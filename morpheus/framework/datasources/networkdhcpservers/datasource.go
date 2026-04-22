// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

// Package networkdhcpservers implements a data source for network_dhcp_server
package networkdhcpservers

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
	summary                         = "read network DHCP server data source"
	ErrorNoValidSearchTerms         = `no valid search terms - an id or name is required`
	ErrorRunningPreApply            = `Error running pre-apply plan: exit status 1`
	ErrorNoNetworkDhcpServerFound   = `no network DHCP server found`
	ErrorMultipleNetworkDhcpServers = `multiple network DHCP servers were returned`
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
	resp.TypeName = req.ProviderTypeName + "_" + constants.SubProviderName + "_network_dhcp_server"
}

// Schema defines the schema for the data source.
func (d *DataSource) Schema(
	ctx context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = NetworkDhcpServersDataSourceSchema(ctx)
}

func dhcpServerAsState(
	ctx context.Context,
	dhcp *sdk.GetNetworkDhcpServer200ResponseNetworkDhcpServer,
	networkServerId int64,
) (NetworkDhcpServersModel, error) {
	state := NetworkDhcpServersModel{
		Id:              convert.Int64ToType(dhcp.Id),
		Name:            convert.StrToType(dhcp.Name),
		LeaseTime:       convert.Int64ToType(dhcp.LeaseTime),
		ServerIpAddress: convert.StrToType(dhcp.ServerIpAddress),
		NetworkServerId: types.Int64Value(networkServerId),
	}

	state.Config = types.DynamicNull()
	state.ConfigNsx = NewConfigNsxValueNull()

	if dhcp.Config != nil {
		if nsx := dhcp.Config.NSXDHCPServerConfiguration2; nsx != nil {
			v, diags := NewConfigNsxValue(
				ConfigNsxValue{}.AttributeTypes(ctx),
				map[string]attr.Value{
					"edge_cluster":      convert.StrToType(nsx.EdgeCluster.Get()),
					"active_edge_node":  convert.StrToType(nsx.PreferredEdgeNode1.Get()),
					"standby_edge_node": convert.StrToType(nsx.PreferredEdgeNode2.Get()),
				},
			)
			if diags.HasError() {
				return NetworkDhcpServersModel{}, fmt.Errorf("error creating config_nsx value")
			}

			state.ConfigNsx = v
		} else if dhcp.Config.MapmapOfStringAny != nil {
			raw, err := json.Marshal(*dhcp.Config.MapmapOfStringAny)
			if err != nil {
				return NetworkDhcpServersModel{}, fmt.Errorf("error marshalling config: %w", err)
			}

			state.Config = types.DynamicValue(types.StringValue(string(raw)))
		}
	}

	return state, nil
}

func getDhcpServerByID(
	ctx context.Context,
	id int64,
	serverId int64,
	apiClient *sdk.APIClient,
) (*sdk.GetNetworkDhcpServer200ResponseNetworkDhcpServer, error) {
	r, hresp, err := apiClient.NetworksAPI.GetNetworkDhcpServer(
		ctx, id, serverId,
	).Execute()
	if r == nil || err != nil || hresp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"GET failed for network DHCP server %d: %s",
			id, providererrors.ErrMsg(err, hresp),
		)
	}

	dhcp := r.GetNetworkDhcpServer()

	return &dhcp, nil
}

func getDhcpServerByName(
	ctx context.Context,
	name string,
	serverId int64,
	apiClient *sdk.APIClient,
) (*sdk.GetNetworkDhcpServer200ResponseNetworkDhcpServer, error) {
	// Phrase is used because the API does not expose an exact Name filter.
	// The subsequent loop performs exact-match filtering on the results.
	rs, hresp, err := apiClient.NetworksAPI.GetNetworkDhcpServers(
		ctx, serverId,
	).Phrase(name).Execute()
	if rs == nil || err != nil || hresp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"GET failed for network DHCP server %s: %s",
			name, providererrors.ErrMsg(err, hresp),
		)
	}

	items := rs.GetNetworkDhcpServers()
	if len(items) == 0 {
		return nil, errors.New(ErrorNoNetworkDhcpServerFound)
	}

	var matchedIDs []int64

	for i := range items {
		if items[i].GetName() != name {
			continue
		}

		matchedIDs = append(matchedIDs, items[i].GetId())
	}

	if len(matchedIDs) == 0 {
		return nil, errors.New(ErrorNoNetworkDhcpServerFound)
	} else if len(matchedIDs) > 1 {
		return nil, errors.New(ErrorMultipleNetworkDhcpServers)
	}

	return getDhcpServerByID(ctx, matchedIDs[0], serverId, apiClient)
}

func getDhcpServer(
	ctx context.Context,
	config *NetworkDhcpServersModel,
	apiClient *sdk.APIClient,
) (*sdk.GetNetworkDhcpServer200ResponseNetworkDhcpServer, error) {
	serverId := config.NetworkServerId.ValueInt64()

	if !config.Id.IsNull() {
		return getDhcpServerByID(ctx, config.Id.ValueInt64(), serverId, apiClient)
	} else if !config.Name.IsNull() {
		return getDhcpServerByName(ctx, config.Name.ValueString(), serverId, apiClient)
	}

	return nil, errors.New(ErrorNoValidSearchTerms)
}

// Read refreshes the Terraform state with the latest data.
func (d *DataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	var config NetworkDhcpServersModel

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

	dhcp, err := getDhcpServer(ctx, &config, apiClient)
	if err != nil {
		resp.Diagnostics.AddError(
			summary,
			err.Error(),
		)

		return
	}

	state, err := dhcpServerAsState(ctx, dhcp, config.NetworkServerId.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError(
			summary,
			err.Error(),
		)

		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
