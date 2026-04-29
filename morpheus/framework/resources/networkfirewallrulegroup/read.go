// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package networkfirewallrulegroup

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/utils/convert"
)

func (r *Resource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	var data NetworkFirewallRuleGroupModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client, err := r.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError("creating client failed", err.Error())

		return
	}

	id := data.Id.ValueInt64()
	serverID := data.NetworkIntegrationId.ValueInt64()

	state, diags := getNetworkFirewallRuleGroupAsState(
		ctx, id, serverID, client, data,
	)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func getNetworkFirewallRuleGroupAsState(
	ctx context.Context,
	id int64,
	serverID int64,
	client *sdk.APIClient,
	prior NetworkFirewallRuleGroupModel,
) (NetworkFirewallRuleGroupModel, diag.Diagnostics) {
	var state NetworkFirewallRuleGroupModel
	var diags diag.Diagnostics

	getResp, httpResp, err := client.NetworksAPI.
		GetNetworkFirewallRuleGroup(ctx, id, serverID).Execute()
	if err != nil || httpResp.StatusCode != http.StatusOK {
		diags.AddError(
			"error reading network firewall rule group",
			fmt.Sprintf("network firewall rule group %d GET failed: ", id)+
				errfmt.ErrMsg(err, httpResp),
		)

		return state, diags
	}

	ruleGroup := getResp.GetRuleGroup()

	state.Id = convert.Int64ToType(ruleGroup.Id)
	state.Name = convert.StrToType(ruleGroup.Name)
	state.Priority = convert.Int64ToType(ruleGroup.Priority)
	state.GroupLayer = convert.StrToType(ruleGroup.GroupLayer)

	if ruleGroup.Description.IsSet() {
		state.Description = convert.StrToType(ruleGroup.Description.Get())
	} else {
		state.Description = types.StringNull()
	}

	// network_integration_id and external_type are not returned in the GET response; preserve from prior state.
	state.NetworkIntegrationId = prior.NetworkIntegrationId
	state.ExternalType = prior.ExternalType

	return state, diags
}
