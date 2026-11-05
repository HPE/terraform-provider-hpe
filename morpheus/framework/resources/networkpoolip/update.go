// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkpoolip

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/clearable"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
)

func (r *Resource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	var plan, currentState NetworkPoolIpModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(req.State.Get(ctx, &currentState)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client, err := r.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError("creating client failed", err.Error())

		return
	}

	networkPoolID := currentState.NetworkPoolId.ValueInt64()
	id := currentState.Id.ValueInt64()

	// Only hostname, ref_type and ref_id are editable; everything else is
	// RequiresReplace or create-only.
	ipReq := &sdk.UpdateNetworkPoolIpRequestNetworkPoolIp{}

	if !plan.Hostname.IsNull() && !plan.Hostname.IsUnknown() {
		ipReq.Hostname = plan.Hostname.ValueStringPointer()
	}

	// Attributes removed from the configuration since the last apply, to be
	// cleared. Read before Record, below, overwrites what it reads.
	pending, pdiags := clearable.Compute(ctx, req.Config, req.State, req.Private, clearableAttributes)
	resp.Diagnostics.Append(pdiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// ref_id is a NullableInt64; cleared by sending an explicit null.
	clearable.NullableInt64(pending, "ref_id", plan.RefId, &ipReq.RefId)

	// ref_type is a plain *string with omitempty, so a nil pointer cannot
	// express JSON null -- the field would simply be omitted and the server
	// would keep the current value. Morpheus binds whichever keys are present
	// (DataBindingUtils.bindObjectToInstance(poolIp, props, ['hostname',
	// 'refType','refId'])), so to clear ref_type we must send the key with a
	// null value. Put refType:null into AdditionalProperties, which the
	// generated model marshals alongside the typed fields.
	switch {
	case pending.Clears("ref_type"):
		if ipReq.AdditionalProperties == nil {
			ipReq.AdditionalProperties = map[string]interface{}{}
		}
		ipReq.AdditionalProperties["refType"] = nil
	case !plan.RefType.IsNull() && !plan.RefType.IsUnknown():
		ipReq.RefType = plan.RefType.ValueStringPointer()
	}

	updateReq := sdk.UpdateNetworkPoolIpRequest{
		NetworkPoolIp: ipReq,
	}

	_, httpResp, err := client.NetworksAPI.
		UpdateNetworkPoolIp(ctx, networkPoolID, id).
		UpdateNetworkPoolIpRequest(updateReq).
		Execute()
	if err != nil || httpResp.StatusCode != http.StatusOK {
		resp.Diagnostics.AddError(
			"error updating network pool IP",
			fmt.Sprintf("network pool IP %d PUT failed: ", id)+
				errfmt.ErrMsg(err, httpResp),
		)

		return
	}

	state, found, diags := getNetworkPoolIpAsState(
		ctx, networkPoolID, id, client, plan,
	)
	if !found && !diags.HasError() {
		diags.AddError("read network pool IP", "reservation not found immediately after update")
	}
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(clearable.Record(ctx, req.Config, resp.Private, clearableAttributes)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
