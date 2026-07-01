// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkfirewallrulegroup

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"

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

	// Detect import: ImportState sets network_integration_id, id, and
	// external_type but not name. On normal refresh, name is always a known
	// string from prior state.
	isImport := data.Name.IsNull()
	priorTenantIds := data.TenantIds

	client, err := r.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError("creating client failed", err.Error())

		return
	}

	id := data.Id.ValueInt64()
	serverID := data.NetworkIntegrationId.ValueInt64()

	state, notFound, diags := getNetworkFirewallRuleGroupAsState(
		ctx, id, serverID, client, data,
	)
	if notFound {
		resp.State.RemoveResource(ctx)

		return
	}

	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	// On normal refresh, preserve tenant_ids from prior state. The API may
	// silently drop IDs that don't exist in the environment. On import there
	// is no prior state, so we keep the API value from getNetworkFirewallRuleGroupAsState.
	if !isImport {
		state.TenantIds = priorTenantIds
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func getNetworkFirewallRuleGroupAsState(
	ctx context.Context,
	id int64,
	serverID int64,
	client *sdk.APIClient,
	prior NetworkFirewallRuleGroupModel,
) (NetworkFirewallRuleGroupModel, bool, diag.Diagnostics) {
	var state NetworkFirewallRuleGroupModel
	var diags diag.Diagnostics

	getResp, httpResp, err := client.NetworksAPI.
		GetNetworkFirewallRuleGroup(ctx, id, serverID).Execute()
	if err != nil || httpResp.StatusCode != http.StatusOK {
		if httpResp != nil && httpResp.StatusCode == http.StatusNotFound {
			return state, true, diags
		}

		diags.AddError(
			"error reading network firewall rule group",
			fmt.Sprintf("network firewall rule group %d GET failed: ", id)+
				errfmt.ErrMsg(err, httpResp),
		)

		return state, false, diags
	}

	ruleGroup := getResp.RuleGroup
	if ruleGroup == nil {
		diags.AddError("API returned nil", "RuleGroup is nil in the response")

		return state, false, diags
	}

	state.Id = convert.Int64ToType(ruleGroup.Id)
	state.Name = convert.StrToType(ruleGroup.Name)
	state.Priority = convert.Int64ToType(ruleGroup.Priority)
	state.GroupLayer = convert.StrToType(ruleGroup.GroupLayer)
	state.Visibility = convert.StrToType(ruleGroup.Visibility)

	if ruleGroup.Description.IsSet() {
		state.Description = convert.StrToType(ruleGroup.Description.Get())
	} else {
		state.Description = types.StringNull()
	}

	// Build tenant_ids from the Tenants array returned by the API.
	if len(ruleGroup.Tenants) > 0 {
		ids := make([]int64, 0, len(ruleGroup.Tenants))
		for _, t := range ruleGroup.Tenants {
			if t.Id != nil {
				ids = append(ids, *t.Id)
			}
		}
		listVal, listDiags := types.SetValueFrom(ctx, types.Int64Type, ids)
		diags.Append(listDiags...)
		state.TenantIds = listVal
	} else {
		state.TenantIds = types.SetValueMust(types.Int64Type, []attr.Value{})
	}

	// network_integration_id and external_type are not returned in the GET response; preserve from prior state.
	state.NetworkIntegrationId = prior.NetworkIntegrationId
	state.ExternalType = prior.ExternalType

	return state, false, diags
}
