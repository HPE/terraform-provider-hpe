// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkrouterfirewallrulegroup

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/utils/cleanup"
)

const createOperation = "create network router firewall rule group resource"

func (r *Resource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	var plan NetworkRouterFirewallRuleGroupModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client, err := r.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError(createOperation, "failed to create client: "+err.Error())

		return
	}

	routerID := plan.RouterId.ValueInt64()

	ruleGroup := sdk.CreateNetworkRouterFirewallRuleGroupRequestRuleGroup{
		Name:         plan.Name.ValueString(),
		ExternalType: "GatewayPolicy",
	}

	if !plan.Description.IsNull() && !plan.Description.IsUnknown() {
		ruleGroup.Description = plan.Description.ValueStringPointer()
	}

	if !plan.Priority.IsNull() && !plan.Priority.IsUnknown() {
		ruleGroup.Priority = plan.Priority.ValueInt64Pointer()
	}

	if !plan.GroupLayer.IsNull() && !plan.GroupLayer.IsUnknown() {
		ruleGroup.GroupLayer = plan.GroupLayer.ValueStringPointer()
	}

	if !plan.Visibility.IsNull() && !plan.Visibility.IsUnknown() {
		ruleGroup.Visibility = plan.Visibility.ValueStringPointer()
	}

	if !plan.TenantIds.IsNull() && !plan.TenantIds.IsUnknown() {
		var ids []int64

		resp.Diagnostics.Append(plan.TenantIds.ElementsAs(ctx, &ids, false)...)
		if resp.Diagnostics.HasError() {
			return
		}

		tenants := make([]sdk.CreateNetworkRouterFirewallRuleGroupRequestRuleGroupTenantsInner, 0, len(ids))
		for i := range ids {
			id := ids[i]
			tenants = append(tenants, sdk.CreateNetworkRouterFirewallRuleGroupRequestRuleGroupTenantsInner{Id: &id})
		}

		ruleGroup.Tenants = tenants
	}

	createReq := sdk.CreateNetworkRouterFirewallRuleGroupRequest{
		RuleGroup: &ruleGroup,
	}

	result, hresp, err := client.NetworksAPI.
		CreateNetworkRouterFirewallRuleGroup(ctx, routerID).
		CreateNetworkRouterFirewallRuleGroupRequest(createReq).Execute()
	if err != nil || hresp.StatusCode != http.StatusOK {
		resp.Diagnostics.AddError(
			createOperation,
			fmt.Sprintf("router %d firewall rule group POST failed: %s", routerID, errfmt.ErrMsg(err, hresp)),
		)

		return
	}

	if !result.Id.IsSet() || result.Id.Get() == nil {
		resp.Diagnostics.AddError("API returned nil", "ID is nil in the response")

		return
	}

	id := *result.Id.Get()
	plan.Id = types.Int64Value(id)

	taintState := func() {
		cleanup.TaintResourceState(ctx, cleanup.TaintResourceStateConfig{
			ResourceType: "network_router_firewall_rule_group",
			ResourceID:   id,
			StateWriter:  &resp.State,
			Diagnostics:  &resp.Diagnostics,
		})
	}

	state, notFound, diags := getRuleGroupAsState(ctx, id, routerID, client, plan)
	if notFound {
		// Unexpected: resource was just created but GET returned 404.
		resp.Diagnostics.AddError(
			createOperation,
			fmt.Sprintf("firewall rule group %d was created but GET returned 404", id),
		)
		taintState()

		return
	}

	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		resp.Diagnostics.AddError(
			createOperation,
			fmt.Sprintf("firewall rule group %d was created but could not be read", id),
		)
		taintState()

		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		taintState()
	}
}
