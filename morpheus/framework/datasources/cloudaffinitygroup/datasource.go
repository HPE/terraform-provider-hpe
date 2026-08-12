// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

// Package cloudaffinitygroup implements a data source for cloud_affinity_group
package cloudaffinitygroup

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"
	"github.com/HPE/terraform-provider-hpe/morpheus/configure"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/constants"
	providererrors "github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/versioncheck"
	"github.com/HPE/terraform-provider-hpe/utils/convert"
)

const (
	summary                          = "read cloud affinity group data source"
	ErrorNoValidSearchTerms          = `no valid search terms - an id or name is required`
	ErrorNoCloudAffinityGroupFound   = `no cloud affinity group found`
	ErrorMultipleCloudAffinityGroups = `multiple cloud affinity groups were returned`

	// gatedFeature names this data source in the appliance version gate
	// diagnostic. Phrased as a plural noun so the message reads "Cloud affinity
	// groups require ...".
	gatedFeature = "Cloud affinity groups"
)

// Ensure the implementation satisfies the expected interfaces.
var _ datasource.DataSource = &DataSource{}

// NewDataSource is a helper function to simplify the provider implementation.
func NewDataSource() datasource.DataSource {
	return &DataSource{}
}

// DataSource is the data source implementation.
type DataSource struct {
	configure.DataSourceWithMorpheusConfigure
	datasource.DataSource
}

// Metadata returns the data source type name.
func (d *DataSource) Metadata(
	_ context.Context,
	req datasource.MetadataRequest,
	resp *datasource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_" + "cloud_affinity_group"
}

// Schema defines the schema for the data source.
func (d *DataSource) Schema(
	ctx context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = CloudAffinityGroupDataSourceSchema(ctx)
}

func affinityGroupAsState(
	ctx context.Context,
	ag *sdk.GetCloudAffinityGroup200ResponseAffinityGroup,
	cloudID int64,
) CloudAffinityGroupModel {
	state := CloudAffinityGroupModel{
		Id:           convert.Int64ToType(ag.Id),
		CloudId:      types.Int64Value(cloudID),
		Name:         convert.StrToType(ag.Name),
		Active:       convert.BoolToType(ag.Active),
		AffinityType: convert.StrToType(ag.AffinityType),
		RefId:        convert.Int64ToType(ag.RefId),
		RefType:      convert.StrToType(ag.RefType),
		Visibility:   convert.StrToType(ag.Visibility),
		Source:       convert.StrToType(ag.Source),
	}

	// Pool — nested {id}
	if ag.Pool != nil {
		state.Pool = NewPoolValueMust(
			PoolValue{}.AttributeTypes(ctx),
			map[string]attr.Value{
				"id": convert.Int64ToType(ag.Pool.Id),
			},
		)
	} else {
		state.Pool = NewPoolValueNull()
	}

	// Servers — [{id, name}] → Set of Int64 (extract IDs only)
	if len(ag.Servers) > 0 {
		serverVals := make([]attr.Value, 0, len(ag.Servers))
		for i := range ag.Servers {
			if ag.Servers[i].Id != nil {
				serverVals = append(serverVals, types.Int64Value(*ag.Servers[i].Id))
			}
		}

		if len(serverVals) > 0 {
			state.Servers, _ = types.SetValue(types.Int64Type, serverVals)
		} else {
			state.Servers = types.SetNull(types.Int64Type)
		}
	} else {
		state.Servers = types.SetNull(types.Int64Type)
	}

	// TenantIds — [{id, name}] → Set of Int64 (extract IDs only)
	if len(ag.Tenants) > 0 {
		tenantVals := make([]attr.Value, 0, len(ag.Tenants))
		for i := range ag.Tenants {
			if ag.Tenants[i].Id != nil {
				tenantVals = append(tenantVals, types.Int64Value(*ag.Tenants[i].Id))
			}
		}

		if len(tenantVals) > 0 {
			state.TenantIds, _ = types.SetValue(types.Int64Type, tenantVals)
		} else {
			state.TenantIds = types.SetNull(types.Int64Type)
		}
	} else {
		state.TenantIds = types.SetNull(types.Int64Type)
	}

	// ResourcePermissions — {all, groups: [{id, default}]}
	// Note: cloud API returns Sites as []map[string]interface{} (untyped)
	if ag.ResourcePermissions != nil {
		groupsVals := make([]attr.Value, 0, len(ag.ResourcePermissions.Sites))
		for _, site := range ag.ResourcePermissions.Sites {
			var id types.Int64
			var dflt types.Bool

			if rawID, ok := site["id"]; ok {
				switch v := rawID.(type) {
				case float64:
					idVal := int64(v)
					id = types.Int64Value(idVal)
				default:
					id = types.Int64Null()
				}
			} else {
				id = types.Int64Null()
			}

			if rawDflt, ok := site["default"]; ok {
				switch v := rawDflt.(type) {
				case bool:
					dflt = types.BoolValue(v)
				default:
					dflt = types.BoolNull()
				}
			} else {
				dflt = types.BoolNull()
			}

			gv := NewGroupsValueMust(
				GroupsValue{}.AttributeTypes(ctx),
				map[string]attr.Value{
					"id":      id,
					"default": dflt,
				},
			)

			groupsVals = append(groupsVals, gv)
		}

		// resource_permissions declares groups with GroupsType as its element
		// type, so the elements must be GroupsValue. Converting them to bare
		// objects first fails the set's element type check.
		groupsSet, _ := types.SetValue(GroupsValue{}.Type(ctx), groupsVals)

		state.ResourcePermissions = NewResourcePermissionsValueMust(
			ResourcePermissionsValue{}.AttributeTypes(ctx),
			map[string]attr.Value{
				"all":    convert.BoolToType(ag.ResourcePermissions.All),
				"groups": groupsSet,
			},
		)
	} else {
		state.ResourcePermissions = NewResourcePermissionsValueNull()
	}

	return state
}

