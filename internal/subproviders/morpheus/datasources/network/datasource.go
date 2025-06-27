// Copyright 2025 Hewlett Packard Enterprise Development LP

package network

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"

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

	net, ok := network.GetNetworkOk()
	if !ok {
		return nil, fmt.Errorf("network %d is nil", id)
	}

	state.Labels = convert.StrSliceToSet(net.Labels)

	state.Id = types.Int64Value(id)
	state.Name = convert.StrToType(net.Name)
	state.DisplayName = convert.StrToType(net.DisplayName)
	state.Description = convert.StrToType(net.Description.Get())
	state.Cidr = convert.StrToType(net.Cidr)
	state.Active = convert.BoolToType(net.Active)
	state.Visibility = convert.StrToType(net.Visibility)

	return state, nil
}

func getNetworkByName(
	ctx context.Context,
	name string,
	client *sdk.APIClient,
) (*NetworkModel, error) {
	networks, hresp, err := client.NetworksAPI.ListNetworks(ctx).Name(name).Execute()
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
