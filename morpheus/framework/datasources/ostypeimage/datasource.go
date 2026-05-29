// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

// Package ostypeimage implements a data source for os_type_image
package ostypeimage

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/HPE/terraform-provider-hpe/morpheus/configure"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/utils/convert"
)

const summary = "read os_type_image data source"

// Ensure the implementation satisfies the expected interfaces.
var _ datasource.DataSource = &DataSource{}

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
	resp.TypeName = req.ProviderTypeName + "_morpheus_os_type_image"
}

// Schema defines the schema for the data source.
func (d *DataSource) Schema(
	ctx context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = OsTypeImageDataSourceSchema(ctx)
}

// Read refreshes the Terraform state with the latest data.
func (d *DataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	var data OsTypeImageModel

	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	client, err := d.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError(summary, "could not create sdk client")

		return
	}

	osTypeID := data.OsTypeId.ValueInt64()

	osTypeResp, httpResp, err := client.LibraryAPI.GetOsType(ctx, osTypeID).Execute()
	if osTypeResp == nil || err != nil || httpResp.StatusCode != http.StatusOK {
		resp.Diagnostics.AddError(
			summary,
			fmt.Sprintf("GET os_type %d failed: %s", osTypeID, errfmt.ErrMsg(err, httpResp)),
		)

		return
	}

	osType := osTypeResp.OsType

	virtualImageName := data.VirtualImageName.ValueString()

	// Prefer images with a tenant_id (custom images) over system images,
	// if multiple images have the same virtual_image_name
	var matchedSystemImageID int64
	var matchedTenantImageID int64
	for _, img := range osType.Images {
		if *img.VirtualImageName == virtualImageName && *img.Account.Get() > 0 {
			matchedTenantImageID = *img.Id

			break
		} else if *img.VirtualImageName == virtualImageName {
			matchedSystemImageID = *img.Id
		}
	}

	if matchedSystemImageID == 0 && matchedTenantImageID == 0 {
		resp.Diagnostics.AddError(
			summary,
			fmt.Sprintf(
				"no image with name '%s' found on os_type %d",
				virtualImageName, osTypeID,
			),
		)

		return
	}

	preferredImageID := matchedSystemImageID
	if matchedTenantImageID > 0 {
		preferredImageID = matchedTenantImageID
	}

	imgResp, imgHTTPResp, imgErr := client.LibraryAPI.GetOsTypeImage(ctx, preferredImageID).Execute()
	if imgResp == nil || imgErr != nil || imgHTTPResp.StatusCode != http.StatusOK {
		resp.Diagnostics.AddError(
			summary,
			fmt.Sprintf("GET os_type_image %d failed: %s",
				matchedSystemImageID, errfmt.ErrMsg(imgErr, imgHTTPResp)),
		)

		return
	}

	img := imgResp.OsTypeImage

	data.Id = convert.Int64ToType(img.Id)
	data.VirtualImageId = convert.Int64ToType(img.VirtualImageId)
	data.VirtualImageName = types.StringValue(virtualImageName)
	data.OsTypeId = types.Int64Value(osTypeID)

	if img.Zone.IsSet() {
		data.CloudId = types.Int64Value(*img.Zone.Get())
	}

	if img.ComputeZoneType.IsSet() {
		data.CloudTypeId = types.Int64Value(*img.ComputeZoneType.Get())
	}

	if img.ProvisionType.IsSet() {
		data.ProvisionTypeId = types.Int64Value(*img.ProvisionType.Get())
	}

	if img.Account.IsSet() {
		data.TenantId = types.Int64Value(*img.Account.Get())
	}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}