func getAffinityGroupByID(
	ctx context.Context,
	id int64,
	cloudID int64,
	apiClient *sdk.APIClient,
) (*sdk.GetCloudAffinityGroup200ResponseAffinityGroup, error) {
	r, hresp, err := apiClient.CloudsAPI.GetCloudAffinityGroup(
		ctx, cloudID, id,
	).Execute()
	if r == nil || err != nil || hresp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"GET failed for cloud affinity group %d: %s",
			id, providererrors.ErrMsg(err, hresp),
		)
	}

	return r.AffinityGroup, nil
}

func getAffinityGroupByName(
	ctx context.Context,
	name string,
	cloudID int64,
	apiClient *sdk.APIClient,
) (*sdk.GetCloudAffinityGroup200ResponseAffinityGroup, error) {
	rs, hresp, err := apiClient.CloudsAPI.ListCloudAffinityGroups(
		ctx, cloudID,
	).Execute()
	if rs == nil || err != nil || hresp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"GET failed for cloud affinity groups with name %s: %s",
			name, providererrors.ErrMsg(err, hresp),
		)
	}

	items := rs.AffinityGroups
	if len(items) == 0 {
		return nil, errors.New(ErrorNoCloudAffinityGroupFound)
	}

	var matchedIDs []int64

	for i := range items {
		if items[i].Name == nil || *items[i].Name != name {
			continue
		}
		if items[i].Id == nil {
			continue
		}

		matchedIDs = append(matchedIDs, *items[i].Id)
	}

	if len(matchedIDs) == 0 {
		return nil, errors.New(ErrorNoCloudAffinityGroupFound)
	} else if len(matchedIDs) > 1 {
		return nil, errors.New(ErrorMultipleCloudAffinityGroups)
	}

	return getAffinityGroupByID(ctx, matchedIDs[0], cloudID, apiClient)
}

func getAffinityGroup(
	ctx context.Context,
	config *CloudAffinityGroupModel,
	apiClient *sdk.APIClient,
) (*sdk.GetCloudAffinityGroup200ResponseAffinityGroup, error) {
	cloudID := config.CloudId.ValueInt64()

	if !config.Id.IsNull() {
		return getAffinityGroupByID(ctx, config.Id.ValueInt64(), cloudID, apiClient)
	} else if !config.Name.IsNull() {
		return getAffinityGroupByName(ctx, config.Name.ValueString(), cloudID, apiClient)
	}

	return nil, errors.New(ErrorNoValidSearchTerms)
}

// Read refreshes the Terraform state with the latest data.
func (d *DataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	var config CloudAffinityGroupModel

	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiClient, err := d.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError(summary, "could not create sdk client")

		return
	}

	// MORPH-15506: refuse to read against an appliance older than the first
	// release with stable affinity group semantics, so the practitioner gets a
	// diagnostic naming the required version instead of an opaque API error.
	//
	// The check sits in Read rather than Configure because the framework calls
	// Configure on every RPC for the type, including ValidateDataSourceConfig,
	// which should not reach the network. One extra request per Read, on an
	// operation that was going to call the API regardless. See
	// versioncheck.Require, including why an unreadable version fails open.
	resp.Diagnostics.Append(versioncheck.Require(
		ctx, apiClient, gatedFeature, constants.AffinityGroupMinVersion,
	)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ag, err := getAffinityGroup(ctx, &config, apiClient)
	if err != nil {
		resp.Diagnostics.AddError(summary, err.Error())

		return
	}

	cloudID := config.CloudId.ValueInt64()
	state := affinityGroupAsState(ctx, ag, cloudID)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
