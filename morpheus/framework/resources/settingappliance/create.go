// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package settingappliance

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/wodigest"
)

// Create applies the configured settings. The settings always exist, so
// creating the resource is an update of them; a configuration that sets
// nothing sends nothing and adopts the settings as they are.
func (r *Resource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	var plan, config resourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)

	if resp.Diagnostics.HasError() {
		return
	}

	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)

		return
	}

	body, send := settingsBody(&plan.SettingApplianceModel, nil)
	if addCredentials(&body, &config.SettingApplianceModel) {
		send = true
	}

	if send {
		if err := updateSettings(ctx, client, body); err != nil {
			errfmt.DiagError(&resp.Diagnostics, errfmt.OpCreate, resourceType, "", err, nil)

			return
		}
	}

	// If the read-back fails, Create reports the error and records nothing:
	// with no state, Terraform runs Create again on the next apply, which
	// sends the settings and the credentials again -- harmless for a
	// singleton. So the digests are recorded only once the read succeeded.
	// The resource is not tainted, as a resource with an object of its own on
	// the appliance would be: there is nothing to replace, and a taint would
	// only plan the same Create as a replacement.
	settings, _, err := getSettings(ctx, client)
	if err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpCreate, resourceType, "", err, nil)

		return
	}

	state := plan
	resp.Diagnostics.Append(mapSettings(&state.SettingApplianceModel, settings)...)
	state.Id = types.StringValue(singletonID)
	state.WriteOnlyRevision = wodigest.NewRevision()

	// Record the digests of the configured credentials and, as the baseline
	// for detecting a change made outside Terraform, the appliance's hashes of
	// the stored ones as read back right after they were sent. If they cannot
	// be recorded, Create records nothing either, as when the read-back fails:
	// the next apply runs Create again rather than starting from a resource
	// with no digests.
	resp.Diagnostics.Append(wo.StoreObserved(ctx, resp.Private, req.Config, settings)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
