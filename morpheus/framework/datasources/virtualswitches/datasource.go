// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

// Package virtualswitches implements a data source listing the virtual switches
// on a cluster.
package virtualswitches

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
	summary = "read virtual switches data source"

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
	resp.TypeName = req.ProviderTypeName + "_" + "virtual_switches"
}

// Schema defines the schema for the data source.
func (d *DataSource) Schema(
	ctx context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = VirtualSwitchesDataSourceSchema(ctx)
}

func segmentsAsState(
	ctx context.Context,
	segments []sdk.ListVirtualSwitches200ResponseVirtualSwitchesInnerSegmentsInner,
) (types.List, error) {
	var convErr error

	list, diags := convert.ToListType(ctx, segments,
		func(s sdk.ListVirtualSwitches200ResponseVirtualSwitchesInnerSegmentsInner) SegmentsValue {
			v, vDiags := NewSegmentsValue(
				SegmentsValue{}.AttributeTypes(ctx),
				map[string]attr.Value{
					"traffic_type": convert.StrToType(s.TrafficType),
				},
			)
			if vDiags.HasError() && convErr == nil {
				convErr = fmt.Errorf("error creating segments value")
			}

			return v
		})
	if diags.HasError() || convErr != nil {
		return types.ListNull(SegmentsValue{}.Type(ctx)), fmt.Errorf("error creating segments list")
	}

	return list, nil
}

func virtualSwitchesAsState(
	ctx context.Context,
	switches []sdk.ListVirtualSwitches200ResponseVirtualSwitchesInner,
) (types.List, error) {
	var convErr error

	list, diags := convert.ToListType(ctx, switches,
		func(sw sdk.ListVirtualSwitches200ResponseVirtualSwitchesInner) VirtualSwitchesValue {
			segments, err := segmentsAsState(ctx, sw.Segments)
			if err != nil && convErr == nil {
				convErr = err
			}

			v, vDiags := NewVirtualSwitchesValue(
				VirtualSwitchesValue{}.AttributeTypes(ctx),
				map[string]attr.Value{
					"active":         convert.BoolToType(sw.Active),
					"base_type":      convert.StrToType(sw.BaseType),
					"bond_mode":      convert.StrToType(sw.BondMode.Get()),
					"cluster_id":     convert.Int64ToType(sw.ClusterId),
					"code":           convert.StrToType(sw.Code.Get()),
					"description":    convert.StrToType(sw.Description.Get()),
					"id":             convert.Int64ToType(sw.Id),
					"mtu":            convert.Int32ToType(sw.Mtu),
					"name":           convert.StrToType(sw.Name),
					"network_count":  convert.Int32ToType(sw.NetworkCount),
					"nic_count":      convert.Int32ToType(sw.NicCount),
					"segments":       segments,
					"status":         convert.StrToType(sw.Status),
					"status_message": convert.StrToType(sw.StatusMessage.Get()),
				},
			)
			if vDiags.HasError() && convErr == nil {
				convErr = fmt.Errorf("error creating virtual switches value")
			}

			return v
		})
	if diags.HasError() || convErr != nil {
		return types.ListNull(VirtualSwitchesValue{}.Type(ctx)), fmt.Errorf("error creating virtual switches list")
	}

	return list, nil
}

// Read refreshes the Terraform state with the latest data.
func (d *DataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	var config VirtualSwitchesModel

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

	rs, hresp, err := apiClient.VirtualSwitchesAPI.ListVirtualSwitches(ctx).ClusterId(clusterID).Execute()
	if hresp != nil && hresp.StatusCode == http.StatusNotFound {
		resp.Diagnostics.AddError(summary, fmt.Sprintf("%s (cluster %d)", ErrorClusterNotFound, clusterID))

		return
	}

	if rs == nil || err != nil || hresp.StatusCode != http.StatusOK {
		resp.Diagnostics.AddError(summary, fmt.Sprintf(
			"GET failed for virtual switches on cluster %d: %s",
			clusterID, providererrors.ErrMsg(err, hresp),
		))

		return
	}

	list, err := virtualSwitchesAsState(ctx, rs.VirtualSwitches)
	if err != nil {
		resp.Diagnostics.AddError(summary, err.Error())

		return
	}

	state := VirtualSwitchesModel{
		ClusterId:       config.ClusterId,
		VirtualSwitches: list,
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
