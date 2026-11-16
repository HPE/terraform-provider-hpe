// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package hostprofile

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/utils/cleanup"
)

func (r *Resource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	var plan HostProfileModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client, err := r.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"create host profile resource",
			"failed to create client: "+err.Error(),
		)

		return
	}

	createReq := sdk.CreateHostProfileRequest{
		Name: plan.Name.ValueString(),
	}
	if !plan.Description.IsNull() && !plan.Description.IsUnknown() {
		createReq.Description.Set(plan.Description.ValueStringPointer())
	}

	cfg, cfgDiags := expandCreateConfig(ctx, plan.Config)
	resp.Diagnostics.Append(cfgDiags...)
	if resp.Diagnostics.HasError() {
		return
	}
	createReq.Config = cfg

	createResp, hresp, err := client.HostProfilesAPI.
		CreateHostProfile(ctx).
		CreateHostProfileRequest(createReq).
		Execute()
	if err != nil || hresp == nil || (hresp.StatusCode != http.StatusOK &&
		hresp.StatusCode != http.StatusCreated) {
		resp.Diagnostics.AddError(
			"create host profile resource",
			"host profile POST failed: "+errfmt.ErrMsg(err, hresp),
		)

		return
	}

	if createResp == nil || createResp.Id == nil {
		resp.Diagnostics.AddError(
			"create host profile resource",
			"host profile POST response missing id",
		)

		return
	}

	id := *createResp.Id

	state, found, diags := getHostProfileAsState(ctx, id, client, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() || !found {
		// Only synthesise a "not found" error when the read itself reported no
		// error: getHostProfileAsState returns found=false both for a genuine
		// 404 and alongside its own diagnostics on a read failure, and stacking
		// a misleading "not found immediately after create" on top of a real
		// error (e.g. a transport failure) would obscure the actual cause.
		if !found && !resp.Diagnostics.HasError() {
			resp.Diagnostics.AddError(
				"create host profile resource",
				fmt.Sprintf(
					"host profile %d not found immediately after create", id,
				),
			)
		}
		// The POST succeeded but the follow-up read failed; taint so the
		// resource is recreated rather than left half-tracked.
		cleanup.TaintResourceState(ctx, cleanup.TaintResourceStateConfig{
			ResourceType: "host_profile",
			ResourceID:   id,
			StateWriter:  &resp.State,
			Diagnostics:  &resp.Diagnostics,
		})

		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
