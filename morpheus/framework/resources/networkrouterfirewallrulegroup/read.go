// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkrouterfirewallrulegroup

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/utils/convert"
)

const readOperation = "read network router firewall rule group resource"

// getRuleGroupAsState fetches the firewall rule group by ID and maps it into a
// NetworkRouterFirewallRuleGroupModel. The plan/prior-state model is required
// so that write-only fields (visibility, tenant_ids) can be preserved — none
// of these are returned by the single GET endpoint.
//
// Returns (state, notFound, diags). notFound is true when the API returned 404;
// the caller should remove the resource from state in that case.
func getRuleGroupAsState(
	ctx context.Context,
	id int64,
	routerID int64,
	client *sdk.APIClient,
	prior NetworkRouterFirewallRuleGroupModel,
) (NetworkRouterFirewallRuleGroupModel, bool, diag.Diagnostics) {
	var state NetworkRouterFirewallRuleGroupModel
	var diags diag.Diagnostics

	resp, hresp, err := client.NetworksAPI.
		GetNetworkRouterFirewallRuleGroup(ctx, id, routerID).Execute()
	if err != nil || hresp.StatusCode != http.StatusOK {
		if hresp != nil && hresp.StatusCode == http.StatusNotFound {
			return state, true, diags
		}

		diags.AddError(
			readOperation,
			fmt.Sprintf("firewall rule group %d GET failed: %s", id, errfmt.ErrMsg(err, hresp)),
		)

		return state, false, diags
	}

	rg := resp.RuleGroup
	if rg == nil {
		diags.AddError("API returned nil", "RuleGroup is nil in the response")

		return state, false, diags
	}

	state.Id = convert.Int64ToType(rg.Id)
	state.RouterId = prior.RouterId // path param — always preserve from prior/plan

	state.Name = convert.StrToType(rg.Name)
	state.Priority = convert.Int64ToType(rg.Priority)
	state.GroupLayer = convert.StrToType(rg.GroupLayer)
	state.ExternalId = convert.StrToType(rg.ExternalId)
	state.Status = convert.StrToType(rg.Status)
	state.Description = convert.StrToType(rg.Description.Get())

	// visibility and tenant_ids are not reliably returned by the single GET
	// endpoint. Fall back to prior state to avoid spurious diffs.
	//
	// visibility: nil-check allows the provider to naturally recover if the
	// API is later patched to return it consistently.
	//
	// tenant_ids: the API always returns the root tenant (id=1) regardless of
	// what was configured, so the returned slice cannot represent the
	// user-configured restriction. Preserve prior/plan value unconditionally.
	if rg.Visibility != nil {
		state.Visibility = convert.StrToType(rg.Visibility)
	} else {
		state.Visibility = prior.Visibility
	}

	state.TenantIds = prior.TenantIds

	return state, false, diags
}

func (r *Resource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	var prior NetworkRouterFirewallRuleGroupModel

	resp.Diagnostics.Append(req.State.Get(ctx, &prior)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client, err := r.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError(readOperation, "failed to create client: "+err.Error())

		return
	}

	state, notFound, diags := getRuleGroupAsState(
		ctx, prior.Id.ValueInt64(), prior.RouterId.ValueInt64(), client, prior,
	)
	if notFound {
		resp.State.RemoveResource(ctx)

		return
	}

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
