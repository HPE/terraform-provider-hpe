// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

// Package networkrouterbgpneighbor implements a data source for network_router_bgp_neighbor
package networkrouterbgpneighbor

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/HPE/terraform-provider-hpe/morpheus/configure"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/constants"
	providererrors "github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/utils/convert"
)

const (
	summary                              = "read network router BGP neighbor data source"
	ErrorNoValidSearchTerms              = `no valid search terms - an id or ip_address is required`
	ErrorNoNetworkRouterBgpNeighborFound = `no network router BGP neighbor found`
	ErrorMultipleBgpNeighbors            = `multiple BGP neighbors were returned`
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
	resp.TypeName = req.ProviderTypeName + "_" + constants.SubProviderName + "_network_router_bgp_neighbor"
}

// Schema defines the schema for the data source.
func (d *DataSource) Schema(
	ctx context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = NetworkRouterBgpNeighborDataSourceSchema(ctx)
}

func bgpNeighborAsState(
	ctx context.Context,
	neighbor *sdk.GetNetworkRouterBgpNeighbor200ResponseNetworkRouterBgpNeighbor,
	routerId int64,
) (NetworkRouterBgpNeighborModel, error) {
	state := NetworkRouterBgpNeighborModel{
		Id:                 convert.Int64ToType(neighbor.Id),
		RouterId:           types.Int64Value(routerId),
		IpAddress:          convert.StrToType(neighbor.IpAddress),
		RemoteAs:           convert.StrToType(neighbor.RemoteAs),
		Weight:             convert.Int64ToType(neighbor.Weight),
		KeepAlive:          convert.Int64ToType(neighbor.KeepAlive),
		HoldDown:           convert.Int64ToType(neighbor.HoldDown),
		RestartMode:        convert.StrToType(neighbor.RestartMode.Get()),
		ProviderId:         convert.StrToType(neighbor.ProviderId.Get()),
		SyncSource:         convert.StrToType(neighbor.SyncSource.Get()),
		InternalId:         convert.StrToType(neighbor.InternalId.Get()),
		ExternalId:         convert.StrToType(neighbor.ExternalId.Get()),
		RefType:            convert.StrToType(neighbor.RefType.Get()),
		Description:        convert.StrToType(neighbor.Description.Get()),
		ForwardingAddress:  convert.StrToType(neighbor.ForwardingAddress.Get()),
		ProtocolAddress:    convert.StrToType(neighbor.ProtocolAddress.Get()),
		RouteFilteringType: convert.StrToType(neighbor.RouteFilteringType.Get()),
		RouteFilteringIn:   convert.StrToType(neighbor.RouteFilteringIn.Get()),
		RouteFilteringOut:  convert.StrToType(neighbor.RouteFilteringOut.Get()),
		BfdEnabled:         convert.BoolToType(neighbor.BfdEnabled.Get()),
		BfdInterval:        convert.Int64ToType(neighbor.BfdInterval.Get()),
		BfdMultiple:        convert.Int64ToType(neighbor.BfdMultiple.Get()),
		AllowAsIn:          convert.BoolToType(neighbor.AllowAsIn.Get()),
		HopLimit:           convert.Int64ToType(neighbor.HopLimit.Get()),
		RefId:              convert.Int64ToType(neighbor.RefId.Get()),
	}

	if neighbor.DateCreated != nil {
		state.DateCreated = types.StringValue(neighbor.DateCreated.String())
	} else {
		state.DateCreated = types.StringNull()
	}

	if neighbor.LastUpdated != nil {
		state.LastUpdated = types.StringValue(neighbor.LastUpdated.String())
	} else {
		state.LastUpdated = types.StringNull()
	}

	// Handle config variants
	state.Config = types.DynamicNull()
	state.ConfigNsxt = NewConfigNsxtValueNull()
	state.ConfigNsxv = NewConfigNsxvValueNull()

	if neighbor.Config != nil {
		if nsxt := neighbor.Config.NSXTBGPNeighborConfig2; nsxt != nil {
			sourceAddrs := make([]attr.Value, 0, len(nsxt.SourceAddresses))
			for _, addr := range nsxt.SourceAddresses {
				sourceAddrs = append(sourceAddrs, types.StringValue(addr))
			}

			sourceAddrsSet, diags := types.SetValue(types.StringType, sourceAddrs)
			if diags.HasError() {
				return NetworkRouterBgpNeighborModel{}, fmt.Errorf("error creating source_addresses set")
			}

			v, diags := NewConfigNsxtValue(
				ConfigNsxtValue{}.AttributeTypes(ctx),
				map[string]attr.Value{
					"source_addresses": sourceAddrsSet,
				},
			)
			if diags.HasError() {
				return NetworkRouterBgpNeighborModel{}, fmt.Errorf("error creating config_nsxt value")
			}

			state.ConfigNsxt = v
		} else if nsxv := neighbor.Config.NSXVBGPNeighborConfig2; nsxv != nil {
			v, diags := NewConfigNsxvValue(
				ConfigNsxvValue{}.AttributeTypes(ctx),
				map[string]attr.Value{
					"router_id": convert.StrToType(nsxv.RouterId),
					"interface": convert.StrToType(nsxv.Interface),
				},
			)
			if diags.HasError() {
				return NetworkRouterBgpNeighborModel{}, fmt.Errorf("error creating config_nsxv value")
			}

			state.ConfigNsxv = v
		} else if neighbor.Config.MapmapOfStringAny != nil {
			raw, err := json.Marshal(*neighbor.Config.MapmapOfStringAny)
			if err != nil {
				return NetworkRouterBgpNeighborModel{}, fmt.Errorf("error marshalling config: %w", err)
			}

			state.Config = types.DynamicValue(types.StringValue(string(raw)))
		}
	}

	return state, nil
}

