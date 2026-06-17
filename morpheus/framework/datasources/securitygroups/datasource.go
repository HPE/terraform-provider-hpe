// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package securitygroups

import (
	"context"
	"fmt"
	"net/http"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/HPE/terraform-provider-hpe/morpheus/configure"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/constants"
	providererrors "github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/utils/convert"
)

const (
	summary = "read security groups data source"
	// listMax bounds the number of records fetched from the API in one call.
	listMax = 250
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
	resp.TypeName = req.ProviderTypeName + "_" + constants.SubProviderName + "_security_groups"
}

// Schema defines the schema for the data source.
func (d *DataSource) Schema(
	ctx context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = SecurityGroupsDataSourceSchema(ctx)
}

// Read refreshes the Terraform state with the latest data.
func (d *DataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	var config SecurityGroupsModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiClient, err := d.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError(summary, "could not create sdk client")

		return
	}

	listReq := apiClient.SecurityGroupsAPI.ListSecurityGroups(ctx).Max(listMax)
	if !config.Name.IsNull() {
		listReq = listReq.Name(config.Name.ValueString())
	}
	if !config.Phrase.IsNull() {
		listReq = listReq.Phrase(config.Phrase.ValueString())
	}

	rs, hresp, err := listReq.Execute()
	if rs == nil || err != nil || hresp.StatusCode != http.StatusOK {
		resp.Diagnostics.AddError(summary, fmt.Sprintf("LIST failed for security groups: %s",
			providererrors.ErrMsg(err, hresp)))

		return
	}

	elemType := types.ObjectType{AttrTypes: securityGroupObjectAttrTypes()}
	objs := make([]attr.Value, 0, len(rs.SecurityGroups))

	for i := range rs.SecurityGroups {
		sg := rs.SecurityGroups[i]
		if !matchesClientFilters(&config, &sg) {
			continue
		}

		obj, diags := securityGroupInnerToObject(&sg)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		objs = append(objs, obj)
	}

	listVal, diags := types.ListValue(elemType, objs)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	config.SecurityGroups = listVal

	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}

// matchesClientFilters applies the filters that the list API does not support
// natively (cloud_id, visibility, active).
func matchesClientFilters(
	config *SecurityGroupsModel,
	sg *sdk.ListSecurityGroups200ResponseAllOfSecurityGroupsInner,
) bool {
	if !config.CloudId.IsNull() {
		if sg.Zone == nil || sg.Zone.Id == nil || *sg.Zone.Id != config.CloudId.ValueInt64() {
			return false
		}
	}

	if !config.Visibility.IsNull() {
		if sg.Visibility == nil || *sg.Visibility != config.Visibility.ValueString() {
			return false
		}
	}

	if !config.Active.IsNull() {
		if sg.Active == nil || *sg.Active != config.Active.ValueBool() {
			return false
		}
	}

	return true
}

func securityGroupInnerToObject(
	sg *sdk.ListSecurityGroups200ResponseAllOfSecurityGroupsInner,
) (types.Object, diag.Diagnostics) {
	cloudID := types.Int64Null()
	if sg.Zone != nil {
		cloudID = convert.Int64ToType(sg.Zone.Id)
	}

	groupsAll := types.BoolNull()
	groupIDs := types.SetNull(types.Int64Type)
	if sg.ResourcePermission != nil {
		groupsAll = convert.BoolToType(sg.ResourcePermission.All)
		groupIDs = siteIDsToSet(sg.ResourcePermission.Sites)
	}

	attrs := map[string]attr.Value{
		"account_id":                     convert.Int64ToType(sg.AccountId),
		"active":                         convert.BoolToType(sg.Active),
		"cloud_id":                       cloudID,
		"description":                    convert.StrToType(sg.Description.Get()),
		"enabled":                        convert.StrToType(sg.Enabled.Get()),
		"external_id":                    convert.StrToType(sg.ExternalId.Get()),
		"group_source":                   convert.StrToType(sg.GroupSource.Get()),
		"id":                             convert.Int64ToType(sg.Id),
		"name":                           convert.StrToType(sg.Name),
		"resource_permission_group_ids":  groupIDs,
		"resource_permission_groups_all": groupsAll,
		"sync_source":                    convert.StrToType(sg.SyncSource),
		"tenant_ids":                     tenantIDsToSet(sg.Tenants),
		"visibility":                     convert.StrToType(sg.Visibility),
	}

	return types.ObjectValue(securityGroupObjectAttrTypes(), attrs)
}

func tenantIDsToSet(
	tenants []sdk.ListSecurityGroups200ResponseAllOfSecurityGroupsInnerTenantsInner,
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
	sites []sdk.ListCloudFolders200ResponseAllOfFoldersInnerResourcePermissionsSitesInner,
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
