// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package settingappliance

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/wodigest"
)

// Update sends the settings that changed and, when ModifyPlan planned it, the
// configured credentials.
func (r *Resource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	var plan, state, config resourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)

	if resp.Diagnostics.HasError() {
		return
	}

	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)

		return
	}

	body, send := settingsBody(&plan.SettingApplianceModel, &state.SettingApplianceModel)

	// ModifyPlan marks write_only_revision unknown when a configured credential
	// changed, was changed on the appliance, or the resource has no record of
	// the credentials at all, as after an upgrade from the SDKv2 resource (an
	// import records digests or a sentinel; see import.go): the configured
	// credentials are re-sent then.
	resend := wodigest.Resend(ctx, req.Plan)
	if resend && addCredentials(&body, &config.SettingApplianceModel) {
		send = true
	}

	// A plan that changes nothing on the appliance -- only write_only_revision,
	// with no credential configured -- sends nothing.
	if send {
		if err := updateSettings(ctx, client, body); err != nil {
			errfmt.DiagError(&resp.Diagnostics, errfmt.OpUpdate, resourceType, "", err, nil)

			return
		}
	}

	// Always record the digests of the configured credentials, sent or not,
	// so that neither an import sentinel nor a missing digest survives an
	// update. They are recorded before the read-back: if it fails, the
	// framework keeps the prior state together with these digests, which are
	// right about what the appliance now holds, and the refresh before the
	// next plan reads the settings themselves back. The appliance's hashes,
	// though, are only observed after the read-back (wo.Observe below), so
	// when it fails after credentials were re-sent, that refresh sees their
	// new hashes against the old baseline, and the next plan re-sends the
	// configured credentials once, which re-baselines them. Recording the
	// digests after the read-back instead would cost the same one re-send,
	// through digests that no longer match the configuration.
	resp.Diagnostics.Append(wo.Store(ctx, resp.Private, req.Config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	settings, _, err := getSettings(ctx, client)
	if err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpUpdate, resourceType, "", err, nil)

		return
	}

	newState := plan
	resp.Diagnostics.Append(mapSettings(&newState.SettingApplianceModel, settings)...)
	newState.Id = types.StringValue(singletonID)

	newState.WriteOnlyRevision = state.WriteOnlyRevision
	if resend {
		newState.WriteOnlyRevision = wodigest.NewRevision()
	}

	// Re-baseline the appliance-side drift detection off the read-back.
	resp.Diagnostics.Append(wo.Observe(ctx, resp.Private, settings)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}
