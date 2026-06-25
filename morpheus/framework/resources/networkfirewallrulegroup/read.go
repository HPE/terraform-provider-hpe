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

	if ruleGroup.Description.IsSet() {
		state.Description = convert.StrToType(ruleGroup.Description.Get())
	} else {
		state.Description = types.StringNull()
	}

	// network_integration_id and external_type are not returned in the GET response; preserve from prior state.
	state.NetworkIntegrationId = prior.NetworkIntegrationId
	state.ExternalType = prior.ExternalType

	// visibility and tenant_ids: in AdditionalProperties on GET response
	state.Visibility = types.StringNull()
	if v, ok := ruleGroup.AdditionalProperties["visibility"]; ok {
		if s, ok := v.(string); ok {
			state.Visibility = types.StringValue(s)
		}
	}

	state.TenantIds = types.ListNull(types.Int64Type)
	if tenantsRaw, ok := ruleGroup.AdditionalProperties["tenants"]; ok {
		if tenantsArr, ok := tenantsRaw.([]interface{}); ok {
			tenantVals := make([]attr.Value, 0, len(tenantsArr))
			for _, t := range tenantsArr {
				if tMap, ok := t.(map[string]interface{}); ok {
					switch id := tMap["id"].(type) {
					case float64:
						tenantVals = append(tenantVals, types.Int64Value(int64(id)))
					case int64:
						tenantVals = append(tenantVals, types.Int64Value(id))
					}
				}
			}
			if len(tenantVals) > 0 {
				tenantList, _ := types.ListValue(types.Int64Type, tenantVals)
				state.TenantIds = tenantList
			}
		}
	}

	return state, false, diags
}
