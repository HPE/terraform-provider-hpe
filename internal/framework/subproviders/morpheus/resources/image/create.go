// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

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

	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/convert"
	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/errors"
	"github.com/HPE/terraform-provider-hpe/internal/framework/utils"
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

	reqImage := sdk.NewAddVirtualImageRequestWithDefaults()
	reqImage.VirtualImage = *sdk.NewAddVirtualImageRequestVirtualImageWithDefaults()

	// auto_join_domain
	if !plan.AutoJoinDomain.IsNull() && !plan.AutoJoinDomain.IsUnknown() {
		reqImage.VirtualImage.SetIsAutoJoinDomain(plan.AutoJoinDomain.ValueBool())
	}

	// cloud_init
	if !plan.CloudInit.IsNull() && !plan.CloudInit.IsUnknown() {
		reqImage.VirtualImage.SetIsCloudInit(plan.CloudInit.ValueBool())
	}

	// config_azure
	if !plan.ConfigAzure.IsNull() && !plan.ConfigAzure.IsUnknown() {
		config := sdk.AddVirtualImageRequestVirtualImageConfig{}
		config.AzureReferenceVirtualImageConfiguration = sdk.NewAzureReferenceVirtualImageConfigurationWithDefaults()
		config.AzureReferenceVirtualImageConfiguration.SetPublisher(plan.ConfigAzure.Publisher.ValueString())
		config.AzureReferenceVirtualImageConfiguration.SetOffer(plan.ConfigAzure.Offer.ValueString())
		config.AzureReferenceVirtualImageConfiguration.SetVersion(plan.ConfigAzure.Version.ValueString())
		config.AzureReferenceVirtualImageConfiguration.SetSku(plan.ConfigAzure.Sku.ValueString())
		reqImage.VirtualImage.SetConfig(config)
	}

	// description
	if !plan.Description.IsNull() && !plan.Description.IsUnknown() {
		reqImage.VirtualImage.SetDescription(plan.Description.ValueString())
	}

	// fips_enabled
	if !plan.FipsEnabled.IsNull() && !plan.FipsEnabled.IsUnknown() {
		reqImage.VirtualImage.SetFipsEnabled(plan.FipsEnabled.ValueBool())
	}

	// is_force_customization
	if !plan.ForceCustomization.IsNull() && !plan.ForceCustomization.IsUnknown() {
		reqImage.VirtualImage.SetIsForceCustomization(plan.ForceCustomization.ValueBool())
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
	reqImage.VirtualImage.SetImageType(plan.ImageType.ValueString())

	// install_agent
	if !plan.InstallAgent.IsNull() && !plan.InstallAgent.IsUnknown() {
		reqImage.VirtualImage.SetInstallAgent(plan.InstallAgent.ValueBool())
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

		reqImage.VirtualImage.SetLabels(labels)
	}

	// min_disk
	if !plan.MinDisk.IsNull() && !plan.MinDisk.IsUnknown() {
		reqImage.VirtualImage.SetMinDisk(plan.MinDisk.ValueInt64() * 1024 * 1024 * 1024)
	}

	// min_ram
	if !plan.MinRam.IsNull() && !plan.MinRam.IsUnknown() {
		reqImage.VirtualImage.SetMinRam(plan.MinRam.ValueInt64() * 1024 * 1024 * 1024)
	}

	// name (required)
	reqImage.VirtualImage.SetName(plan.Name.ValueString())

	// os_type_id
	if !plan.OsTypeId.IsNull() && !plan.OsTypeId.IsUnknown() {
		reqImage.VirtualImage.SetOsType(plan.OsTypeId.ValueInt64())
	}

	// ssh_password_wo
	if !config.SshPasswordWo.IsNull() && !config.SshPasswordWo.IsUnknown() {
		reqImage.VirtualImage.SetSshPassword(config.SshPasswordWo.ValueString())
	}

	// ssh_username
	if !plan.SshUsername.IsNull() && !plan.SshUsername.IsUnknown() {
		reqImage.VirtualImage.SetSshUsername(plan.SshUsername.ValueString())
	}

	// storage_provider_id
	if !plan.StorageProviderId.IsNull() && !plan.StorageProviderId.IsUnknown() {
		storageProvider := sdk.NewAddVirtualImageRequestVirtualImageStorageProviderWithDefaults()
		storageProvider.SetId(plan.StorageProviderId.ValueInt64())

		reqImage.VirtualImage.SetStorageProvider(*storageProvider)
	}

	// sysprep
	if !plan.Sysprep.IsNull() && !plan.Sysprep.IsUnknown() {
		reqImage.VirtualImage.SetIsSysprep(plan.Sysprep.ValueBool())
	}

	// tags
	tags, diags := convert.FromSetType(ctx, plan.Tags, tagMapper)
	if diags.HasError() {
		tflog.Error(ctx, "cannot convert volumes")
		resp.Diagnostics.Append(diags...)

		return
	}
	reqImage.VirtualImage.SetTags(tags)

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
		reqImage.VirtualImage.SetTrialVersion(plan.TrialVersion.ValueBool())
	}

	// uefi
	if !plan.Uefi.IsNull() && !plan.Uefi.IsUnknown() {
		reqImage.VirtualImage.SetUefi(plan.Uefi.ValueBool())
	}

	// url
	if !plan.Url.IsNull() && !plan.Url.IsUnknown() {
		reqImage.VirtualImage.SetUrl(plan.Url.ValueString())
	}

	// user_data
	if !plan.UserData.IsNull() && !plan.UserData.IsUnknown() {
		reqImage.VirtualImage.UserData.Set(plan.UserData.ValueStringPointer())
	}

	// virtio_supported
	if !plan.VirtioSupported.IsNull() && !plan.VirtioSupported.IsUnknown() {
		reqImage.VirtualImage.SetVirtioSupported(plan.VirtioSupported.ValueBool())
	}

	// visibility
	if !plan.Visibility.IsNull() && !plan.Visibility.IsUnknown() {
		reqImage.VirtualImage.SetVisibility(plan.Visibility.ValueString())
	}

	// vm_tools_installed
	if !plan.VmToolsInstalled.IsNull() && !plan.VmToolsInstalled.IsUnknown() {
		reqImage.VirtualImage.SetVmToolsInstalled(plan.VmToolsInstalled.ValueBool())
	}

	image, httpResp, err := client.LibraryAPI.AddVirtualImage(ctx).
		AddVirtualImageRequest(*reqImage).
		Execute()

	if err != nil || httpResp.StatusCode != http.StatusOK {
		resp.Diagnostics.AddError("error creating image", errors.ErrMsg(err, httpResp))

		return
	}

	imageId := image.VirtualImage.GetId()

	// Helper to set partial state on error
	setPartialState := func(id int64) {
		utils.SetPartialState(ctx, utils.SetPartialStateConfig{
			ResourceType: "image",
			ResourceID:   id,
			StateWriter:  &resp.State,
			Diagnostics:  &resp.Diagnostics,
		})
	}

	// Set state
	plan.Id = convert.Int64ToType(&imageId)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError(
			"failed to set initial image state",
			fmt.Sprintf("Image %d was created but state could not be saved", imageId),
		)
		setPartialState(imageId)

		return
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
		if err != nil || httpResp.StatusCode != http.StatusOK {
			return "", backoff.Permanent(err)
		}

		status := resp.VirtualImage.GetStatus()

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
		resp.Diagnostics.AddError(
			"create image resource",
			fmt.Sprintf("image %s: creation failed current status is: %s", plan.Name.ValueString(), status),
		)
		utils.SetPartialState(ctx, utils.SetPartialStateConfig{
			ResourceType: "image",
			ResourceID:   plan.Id.ValueInt64(),
			StateWriter:  &resp.State,
			Diagnostics:  &resp.Diagnostics,
		})

		return
	}

	state, diag := getImageAsState(ctx, plan.Id.ValueInt64(), client, plan)
	if resp.Diagnostics.Append(diag...); resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError(
			"failed to read image state",
			fmt.Sprintf("Image %d was created but could not be read", plan.Id.ValueInt64()),
		)
		utils.SetPartialState(ctx, utils.SetPartialStateConfig{
			ResourceType: "image",
			ResourceID:   plan.Id.ValueInt64(),
			StateWriter:  &resp.State,
			Diagnostics:  &resp.Diagnostics,
		})

		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError(
			"failed to set image state",
			fmt.Sprintf("Image %d was created but state could not be saved", plan.Id.ValueInt64()),
		)
		utils.SetPartialState(ctx, utils.SetPartialStateConfig{
			ResourceType: "image",
			ResourceID:   plan.Id.ValueInt64(),
			StateWriter:  &resp.State,
			Diagnostics:  &resp.Diagnostics,
		})

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

func tagMapper(
	in TagsValue,
) sdk.AddVirtualImageRequestVirtualImageTagsInner {
	return sdk.AddVirtualImageRequestVirtualImageTagsInner{
		Name:  in.Name.ValueString(),
		Value: in.Value.ValueString(),
	}
}
