// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package image

import (
	"context"
	"fmt"
	"net/http"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/convert"
	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/errors"
)

func getImageAsState(
	ctx context.Context,
	id int64,
	client *sdk.APIClient,
	plan ImageModel,
) (ImageModel, diag.Diagnostics) {
	var state ImageModel
	var diags diag.Diagnostics

	resp, httpResp, err := client.LibraryAPI.GetVirtualImage(ctx, id).Execute()
	if err != nil || httpResp.StatusCode != http.StatusOK {
		diags.AddError(
			"populate image resource",
			fmt.Sprintf("image %d GET failed: ", id)+errors.ErrMsg(err, httpResp),
		)

		return state, diags
	}

	image := resp.GetVirtualImage()

	// auto_join_domain
	state.AutoJoinDomain = convert.BoolToType(image.IsAutoJoinDomain)

	// config_azure
	if image.GetImageType() == "azure-reference" {
		if image.GetConfig().AzureReferenceVirtualImageConfiguration != nil {
			state.ConfigAzure.Publisher = convert.StrToType(&image.GetConfig().AzureReferenceVirtualImageConfiguration.Publisher)
			state.ConfigAzure.Offer = convert.StrToType(&image.GetConfig().AzureReferenceVirtualImageConfiguration.Offer)
			state.ConfigAzure.Version = convert.StrToType(&image.GetConfig().AzureReferenceVirtualImageConfiguration.Version)
			state.ConfigAzure.Sku = convert.StrToType(&image.GetConfig().AzureReferenceVirtualImageConfiguration.Sku)
		}

		state.ConfigAzure.state = attr.ValueStateKnown
	} else {
		state.ConfigAzure = NewConfigAzureValueNull()
	}

	// cloud_init
	state.CloudInit = convert.BoolToType(image.IsCloudInit)

	// description
	state.Description = convert.StrToType(image.Description.Get())

	// fips_enabled
	state.FipsEnabled = convert.BoolToType(image.FipsEnabled)

	// force_customization
	state.ForceCustomization = convert.BoolToType(image.IsForceCustomization)

	// external_id
	state.ExternalId = convert.StrToType(image.ExternalId.Get())

	// id
	state.Id = convert.Int64ToType(image.Id)

	// image_type
	state.ImageType = convert.StrToType(image.ImageType)

	// install_agent
	state.InstallAgent = convert.BoolToType(image.InstallAgent)

	// labels
	state.Labels = convert.StrSliceToSet(image.GetLabels())

	// min_disk
	if image.MinDisk.Get() != nil {
		convertToGb := int64(*image.MinDisk.Get() / 1024 / 1024 / 1024)
		state.MinDisk = convert.Int64ToType(&convertToGb)
	} else {
		state.MinDisk = basetypes.NewInt64Null()
	}

	// min_ram
	if image.MinRam.Get() != nil {
		convertToGb := int64(*image.MinRam.Get() / 1024 / 1024 / 1024)
		state.MinRam = convert.Int64ToType(&convertToGb)
	} else {
		state.MinRam = basetypes.NewInt64Null()
	}

	// name
	state.Name = convert.StrToType(image.Name)

	// os_type_id
	state.OsTypeId = convert.Int64ToType(image.GetOsType().Id)

	// owner_id
	state.OwnerId = convert.Int64ToType(image.OwnerId)

	// raw_size
	state.RawSize = convert.Int64ToType(image.RawSize.Get())

	// ssh_username
	state.SshUsername = convert.StrToType(image.SshUsername.Get())

	// ssh_key_wo_version
	state.SshKeyWoVersion = plan.SshKeyWoVersion

	// ssh_password_wo_version
	state.SshPasswordWoVersion = plan.SshPasswordWoVersion

	// storage_provider_id
	state.StorageProviderId = convert.Int64ToType(image.GetStorageProvider().Id)

	// sysprep
	state.Sysprep = convert.BoolToType(image.IsSysprep)

	// system_image
	state.SystemImage = convert.BoolToType(image.SystemImage)

	// tags
	tags, d := convert.ToSetType(
		ctx,
		image.Tags,
		func(
			in sdk.AddCatalogItemTypeRequestCatalogItemTypeOneOfConfigEvarsInner,
		) TagsValue {
			return TagsValue{
				Name:  convert.StrToType(in.Name),
				Value: convert.StrToType(in.Value),
				state: attr.ValueStateKnown,
			}
		},
	)
	diags.Append(d...)
	state.Tags = tags

	// tenant_ids
	state.TenantIds = types.SetNull(types.Int64Type)
	if len(image.Accounts) > 0 {
		var tenantValues []attr.Value
		for _, tenant := range image.Accounts {
			if tenant.Id != nil {
				tenantValues = append(tenantValues, types.Int64Value(*tenant.Id))
			}
		}
		if len(tenantValues) > 0 {
			tenantSet, d := types.SetValue(
				types.Int64Type, tenantValues,
			)
			diags.Append(d...)
			if diags.HasError() {
				return state, diags
			}
			state.TenantIds = tenantSet
		}
	}

	// trial_version
	state.TrialVersion = convert.BoolToType(image.TrialVersion)

	// uefi
	state.Uefi = convert.BoolToType(image.Uefi.Get())

	// url
	state.Url = plan.Url

	// user_data
	state.UserData = convert.StrToType(image.UserData.Get())

	// virtio_supported
	state.VirtioSupported = convert.BoolToType(image.VirtioSupported)

	// visibility
	state.Visibility = convert.StrToType(image.Visibility)

	// vm_tools_installed
	state.VmToolsInstalled = convert.BoolToType(image.VmToolsInstalled)

	return state, diags
}

// Read implements resource.Resource.
func (r *Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data ImageModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client, err := r.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError("creating client failed", err.Error())

		return
	}

	state, diag := getImageAsState(ctx, data.Id.ValueInt64(), client, data)
	if resp.Diagnostics.Append(diag...); resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}
