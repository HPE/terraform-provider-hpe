// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networktransportzone

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
	summary                 = "read network transport zone data source"
	ErrorNoValidSearchTerms = `no valid search terms - an id or name is required`
	ErrorNoTransportZone    = `no network transport zone found`
	ErrorMultipleFound      = `multiple network transport zones were returned`
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
	resp.TypeName = req.ProviderTypeName + "_" + constants.SubProviderName + "_network_transport_zone"
}

// Schema defines the schema for the data source.
func (d *DataSource) Schema(
	ctx context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = NetworkTransportZoneDataSourceSchema(ctx)
}

// Read refreshes the Terraform state with the latest data.
func (d *DataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	var config NetworkTransportZoneModel

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

	tz, err := getTransportZone(ctx, &config, serverID, client)
	if err != nil {
		resp.Diagnostics.AddError(summary, err.Error())

		return
	}

	state := transportZoneAsState(ctx, tz, serverID)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func getTransportZone(
	ctx context.Context,
	config *NetworkTransportZoneModel,
	serverID int64,
	client *sdk.APIClient,
) (*sdk.GetNetworkTransportZone200ResponseNetworkScope, error) {
	if !config.Id.IsNull() {
		return getTransportZoneByID(ctx, config.Id.ValueInt64(), serverID, client)
	} else if !config.Name.IsNull() {
		return getTransportZoneByName(ctx, config.Name.ValueString(), serverID, client)
	}

	return nil, errors.New(ErrorNoValidSearchTerms)
}

func getTransportZoneByID(
	ctx context.Context,
	id int64,
	serverID int64,
	client *sdk.APIClient,
) (*sdk.GetNetworkTransportZone200ResponseNetworkScope, error) {
	r, hresp, err := client.NetworksAPI.GetNetworkTransportZone(ctx, id, serverID).Execute()
	if err != nil || hresp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"GET failed for network transport zone %d (server %d): %s",
			id, serverID, errfmt.ErrMsg(err, hresp),
		)
	}

	if r.NetworkScope == nil {
		return nil, fmt.Errorf(
			"GET failed for network transport zone %d (server %d): response missing networkScope",
			id, serverID,
		)
	}

	return r.NetworkScope, nil
}

func getTransportZoneByName(
	ctx context.Context,
	name string,
	serverID int64,
	client *sdk.APIClient,
) (*sdk.GetNetworkTransportZone200ResponseNetworkScope, error) {
	rs, hresp, err := client.NetworksAPI.GetNetworkTransportZones(ctx, serverID).Phrase(name).Execute()
	if err != nil || hresp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"LIST failed for network transport zones (server %d): %s",
			serverID, errfmt.ErrMsg(err, hresp),
		)
	}

	// The LIST response returns typed networkScopes; match by exact name.
	var matched []sdk.GetNetworkTransportZones200ResponseAllOfNetworkScopesInner
	for _, item := range rs.NetworkScopes {
		if item.Name != nil && *item.Name == name {
			matched = append(matched, item)
		}
	}

	if len(matched) == 0 {
		return nil, errors.New(ErrorNoTransportZone)
	}

	if len(matched) > 1 {
		return nil, errors.New(ErrorMultipleFound)
	}

	if matched[0].Id == nil {
		return nil, fmt.Errorf("transport zone %q has nil id", name)
	}

	return getTransportZoneByID(ctx, *matched[0].Id, serverID, client)
}

func transportZoneAsState(
	ctx context.Context,
	tz *sdk.GetNetworkTransportZone200ResponseNetworkScope,
	serverID int64,
) NetworkTransportZoneModel {
	state := NetworkTransportZoneModel{}

	state.Id = convert.Int64ToType(tz.Id)
	state.Name = convert.StrToType(tz.Name)
	state.Description = convert.StrToType(tz.Description)
	state.NetworkServerId = types.Int64Value(serverID)
	state.ProviderId = convert.StrToType(tz.ProviderId)
	state.ExternalId = convert.StrToType(tz.ExternalId)
	state.InternalId = convert.StrToType(tz.InternalId)
	state.StreamType = convert.StrToType(tz.StreamType)
	state.Active = convert.BoolToType(tz.Active)
	state.DisplayName = convert.StrToType(tz.DisplayName)
	state.Enabled = convert.BoolToType(tz.Enabled)
	state.Status = convert.StrToType(tz.Status)
	state.Visibility = convert.StrToType(tz.Visibility)

	// Config — convert struct to dynamic via ToMap.
	if tz.Config != nil {
		configMap, mapErr := tz.Config.ToMap()
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
