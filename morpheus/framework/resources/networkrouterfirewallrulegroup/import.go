// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkrouterfirewallrulegroup

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (r *Resource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	parts := strings.SplitN(req.ID, ".", 2)
	if len(parts) != 2 {
		resp.Diagnostics.AddError(
			"import network router firewall rule group resource",
			"provided import ID '"+req.ID+"' is invalid, expected format 'router_id.group_id'",
		)

		return
	}

	routerID, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		resp.Diagnostics.AddError(
			"import network router firewall rule group resource",
			"provided router_id '"+parts[0]+"' is invalid (non-number)",
		)

		return
	}

	id, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		resp.Diagnostics.AddError(
			"import network router firewall rule group resource",
			"provided group_id '"+parts[1]+"' is invalid (non-number)",
		)

		return
	}

	// Seed a minimal prior state so getRuleGroupAsState can preserve write-only
	// fields. On import, visibility is null and tenant_ids is an empty set —
	// users must re-apply to restore actual values.
	prior := NetworkRouterFirewallRuleGroupModel{
		RouterId:   types.Int64Value(routerID),
		Visibility: types.StringNull(),
	}

	var setDiags diag.Diagnostics
	prior.TenantIds, setDiags = types.SetValue(types.Int64Type, []attr.Value{})
	resp.Diagnostics.Append(setDiags...)

	if resp.Diagnostics.HasError() {
		return
	}

	client, err := r.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"import network router firewall rule group resource",
			"failed to create client: "+err.Error(),
		)

		return
	}

	state, notFound, diags := getRuleGroupAsState(ctx, id, routerID, client, prior)
	if notFound {
		resp.Diagnostics.AddError(
			"import network router firewall rule group resource",
			fmt.Sprintf("firewall rule group %d not found on router %d", id, routerID),
		)

		return
	}

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
