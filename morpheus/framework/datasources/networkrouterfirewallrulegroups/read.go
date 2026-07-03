// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkrouterfirewallrulegroups

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	providererrors "github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/utils/convert"
)

const summary = "read network router firewall rule groups data source"

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

		// Description is NullableString — guard with IsSet() and nil-check Get().
		description := types.StringNull()
		if rg.Description.IsSet() && rg.Description.Get() != nil {
			description = types.StringValue(*rg.Description.Get())
		}

		attrs := map[string]attr.Value{
			"id":          convert.Int64ToType(rg.Id),
			"name":        convert.StrToType(rg.Name),
			"description": description,
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
