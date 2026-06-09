// (C) Copyright 2025-2026 Hewlett Packard Enterprise Development LP

package image

import (
	"context"
	stderr "errors"
	"fmt"
	"net/http"
	"slices"
	"time"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"
	"github.com/cenkalti/backoff/v5"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/utils/cleanup"
	"github.com/HPE/terraform-provider-hpe/utils/convert"
)

// Create implements resource.Resource.
func (r *Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan, config ImageModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
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

	reqImage := &sdk.AddVirtualImageRequest{}
	reqImage.VirtualImage = sdk.AddVirtualImageRequestVirtualImage{}

	// auto_join_domain
	if !plan.AutoJoinDomain.IsNull() && !plan.AutoJoinDomain.IsUnknown() {
		reqImage.VirtualImage.IsAutoJoinDomain = plan.AutoJoinDomain.ValueBoolPointer()
	}

	// cloud_init
	if !plan.CloudInit.IsNull() && !plan.CloudInit.IsUnknown() {
		reqImage.VirtualImage.IsCloudInit = plan.CloudInit.ValueBoolPointer()
	}

	// config_azure
	if !plan.ConfigAzure.IsNull() && !plan.ConfigAzure.IsUnknown() {
		azureConfig := &sdk.AzureReferenceVirtualImageConfiguration1{}
		azureConfig.Publisher = plan.ConfigAzure.Publisher.ValueString()
		azureConfig.Offer = plan.ConfigAzure.Offer.ValueString()
		azureConfig.Version = plan.ConfigAzure.Version.ValueString()
		azureConfig.Sku = plan.ConfigAzure.Sku.ValueString()

		config := sdk.AddVirtualImageRequestVirtualImageConfig{
			AzureReferenceVirtualImageConfiguration1: azureConfig,
		}
		reqImage.VirtualImage.Config = &config
	}

	// description
	if !plan.Description.IsNull() && !plan.Description.IsUnknown() {
		reqImage.VirtualImage.Description = plan.Description.ValueStringPointer()
	}

	// fips_enabled
	if !plan.FipsEnabled.IsNull() && !plan.FipsEnabled.IsUnknown() {
		reqImage.VirtualImage.FipsEnabled = plan.FipsEnabled.ValueBoolPointer()
	}

	// is_force_customization
	if !plan.ForceCustomization.IsNull() && !plan.ForceCustomization.IsUnknown() {
		reqImage.VirtualImage.IsForceCustomization = plan.ForceCustomization.ValueBoolPointer()
	}

	// uncomment when adding support for uploading files
	// file upload will need to be after image creation, we only want to
	// verify here that the file exists before we submit the API request
	// if !plan.File.IsNull() && !plan.File.IsUnknown() {
	// 	if _, err := os.Stat(plan.File.ValueString()); err != nil {
	// 		resp.Diagnostics.AddError("could not read the specified image file", err.Error())

	// 		return
	// 	}
	// }

	// image_type (required)
	reqImage.VirtualImage.ImageType = plan.ImageType.ValueStringPointer()

	// install_agent
	if !plan.InstallAgent.IsNull() && !plan.InstallAgent.IsUnknown() {
		reqImage.VirtualImage.InstallAgent = plan.InstallAgent.ValueBoolPointer()
	}

	// labels
	if !plan.Labels.IsNull() && !plan.Labels.IsUnknown() {
		labels, err := convert.SetToStrSlice(plan.Labels)
		if err != nil {
			resp.Diagnostics.AddError(
				"create image resource",
				"image "+plan.Name.ValueString()+": failed to parse label: "+err.Error(),
			)

			return
		}

		reqImage.VirtualImage.Labels = labels
	}

	// min_disk
	if !plan.MinDisk.IsNull() && !plan.MinDisk.IsUnknown() {
		val := plan.MinDisk.ValueInt64() * 1024 * 1024 * 1024
		reqImage.VirtualImage.MinDisk = *sdk.NewNullableInt64(&val)
	}

	// min_ram
	if !plan.MinRam.IsNull() && !plan.MinRam.IsUnknown() {
		val := plan.MinRam.ValueInt64() * 1024 * 1024 * 1024
		reqImage.VirtualImage.MinRam = *sdk.NewNullableInt64(&val)
	}

	// name (required)
	reqImage.VirtualImage.Name = plan.Name.ValueStringPointer()

	// os_type_id
	if !plan.OsTypeId.IsNull() && !plan.OsTypeId.IsUnknown() {
		reqImage.VirtualImage.OsType.Set(plan.OsTypeId.ValueInt64Pointer())
	}

	// ssh_password_wo
	if !config.SshPasswordWo.IsNull() && !config.SshPasswordWo.IsUnknown() {
		sshPwd := config.SshPasswordWo.ValueString()
		reqImage.VirtualImage.SshPassword.Set(&sshPwd)
	}

	// ssh_username
	if !plan.SshUsername.IsNull() && !plan.SshUsername.IsUnknown() {
		sshUser := plan.SshUsername.ValueString()
		reqImage.VirtualImage.SshUsername.Set(&sshUser)
	}

	// storage_provider_id
	if !plan.StorageProviderId.IsNull() && !plan.StorageProviderId.IsUnknown() {
		storageProvider := &sdk.AddVirtualImageRequestVirtualImageStorageProvider{}
		storageProvider.Id = plan.StorageProviderId.ValueInt64Pointer()

		reqImage.VirtualImage.StorageProvider = storageProvider
	}

	// sysprep
	if !plan.Sysprep.IsNull() && !plan.Sysprep.IsUnknown() {
		reqImage.VirtualImage.IsSysprep = plan.Sysprep.ValueBoolPointer()
	}

	// tags
	tags, diags := convert.FromSetType(ctx, plan.Tags, createTagMapper)
	if diags.HasError() {
		tflog.Error(ctx, "cannot convert volumes")
		resp.Diagnostics.Append(diags...)

		return
	}
	reqImage.VirtualImage.Tags = tags

	// tenant_id
	if !plan.TenantIds.IsNull() && !plan.TenantIds.IsUnknown() {
		var tenantIDs []types.Int64
		diags := plan.TenantIds.ElementsAs(ctx, &tenantIDs, false)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		for _, idVal := range tenantIDs {
			if !idVal.IsNull() {
				reqImage.VirtualImage.Accounts = append(
					reqImage.VirtualImage.Accounts,
					idVal.ValueInt64(),
				)
			}
		}
	}

	// trial_version
	if !plan.TrialVersion.IsNull() && !plan.TrialVersion.IsUnknown() {
		reqImage.VirtualImage.TrialVersion = plan.TrialVersion.ValueBoolPointer()
	}

	// uefi
	if !plan.Uefi.IsNull() && !plan.Uefi.IsUnknown() {
		reqImage.VirtualImage.Uefi = plan.Uefi.ValueBoolPointer()
	}

	// url
	if !plan.Url.IsNull() && !plan.Url.IsUnknown() {
		reqImage.VirtualImage.Url = plan.Url.ValueStringPointer()
	}

	// user_data
	if !plan.UserData.IsNull() && !plan.UserData.IsUnknown() {
		reqImage.VirtualImage.UserData.Set(plan.UserData.ValueStringPointer())
	}

	// virtio_supported
	if !plan.VirtioSupported.IsNull() && !plan.VirtioSupported.IsUnknown() {
		reqImage.VirtualImage.VirtioSupported = plan.VirtioSupported.ValueBoolPointer()
	}

	// visibility
	if !plan.Visibility.IsNull() && !plan.Visibility.IsUnknown() {
		reqImage.VirtualImage.Visibility = plan.Visibility.ValueStringPointer()
	}

	// vm_tools_installed
	if !plan.VmToolsInstalled.IsNull() && !plan.VmToolsInstalled.IsUnknown() {
		reqImage.VirtualImage.VmToolsInstalled = plan.VmToolsInstalled.ValueBoolPointer()
	}

	image, httpResp, err := client.LibraryAPI.AddVirtualImage(ctx).
		AddVirtualImageRequest(*reqImage).
		Execute()

	if err != nil || httpResp.StatusCode != http.StatusOK {
		resp.Diagnostics.AddError("error creating image", errfmt.ErrMsg(err, httpResp))

		return
	}

	if image.VirtualImage == nil {
		resp.Diagnostics.AddError("API returned nil", "VirtualImage is nil in the response")

		return
	}

	if image.VirtualImage.Id == nil {
		resp.Diagnostics.AddError("API returned nil", "VirtualImage ID is nil in the response")

		return
	}

	plan.Id = convert.Int64ToType(image.VirtualImage.Id)

	// Helper to taint the resource state on an error after the POST request
	taintResourceState := func(id int64) {
		cleanup.TaintResourceState(ctx, cleanup.TaintResourceStateConfig{
			ResourceType: "image",
			ResourceID:   id,
			StateWriter:  &resp.State,
			Diagnostics:  &resp.Diagnostics,
		})
	}

	// rough template for what needs to be done to support file uploading
	// if !plan.File.IsNull() && !plan.File.IsUnknown() {
	// 	if diags := uploadImageFile(ctx, client, plan.Id.ValueInt64(), plan.File.ValueString()); diags.HasError() {
	// 		resp.Diagnostics = append(resp.Diagnostics, diags...)

	// 		return
	// 	}
	// file, err := os.Open(plan.File.ValueString())
	// if err != nil {
	// 	resp.Diagnostics.AddError("could not read the specified image file", err.Error())

	// 	return
	// }

	// _, httpResp, err := client.LibraryAPI.AddVirtualImageFile(
	// 	ctx,
	// 	plan.Id.ValueInt64(),
	// ).Filename(plan.File.ValueString()).Body(file).Execute()
	// if err != nil || httpResp.StatusCode != http.StatusOK {
	// 	resp.Diagnostics.AddError("error uploading image file", errors.ErrMsg(err, httpResp))

	// 	return
	// }
	// }

	waitForReady := func() (string, error) {
		resp, httpResp, err := client.LibraryAPI.GetVirtualImage(ctx, plan.Id.ValueInt64()).Execute()
		if err != nil {
			if httpResp == nil || httpResp.StatusCode != http.StatusOK {
				return "", backoff.Permanent(err)
			}
		}

		status := ""
		if resp.VirtualImage != nil && resp.VirtualImage.Status != nil {
			status = *resp.VirtualImage.Status
		}

		return status, checkStatusDone(
			status,
			CreateTargetStatuses,
			CreateErrorStatuses,
		)
	}

	if status, err := backoff.Retry(
		ctx,
		waitForReady,
		backoff.WithBackOff(backoff.NewConstantBackOff(5*time.Second)),
		backoff.WithMaxElapsedTime(45*time.Minute),
	); err != nil {
		if status == "" {
			errUnwrapped := stderr.Unwrap(err)
			if errUnwrapped != nil {
				resp.Diagnostics.AddError(
					"image provisioning failed",
					fmt.Sprintf("Image %d failed to reach provisioned status: %v", plan.Id.ValueInt64(), errUnwrapped),
				)
			} else {
				resp.Diagnostics.AddError(
					"image provisioning failed",
					fmt.Sprintf("Image %d failed to reach provisioned status - unknown error.", plan.Id.ValueInt64()),
				)
			}
		} else {
			resp.Diagnostics.AddError(
				"create image resource",
				fmt.Sprintf("image %s: creation failed current status is: %s", plan.Name.ValueString(), status),
			)
		}
		taintResourceState(plan.Id.ValueInt64())

		return
	}

	state, diag := getImageAsState(ctx, plan.Id.ValueInt64(), client, plan)
	if resp.Diagnostics.Append(diag...); resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError(
			"failed to read image state",
			fmt.Sprintf("Image %d was created but could not be read", plan.Id.ValueInt64()),
		)
		taintResourceState(plan.Id.ValueInt64())

		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError(
			"failed to set image state",
			fmt.Sprintf("Image %d was created but state could not be saved", plan.Id.ValueInt64()),
		)
		taintResourceState(plan.Id.ValueInt64())

		return
	}
}

func checkStatusDone(status string, targetStatuses []string, errorStatuses []string) error {
	switch {
	case slices.Contains(errorStatuses, status):
		return backoff.Permanent(stderr.New("reached error status: " + status))
	case slices.Contains(targetStatuses, status):
		return nil
	default:
		return backoff.RetryAfter(5)
	}
}

func createTagMapper(
	in TagsValue,
) sdk.AddVirtualImageRequestVirtualImageTagsInner {
	return sdk.AddVirtualImageRequestVirtualImageTagsInner{
		Name:  in.Name.ValueString(),
		Value: in.Value.ValueString(),
	}
}

func updateTagMapper(
	in TagsValue,
) sdk.UpdateVirtualImageRequestVirtualImageTagsInner {
	return sdk.UpdateVirtualImageRequestVirtualImageTagsInner{
		Name:  in.Name.ValueString(),
		Value: in.Value.ValueString(),
	}
}
