// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

// Package virtualswitch implements a data source for a single virtual switch.
package virtualswitch

import (
	"context"
	"errors"
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
	summary = "read virtual switch data source"

	// ErrorNoValidSearchTerms is returned when neither id nor (cluster_id +
	// name) were supplied.
	ErrorNoValidSearchTerms = `no valid search terms - an id, or a cluster_id and name, is required`
	// ErrorNoVirtualSwitchFound is returned when a cluster_id + name lookup
	// matched nothing.
	ErrorNoVirtualSwitchFound = `no virtual switch found`
	// ErrorMultipleVirtualSwitches is returned when a cluster_id + name lookup
	// matched more than one switch.
	ErrorMultipleVirtualSwitches = `multiple virtual switches were returned`
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
	resp.TypeName = req.ProviderTypeName + "_" + "virtual_switch"
}

// Schema defines the schema for the data source.
func (d *DataSource) Schema(
	ctx context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = VirtualSwitchDataSourceSchema(ctx)
}

func hostIpsAsState(
	ctx context.Context,
	hostIps []sdk.GetVirtualSwitch200ResponseVirtualSwitchAllOfSegmentsInner1HostIpsInner,
) (types.List, error) {
	var convErr error

	list, diags := convert.ToListType(ctx, hostIps,
		func(h sdk.GetVirtualSwitch200ResponseVirtualSwitchAllOfSegmentsInner1HostIpsInner) HostIpsValue {
			v, vDiags := NewHostIpsValue(
				HostIpsValue{}.AttributeTypes(ctx),
				map[string]attr.Value{
					"ip":        convert.StrToType(h.Ip),
					"server_id": convert.Int64ToType(h.ServerId),
				},
			)
			if vDiags.HasError() && convErr == nil {
				convErr = fmt.Errorf("error creating host_ips value")
			}

			return v
		})
	if diags.HasError() || convErr != nil {
		return types.ListNull(HostIpsValue{}.Type(ctx)), fmt.Errorf("error creating host_ips list")
	}

	return list, nil
}

