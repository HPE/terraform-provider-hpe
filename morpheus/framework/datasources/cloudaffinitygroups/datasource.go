// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

// Package cloudaffinitygroups implements a data source for cloud_affinity_groups
package cloudaffinitygroups

import (
	"context"
	"fmt"
	"net/http"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"
	"github.com/HPE/terraform-provider-hpe/morpheus/configure"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/constants"
	providererrors "github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/versioncheck"
	"github.com/HPE/terraform-provider-hpe/utils/convert"
)

const (
	summary = "read cloud affinity groups data source"

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
	resp.TypeName = req.ProviderTypeName + "_" + "cloud_affinity_groups"
}

// Schema defines the schema for the data source.
func (d *DataSource) Schema(
	ctx context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = CloudAffinityGroupsDataSourceSchema(ctx)
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
	var config CloudAffinityGroupsModel

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

	cloudID := config.CloudId.ValueInt64()

	// MORPH-15506: refuse to read against an appliance older than the first
	// release with stable affinity group semantics, so the practitioner gets a
	// diagnostic naming the required version instead of an opaque API error.
	//
	// The check sits in Read rather than Configure because the framework calls
	// Configure on every RPC for the type, including ValidateDataSourceConfig,
	// which should not reach the network. That does mean one extra request per
	// Read of this plural data source, but a Read is a single LIST call so the
	// gate at worst doubles a cheap operation, and it is per data source block
	// rather than per group returned. See versioncheck.Require, including why
	// an unreadable version fails open.
	resp.Diagnostics.Append(versioncheck.Require(
		ctx, apiClient, gatedFeature, constants.AffinityGroupMinVersion,
	)...)
	if resp.Diagnostics.HasError() {
		return
	}

	rs, hresp, err := apiClient.CloudsAPI.ListCloudAffinityGroups(ctx, cloudID).Execute()
	if rs == nil || err != nil || hresp.StatusCode != http.StatusOK {
		resp.Diagnostics.AddError(summary, fmt.Sprintf(
			"LIST failed for cloud affinity groups: %s",
			providererrors.ErrMsg(err, hresp)))

		return
	}

	objs := make([]attr.Value, 0, len(rs.AffinityGroups))

	for i := range rs.AffinityGroups {
		ag := &rs.AffinityGroups[i]
		if !matchesFilters(ag, filters) {
			continue
		}

		v, diags := affinityGroupToValue(ctx, ag)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		objs = append(objs, v)
	}

	// The generated schema declares the element type of affinity_groups as the
	// custom AffinityGroupsType, so both the element type used here and every
	// element must be the generated type. A bare types.Object element fails the
	// set's element type check with "Invalid Set Element Type".
	setVal, diags := types.SetValue(AffinityGroupsValue{}.Type(ctx), objs)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	config.AffinityGroups = setVal

	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}

