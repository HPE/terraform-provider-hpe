// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package image

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"

	"github.com/HPE/terraform-provider-hpe/morpheus/configure"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/constants"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/utils/convert"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource = &DataSource{}
)

// NewDataSource is a helper function to simplify the provider implementation.
func NewDataSource() datasource.DataSource {
	return &DataSource{}
}

// DataSource is the data source implementation.
type DataSource struct {
	configure.DataSourceWithMorpheusConfigure
	datasource.DataSource
}

// Metadata returns the data source type name.
func (d *DataSource) Metadata(
	_ context.Context,
	req datasource.MetadataRequest,
	resp *datasource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_" + constants.SubProviderName + "_image"
}

// Schema defines the schema for the data source.
func (d *DataSource) Schema(
	ctx context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = ImageDataSourceSchema(ctx)
}

func (d *DataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	var data ImageModel

	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	client, err := d.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError("creating client failed", err.Error())

		return
	}

	if !data.Name.IsNull() {
		if diags := getImageByName(ctx, &data, client); diags.HasError() {
			resp.Diagnostics.AddError(fmt.Sprintf("failed to get image by name '%s'", data.Name.ValueString()), "")
			resp.Diagnostics = append(resp.Diagnostics, diags...)

			return
		}
	}

	if !data.Id.IsNull() {
		if err := getImageByID(ctx, &data, client, data.Id.ValueInt64()); err != nil {
			resp.Diagnostics.AddError(fmt.Sprintf("failed to get image by id '%d'", data.Id.ValueInt64()), "")
			resp.Diagnostics = append(resp.Diagnostics, diags...)

			return
		}
	}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func getImageByID(
	ctx context.Context,
	data *ImageModel,
	client *sdk.APIClient,
	id int64,
) diag.Diagnostics {
	var diags diag.Diagnostics

	imageResp, httpResp, err := client.LibraryAPI.GetVirtualImage(ctx, id).Execute()
	if imageResp == nil || err != nil || httpResp.StatusCode != http.StatusOK {
		diags.AddError(
			fmt.Sprintf(
				"GET failed for image '%d'", id,
			),
			errfmt.ErrMsg(err, httpResp),
		)

		return diags
	}

	if imageResp.VirtualImage == nil {
		diags.AddError(
			fmt.Sprintf(
				"GET failed for image '%d'", id,
			),
			"image response did not include virtualImage",
		)

		return diags
	}

	image := *imageResp.VirtualImage

	d := parseAsData(ctx, image, data)
	diags = append(diags, d...)

	return diags
}

func getImageByName(
	ctx context.Context,
	data *ImageModel,
	client *sdk.APIClient,
) diag.Diagnostics {
	var diags diag.Diagnostics

	imageListReq := client.LibraryAPI.ListVirtualImages(ctx)

	name := data.Name.ValueString()
	imageListReq = imageListReq.Name(name)

	if !data.ImageType.IsNull() {
		imageListReq = imageListReq.ImageType(data.ImageType.ValueString())
	}

	imageListResp, httpResp, err := imageListReq.Execute()
	if imageListResp == nil || err != nil || httpResp.StatusCode != http.StatusOK {
		diags.AddError(
			fmt.Sprintf("GET failed for image '%s'", name),
			errfmt.ErrMsg(err, httpResp),
		)

		return diags
	}

	var images []sdk.ListVirtualImages200ResponseAllOfVirtualImagesInner

	for _, image := range imageListResp.VirtualImages {
		if image.Name != nil && *image.Name == name {
			if !data.ImageType.IsNull() {
				// skip if image type doesn't match
				if image.ImageType == nil || *image.ImageType != data.ImageType.ValueString() {
					continue
				}
			}
			images = append(images, image)
		}
	}

	if len(images) > 1 {
		diags.AddError(
			"multiple images were returned",
			fmt.Sprintf("datasource expected 1 image, but %d were found", len(images)),
		)

		return diags
	} else if len(images) == 0 {
		diags.AddError(
			"no image found",
			fmt.Sprintf(
				"could not find image by name '%s' and image_type '%s'",
				name,
				data.ImageType.ValueString(),
			),
		)

		return diags
	}

	image := images[0]

	// okay, we need to get a bit silly here
	// despite the structure of both requests being the same our SDK generates
	// two different structs. To avoid code duplication we marshal the image from
	// the list response into JSON, and then unmarshal it into the struct used by
	// the get by ID endpoint.
	imgJSON, err := image.MarshalJSON()
	if err != nil {
		// it is very unlikely for this to error out, if it does it likely indicates
		// an SDK error
		diags.AddError("could not convert image data from API", err.Error())

		return diags
	}

	ma := sdk.GetVirtualImage200ResponseVirtualImage{}
	if err := ma.UnmarshalJSON(imgJSON); err != nil {
		// it is very unlikely for this to error out, if it does it likely indicates
		// an SDK error
		diags.AddError("could not convert image data from API", err.Error())

		return diags
	}

	d := parseAsData(ctx, ma, data)
	diags = append(diags, d...)

	return nil
}