func segmentsAsState(
	ctx context.Context,
	segments []sdk.GetVirtualSwitch200ResponseVirtualSwitchAllOfSegmentsInner1,
) (types.List, error) {
	var convErr error

	list, diags := convert.ToListType(ctx, segments,
		func(s sdk.GetVirtualSwitch200ResponseVirtualSwitchAllOfSegmentsInner1) SegmentsValue {
			hostIps, err := hostIpsAsState(ctx, s.HostIps)
			if err != nil && convErr == nil {
				convErr = err
			}

			v, vDiags := NewSegmentsValue(
				SegmentsValue{}.AttributeTypes(ctx),
				map[string]attr.Value{
					"gateway":             convert.StrToType(s.Gateway.Get()),
					"host_ips":            hostIps,
					"id":                  convert.Int64ToType(s.Id),
					"interface_name":      convert.StrToType(s.InterfaceName.Get()),
					"ip_address":          convert.StrToType(s.IpAddress.Get()),
					"mtu":                 convert.Int32ToType(s.Mtu.Get()),
					"netmask":             convert.StrToType(s.Netmask.Get()),
					"ovs_bridge":          convert.StrToType(s.OvsBridge.Get()),
					"ovs_libvirt_network": convert.StrToType(s.OvsLibvirtNetwork.Get()),
					"ovs_uuid":            convert.StrToType(s.OvsUuid.Get()),
					"status":              convert.StrToType(s.Status),
					"status_message":      convert.StrToType(s.StatusMessage.Get()),
					"traffic_type":        convert.StrToType(s.TrafficType),
					"vlan_id":             convert.Int32ToType(s.VlanId.Get()),
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

func uplinksAsState(
	ctx context.Context,
	uplinks []sdk.GetVirtualSwitch200ResponseVirtualSwitchAllOfUplinksInner,
) (types.List, error) {
	var convErr error

	list, diags := convert.ToListType(ctx, uplinks,
		func(u sdk.GetVirtualSwitch200ResponseVirtualSwitchAllOfUplinksInner) UplinksValue {
			v, vDiags := NewUplinksValue(
				UplinksValue{}.AttributeTypes(ctx),
				map[string]attr.Value{
					"id":             convert.Int64ToType(u.Id),
					"nic1":           convert.StrToType(u.Nic1),
					"nic2":           convert.StrToType(u.Nic2.Get()),
					"server_id":      convert.Int64ToType(u.ServerId),
					"server_name":    convert.StrToType(u.ServerName.Get()),
					"sync_message":   convert.StrToType(u.SyncMessage.Get()),
					"sync_status":    convert.StrToType(u.SyncStatus),
					"last_sync_date": convert.TimeToType(u.LastSyncDate.Get()),
				},
			)
			if vDiags.HasError() && convErr == nil {
				convErr = fmt.Errorf("error creating uplinks value")
			}

			return v
		})
	if diags.HasError() || convErr != nil {
		return types.ListNull(UplinksValue{}.Type(ctx)), fmt.Errorf("error creating uplinks list")
	}

	return list, nil
}

func virtualSwitchAsState(
	ctx context.Context,
	sw *sdk.GetVirtualSwitch200ResponseVirtualSwitch,
) (VirtualSwitchModel, error) {
	segments, err := segmentsAsState(ctx, sw.Segments)
	if err != nil {
		return VirtualSwitchModel{}, err
	}

	uplinks, err := uplinksAsState(ctx, sw.Uplinks)
	if err != nil {
		return VirtualSwitchModel{}, err
	}

	return VirtualSwitchModel{
		Active:        convert.BoolToType(sw.Active),
		BaseType:      convert.StrToType(sw.BaseType),
		BondMode:      convert.StrToType(sw.BondMode.Get()),
		ClusterId:     convert.Int64ToType(sw.ClusterId),
		Code:          convert.StrToType(sw.Code.Get()),
		Description:   convert.StrToType(sw.Description.Get()),
		Id:            convert.Int64ToType(sw.Id),
		Mtu:           convert.Int32ToType(sw.Mtu),
		Name:          convert.StrToType(sw.Name),
		NetworkCount:  convert.Int32ToType(sw.NetworkCount),
		NicCount:      convert.Int32ToType(sw.NicCount),
		Segments:      segments,
		Status:        convert.StrToType(sw.Status),
		StatusMessage: convert.StrToType(sw.StatusMessage.Get()),
		Uplinks:       uplinks,
	}, nil
}

func getVirtualSwitchByID(
	ctx context.Context,
	id int64,
	apiClient *sdk.APIClient,
) (*sdk.GetVirtualSwitch200ResponseVirtualSwitch, error) {
	r, hresp, err := apiClient.VirtualSwitchesAPI.GetVirtualSwitch(ctx, id).Execute()
	if hresp != nil && hresp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("%s with id %d", ErrorNoVirtualSwitchFound, id)
	}

	if r == nil || r.VirtualSwitch == nil || err != nil || hresp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"GET failed for virtual switch %d: %s",
			id, providererrors.ErrMsg(err, hresp),
		)
	}

	return r.VirtualSwitch, nil
}

func getVirtualSwitchByName(
	ctx context.Context,
	clusterID int64,
	name string,
	apiClient *sdk.APIClient,
) (*sdk.GetVirtualSwitch200ResponseVirtualSwitch, error) {
	rs, hresp, err := apiClient.VirtualSwitchesAPI.ListVirtualSwitches(ctx).ClusterId(clusterID).Execute()
	if hresp != nil && hresp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("%s (cluster %d)", ErrorClusterNotFound, clusterID)
	}

	if rs == nil || err != nil || hresp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"GET failed for virtual switches on cluster %d: %s",
			clusterID, providererrors.ErrMsg(err, hresp),
		)
	}

	var matchedIDs []int64

	for i := range rs.VirtualSwitches {
		sw := rs.VirtualSwitches[i]
		if sw.Name == nil || *sw.Name != name {
			continue
		}

		if sw.Id != nil {
			matchedIDs = append(matchedIDs, *sw.Id)
		}
	}

	if len(matchedIDs) == 0 {
		return nil, fmt.Errorf("%s matching name %q on cluster %d", ErrorNoVirtualSwitchFound, name, clusterID)
	} else if len(matchedIDs) > 1 {
		return nil, fmt.Errorf("%s matching name %q on cluster %d", ErrorMultipleVirtualSwitches, name, clusterID)
	}

	return getVirtualSwitchByID(ctx, matchedIDs[0], apiClient)
}

func getVirtualSwitch(
	ctx context.Context,
	config *VirtualSwitchModel,
	apiClient *sdk.APIClient,
) (*sdk.GetVirtualSwitch200ResponseVirtualSwitch, error) {
	if !config.Id.IsNull() {
		return getVirtualSwitchByID(ctx, config.Id.ValueInt64(), apiClient)
	}

	if !config.Name.IsNull() && !config.ClusterId.IsNull() {
		return getVirtualSwitchByName(
			ctx, config.ClusterId.ValueInt64(), config.Name.ValueString(), apiClient,
		)
	}

	return nil, errors.New(ErrorNoValidSearchTerms)
}

// Read refreshes the Terraform state with the latest data.
func (d *DataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	var config VirtualSwitchModel

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

	sw, err := getVirtualSwitch(ctx, &config, apiClient)
	if err != nil {
		resp.Diagnostics.AddError(summary, err.Error())

		return
	}

	state, err := virtualSwitchAsState(ctx, sw)
	if err != nil {
		resp.Diagnostics.AddError(summary, err.Error())

		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
