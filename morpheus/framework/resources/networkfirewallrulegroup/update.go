// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkfirewallrulegroup

import (
	"context"
	"fmt"
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
	var plan, currentState NetworkFirewallRuleGroupModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &currentState)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client, err := r.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError("creating client failed", err.Error())

		return
	}

	id := currentState.Id.ValueInt64()
	serverID := currentState.NetworkIntegrationId.ValueInt64()

	ruleGroup := sdk.NewUpdateNetworkFirewallRuleGroupRequestRuleGroupWithDefaults()

	if !plan.Name.IsNull() && !plan.Name.IsUnknown() {
		ruleGroup.SetName(plan.Name.ValueString())
	}

	if plan.Description.IsNull() {
		ruleGroup.SetDescriptionNil()
	} else if !plan.Description.IsUnknown() {
		ruleGroup.SetDescription(plan.Description.ValueString())
	}

	if plan.Priority.IsNull() {
		ruleGroup.SetPriorityNil()
	} else if !plan.Priority.IsUnknown() {
		ruleGroup.SetPriority(plan.Priority.ValueInt64())
	}

	updateReq := sdk.NewUpdateNetworkFirewallRuleGroupRequestWithDefaults()
	updateReq.SetRuleGroup(*ruleGroup)

	_, httpResp, err := client.NetworksAPI.
		UpdateNetworkFirewallRuleGroup(ctx, id, serverID).
		UpdateNetworkFirewallRuleGroupRequest(*updateReq).Execute()
	if err != nil || httpResp.StatusCode != http.StatusOK {
		resp.Diagnostics.AddError(
			"error updating network firewall rule group",
			fmt.Sprintf("network firewall rule group %d PUT failed: ", id)+
				errfmt.ErrMsg(err, httpResp),
		)

		return
	}

	state, _, diags := getNetworkFirewallRuleGroupAsState(
		ctx, id, serverID, client, plan,
	)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