func parseAsData(
	ctx context.Context,
	image sdk.GetVirtualImage200ResponseVirtualImage,
	data *ImageModel,
) diag.Diagnostics {
	var diags diag.Diagnostics

	// auto_join_domain
	data.AutoJoinDomain = convert.BoolToType(image.IsAutoJoinDomain)

	// cloud_init
	data.CloudInit = convert.BoolToType(image.IsCloudInit)

	// config_azure
	if image.ImageType != nil && *image.ImageType == "azure-reference" {
		if image.Config != nil && image.Config.AzureReferenceVirtualImageConfiguration3 != nil {
			data.ConfigAzure.Publisher = convert.StrToType(&image.Config.AzureReferenceVirtualImageConfiguration3.Publisher)
			data.ConfigAzure.Offer = convert.StrToType(&image.Config.AzureReferenceVirtualImageConfiguration3.Offer)
			data.ConfigAzure.Version = convert.StrToType(&image.Config.AzureReferenceVirtualImageConfiguration3.Version)
			data.ConfigAzure.Sku = convert.StrToType(&image.Config.AzureReferenceVirtualImageConfiguration3.Sku)
		}

		data.ConfigAzure.state = attr.ValueStateKnown
	} else {
		data.ConfigAzure = NewConfigAzureValueNull()
	}

	// console_keymap
	data.ConsoleKeymap = convert.StrToType(image.ConsoleKeymap.Get())

	// description
	data.Description = convert.StrToType(image.Description.Get())

	// external_id
	data.ExternalId = convert.StrToType(image.ExternalId.Get())

	// fips_enabled
	data.FipsEnabled = convert.BoolToType(image.FipsEnabled)

	// force_customization
	data.ForceCustomization = convert.BoolToType(image.IsForceCustomization)

	// id
	data.Id = convert.Int64ToType(image.Id)

	// image_type
	data.ImageType = convert.StrToType(image.ImageType)

	// install_agent
	data.InstallAgent = convert.BoolToType(image.InstallAgent)

	// lables
	data.Labels = convert.StrSliceToSet(image.Labels)

	// min_disk
	if image.MinDisk.Get() != nil {
		convertToGb := *image.MinDisk.Get() / 1024 / 1024 / 1024
		data.MinDisk = convert.Int64ToType(&convertToGb)
	} else {
		data.MinDisk = basetypes.NewInt64Null()
	}

	// min_ram
	if image.MinRam.Get() != nil {
		convertToGb := *image.MinRam.Get() / 1024 / 1024 / 1024
		data.MinRam = convert.Int64ToType(&convertToGb)
	} else {
		data.MinRam = basetypes.NewInt64Null()
	}

	// name
	data.Name = convert.StrToType(image.Name)

	// os_type_id
	var osTypeID *int64
	if image.OsType != nil {
		osTypeID = image.OsType.Id
	}
	data.OsTypeId = convert.Int64ToType(osTypeID)

	// owner_id
	data.OwnerId = convert.Int64ToType(image.OwnerId)

	// raw_size
	data.RawSize = convert.Int64ToType(image.RawSize.Get())

	// ssh_key
	data.SshKey = convert.StrToType(image.SshKey.Get())

	// ssh_username
	data.SshUsername = convert.StrToType(image.SshUsername.Get())

	// status
	data.Status = convert.StrToType(image.Status)

	// storage_provider_id
	var storageProviderID *int64
	if image.StorageProvider != nil {
		storageProviderID = image.StorageProvider.Id
	}
	data.StorageProviderId = convert.Int64ToType(storageProviderID)

	// sysprep
	data.Sysprep = convert.BoolToType(image.IsSysprep)

	// system_image
	data.SystemImage = convert.BoolToType(image.SystemImage)

	// tags
	tags, d := convert.ToSetType(
		ctx,
		image.Tags,
		func(
			in sdk.GetVirtualImage200ResponseVirtualImageTagsInner,
		) TagsValue {
			return TagsValue{
				Name:  convert.StrToType(in.Name),
				Value: convert.StrToType(in.Value),
				state: attr.ValueStateKnown,
			}
		},
	)
	diags.Append(d...)
	data.Tags = tags

	// tenants
	tenants, d := convert.ToSetType(
		ctx, image.Accounts,
		func(in sdk.GetVirtualImage200ResponseVirtualImageAccountsInner) TenantsValue {
			return TenantsValue{
				Name: convert.StrToType(in.Name),
				Id:   convert.Int64ToType(in.Id),
			}
		},
	)
	diags.Append(d...)
	data.Tenants = tenants

	// trial_version
	data.TrialVersion = convert.BoolToType(image.TrialVersion)

	// uefi
	data.Uefi = convert.BoolToType(image.Uefi.Get())

	// user_data
	data.UserData = convert.StrToType(image.UserData.Get())

	// user_defined
	data.UserDefined = convert.BoolToType(image.UserDefined)

	// user_uploaded
	data.UserUploaded = convert.BoolToType(image.UserUploaded)

	// virtio_supported
	data.UserData = convert.StrToType(image.UserData.Get())

	// visibility
	data.Visibility = convert.StrToType(image.Visibility)

	// vm_tools_installed
	data.VmToolsInstalled = convert.BoolToType(image.VmToolsInstalled)

	return diags
}
