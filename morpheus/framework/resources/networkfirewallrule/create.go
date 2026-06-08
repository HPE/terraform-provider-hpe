// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkfirewallrule

import (
	"context"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/utils/cleanup"
)

func (r *Resource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	var plan NetworkFirewallRuleModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client, err := r.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError("creating client failed", err.Error())

		return
	}

	rule := &sdk.CreateNetworkFirewallRuleRequestRule{
		Name: plan.Name.ValueString(),
	}

	if !plan.Description.IsNull() && !plan.Description.IsUnknown() {
		desc := plan.Description.ValueString()
		rule.Description.Set(&desc)
	}

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
		sources := &sdk.CreateNetworkFirewallRuleRequestRuleSources{
			Id: setValueToStringSlice(plan.Sources.Id),
		}
		rule.Sources = sources
	}

	if !plan.Destinations.IsNull() && !plan.Destinations.IsUnknown() {
		destinations := &sdk.CreateNetworkFirewallRuleRequestRuleDestinations{
			Id: setValueToStringSlice(plan.Destinations.Id),
		}
		rule.Destinations = destinations
	}

	if !plan.Scopes.IsNull() && !plan.Scopes.IsUnknown() {
		scopes := &sdk.CreateNetworkFirewallRuleRequestRuleScopes{
			Id: setValueToStringSlice(plan.Scopes.Id),
		}
		rule.Scopes = scopes
	}

	config := &sdk.CreateNetworkFirewallRuleRequestRuleConfig{}
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
		ruleGroup := &sdk.CreateNetworkFirewallRuleRequestRuleRuleGroup{}
		if !plan.RuleGroupId.Id.IsNull() && !plan.RuleGroupId.Id.IsUnknown() {
			ruleGroupID := int32(plan.RuleGroupId.Id.ValueInt64()) //nolint:gosec // API uses int32
			ruleGroup.Id = &ruleGroupID
		}
		rule.RuleGroup = ruleGroup
	}

	addReq := &sdk.CreateNetworkFirewallRuleRequest{
		Rule: rule,
	}

	networkIntegrationId := plan.NetworkIntegrationId.ValueInt64()

	createResp, httpResp, err := client.NetworksAPI.
		CreateNetworkFirewallRule(ctx, networkIntegrationId).
		CreateNetworkFirewallRuleRequest(*addReq).Execute()
	if err != nil || (httpResp != nil && httpResp.StatusCode != http.StatusOK) {
		resp.Diagnostics.AddError(
			"error creating network firewall rule",
			"network firewall rule "+plan.Name.ValueString()+" POST failed: "+errfmt.ErrMsg(err, httpResp),
		)

		return
	}

	createdID := *createResp.Id.Get()

	// Taint the resource state on post-create failures so Terraform can
	// destroy and recreate the resource on the next apply.
	taintResourceState := func(id int64) {
		cleanup.TaintResourceState(ctx, cleanup.TaintResourceStateConfig{
			ResourceType: "network_firewall_rule",
			ResourceID:   id,
			StateWriter:  &resp.State,
			Diagnostics:  &resp.Diagnostics,
		})
	}

	state, _, diag := getNetworkFirewallRuleAsState(ctx, createdID, networkIntegrationId, client)
	if resp.Diagnostics.Append(diag...); resp.Diagnostics.HasError() {
		taintResourceState(createdID)

		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