func getBgpNeighborByID(
	ctx context.Context,
	id int64,
	routerId int64,
	apiClient *sdk.APIClient,
) (*sdk.GetNetworkRouterBgpNeighbor200ResponseNetworkRouterBgpNeighbor, error) {
	r, hresp, err := apiClient.NetworksAPI.GetNetworkRouterBgpNeighbor(
		ctx, id, routerId,
	).Execute()
	if r == nil || err != nil || hresp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"GET failed for network router BGP neighbor %d: %s",
			id, providererrors.ErrMsg(err, hresp),
		)
	}

	neighbor := r.GetNetworkRouterBgpNeighbor()

	return &neighbor, nil
}

func getBgpNeighborByIpAddress(
	ctx context.Context,
	ipAddress string,
	routerId int64,
	apiClient *sdk.APIClient,
) (*sdk.GetNetworkRouterBgpNeighbor200ResponseNetworkRouterBgpNeighbor, error) {
	rs, hresp, err := apiClient.NetworksAPI.GetNetworkRoutersBgpNeighbors(
		ctx, routerId,
	).Execute()
	if rs == nil || err != nil || hresp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"GET failed for network router BGP neighbors with IP %s: %s",
			ipAddress, providererrors.ErrMsg(err, hresp),
		)
	}

	items := rs.GetNetworkRouterBgpNeighbors()
	if len(items) == 0 {
		return nil, errors.New(ErrorNoNetworkRouterBgpNeighborFound)
	}

	var matchedIDs []int64

	for i := range items {
		if items[i].GetIpAddress() != ipAddress {
			continue
		}

		matchedIDs = append(matchedIDs, items[i].GetId())
	}

	if len(matchedIDs) == 0 {
		return nil, errors.New(ErrorNoNetworkRouterBgpNeighborFound)
	} else if len(matchedIDs) > 1 {
		return nil, errors.New(ErrorMultipleBgpNeighbors)
	}

	return getBgpNeighborByID(ctx, matchedIDs[0], routerId, apiClient)
}

func getBgpNeighbor(
	ctx context.Context,
	config *NetworkRouterBgpNeighborModel,
	apiClient *sdk.APIClient,
) (*sdk.GetNetworkRouterBgpNeighbor200ResponseNetworkRouterBgpNeighbor, error) {
	routerId := config.RouterId.ValueInt64()

	if !config.Id.IsNull() {
		return getBgpNeighborByID(ctx, config.Id.ValueInt64(), routerId, apiClient)
	} else if !config.IpAddress.IsNull() {
		return getBgpNeighborByIpAddress(ctx, config.IpAddress.ValueString(), routerId, apiClient)
	}

	return nil, errors.New(ErrorNoValidSearchTerms)
}

// Read refreshes the Terraform state with the latest data.
func (d *DataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	var config NetworkRouterBgpNeighborModel

	// Read config
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiClient, err := d.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			summary,
			"could not create sdk client",
		)

		return
	}

	neighbor, err := getBgpNeighbor(ctx, &config, apiClient)
	if err != nil {
		resp.Diagnostics.AddError(
			summary,
			err.Error(),
		)

		return
	}

	routerId := config.RouterId.ValueInt64()
	state, err := bgpNeighborAsState(ctx, neighbor, routerId)
	if err != nil {
		resp.Diagnostics.AddError(
			summary,
			err.Error(),
		)

		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
