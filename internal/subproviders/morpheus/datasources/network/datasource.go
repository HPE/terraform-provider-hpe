// Copyright 2025 Hewlett Packard Enterprise Development LP

package network

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/datasource"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/sdk"

	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/configure"
	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/convert"
	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/errors"
)

var _ datasource.DataSource = &DataSource{}

func NewDataSource() datasource.DataSource {
	return &DataSource{}
}

type DataSource struct {
	configure.DataSourceWithMorpheusConfigure
	datasource.DataSource
}

func (d *DataSource) Metadata(
	_ context.Context,
	req datasource.MetadataRequest,
	resp *datasource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_morpheus_network"
}

func (d *DataSource) Schema(
	ctx context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = NetworkDataSourceSchema(ctx)
}

func getNetwork(
	ctx context.Context,
	config NetworkModel,
	client *sdk.APIClient,
) (*NetworkModel, error) {
	if !config.Id.IsNull() {
		return getNetworkByID(ctx, config.Id.ValueInt64(), client)
	}
	if !config.Name.IsNull() {
		return getNetworkByName(ctx, config.Name.ValueString(), client)
	}

	return nil, fmt.Errorf("either id or name must be specified")
}

func getNetworkByID(
	ctx context.Context,
	id int64,
	client *sdk.APIClient,
) (*NetworkModel, error) {
	network, hresp, err := client.NetworksAPI.GetNetwork(ctx, id).Execute()
	if err != nil || hresp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("network %d GET failed: %s", id, errors.ErrMsg(err, hresp))
	}

	state := &NetworkModel{}

	state.Labels = convert.StrSliceToSet(network.GetNetwork().Labels)

	if id, ok := network.Network.GetIdOk(); ok {
		state.Id = convert.Int64ToType(id)
	}
	if name, ok := network.Network.GetNameOk(); ok {
		state.Name = convert.StrToType(name)
	}
	if displayName, ok := network.Network.GetDisplayNameOk(); ok {
		state.DisplayName = convert.StrToType(displayName)
	}
	if description, ok := network.Network.GetDescriptionOk(); ok {
		state.Description = convert.StrToType(description)
	}
	if cidr, ok := network.Network.GetCidrOk(); ok {
		state.Cidr = convert.StrToType(cidr)
	}
	if active, ok := network.Network.GetActiveOk(); ok {
		state.Active = convert.BoolToType(active)
	}
	if visibility, ok := network.Network.GetVisibilityOk(); ok {
		state.Visibility = convert.StrToType(visibility)
	}

	return state, nil
}

func getNetworkByName(
	ctx context.Context,
	name string,
	client *sdk.APIClient,
) (*NetworkModel, error) {
	networks, hresp, err := client.NetworksAPI.ListNetworks(ctx).Execute()
	if err != nil || hresp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("network %s list failed: %s", name, errors.ErrMsg(err, hresp))
	}

	var matchingNetworks []sdk.ListNetworks200ResponseAllOfNetworksInner
	for _, network := range networks.GetNetworks() {
		if networkName, ok := network.GetNameOk(); ok && *networkName == name {
			matchingNetworks = append(matchingNetworks, network)
		}
	}

	if len(matchingNetworks) == 0 {
		return nil, fmt.Errorf("network %s not found", name)
	}

	if len(matchingNetworks) > 1 {
		var networkIDs []string
		for _, n := range matchingNetworks {
			if id, ok := n.GetIdOk(); ok {
				networkIDs = append(networkIDs, fmt.Sprintf("%d", *id))
			}
		}

		return nil, fmt.Errorf(
			"multiple networks found with name %s. Network IDs: %s. "+
				"Please specify an ID instead",
			name,
			strings.Join(networkIDs, ", "),
		)
	}

	id, ok := matchingNetworks[0].GetIdOk()
	if !ok {
		return nil, fmt.Errorf("network %s has missing ID", name)
	}

	return getNetworkByID(ctx, *id, client)
}

func (d *DataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	var config NetworkModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client, err := d.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"read network data source",
			fmt.Sprintf("failed to create client: %s", err.Error()),
		)

		return
	}

	state, err := getNetwork(ctx, config, client)
	if err != nil {
		resp.Diagnostics.AddError(
			"read network data source",
			err.Error(),
		)

		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}