func compileFilters(
	ctx context.Context,
	filterSet types.Set,
	diags *diag.Diagnostics,
) []compiledFilter {
	if filterSet.IsNull() || filterSet.IsUnknown() {
		return nil
	}

	var filterBlocks []FilterValue
	diags.Append(filterSet.ElementsAs(ctx, &filterBlocks, false)...)
	if diags.HasError() {
		return nil
	}

	filters := make([]compiledFilter, 0, len(filterBlocks))

	for _, b := range filterBlocks {
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

func matchesFilters(
	ag *sdk.ListCloudAffinityGroups200ResponseAllOfAffinityGroupsInner,
	filters []compiledFilter,
) bool {
	for _, f := range filters {
		val, ok := fieldValue(ag, f.field)
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

func fieldValue(
	ag *sdk.ListCloudAffinityGroups200ResponseAllOfAffinityGroupsInner,
	field string,
) (string, bool) {
	switch field {
	case "name":
		if ag.Name != nil {
			return *ag.Name, true
		}
	case "affinity_type":
		if ag.AffinityType != nil {
			return *ag.AffinityType, true
		}
	case "source":
		if ag.Source != nil {
			return *ag.Source, true
		}
	case "visibility":
		if ag.Visibility != nil {
			return *ag.Visibility, true
		}
	}

	return "", false
}

// affinityGroupToValue maps an API affinity group into the generated custom
// object value used as the element of the affinity_groups set. It must return
// AffinityGroupsValue rather than a bare object: the schema declares the set's
// element type as AffinityGroupsType, and a set rejects any element whose type
// is not exactly the declared element type.
func affinityGroupToValue(
	ctx context.Context,
	ag *sdk.ListCloudAffinityGroups200ResponseAllOfAffinityGroupsInner,
) (AffinityGroupsValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Pool. AffinityGroupsValue declares pool as a plain object type, so this
	// stays a types.Object rather than a PoolValue.
	poolVal := types.ObjectNull(PoolValue{}.AttributeTypes(ctx))

	if ag.Pool != nil {
		v, d := types.ObjectValue(
			PoolValue{}.AttributeTypes(ctx),
			map[string]attr.Value{
				"id": convert.Int64ToType(ag.Pool.Id),
			},
		)
		diags.Append(d...)
		poolVal = v
	}

	// Servers
	serversVal := types.SetNull(types.Int64Type)

	if len(ag.Servers) > 0 {
		vals := make([]attr.Value, 0, len(ag.Servers))
		for i := range ag.Servers {
			if ag.Servers[i].Id != nil {
				vals = append(vals, types.Int64Value(*ag.Servers[i].Id))
			}
		}

		if len(vals) > 0 {
			v, d := types.SetValue(types.Int64Type, vals)
			diags.Append(d...)
			serversVal = v
		}
	}

	// TenantIds
	tenantIdsVal := types.SetNull(types.Int64Type)

	if len(ag.Tenants) > 0 {
		vals := make([]attr.Value, 0, len(ag.Tenants))
		for i := range ag.Tenants {
			if ag.Tenants[i].Id != nil {
				vals = append(vals, types.Int64Value(*ag.Tenants[i].Id))
			}
		}

		if len(vals) > 0 {
			v, d := types.SetValue(types.Int64Type, vals)
			diags.Append(d...)
			tenantIdsVal = v
		}
	}

	// ResourcePermissions — cloud list uses untyped Sites []map[string]interface{}
	rpVal := types.ObjectNull(ResourcePermissionsValue{}.AttributeTypes(ctx))

	if ag.ResourcePermissions != nil {
		v, d := resourcePermissionsToObject(ctx, ag.ResourcePermissions)
		diags.Append(d...)
		rpVal = v
	}

	v, d := NewAffinityGroupsValue(
		AffinityGroupsValue{}.AttributeTypes(ctx),
		map[string]attr.Value{
			"active":               convert.BoolToType(ag.Active),
			"affinity_type":        convert.StrToType(ag.AffinityType),
			"id":                   convert.Int64ToType(ag.Id),
			"name":                 convert.StrToType(ag.Name),
			"pool":                 poolVal,
			"ref_id":               convert.Int64ToType(ag.RefId),
			"ref_type":             convert.StrToType(ag.RefType),
			"resource_permissions": rpVal,
			"servers":              serversVal,
			"source":               convert.StrToType(ag.Source),
			"tenant_ids":           tenantIdsVal,
			"visibility":           convert.StrToType(ag.Visibility),
		},
	)
	diags.Append(d...)

	return v, diags
}

// resourcePermissionsToObject maps the resource permissions of a cloud affinity
// group. The groups set is declared with GroupsType as its element type, so its
// elements are built with NewGroupsValue rather than as bare objects.
func resourcePermissionsToObject(
	ctx context.Context,
	rp *sdk.ListCloudAffinityGroups200ResponseAllOfAffinityGroupsInnerResourcePermissions,
) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	groupVals := make([]attr.Value, 0, len(rp.Sites))

	for _, site := range rp.Sites {
		id := types.Int64Null()
		if v, ok := site["id"].(float64); ok {
			id = types.Int64Value(int64(v))
		}

		dflt := types.BoolNull()
		if v, ok := site["default"].(bool); ok {
			dflt = types.BoolValue(v)
		}

		gv, d := NewGroupsValue(
			GroupsValue{}.AttributeTypes(ctx),
			map[string]attr.Value{
				"id":      id,
				"default": dflt,
			},
		)
		diags.Append(d...)
		groupVals = append(groupVals, gv)
	}

	groupsSet, d := types.SetValue(GroupsValue{}.Type(ctx), groupVals)
	diags.Append(d...)

	obj, d := types.ObjectValue(
		ResourcePermissionsValue{}.AttributeTypes(ctx),
		map[string]attr.Value{
			"all":    convert.BoolToType(rp.All),
			"groups": groupsSet,
		},
	)
	diags.Append(d...)

	return obj, diags
}
