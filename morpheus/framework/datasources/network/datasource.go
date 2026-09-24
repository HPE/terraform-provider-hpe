// Copyright 2025 Hewlett Packard Enterprise Development LP

package network

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"

	"github.com/HPE/terraform-provider-hpe/morpheus/configure"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/utils/convert"
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
	resp.TypeName = req.ProviderTypeName + "_" + "network"
}

func (d *DataSource) Schema(
	ctx context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = NetworkDataSourceSchema(ctx)
}

// getNetwork finds the network by id or by name. Morpheus keeps one network
// record per cloud, so a name may match a record in each of several clouds
// that cover the same account; cloud_id selects among them, and with an id it
// is checked.
func getNetwork(
	ctx context.Context,
	config NetworkModel,
	client *sdk.APIClient,
) (*NetworkModel, error) {
	var cloudID *int64
	if !config.CloudId.IsNull() && !config.CloudId.IsUnknown() {
		v := config.CloudId.ValueInt64()
		cloudID = &v
	}

	if !config.Id.IsNull() {
		state, err := getNetworkByID(ctx, config.Id.ValueInt64(), client)
		if err != nil {
			return nil, err
		}

		if cloudID != nil && (state.CloudId.IsNull() || state.CloudId.ValueInt64() != *cloudID) {
			return nil, fmt.Errorf("network %d belongs to cloud %s, not cloud %d",
				config.Id.ValueInt64(), state.CloudId.String(), *cloudID)
		}

		return state, nil
	}

	if !config.Name.IsNull() {
		return getNetworkByName(ctx, config.Name.ValueString(), cloudID, client)
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
		return nil, fmt.Errorf("network %d GET failed: %s", id, errfmt.ErrMsg(err, hresp))
	}

	state := &NetworkModel{}

	if network.Network == nil {
		return nil, fmt.Errorf("network %d is nil", id)
	}
	net := network.Network

	state.Labels = convert.StrSliceToSet(net.Labels)

	state.Id = types.Int64Value(id)
	state.Name = convert.StrToType(net.Name)
	state.DisplayName = convert.StrToType(net.DisplayName.Get())
	state.Description = convert.StrToType(net.Description.Get())
	state.Cidr = convert.StrToType(net.Cidr.Get())
	state.Active = convert.BoolToType(net.Active)
	state.Visibility = convert.StrToType(net.Visibility)

	state.CloudId = types.Int64Null()
	if net.Zone != nil {
		state.CloudId = convert.Int64ToType(net.Zone.Id)
	}

	return state, nil
}

func getNetworkByName(
	ctx context.Context,
	name string,
	cloudID *int64,
	client *sdk.APIClient,
) (*NetworkModel, error) {
	request := client.NetworksAPI.ListNetworks(ctx).Name(name)
	if cloudID != nil {
		request = request.ZoneId(*cloudID)
	}

	networks, hresp, err := request.Execute()
	if err != nil || hresp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("network %s list failed: %s", name, errfmt.ErrMsg(err, hresp))
	}

	var matchingNetworks []sdk.ListNetworks200ResponseAllOfNetworksInner
	for _, network := range networks.Networks {
		if network.Name != nil && *network.Name == name {
			matchingNetworks = append(matchingNetworks, network)
		}
	}

	if len(matchingNetworks) == 0 {
		if cloudID != nil {
			return nil, fmt.Errorf("network %s not found in cloud %d", name, *cloudID)
		}

		return nil, fmt.Errorf("network %s not found", name)
	}

	if len(matchingNetworks) > 1 {
		// Usually one record per cloud for the same network: Morpheus keeps a
		// network record in every cloud that sees it. Say which, so cloud_id
		// can pick.
		var found []string
		for _, n := range matchingNetworks {
			if n.Id == nil {
				continue
			}

			if n.Zone != nil && n.Zone.Id != nil {
				found = append(found, fmt.Sprintf("%d (cloud %d)", *n.Id, *n.Zone.Id))
			} else {
				found = append(found, fmt.Sprintf("%d", *n.Id))
			}
		}

		return nil, fmt.Errorf(
			"multiple networks found with name %s: %s. "+
				"Set cloud_id to select the cloud, or use id",
			name,
			strings.Join(found, ", "),
		)
	}

	id := matchingNetworks[0].Id
	if id == nil {
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
