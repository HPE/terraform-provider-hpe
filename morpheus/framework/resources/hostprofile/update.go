// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package hostprofile

import (
	"context"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
)

func (r *Resource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	var plan HostProfileModel
	var state HostProfileModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// The seeded HVM default profile rejects modification on the API; refuse
	// up front with a clear message rather than surfacing a 400.
	if state.Code.ValueString() == seedCode {
		resp.Diagnostics.AddError(
			"update host profile resource",
			"the seeded default host profile '"+seedCode+
				"' is immutable and cannot be updated",
		)

		return
	}

	client, err := r.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"update host profile resource",
			"failed to create client: "+err.Error(),
		)

		return
	}

	id := state.Id.ValueInt64()

	updateReq := sdk.UpdateHostProfileRequest{
		Name: plan.Name.ValueStringPointer(),
	}
	if !plan.Description.IsNull() && !plan.Description.IsUnknown() {
		updateReq.Description.Set(plan.Description.ValueStringPointer())
	} else {
		// The host-profile PUT ignores an explicit null: its update handler
		// only assigns the description when body.description != null, so a
		// removed description would otherwise persist unchanged. (Only the PUT
		// path is reachable here -- the SDK issues UpdateHostProfile as an HTTP
		// PUT; the server's PATCH handler, which does clear on null, is never
		// called.) Send an empty string to clear it instead; the read-back
		// coerces "" back to null so the post-apply state matches the now-null
		// configuration. This clear-on-removal behaviour is pinned by the
		// "remove description" step in TestAccMorpheusHostProfileResourceUpdateOk.
		empty := ""
		updateReq.Description.Set(&empty)
	}

	cfg, cfgDiags := expandUpdateConfig(ctx, plan.Config)
	resp.Diagnostics.Append(cfgDiags...)
	if resp.Diagnostics.HasError() {
		return
	}
	updateReq.Config = cfg

	// Clearing the whole config document: when the plan drops config entirely
	// but state still carries one, expandUpdateConfig returns nil, which the
	// SDK omits from the PUT body -- and the API's update handler only writes
	// config when it is present, so an omitted config would silently preserve
	// the old document. Send an explicit empty object instead: setConfigMap
	// does a full replace and clears the stored config on an empty map. (With
	// every config field now plain-optional, individual leaves are cleared the
	// same way -- the full map we send fully replaces the stored one.)
	if cfg == nil && !state.Config.IsNull() && !plan.Config.IsUnknown() {
		updateReq.Config = &sdk.UpdateHostProfileRequestConfig{}
	}

	_, hresp, err := client.HostProfilesAPI.
		UpdateHostProfile(ctx, id).
		UpdateHostProfileRequest(updateReq).
		Execute()
	if err != nil || hresp == nil || hresp.StatusCode != http.StatusOK {
		resp.Diagnostics.AddError(
			"update host profile resource",
			"host profile PUT failed: "+errfmt.ErrMsg(err, hresp),
		)

		return
	}

	newState, found, diags := getHostProfileAsState(ctx, id, client, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	if !found {
		resp.State.RemoveResource(ctx)

		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}
