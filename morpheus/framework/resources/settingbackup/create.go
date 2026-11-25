// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package settingbackup

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
)

// Create applies the configured settings. The settings always exist, so
// creating the resource is an update of them; a configuration that sets
// nothing sends nothing and adopts the settings as they are.
func (r *Resource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	var plan resourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)

		return
	}

	body, send := settingsBody(&plan.SettingBackupModel, nil)

	if send {
		if err := updateSettings(ctx, client, body); err != nil {
			errfmt.DiagError(&resp.Diagnostics, errfmt.OpCreate, resourceType, "", err, nil)

			return
		}
	}

	// If the read-back fails, Create reports the error and records nothing:
	// with no state, Terraform runs Create again on the next apply, which
	// sends the settings again -- harmless for a singleton. The resource is not
	// tainted, as a resource with an object of its own on the appliance would
	// be: there is nothing to replace, and a taint would only plan the same
	// Create as a replacement.
	settings, _, err := getSettings(ctx, client)
	if err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpCreate, resourceType, "", err, nil)

		return
	}

	state := plan
	resp.Diagnostics.Append(mapSettings(&state.SettingBackupModel, settings)...)
	state.Id = types.StringValue(singletonID)

	// A setting the appliance accepted but did not apply is reported with its
	// reason here rather than as an inconsistent result by Terraform. Nothing
	// is recorded: as above, the next apply runs Create again.
	if send {
		resp.Diagnostics.Append(checkApplied(body, &plan.SettingBackupModel, &state.SettingBackupModel,
			backupApplianceHint(ctx, client))...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
