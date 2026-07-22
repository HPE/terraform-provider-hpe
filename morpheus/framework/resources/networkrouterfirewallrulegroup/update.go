// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkrouterfirewallrulegroup

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
)

const updateOperation = "update network router firewall rule group resource"

func (r *Resource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	var plan NetworkRouterFirewallRuleGroupModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client, err := r.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			updateOperation,
			"failed to create client: "+err.Error(),
		)

		return
	}

	id := plan.Id.ValueInt64()
	routerID := plan.RouterId.ValueInt64()

	// The SDK's UpdateNetworkRouterFirewallRuleGroupRequest.RuleGroup is typed
	// as map[string]interface{} — build it manually.
	ruleGroup := map[string]interface{}{
		"name":         plan.Name.ValueString(),
		"externalType": "GatewayPolicy",
	}

	if !plan.Description.IsNull() && !plan.Description.IsUnknown() {
		ruleGroup["description"] = plan.Description.ValueString()
	}

	if !plan.Priority.IsUnknown() {
		ruleGroup["priority"] = plan.Priority.ValueInt64()
	}

	if !plan.GroupLayer.IsNull() && !plan.GroupLayer.IsUnknown() {
		ruleGroup["groupLayer"] = plan.GroupLayer.ValueString()
	}

	if !plan.Visibility.IsNull() && !plan.Visibility.IsUnknown() {
		ruleGroup["visibility"] = plan.Visibility.ValueString()
	}

	if !plan.TenantIds.IsUnknown() {
		var ids []int64
		if !plan.TenantIds.IsNull() {
			resp.Diagnostics.Append(plan.TenantIds.ElementsAs(ctx, &ids, false)...)
			if resp.Diagnostics.HasError() {
				return
			}
		}

		tenants := make([]map[string]interface{}, 0, len(ids))
		for _, tid := range ids {
			tenants = append(tenants, map[string]interface{}{"id": tid})
		}

		ruleGroup["tenants"] = tenants
	}

	updateReq := sdk.UpdateNetworkRouterFirewallRuleGroupRequest{
		RuleGroup: ruleGroup,
	}

	_, hresp, err := client.NetworksAPI.
		UpdateNetworkRouterFirewallRuleGroup(ctx, id, routerID).
		UpdateNetworkRouterFirewallRuleGroupRequest(updateReq).Execute()
	if err != nil || hresp.StatusCode != http.StatusOK {
		resp.Diagnostics.AddError(
			updateOperation,
			fmt.Sprintf("firewall rule group %d PUT failed: %s", id, errfmt.ErrMsg(err, hresp)),
		)

		return
	}

	state, notFound, diags := getRuleGroupAsState(ctx, id, routerID, client, plan)
	if notFound {
		// Unexpected: resource existed during PUT but GET returned 404.
		resp.Diagnostics.AddError(
			updateOperation,
			fmt.Sprintf("firewall rule group %d was updated but GET returned 404", id),
		)

		return
	}

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
