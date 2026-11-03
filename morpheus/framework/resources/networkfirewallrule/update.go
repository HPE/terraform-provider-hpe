// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkfirewallrule

import (
	"context"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/clearable"
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

	// Attributes removed from the configuration since the last apply, to be
	// cleared. Read before Record, below, overwrites what it reads.
	pending, pdiags := clearable.Compute(ctx, req.Config, req.State, req.Private, clearableAttributes)
	resp.Diagnostics.Append(pdiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	rule := &sdk.UpdateNetworkFirewallRuleRequestRule{}

	if !plan.Name.IsNull() && !plan.Name.IsUnknown() {
		rule.Name = plan.Name.ValueStringPointer()
	}

	clearable.NullableString(pending, "description", plan.Description, &rule.Description)

	if !plan.Enabled.IsNull() && !plan.Enabled.IsUnknown() {
		rule.Enabled = plan.Enabled.ValueBoolPointer()
	}

	if !plan.Direction.IsNull() && !plan.Direction.IsUnknown() {
		rule.Direction = plan.Direction.ValueStringPointer()
	}

	if !plan.Policy.IsNull() && !plan.Policy.IsUnknown() {
		rule.Policy = plan.Policy.ValueStringPointer()
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

	// application and profile: the planned list, or an empty list when being
	// cleared (an empty list matches any application or profile)
	if apps, ok, ldiags := pending.Strings(ctx, "application", plan.Application); ok {
		resp.Diagnostics.Append(ldiags...)
		if resp.Diagnostics.HasError() {
			return
		}
		config.Application = apps
		configSet = true
	}

	if profiles, ok, ldiags := pending.Strings(ctx, "profile", plan.Profile); ok {
		resp.Diagnostics.Append(ldiags...)
		if resp.Diagnostics.HasError() {
			return
		}
		config.Profile = profiles
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

	resp.Diagnostics.Append(clearable.Record(ctx, req.Config, resp.Private, clearableAttributes)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
