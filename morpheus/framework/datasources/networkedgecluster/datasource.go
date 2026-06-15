// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkedgecluster

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/HPE/terraform-provider-hpe/morpheus/configure"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/constants"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/utils/convert"
)

const (
	summary                 = "read network edge cluster data source"
	ErrorNoValidSearchTerms = `no valid search terms - an id or name is required`
	ErrorNoEdgeClusterFound = `no network edge cluster found`
	ErrorMultipleFound      = `multiple network edge clusters were returned`
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
	resp.TypeName = req.ProviderTypeName + "_" + constants.SubProviderName + "_network_edge_cluster"
}

// Schema defines the schema for the data source.
func (d *DataSource) Schema(
	ctx context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = NetworkEdgeClusterDataSourceSchema(ctx)
}

// Read refreshes the Terraform state with the latest data.
func (d *DataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	var config NetworkEdgeClusterModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client, err := d.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError(summary, fmt.Sprintf("failed to create client: %s", err.Error()))

		return
	}

	serverID := config.NetworkServerId.ValueInt64()

	edgeCluster, err := getEdgeCluster(ctx, &config, serverID, client)
	if err != nil {
		resp.Diagnostics.AddError(summary, err.Error())

		return
	}

	state := edgeClusterAsState(ctx, edgeCluster, serverID)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func getEdgeCluster(
	ctx context.Context,
	config *NetworkEdgeClusterModel,
	serverID int64,
	client *sdk.APIClient,
) (*sdk.GetNetworkEdgeCluster200ResponseNetworkEdgeCluster, error) {
	if !config.Id.IsNull() {
		return getEdgeClusterByID(ctx, config.Id.ValueInt64(), serverID, client)
	} else if !config.Name.IsNull() {
		return getEdgeClusterByName(ctx, config.Name.ValueString(), serverID, client)
	}

	return nil, errors.New(ErrorNoValidSearchTerms)
}

func getEdgeClusterByID(
	ctx context.Context,
	id int64,
	serverID int64,
	client *sdk.APIClient,
) (*sdk.GetNetworkEdgeCluster200ResponseNetworkEdgeCluster, error) {
	r, hresp, err := client.NetworksAPI.GetNetworkEdgeCluster(ctx, id, serverID).Execute()
	if err != nil || hresp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"GET failed for network edge cluster %d (server %d): %s",
			id, serverID, errfmt.ErrMsg(err, hresp),
		)
	}

	if r.NetworkEdgeCluster == nil {
		return nil, fmt.Errorf(
			"GET failed for network edge cluster %d (server %d): response missing networkEdgeCluster",
			id, serverID,
		)
	}

	return r.NetworkEdgeCluster, nil
}

func getEdgeClusterByName(
	ctx context.Context,
	name string,
	serverID int64,
	client *sdk.APIClient,
) (*sdk.GetNetworkEdgeCluster200ResponseNetworkEdgeCluster, error) {
	rs, hresp, err := client.NetworksAPI.GetNetworkEdgeClusters(ctx, serverID).Phrase(name).Execute()
	if err != nil || hresp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"LIST failed for network edge clusters (server %d): %s",
			serverID, errfmt.ErrMsg(err, hresp),
		)
	}

	// The LIST response returns typed networkEdgeClusters; match by exact name.
	var matched []sdk.GetNetworkEdgeClusters200ResponseAllOfNetworkEdgeClustersInner
	for _, item := range rs.NetworkEdgeClusters {
		if item.Name != nil && *item.Name == name {
			matched = append(matched, item)
		}
	}

	if len(matched) == 0 {
		return nil, errors.New(ErrorNoEdgeClusterFound)
	}

	if len(matched) > 1 {
		return nil, errors.New(ErrorMultipleFound)
	}

	if matched[0].Id == nil {
		return nil, fmt.Errorf("edge cluster %q has nil id", name)
	}

	return getEdgeClusterByID(ctx, *matched[0].Id, serverID, client)
}

func edgeClusterAsState(
	ctx context.Context,
	ec *sdk.GetNetworkEdgeCluster200ResponseNetworkEdgeCluster,
	serverID int64,
) NetworkEdgeClusterModel {
	state := NetworkEdgeClusterModel{}

	state.Id = convert.Int64ToType(ec.Id)
	state.Name = convert.StrToType(ec.Name)
	state.Description = convert.StrToType(ec.Description)
	state.NetworkServerId = types.Int64Value(serverID)
	state.ProviderId = convert.StrToType(ec.ProviderId)
	state.ExternalId = convert.StrToType(ec.ExternalId)
	state.InternalId = convert.StrToType(ec.InternalId)
	state.Active = convert.BoolToType(ec.Active)
	state.DisplayName = convert.StrToType(ec.DisplayName)
	state.Enabled = convert.BoolToType(ec.Enabled)
	state.Visibility = convert.StrToType(ec.Visibility)

	// Config — convert struct to dynamic via JSON roundtrip.
	if ec.Config != nil {
		configMap, mapErr := ec.Config.ToMap()
		if mapErr == nil {
			dyn, dynErr := convert.MapToDynamic(ctx, configMap)
			if dynErr == nil {
				state.Config = dyn
			} else {
				state.Config = types.DynamicNull()
			}
		} else {
			state.Config = types.DynamicNull()
		}
	} else {
		state.Config = types.DynamicNull()
	}

	return state
}
