// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package securitygroup

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/HPE/terraform-provider-hpe/morpheus/configure"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/constants"
	providererrors "github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/utils/convert"
)

const (
	summary                     = "read security group data source"
	ErrorNoValidSearchTerms     = `no valid search terms - an id or name is required`
	ErrorRunningPreApply        = `Error running pre-apply plan: exit status 1`
	ErrorNoSecurityGroupFound   = `no security group found`
	ErrorMultipleSecurityGroups = `multiple security groups were returned`
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
	resp.TypeName = req.ProviderTypeName + "_" + constants.SubProviderName + "_security_group"
}

// Schema defines the schema for the data source.
func (d *DataSource) Schema(
	ctx context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = SecurityGroupDataSourceSchema(ctx)
	resp.Schema.Description = "Retrieves information about a single Morpheus security group by id or name."
	resp.Schema.MarkdownDescription = "Retrieves information about a single Morpheus security group by id or name."
}

// securityGroupAsState maps an API security group into the datasource model.
func securityGroupAsState(
	sg *sdk.GetSecurityGroups200ResponseSecurityGroup,
) SecurityGroupModel {
	state := SecurityGroupModel{
		Id:          convert.Int64ToType(sg.Id),
		Name:        convert.StrToType(sg.Name),
		Description: convert.StrToType(sg.Description.Get()),
		AccountId:   convert.Int64ToType(sg.AccountId),
		GroupSource: convert.StrToType(sg.GroupSource.Get()),
		ExternalId:  convert.StrToType(sg.ExternalId.Get()),
		Enabled:     convert.StrToType(sg.Enabled.Get()),
		SyncSource:  convert.StrToType(sg.SyncSource),
		Visibility:  convert.StrToType(sg.Visibility),
		Active:      convert.BoolToType(sg.Active),
	}

	if sg.Zone != nil {
		state.CloudId = convert.Int64ToType(sg.Zone.Id)
	} else {
		state.CloudId = types.Int64Null()
	}

	state.TenantIds = tenantIDsToSet(sg.Tenants)

	if sg.ResourcePermission != nil {
		state.ResourcePermissionGroupsAll = convert.BoolToType(sg.ResourcePermission.All)
		state.ResourcePermissionGroupIds = siteIDsToSet(sg.ResourcePermission.Sites)
	} else {
		state.ResourcePermissionGroupsAll = types.BoolNull()
		state.ResourcePermissionGroupIds = types.SetNull(types.Int64Type)
	}

	return state
}

func tenantIDsToSet(
	tenants []sdk.AddSecurityGroups200ResponseSecurityGroupAllOfTenantsInner,
) types.Set {
	if len(tenants) == 0 {
		return types.SetNull(types.Int64Type)
	}

	values := make([]attr.Value, 0, len(tenants))
	for _, t := range tenants {
		if t.Id != nil {
			values = append(values, types.Int64Value(*t.Id))
		}
	}
	if len(values) == 0 {
		return types.SetNull(types.Int64Type)
	}

	set, diags := types.SetValue(types.Int64Type, values)
	if diags.HasError() {
		return types.SetNull(types.Int64Type)
	}

	return set
}

func siteIDsToSet(
	sites []sdk.AddSecurityGroups200ResponseSecurityGroupAllOfResourcePermissionSitesInner,
) types.Set {
	if len(sites) == 0 {
		return types.SetNull(types.Int64Type)
	}

	values := make([]attr.Value, 0, len(sites))
	for _, s := range sites {
		if s.Id != nil {
			values = append(values, types.Int64Value(*s.Id))
		}
	}
	if len(values) == 0 {
		return types.SetNull(types.Int64Type)
	}

	set, diags := types.SetValue(types.Int64Type, values)
	if diags.HasError() {
		return types.SetNull(types.Int64Type)
	}

	return set
}

func getSecurityGroupByID(
	ctx context.Context,
	id int64,
	apiClient *sdk.APIClient,
) (*sdk.GetSecurityGroups200ResponseSecurityGroup, error) {
	r, hresp, err := apiClient.SecurityGroupsAPI.GetSecurityGroups(ctx, id).Execute()
	if r == nil || err != nil || hresp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GET failed for security group %d: %s", id, providererrors.ErrMsg(err, hresp))
	}

	if r.SecurityGroup == nil {
		return nil, fmt.Errorf("GET failed for security group %d: response missing securityGroup", id)
	}

	return r.SecurityGroup, nil
}

func getSecurityGroupByName(
	ctx context.Context,
	name string,
	apiClient *sdk.APIClient,
) (*sdk.GetSecurityGroups200ResponseSecurityGroup, error) {
	rs, hresp, err := apiClient.SecurityGroupsAPI.ListSecurityGroups(ctx).Name(name).Execute()
	if rs == nil || err != nil || hresp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GET failed for security group %s: %s", name, providererrors.ErrMsg(err, hresp))
	}

	var matchedID int64

	var matchCount int

	for _, sg := range rs.SecurityGroups {
		if sg.Name != nil && *sg.Name == name {
			if sg.Id != nil {
				matchedID = *sg.Id
			}

			matchCount++
		}
	}

	if matchCount == 0 {
		return nil, errors.New(ErrorNoSecurityGroupFound)
	} else if matchCount > 1 {
		return nil, errors.New(ErrorMultipleSecurityGroups)
	}

	return getSecurityGroupByID(ctx, matchedID, apiClient)
}

func getSecurityGroup(
	ctx context.Context,
	config *SecurityGroupModel,
	apiClient *sdk.APIClient,
) (*sdk.GetSecurityGroups200ResponseSecurityGroup, error) {
	if !config.Id.IsNull() {
		return getSecurityGroupByID(ctx, config.Id.ValueInt64(), apiClient)
	} else if !config.Name.IsNull() {
		return getSecurityGroupByName(ctx, config.Name.ValueString(), apiClient)
	}

	return nil, errors.New(ErrorNoValidSearchTerms)
}

// Read refreshes the Terraform state with the latest data.
func (d *DataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	var config SecurityGroupModel

	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiClient, err := d.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			summary,
			"could not create sdk client",
		)

		return
	}

	sg, err := getSecurityGroup(ctx, &config, apiClient)
	if err != nil {
		resp.Diagnostics.AddError(
			summary,
			err.Error(),
		)

		return
	}

	state := securityGroupAsState(sg)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
