package image

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
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

	// labels
	if !plan.Labels.IsNull() {
		labels, err := convert.SetToStrSlice(plan.Labels)
		if err != nil {
			resp.Diagnostics.AddError(
				"update image resource",
				"image "+name+": failed to parse label: "+err.Error(),
			)

			return
		}

		updateImage.Labels = labels
	}

	// min_disk
	if !plan.MinDisk.IsNull() {
		val := plan.MinDisk.ValueInt64() * 1024 * 1024 * 1024
		updateImage.MinDisk.Set(&val)
	}

	// min_ram
	if !plan.MinRam.IsNull() {
		val := plan.MinRam.ValueInt64()
		updateImage.MinRamGB.Set(&val)
	}

	// name
	if !plan.Name.IsNull() {
		updateImage.Name = &name
	}

	// os_type_id
	if !plan.OsTypeId.IsNull() {
		updateImage.OsType.Set(plan.OsTypeId.ValueInt64Pointer())
	}

	// ssh_password_wo
	if !plan.SshPasswordWoVersion.Equal(state.SshPasswordWoVersion) {
		if config.SshPasswordWo.IsUnknown() {
			resp.Diagnostics.AddError(
				"update image resource",
				fmt.Sprintf("image %s: 'ssh_password_wo_version' changed, "+
					"but 'ssh_password_wo' is not set", name),
			)

			return
		}
		updateImage.SshPassword.Set(config.SshPasswordWo.ValueStringPointer())
	}

	// ssh_username
	if !plan.SshUsername.IsNull() {
		updateImage.SshUsername.Set(plan.SshUsername.ValueStringPointer())
	}

	// storage_provider_id
	if !plan.StorageProviderId.IsNull() {
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

	// tenant_ids
	if !plan.TenantIds.IsNull() {
		tenantIDs, err := convert.SetToSlice[types.Int64](ctx, plan.TenantIds)
		if err != nil {
			log.Fatal(err)
		}

		for _, idVal := range tenantIDs {
			if !idVal.IsNull() {
				updateImage.Accounts = append(
					updateImage.Accounts,
					idVal.ValueInt64(),
				)
			}
		}
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
	if !plan.UserData.IsNull() {
		updateImage.UserData.Set(plan.UserData.ValueStringPointer())
	}

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
	if image.VirtualImage != nil {
		plan.Id = convert.Int64ToType(image.VirtualImage.Id)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	state, diag := getImageAsState(ctx, plan.Id.ValueInt64(), client, plan)
	if resp.Diagnostics.Append(diag...); resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}
