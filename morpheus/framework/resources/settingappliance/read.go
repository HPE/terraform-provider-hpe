// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package settingappliance

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/utils/schemadefaults"
)

func (r *Resource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	var state resourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)

		return
	}

	settings, httpResp, err := getSettings(ctx, client)
	if err != nil {
		// As the SDKv2 resource did: an appliance that has no settings to
		// return has nothing for the resource to manage, so it is forgotten
		// rather than failing every plan.
		if errfmt.IsNotFound(httpResp) {
			tflog.Warn(ctx, "appliance settings not found, removing the resource from state")
			resp.State.RemoveResource(ctx)

			return
		}

		errfmt.DiagError(&resp.Diagnostics, errfmt.OpRead, resourceType, "", err, nil)

		return
	}

	resp.Diagnostics.Append(mapSettings(&state.SettingApplianceModel, settings)...)
	state.Id = types.StringValue(singletonID)

	// Compare the appliance's hash of each stored credential with the baseline
	// recorded when the provider last sent it: an inequality means it was
	// changed outside Terraform, and the next plan re-sends the configured
	// value.
	drifted, diags := wo.CheckDrift(ctx, resp.Private, settings)
	resp.Diagnostics.Append(diags...)

	if len(drifted) > 0 {
		tflog.Info(ctx, "write-only value changed on the appliance",
			map[string]any{"attributes": drifted})
	}

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)

	// The schema declares no defaults today, so this is a no-op; it keeps an
	// imported resource's plan clean should one be added to the spec.
	resp.Diagnostics.Append(schemadefaults.Apply(ctx, resourceSchema(ctx), &resp.State)...)
}
