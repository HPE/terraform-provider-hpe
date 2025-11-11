package image

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/convert"
	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/errors"
	"github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
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

	updateRequest := sdk.NewUpdateVirtualImageRequestWithDefaults()
	updateImage := sdk.NewUpdateVirtualImageRequestVirtualImageWithDefaults()

	// auto_join_domain
	if !plan.AutoJoinDomain.IsNull() {
		updateImage.SetIsAutoJoinDomain(plan.AutoJoinDomain.ValueBool())
	}

	// cloud_init
	if !plan.CloudInit.IsNull() {
		updateImage.SetIsCloudInit(plan.CloudInit.ValueBool())
	}

	// config_azure
	// keep unimplemented for now, we can decide in the future if this
	// should continue being a replace or if we want to support in-place updates
	// of the Azure image config

	// description
	// description for whatever reason cannot be updated at the moment

	// fips_enabled
	if !plan.FipsEnabled.IsNull() {
		updateImage.SetFipsEnabled(plan.FipsEnabled.ValueBool())
	}

	// force_customization
	if !plan.ForceCustomization.IsNull() {
		updateImage.SetIsForceCustomization(plan.ForceCustomization.ValueBool())
	}

	// image_type
	if !plan.ImageType.IsNull() {
		updateImage.SetImageType(plan.ImageType.ValueString())
	}

	// install_agent
	if !plan.InstallAgent.IsNull() {
		updateImage.SetInstallAgent(plan.InstallAgent.ValueBool())
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

		updateImage.SetLabels(labels)
	}

	// min_disk
	if !plan.MinDisk.IsNull() {
		updateImage.SetMinDisk(plan.MinDisk.ValueInt64() * 1024 * 1024 * 1024)

	}

	// min_ram
	if !plan.MinRam.IsNull() {
		updateImage.SetMinRamGB(plan.MinRam.ValueInt64())
	}

	// name
	if !plan.Name.IsNull() {
		updateImage.SetName(name)
	}

	// os_type_id
	if !plan.OsTypeId.IsNull() {
		updateImage.SetOsType(plan.OsTypeId.ValueInt64())
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
		updateImage.SetSshPassword(config.SshPasswordWo.ValueString())
	}

	// ssh_username
	if !plan.SshUsername.IsNull() {
		updateImage.SetSshUsername(plan.SshUsername.ValueString())
	}

	// storage_provider_id
	if !plan.StorageProviderId.IsNull() {
		storageProvider := sdk.NewAddVirtualImageRequestVirtualImageStorageProviderWithDefaults()
		storageProvider.SetId(plan.StorageProviderId.ValueInt64())

		updateImage.SetStorageProvider(*storageProvider)
	}

	// sysprep
	if !plan.Sysprep.IsNull() {
		updateImage.SetIsSysprep(plan.Sysprep.ValueBool())
	}

	// tags
	tags, diags := convert.FromSetType(ctx, plan.Tags, tagMapper)
	if diags.HasError() {
		tflog.Error(ctx, "cannot convert tags")
		resp.Diagnostics.Append(diags...)

		return
	}
	updateImage.SetTags(tags)

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
		updateImage.SetTrialVersion(plan.TrialVersion.ValueBool())
	}

	// uefi
	if !plan.Uefi.IsNull() {
		updateImage.SetUefi(plan.Uefi.ValueBool())
	}

	// url
	// keep unimplemented for now, we can decide in the future if this
	// should continue being a replace or if we want to support in-place updates
	// of the image URL

	// user_data
	if !plan.UserData.IsNull() {
		updateImage.SetUserData(plan.UserData.ValueString())
	}

	// virtio_supported
	if !plan.VirtioSupported.IsNull() {
		updateImage.SetVirtioSupported(plan.VirtioSupported.ValueBool())
	}

	// visibility
	if !plan.Visibility.IsNull() {
		updateImage.SetVisibility(plan.Visibility.ValueString())
	}

	// vm_tools_installed
	if !plan.VmToolsInstalled.IsNull() {
		updateImage.SetVmToolsInstalled(plan.VmToolsInstalled.ValueBool())
	}

	// send the API request here
	updateRequest.SetVirtualImage(*updateImage)
	image, httpResp, err := client.LibraryAPI.UpdateVirtualImage(ctx, state.Id.ValueInt64()).
		UpdateVirtualImageRequest(*updateRequest).
		Execute()
	if err != nil || httpResp.StatusCode != http.StatusOK {
		resp.Diagnostics.AddError("error updating image", errors.ErrMsg(err, httpResp))

		return
	}
	// set the ID value in state
	plan.Id = convert.Int64ToType(image.GetVirtualImage().Id)

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
