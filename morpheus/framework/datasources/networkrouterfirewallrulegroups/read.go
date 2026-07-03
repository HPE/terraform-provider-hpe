// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkrouterfirewallrulegroups

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"
	providererrors "github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/utils/convert"
)

const summary = "read network router firewall rule groups data source"

// filterModel decodes a single filter block.
type filterModel struct {
	Name   types.String `tfsdk:"name"`
	Values types.Set    `tfsdk:"values"`
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
	var model NetworkRouterFirewallRuleGroupsModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &model)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var blocks []filterModel
	resp.Diagnostics.Append(model.Filter.ElementsAs(ctx, &blocks, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	filters := compileFilters(ctx, blocks, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	routerID := model.RouterId.ValueInt64()

	apiClient, err := d.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError(summary, "could not create sdk client")

		return
	}

	rs, hresp, err := apiClient.NetworksAPI.GetNetworkRouterFirewallRuleGroups(ctx, routerID).Execute()
	if rs == nil || err != nil || hresp.StatusCode != http.StatusOK {
		resp.Diagnostics.AddError(summary, fmt.Sprintf(
			"LIST failed for network router firewall rule groups (router_id=%d): %s",
			routerID, providererrors.ErrMsg(err, hresp)))

		return
	}

	objs := make([]attr.Value, 0, len(rs.RuleGroups))

	for i := range rs.RuleGroups {
		rg := rs.RuleGroups[i]
		if !ruleGroupMatchesFilters(&rg, filters) {
			continue
		}

		attrs := map[string]attr.Value{
			"id":          convert.Int64ToType(rg.Id),
			"name":        convert.StrToType(rg.Name),
			"description": convert.StrToType(rg.Description.Get()),
			"external_id": convert.StrToType(rg.ExternalId),
			"status":      convert.StrToType(rg.Status),
			"priority":    convert.Int64ToType(rg.Priority),
			"group_layer": convert.StrToType(rg.GroupLayer),
		}

		v, diags := NewRuleGroupsValue(RuleGroupsValue{}.AttributeTypes(ctx), attrs)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		objs = append(objs, v)
	}

	setVal, diags := types.SetValue(RuleGroupsValue{}.Type(ctx), objs)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	model.RuleGroups = setVal

	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
}

// compileFilters converts the filter blocks from configuration into compiled
// regular expressions. Invalid patterns are reported as diagnostics.
func compileFilters(
	ctx context.Context,
	blocks []filterModel,
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

// ruleGroupMatchesFilters reports whether rg satisfies every filter block.
// Within a block, the field must match ANY value (OR); across blocks all must
// match (AND).
func ruleGroupMatchesFilters(
	rg *sdk.GetNetworkRouterFirewallRuleGroups200ResponseRuleGroupsInner,
	filters []compiledFilter,
) bool {
	for _, f := range filters {
		val, ok := ruleGroupFieldValue(rg, f.field)
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

// ruleGroupFieldValue returns the string representation of the named filter
// field for regex matching, and whether the field is present.
func ruleGroupFieldValue(
	rg *sdk.GetNetworkRouterFirewallRuleGroups200ResponseRuleGroupsInner,
	field string,
) (string, bool) {
	switch field {
	case "name":
		if rg.Name != nil {
			return *rg.Name, true
		}
	case "id":
		if rg.Id != nil {
			return strconv.FormatInt(*rg.Id, 10), true
		}
	case "external_id":
		if rg.ExternalId != nil {
			return *rg.ExternalId, true
		}
	case "status":
		if rg.Status != nil {
			return *rg.Status, true
		}
	case "priority":
		if rg.Priority != nil {
			return strconv.FormatInt(*rg.Priority, 10), true
		}
	case "group_layer":
		if rg.GroupLayer != nil {
			return *rg.GroupLayer, true
		}
	}

	return "", false
}
