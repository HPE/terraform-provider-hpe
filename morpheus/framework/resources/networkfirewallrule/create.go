// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkfirewallrule

import (
	"context"
	"net/http"
	"strconv"

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

	rule := sdk.NewCreateNetworkFirewallRuleRequestRule(plan.Name.ValueString())

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
		priorityVal, err := strconv.ParseInt(plan.Priority.ValueString(), 10, 64)
		if err != nil {
			resp.Diagnostics.AddError("invalid priority value", "priority must be a valid integer: "+err.Error())

			return
		}

		rule.SetPriority(priorityVal)
	}

	if !plan.Sources.IsNull() && !plan.Sources.IsUnknown() {
		sources := sdk.NewCreateNetworkFirewallRuleRequestRuleSourcesWithDefaults()
		sources.SetId(setValueToStringSlice(plan.Sources.Id))
		rule.SetSources(*sources)
	}

	if !plan.Destinations.IsNull() && !plan.Destinations.IsUnknown() {
		destinations := sdk.NewCreateNetworkFirewallRuleRequestRuleDestinationsWithDefaults()
		destinations.SetId(setValueToStringSlice(plan.Destinations.Id))
		rule.SetDestinations(*destinations)
	}

	if !plan.Scopes.IsNull() && !plan.Scopes.IsUnknown() {
		scopes := sdk.NewCreateNetworkFirewallRuleRequestRuleScopesWithDefaults()
		scopes.SetId(setValueToStringSlice(plan.Scopes.Id))
		rule.SetScopes(*scopes)
	}

	config := sdk.NewCreateNetworkFirewallRuleRequestRuleConfigWithDefaults()
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
		ruleGroup := sdk.NewCreateNetworkFirewallRuleRequestRuleRuleGroupWithDefaults()
		if !plan.RuleGroupId.Id.IsNull() && !plan.RuleGroupId.Id.IsUnknown() {
			ruleGroupID := int32(plan.RuleGroupId.Id.ValueInt64()) //nolint:gosec // API uses int32
			ruleGroup.SetId(ruleGroupID)
		}
		rule.SetRuleGroup(*ruleGroup)
	}

	addReq := sdk.NewCreateNetworkFirewallRuleRequestWithDefaults()
	addReq.SetRule(*rule)

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

	createdID := createResp.GetId()

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
