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

	rule := sdk.NewUpdateNetworkFirewallRuleRequestRule()

	if !plan.Name.IsNull() && !plan.Name.IsUnknown() {
		rule.SetName(plan.Name.ValueString())
	}

	if !plan.Description.IsNull() && !plan.Description.IsUnknown() {
		rule.SetDescription(plan.Description.ValueString())
	}

	if !plan.Enabled.IsNull() && !plan.Enabled.IsUnknown() {
		rule.SetEnabled(plan.Enabled.ValueBool())
	}

	if !plan.Direction.IsNull() && !plan.Direction.IsUnknown() {
		rule.SetDirection(plan.Direction.ValueString())
	}

	if !plan.Policy.IsNull() && !plan.Policy.IsUnknown() {
		rule.SetPolicy(plan.Policy.ValueString())
	}

	if !plan.Priority.IsNull() && !plan.Priority.IsUnknown() {
		rule.SetPriority(plan.Priority.ValueInt64())
	}

	if !plan.Sources.IsNull() && !plan.Sources.IsUnknown() {
		sources := sdk.NewUpdateNetworkFirewallRuleRequestRuleSources()
		sources.SetId(setValueToStringSlice(plan.Sources.Id))
		rule.SetSources(*sources)
	}

	if !plan.Destinations.IsNull() && !plan.Destinations.IsUnknown() {
		destinations := sdk.NewUpdateNetworkFirewallRuleRequestRuleDestinations()
		destinations.SetId(setValueToStringSlice(plan.Destinations.Id))
		rule.SetDestinations(*destinations)
	}

	if !plan.Scopes.IsNull() && !plan.Scopes.IsUnknown() {
		scopes := sdk.NewUpdateNetworkFirewallRuleRequestRuleScopes()
		scopes.SetId(setValueToStringSlice(plan.Scopes.Id))
		rule.SetScopes(*scopes)
	}

	config := sdk.NewUpdateNetworkFirewallRuleRequestRuleConfig()
	configSet := false

	if !plan.Application.IsNull() && !plan.Application.IsUnknown() {
		config.SetApplication(setValueToStringSlice(plan.Application))
		configSet = true
	}

	if !plan.Profile.IsNull() && !plan.Profile.IsUnknown() {
		config.SetProfile(setValueToStringSlice(plan.Profile))
		configSet = true
	}

	if configSet {
		rule.SetConfig(*config)
	}

	if !plan.RuleGroupId.IsNull() && !plan.RuleGroupId.IsUnknown() {
		ruleGroup := sdk.NewUpdateNetworkFirewallRuleRequestRuleRuleGroup()
		if !plan.RuleGroupId.Id.IsNull() && !plan.RuleGroupId.Id.IsUnknown() {
			ruleGroupID := int32(plan.RuleGroupId.Id.ValueInt64()) //nolint:gosec // API uses int32
			ruleGroup.SetId(ruleGroupID)
		}

		rule.SetRuleGroup(*ruleGroup)
	}

	updateReq := sdk.NewUpdateNetworkFirewallRuleRequestWithDefaults()
	updateReq.SetRule(*rule)

	id := currentState.Id.ValueInt64()
	networkIntegrationId := currentState.NetworkIntegrationId.ValueInt64()

	_, httpResp, err := client.NetworksAPI.
		UpdateNetworkFirewallRule(ctx, id, networkIntegrationId).
		UpdateNetworkFirewallRuleRequest(*updateReq).Execute()
	if err != nil || (httpResp != nil && httpResp.StatusCode != http.StatusOK) {
		resp.Diagnostics.AddError(
			"error updating network firewall rule",
			"network firewall rule "+plan.Name.ValueString()+" PUT failed: "+errfmt.ErrMsg(err, httpResp),
		)

		return
	}

	state, _, diag := getNetworkFirewallRuleAsState(ctx, id, networkIntegrationId, client)
	if resp.Diagnostics.Append(diag...); resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
