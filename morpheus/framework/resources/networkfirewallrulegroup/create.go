// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkfirewallrulegroup

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/utils/cleanup"
)

func (r *Resource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	var plan NetworkFirewallRuleGroupModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client, err := r.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError("creating client failed", err.Error())

		return
	}

	name := plan.Name.ValueString()

	ruleGroup := &sdk.CreateNetworkFirewallRuleGroupRequestRuleGroup{
		Name:         name,
		ExternalType: plan.ExternalType.ValueString(),
	}

	if !plan.Description.IsNull() && !plan.Description.IsUnknown() {
		desc := plan.Description.ValueString()
		ruleGroup.Description = &desc
	}

	if !plan.Priority.IsNull() && !plan.Priority.IsUnknown() {
		priority := plan.Priority.ValueInt64()
		ruleGroup.Priority = &priority
	}

	if !plan.GroupLayer.IsNull() && !plan.GroupLayer.IsUnknown() {
		layer := plan.GroupLayer.ValueString()
		ruleGroup.GroupLayer = &layer
	}

	createReq := &sdk.CreateNetworkFirewallRuleGroupRequest{
		RuleGroup: ruleGroup,
	}

	serverID := plan.NetworkIntegrationId.ValueInt64()

	createResp, httpResp, err := client.NetworksAPI.
		CreateNetworkFirewallRuleGroup(ctx, serverID).
		CreateNetworkFirewallRuleGroupRequest(*createReq).Execute()
	if err != nil || httpResp.StatusCode != http.StatusOK {
		resp.Diagnostics.AddError(
			"error creating network firewall rule group",
			"network firewall rule group "+name+" POST failed: "+
				errfmt.ErrMsg(err, httpResp),
		)

		return
	}

	if createResp == nil {
		resp.Diagnostics.AddError(
			"error creating network firewall rule group",
			"network firewall rule group "+name+" POST returned an empty response",
		)

		return
	}

	createdID := *createResp.Id.Get()
	plan.Id = types.Int64Value(createdID)

	taintResourceState := func(id int64) {
		cleanup.TaintResourceState(ctx, cleanup.TaintResourceStateConfig{
			ResourceType: "network_firewall_rule_group",
			ResourceID:   id,
			StateWriter:  &resp.State,
			Diagnostics:  &resp.Diagnostics,
		})
	}

	state, _, diags := getNetworkFirewallRuleGroupAsState(
		ctx, createdID, serverID, client, plan,
	)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		resp.Diagnostics.AddError(
			"failed to read network firewall rule group state",
			fmt.Sprintf(
				"Network firewall rule group %d was created but could not be read",
				createdID,
			),
		)
		taintResourceState(createdID)

		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError(
			"failed to set network firewall rule group state",
			fmt.Sprintf(
				"Network firewall rule group %d was created but state could not be saved",
				createdID,
			),
		)
		taintResourceState(createdID)

		return
	}
}
