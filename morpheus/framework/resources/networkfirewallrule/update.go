// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkfirewallrule

import (
	"context"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
)

func (r *Resource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	var plan, currentState NetworkFirewallRuleModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(req.State.Get(ctx, &currentState)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client, err := r.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError("creating client failed", err.Error())

		return
	}

	ruleMap := map[string]interface{}{}

	if !plan.Name.IsNull() && !plan.Name.IsUnknown() {
		ruleMap["name"] = plan.Name.ValueString()
	}

	if !plan.Description.IsNull() && !plan.Description.IsUnknown() {
		ruleMap["description"] = plan.Description.ValueString()
	}

	if !plan.Enabled.IsNull() && !plan.Enabled.IsUnknown() {
		ruleMap["enabled"] = plan.Enabled.ValueBool()
	}

	if !plan.Direction.IsNull() && !plan.Direction.IsUnknown() {
		ruleMap["direction"] = plan.Direction.ValueString()
	}

	if !plan.Policy.IsNull() && !plan.Policy.IsUnknown() {
		ruleMap["policy"] = plan.Policy.ValueString()
	}

	if !plan.Priority.IsNull() && !plan.Priority.IsUnknown() {
		ruleMap["priority"] = plan.Priority.ValueString()
	}

	if !plan.Sources.IsNull() && !plan.Sources.IsUnknown() {
		ruleMap["sources"] = map[string]interface{}{
			"id": setValueToStringSlice(plan.Sources.Id),
		}
	}

	if !plan.Destinations.IsNull() && !plan.Destinations.IsUnknown() {
		ruleMap["destinations"] = map[string]interface{}{
			"id": setValueToStringSlice(plan.Destinations.Id),
		}
	}

	if !plan.Scopes.IsNull() && !plan.Scopes.IsUnknown() {
		ruleMap["scopes"] = map[string]interface{}{
			"id": setValueToStringSlice(plan.Scopes.Id),
		}
	}

	if !plan.Config.IsNull() && !plan.Config.IsUnknown() {
		configMap := map[string]interface{}{}
		if !plan.Config.Application.IsNull() && !plan.Config.Application.IsUnknown() {
			configMap["application"] = setValueToStringSlice(plan.Config.Application)
		}
		if !plan.Config.Profile.IsNull() && !plan.Config.Profile.IsUnknown() {
			configMap["profile"] = setValueToStringSlice(plan.Config.Profile)
		}

		ruleMap["config"] = configMap
	}

	if !plan.RuleGroupId.IsNull() && !plan.RuleGroupId.IsUnknown() {
		ruleGroupMap := map[string]interface{}{}
		if !plan.RuleGroupId.Id.IsNull() && !plan.RuleGroupId.Id.IsUnknown() {
			ruleGroupMap["id"] = plan.RuleGroupId.Id.ValueInt64()
		}

		ruleMap["ruleGroup"] = ruleGroupMap
	}

	updateReq := sdk.NewUpdateNetworkFirewallRuleRequestWithDefaults()
	updateReq.SetRule(ruleMap)

	id := currentState.Id.ValueInt64()
	serverId := currentState.NetworkServerId.ValueInt64()

	_, httpResp, err := client.NetworksAPI.
		UpdateNetworkFirewallRule(ctx, id, serverId).
		UpdateNetworkFirewallRuleRequest(*updateReq).Execute()
	if err != nil || (httpResp != nil && httpResp.StatusCode != http.StatusOK) {
		resp.Diagnostics.AddError(
			"error updating network firewall rule",
			"network firewall rule "+plan.Name.ValueString()+" PUT failed: "+errfmt.ErrMsg(err, httpResp),
		)

		return
	}

	state, diag := getNetworkFirewallRuleAsState(ctx, id, serverId, client)
	if resp.Diagnostics.Append(diag...); resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
