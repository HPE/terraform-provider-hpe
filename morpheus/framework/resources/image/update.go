package image

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/clearable"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/tenancy"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/wodigest"
	"github.com/HPE/terraform-provider-hpe/utils/convert"
)

func (r *Resource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	var plan, state, config ImageModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client, err := r.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError("creating client failed", err.Error())

		return
	}

	name := plan.Name.ValueString()

	updateRequest := &sdk.UpdateVirtualImageRequest{}
	updateImage := &sdk.UpdateVirtualImageRequestVirtualImage{}

	// auto_join_domain
	if !plan.AutoJoinDomain.IsNull() {
		updateImage.IsAutoJoinDomain = plan.AutoJoinDomain.ValueBoolPointer()
	}

	// cloud_init
	if !plan.CloudInit.IsNull() {
		updateImage.IsCloudInit = plan.CloudInit.ValueBoolPointer()
	}

	// config_azure
	// keep unimplemented for now, we can decide in the future if this
	// should continue being a replace or if we want to support in-place updates
	// of the Azure image config

	// description
	// description for whatever reason cannot be updated at the moment

	// fips_enabled
	if !plan.FipsEnabled.IsNull() {
		updateImage.FipsEnabled = plan.FipsEnabled.ValueBoolPointer()
	}

	// force_customization
	if !plan.ForceCustomization.IsNull() {
		updateImage.IsForceCustomization = plan.ForceCustomization.ValueBoolPointer()
	}

	// image_type
	if !plan.ImageType.IsNull() {
		updateImage.ImageType = plan.ImageType.ValueStringPointer()
	}

	// install_agent
	if !plan.InstallAgent.IsNull() {
		updateImage.InstallAgent = plan.InstallAgent.ValueBoolPointer()
	}

	// Attributes removed from the configuration since the last apply, to be
	// cleared. Read before Record, below, overwrites what it reads.
	pending, pdiags := clearable.Compute(ctx, req.Config, req.State, req.Private, clearableAttributes)
	resp.Diagnostics.Append(pdiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// labels: the planned list, or an empty list when being cleared
	labels, sendLabels, ldiags := pending.Strings(ctx, "labels", plan.Labels)
	resp.Diagnostics.Append(ldiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if sendLabels {
		updateImage.Labels = labels
	}

	// min_disk: in bytes on the wire; null when being cleared
	switch {
	case pending.Clears("min_disk"):
		updateImage.MinDisk.Set(nil)
	case !plan.MinDisk.IsNull() && !plan.MinDisk.IsUnknown():
		val := plan.MinDisk.ValueInt64() * 1024 * 1024 * 1024
		updateImage.MinDisk.Set(&val)
	}

	// min_ram: not clearable — Morpheus rejects a null with a server error
	// and ignores zero — so it is only ever sent when known.
	if !plan.MinRam.IsNull() && !plan.MinRam.IsUnknown() {
		val := plan.MinRam.ValueInt64()
		updateImage.MinRamGB.Set(&val)
	}

	// name
	if !plan.Name.IsNull() {
		updateImage.Name = &name
	}

	// os_type_id
	clearable.NullableInt64(pending, "os_type_id", plan.OsTypeId, &updateImage.OsType)

	// ssh_password_wo: re-sent when ModifyPlan marked write_only_revision
	// unknown (its digest changed, or the deprecated marker did). See
	// write_only.go.
	resend := wodigest.Resend(ctx, req.Plan)
	if resend && !config.SshPasswordWo.IsNull() && !config.SshPasswordWo.IsUnknown() {
		updateImage.SshPassword.Set(config.SshPasswordWo.ValueStringPointer())
	}

	// ssh_username
	clearable.NullableString(pending, "ssh_username", plan.SshUsername, &updateImage.SshUsername)

	// storage_provider_id: not clearable — Morpheus ignores it on update and a
	// change replaces the image — so it is only ever sent when known.
	if !plan.StorageProviderId.IsNull() && !plan.StorageProviderId.IsUnknown() {
		storageProvider := &sdk.UpdateVirtualImageRequestVirtualImageStorageProvider{}
		storageProvider.Id = plan.StorageProviderId.ValueInt64Pointer()

		updateImage.StorageProvider = storageProvider
	}

	// sysprep
	if !plan.Sysprep.IsNull() {
		updateImage.IsSysprep = plan.Sysprep.ValueBoolPointer()
	}

	// tags
	tags, diags := convert.FromSetType(ctx, plan.Tags, updateTagMapper)
	if diags.HasError() {
		tflog.Error(ctx, "cannot convert tags")
		resp.Diagnostics.Append(diags...)

		return
	}
	updateImage.Tags = tags

	// tenant_ids: the planned tenants, or an empty list when being cleared
	tenantIDs, sendTenants, tdiags := pending.Int64s(ctx, "tenant_ids", plan.TenantIds)
	resp.Diagnostics.Append(tdiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if sendTenants {
		updateImage.Accounts = tenantIDs
	}

	// trial_version
	if !plan.TrialVersion.IsNull() {
		updateImage.TrialVersion = plan.TrialVersion.ValueBoolPointer()
	}

	// uefi
	if !plan.Uefi.IsNull() {
		updateImage.Uefi = plan.Uefi.ValueBoolPointer()
	}

	// url
	// keep unimplemented for now, we can decide in the future if this
	// should continue being a replace or if we want to support in-place updates
	// of the image URL

	// user_data
	clearable.NullableString(pending, "user_data", plan.UserData, &updateImage.UserData)

	// virtio_supported
	if !plan.VirtioSupported.IsNull() {
		updateImage.VirtioSupported = plan.VirtioSupported.ValueBoolPointer()
	}

	// visibility
	if !plan.Visibility.IsNull() {
		updateImage.Visibility = plan.Visibility.ValueStringPointer()
	}

	// vm_tools_installed
	if !plan.VmToolsInstalled.IsNull() {
		updateImage.VmToolsInstalled = plan.VmToolsInstalled.ValueBoolPointer()
	}

	// send the API request here
	updateRequest.VirtualImage = *updateImage
	image, httpResp, err := client.LibraryAPI.UpdateVirtualImage(ctx, state.Id.ValueInt64()).
		UpdateVirtualImageRequest(*updateRequest).
		Execute()
	if err != nil || httpResp.StatusCode != http.StatusOK {
		resp.Diagnostics.AddError("error updating image", errfmt.ErrMsg(err, httpResp))

		return
	}
	// set the ID value in state
	if image.VirtualImage == nil {
		resp.Diagnostics.AddError("API returned nil", "VirtualImage is nil in the response")

		return
	}

	if image.VirtualImage.Id == nil {
		resp.Diagnostics.AddError("API returned nil", "VirtualImage ID is nil in the response")

		return
	}

	plan.Id = convert.Int64ToType(image.VirtualImage.Id)

	plan.WriteOnlyRevision = state.WriteOnlyRevision
	if resend {
		plan.WriteOnlyRevision = wodigest.NewRevision()
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	state, apiObj, found, diag := getImageAsState(ctx, plan.Id.ValueInt64(), client, plan)
	if resp.Diagnostics.Append(diag...); resp.Diagnostics.HasError() {
		return
	}
	if !found {
		resp.Diagnostics.AddError(
			"failed to read image state",
			fmt.Sprintf("image %d not found after update", plan.Id.ValueInt64()),
		)

		return
	}

	if d := tenancy.CheckVisibilityApplied(plan.Visibility, state.Visibility); d != nil {
		resp.Diagnostics.Append(d)

		return
	}

	// Record the digests of the configured values, sent or not, once the
	// read-back succeeded, so private state and state move together, and
	// re-baseline the appliance-side drift hash from the read-back object.
	resp.Diagnostics.Append(wo.StoreObserved(ctx, resp.Private, req.Config, apiObj)...)

	resp.Diagnostics.Append(clearable.Record(ctx, req.Config, resp.Private, clearableAttributes)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}
