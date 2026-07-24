// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package morpheusdetails

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	vmaascmpclient "github.com/HPE/terraform-provider-hpe/greenlake/cloud/sdk/vmaascmp/client"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &DataSource{}
	_ datasource.DataSourceWithConfigure = &DataSource{}
)

// NewDataSource is a helper function to simplify the provider implementation.
func NewDataSource() datasource.DataSource {
	return &DataSource{}
}

// DataSource is the morpheus_details data source implementation.
type DataSource struct {
	client *vmaascmpclient.APIClient
}

// MorpheusDetailsModel maps the data source schema to a Go struct.
type MorpheusDetailsModel struct {
	ID          types.String `tfsdk:"id"`
	AccessToken types.String `tfsdk:"access_token"`
	Expiry      types.Int64  `tfsdk:"expiry"`
	URL         types.String `tfsdk:"url"`
}

// Metadata returns the data source type name.
func (d *DataSource) Metadata(
	_ context.Context,
	req datasource.MetadataRequest,
	resp *datasource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_morpheus_details"
}

// Schema defines the schema for the data source.
func (d *DataSource) Schema(
	_ context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Description: "Retrieves details of the Morpheus instance used by GreenLake " +
			"VMaaS, including the access token, its expiry, and the Morpheus URL.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Service instance ID of the Morpheus instance.",
				Computed:    true,
			},
			"access_token": schema.StringAttribute{
				Description: "Morpheus access token.",
				Computed:    true,
				Sensitive:   true,
			},
			"expiry": schema.Int64Attribute{
				Description: "Time until the access token expires, in seconds.",
				Computed:    true,
			},
			"url": schema.StringAttribute{
				Description: "Morpheus instance URL.",
				Computed:    true,
			},
		},
	}
}

// Configure adds the configured vmaascmp client to the data source.
func (d *DataSource) Configure(
	ctx context.Context,
	req datasource.ConfigureRequest,
	resp *datasource.ConfigureResponse,
) {
	// provider.Configure is not guaranteed to have run yet.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*vmaascmpclient.APIClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected provider data type",
			"Expected *vmaascmpclient.APIClient from the greenlake_cloud provider "+
				"configuration. This is a bug in the provider.",
		)

		return
	}

	d.client = client
}

// Read refreshes the Terraform state with the latest Morpheus details.
func (d *DataSource) Read(
	ctx context.Context,
	_ datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	morpheusDetails, err := d.client.GetCMPDetails(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to read Morpheus details",
			err.Error(),
		)

		return
	}

	state := MorpheusDetailsModel{
		ID:          types.StringValue(morpheusDetails.ID),
		AccessToken: types.StringValue(morpheusDetails.AccessToken),
		Expiry:      types.Int64Value(morpheusDetails.ValidTill),
		URL:         types.StringValue(morpheusDetails.URL),
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
