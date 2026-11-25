// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package settingbackup

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
)

// Update sends the settings that changed.
func (r *Resource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	var plan, state resourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)

		return
	}

	body, send := settingsBody(&plan.SettingBackupModel, &state.SettingBackupModel)

	// A plan that changes nothing on the appliance sends nothing.
	if send {
		if err := updateSettings(ctx, client, body); err != nil {
			errfmt.DiagError(&resp.Diagnostics, errfmt.OpUpdate, resourceType, "", err, nil)

			return
		}
	}

	settings, _, err := getSettings(ctx, client)
	if err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpUpdate, resourceType, "", err, nil)

		return
	}

	newState := plan
	resp.Diagnostics.Append(mapSettings(&newState.SettingBackupModel, settings)...)
	newState.Id = types.StringValue(singletonID)

	// A setting the appliance accepted but did not apply is reported with its
	// reason here rather than as an inconsistent result by Terraform. The
	// prior state is kept; the refresh before the next plan reads the rest.
	if send {
		resp.Diagnostics.Append(checkApplied(body, &plan.SettingBackupModel, &newState.SettingBackupModel,
			backupApplianceHint(ctx, client))...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}
