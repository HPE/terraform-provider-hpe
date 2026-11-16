// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

// Package clusteravailablenics implements a data source listing the candidate
// physical NICs per host in a cluster.
package clusteravailablenics

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"

	"github.com/HPE/terraform-provider-hpe/morpheus/configure"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/constants"
	providererrors "github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/versioncheck"
	"github.com/HPE/terraform-provider-hpe/utils/convert"
)

const (
	summary = "read cluster available nics data source"

	// ErrorClusterNotFound explains the 404 the API returns both for a missing
	// cluster and for a cluster that is not on a virtual-switch-capable layout.
	ErrorClusterNotFound = `cluster not found or not on a virtual-switch-capable (2.0+) layout`
)

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
	resp.TypeName = req.ProviderTypeName + "_" + "cluster_available_nics"
}

// Schema defines the schema for the data source.
func (d *DataSource) Schema(
	ctx context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = ClusterAvailableNicsDataSourceSchema(ctx)
}

func nicsAsState(
	ctx context.Context,
	nics []sdk.ListClusterAvailableNics200ResponseAvailableNicsInnerNicsInner,
) (types.List, error) {
	var convErr error

	list, diags := convert.ToListType(ctx, nics,
		func(n sdk.ListClusterAvailableNics200ResponseAvailableNicsInnerNicsInner) NicsValue {
			v, vDiags := NewNicsValue(
				NicsValue{}.AttributeTypes(ctx),
				map[string]attr.Value{
					"active":      convert.BoolToType(n.Active.Get()),
					"consumed":    convert.BoolToType(n.Consumed),
					"mac_address": convert.StrToType(n.MacAddress.Get()),
					"name":        convert.StrToType(n.Name),
				},
			)
			if vDiags.HasError() && convErr == nil {
				convErr = fmt.Errorf("error creating nics value")
			}

			return v
		})
	if diags.HasError() || convErr != nil {
		return types.ListNull(NicsValue{}.Type(ctx)), fmt.Errorf("error creating nics list")
	}

	return list, nil
}

func availableNicsAsState(
	ctx context.Context,
	hosts []sdk.ListClusterAvailableNics200ResponseAvailableNicsInner,
) (types.List, error) {
	var convErr error

	list, diags := convert.ToListType(ctx, hosts,
		func(h sdk.ListClusterAvailableNics200ResponseAvailableNicsInner) AvailableNicsValue {
			nics, err := nicsAsState(ctx, h.Nics)
			if err != nil && convErr == nil {
				convErr = err
			}

			v, vDiags := NewAvailableNicsValue(
				AvailableNicsValue{}.AttributeTypes(ctx),
				map[string]attr.Value{
					"nics":        nics,
					"server_id":   convert.Int64ToType(h.ServerId),
					"server_name": convert.StrToType(h.ServerName),
				},
			)
			if vDiags.HasError() && convErr == nil {
				convErr = fmt.Errorf("error creating available_nics value")
			}

			return v
		})
	if diags.HasError() || convErr != nil {
		return types.ListNull(AvailableNicsValue{}.Type(ctx)), fmt.Errorf("error creating available_nics list")
	}

	return list, nil
}

// Read refreshes the Terraform state with the latest data.
func (d *DataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	var config ClusterAvailableNicsModel

	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiClient, err := d.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError(summary, "could not create sdk client")

		return
	}

	resp.Diagnostics.Append(versioncheck.Require(
		ctx, apiClient, "Virtual switches", constants.VirtualSwitchMinVersion,
	)...)
	if resp.Diagnostics.HasError() {
		return
	}

	clusterID := config.ClusterId.ValueInt64()

	rs, hresp, err := apiClient.VirtualSwitchesAPI.ListClusterAvailableNics(ctx, clusterID).Execute()
	if hresp != nil && hresp.StatusCode == http.StatusNotFound {
		resp.Diagnostics.AddError(summary, fmt.Sprintf("%s (cluster %d)", ErrorClusterNotFound, clusterID))

		return
	}

	if rs == nil || err != nil || hresp.StatusCode != http.StatusOK {
		resp.Diagnostics.AddError(summary, fmt.Sprintf(
			"GET failed for available NICs on cluster %d: %s",
			clusterID, providererrors.ErrMsg(err, hresp),
		))

		return
	}

	list, err := availableNicsAsState(ctx, rs.AvailableNics)
	if err != nil {
		resp.Diagnostics.AddError(summary, err.Error())

		return
	}

	state := ClusterAvailableNicsModel{
		AvailableNics: list,
		ClusterId:     config.ClusterId,
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
