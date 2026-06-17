// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package securitygroups

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strconv"

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

// compiledFilter is a filter block with its values pre-compiled as regular
// expressions.
type compiledFilter struct {
	field string
	res   []*regexp.Regexp
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

	filters := compileFilters(ctx, config.Filter, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	apiClient, err := d.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError(summary, "could not create sdk client")

		return
	}

	// Sort ascending by default; only descending when explicitly set to false.
	direction := "asc"
	if !config.SortAscending.IsNull() && !config.SortAscending.ValueBool() {
		direction = "desc"
	}

	rs, hresp, err := apiClient.SecurityGroupsAPI.ListSecurityGroups(ctx).
		Max(listMax).
		Sort("id").
		Direction(direction).
		Execute()
	if rs == nil || err != nil || hresp.StatusCode != http.StatusOK {
		resp.Diagnostics.AddError(summary, fmt.Sprintf("LIST failed for security groups: %s",
			providererrors.ErrMsg(err, hresp)))

		return
	}

	elemType := types.ObjectType{AttrTypes: securityGroupObjectAttrTypes()}
	objs := make([]attr.Value, 0, len(rs.SecurityGroups))

	for i := range rs.SecurityGroups {
		sg := rs.SecurityGroups[i]
		if !securityGroupMatchesFilters(&sg, filters) {
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

// compileFilters converts the filter blocks from configuration into compiled
// regular expressions. Invalid patterns are reported as diagnostics.
func compileFilters(
	ctx context.Context,
	blocks []securityGroupsFilterModel,
	diags *diag.Diagnostics,
) []compiledFilter {
	filters := make([]compiledFilter, 0, len(blocks))

	for _, b := range blocks {
		field := b.Name.ValueString()

		var values []string
		diags.Append(b.Values.ElementsAs(ctx, &values, false)...)
		if diags.HasError() {
			return nil
		}

		res := make([]*regexp.Regexp, 0, len(values))
		for _, v := range values {
			re, err := regexp.Compile(v)
			if err != nil {
				diags.AddError(summary,
					fmt.Sprintf("invalid regular expression %q for filter %q: %s", v, field, err))

				return nil
			}
			res = append(res, re)
		}

		filters = append(filters, compiledFilter{field: field, res: res})
	}

	return filters
}

// securityGroupMatchesFilters reports whether sg satisfies every filter block.
// Within a block, the field must match ANY value (OR); across blocks all must
// match (AND).
func securityGroupMatchesFilters(
	sg *sdk.ListSecurityGroups200ResponseAllOfSecurityGroupsInner,
	filters []compiledFilter,
) bool {
	for _, f := range filters {
		val, ok := securityGroupFieldValue(sg, f.field)
		if !ok {
			return false
		}

		matched := false
		for _, re := range f.res {
			if re.MatchString(val) {
				matched = true

				break
			}
		}
		if !matched {
			return false
		}
	}

	return true
}

// securityGroupFieldValue returns the string representation of the named field
// for regex matching, and whether the field is present.
func securityGroupFieldValue(
	sg *sdk.ListSecurityGroups200ResponseAllOfSecurityGroupsInner,
	field string,
) (string, bool) {
	switch field {
	case filterFieldName:
		if sg.Name != nil {
			return *sg.Name, true
		}
	case filterFieldVisibility:
		if sg.Visibility != nil {
			return *sg.Visibility, true
		}
	case filterFieldCloudID:
		if sg.Zone != nil && sg.Zone.Id != nil {
			return strconv.FormatInt(*sg.Zone.Id, 10), true
		}
	case filterFieldActive:
		if sg.Active != nil {
			return strconv.FormatBool(*sg.Active), true
		}
	}

	return "", false
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
