// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

// Package hostprofile implements a data source for a single host profile.
package hostprofile

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/datasource"

	"github.com/HPE/terraform-provider-hpe/morpheus/configure"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
)

const summary = "read host_profile data source"

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
	resp.TypeName = req.ProviderTypeName + "_" + "host_profile"
}

// Schema defines the schema for the data source.
func (d *DataSource) Schema(
	ctx context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = HostProfileDataSourceSchema(ctx)
}

// Read refreshes the Terraform state with the latest data.
func (d *DataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	var config HostProfileModel

	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	client, err := d.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError(summary, "could not create sdk client")

		return
	}

	id := config.Id.ValueInt64()

	r, httpResp, err := client.HostProfilesAPI.GetHostProfile(ctx, id).Execute()
	if r == nil || httpResp == nil || err != nil || httpResp.StatusCode != http.StatusOK {
		resp.Diagnostics.AddError(
			summary,
			fmt.Sprintf("GET host_profile %d failed: %s", id, errfmt.ErrMsg(err, httpResp)),
		)

		return
	}

	if r.HostProfile == nil {
		resp.Diagnostics.AddError(summary, fmt.Sprintf("host_profile %d: response missing hostProfile", id))

		return
	}

	state, mapDiags := hostProfileAsState(ctx, r.HostProfile)
	resp.Diagnostics.Append(mapDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
