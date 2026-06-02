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

	rule := &sdk.UpdateNetworkFirewallRuleRequestRule{}

	if !plan.Name.IsNull() && !plan.Name.IsUnknown() {
		name := plan.Name.ValueString()
		rule.Name = &name
	}

	if !plan.Description.IsNull() && !plan.Description.IsUnknown() {
		desc := plan.Description.ValueString()
		rule.Description.Set(&desc)
	}

	if !plan.Enabled.IsNull() && !plan.Enabled.IsUnknown() {
		enabled := plan.Enabled.ValueBool()
		rule.Enabled = &enabled
	}

	if !plan.Direction.IsNull() && !plan.Direction.IsUnknown() {
		direction := plan.Direction.ValueString()
		rule.Direction = &direction
	}

	if !plan.Policy.IsNull() && !plan.Policy.IsUnknown() {
		policy := plan.Policy.ValueString()
		rule.Policy = &policy
	}

	if !plan.Priority.IsNull() && !plan.Priority.IsUnknown() {
		priority := plan.Priority.ValueInt64()
		rule.Priority.Set(&priority)
	}

	if !plan.Sources.IsNull() && !plan.Sources.IsUnknown() {
		sources := &sdk.UpdateNetworkFirewallRuleRequestRuleSources{
			Id: setValueToStringSlice(plan.Sources.Id),
		}
		rule.Sources = sources
	}

	if !plan.Destinations.IsNull() && !plan.Destinations.IsUnknown() {
		destinations := &sdk.UpdateNetworkFirewallRuleRequestRuleDestinations{
			Id: setValueToStringSlice(plan.Destinations.Id),
		}
		rule.Destinations = destinations
	}

	if !plan.Scopes.IsNull() && !plan.Scopes.IsUnknown() {
		scopes := &sdk.UpdateNetworkFirewallRuleRequestRuleScopes{
			Id: setValueToStringSlice(plan.Scopes.Id),
		}
		rule.Scopes = scopes
	}

	config := &sdk.UpdateNetworkFirewallRuleRequestRuleConfig{}
	configSet := false

	if !plan.Application.IsNull() && !plan.Application.IsUnknown() {
		config.Application = setValueToStringSlice(plan.Application)
		configSet = true
	}

	if !plan.Profile.IsNull() && !plan.Profile.IsUnknown() {
		config.Profile = setValueToStringSlice(plan.Profile)
		configSet = true
	}

	if configSet {
		rule.Config = config
	}

	if !plan.RuleGroupId.IsNull() && !plan.RuleGroupId.IsUnknown() {
		ruleGroup := &sdk.UpdateNetworkFirewallRuleRequestRuleRuleGroup{}
		if !plan.RuleGroupId.Id.IsNull() && !plan.RuleGroupId.Id.IsUnknown() {
			ruleGroupID := int32(plan.RuleGroupId.Id.ValueInt64()) //nolint:gosec // API uses int32
			ruleGroup.Id = &ruleGroupID
		}

		rule.RuleGroup = ruleGroup
	}

	updateReq := &sdk.UpdateNetworkFirewallRuleRequest{
		Rule: rule,
	}

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
