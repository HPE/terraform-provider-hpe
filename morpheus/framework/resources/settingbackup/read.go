// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package settingbackup

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
			tflog.Warn(ctx, "backup settings not found, removing the resource from state")
			resp.State.RemoveResource(ctx)

			return
		}

		errfmt.DiagError(&resp.Diagnostics, errfmt.OpRead, resourceType, "", err, nil)

		return
	}

	resp.Diagnostics.Append(mapSettings(&state.SettingBackupModel, settings)...)
	state.Id = types.StringValue(singletonID)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)

	// The schema declares no defaults today, so this is a no-op; it keeps an
	// imported resource's plan clean should one be added to the spec.
	resp.Diagnostics.Append(schemadefaults.Apply(ctx, resourceSchema(ctx), &resp.State)...)
}
