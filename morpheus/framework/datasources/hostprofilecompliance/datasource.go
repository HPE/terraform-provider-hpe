// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

// Package hostprofilecompliance implements a data source reporting host
// profile compliance for a cluster or a single server.
package hostprofilecompliance

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework-validators/datasourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"

	"github.com/HPE/terraform-provider-hpe/morpheus/configure"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
)

const summary = "read host_profile_compliance data source"

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource                     = &DataSource{}
	_ datasource.DataSourceWithConfigValidators = &DataSource{}
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
	resp.TypeName = req.ProviderTypeName + "_" + "host_profile_compliance"
}

// Schema defines the schema for the data source.
func (d *DataSource) Schema(
	ctx context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = HostProfileComplianceDataSourceSchema(ctx)
}

// ConfigValidators enforces that exactly one of cluster_id or server_id is set.
func (d *DataSource) ConfigValidators(
	_ context.Context,
) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{
		datasourcevalidator.ExactlyOneOf(
			path.MatchRoot("cluster_id"),
			path.MatchRoot("server_id"),
		),
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *DataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	var config HostProfileComplianceModel

	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Defensive runtime check mirroring the ExactlyOneOf config validator.
	hasCluster := !config.ClusterId.IsNull() && !config.ClusterId.IsUnknown()
	hasServer := !config.ServerId.IsNull() && !config.ServerId.IsUnknown()

	if hasCluster == hasServer {
		resp.Diagnostics.AddError(
			summary,
			"exactly one of cluster_id or server_id must be set",
		)

		return
	}

	client, err := d.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError(summary, "could not create sdk client")

		return
	}

	if hasCluster {
		clusterID := config.ClusterId.ValueInt64()

		r, httpResp, err := client.HostProfilesAPI.
			GetClusterHostProfileCompliance(ctx, clusterID).Execute()
		if r == nil || err != nil || httpResp == nil ||
			httpResp.StatusCode != http.StatusOK {
			resp.Diagnostics.AddError(
				summary,
				fmt.Sprintf("GET cluster %d host profile compliance failed: %s",
					clusterID, errfmt.ErrMsg(err, httpResp)),
			)

			return
		}

		state, mapDiags := clusterComplianceAsState(ctx, &config, r)
		resp.Diagnostics.Append(mapDiags...)
		if resp.Diagnostics.HasError() {
			return
		}

		resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)

		return
	}

	serverID := config.ServerId.ValueInt64()

	r, httpResp, err := client.HostProfilesAPI.
		GetServerHostProfileCompliance(ctx, serverID).Execute()
	if r == nil || err != nil || httpResp == nil ||
		httpResp.StatusCode != http.StatusOK {
		resp.Diagnostics.AddError(
			summary,
			fmt.Sprintf("GET server %d host profile compliance failed: %s",
				serverID, errfmt.ErrMsg(err, httpResp)),
		)

		return
	}

	state, mapDiags := serverComplianceAsState(ctx, &config, r)
	resp.Diagnostics.Append(mapDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
