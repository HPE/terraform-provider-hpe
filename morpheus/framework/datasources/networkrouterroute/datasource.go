// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

// Package networkrouterroute implements a data source for network_router_route
package networkrouterroute

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
	providererrors "github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/utils/convert"
)

const (
	summary                          = "read network router route data source"
	ErrorNoValidSearchTerms          = `no valid search terms - an id or name is required`
	ErrorNoNetworkRouterRouteFound   = `no network router route found`
	ErrorMultipleNetworkRouterRoutes = `multiple network router routes were returned`
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
	resp.TypeName = req.ProviderTypeName + "_" + constants.SubProviderName + "_network_router_route"
}

// Schema defines the schema for the data source.
func (d *DataSource) Schema(
	ctx context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = NetworkRouterRouteDataSourceSchema(ctx)
}

func routeAsState(
	route *sdk.GetNetworkRouterRoute200ResponseNetworkRoute,
	routerID int64,
) NetworkRouterRouteModel {

	// We're only doing this for now because SDK generated with float32.
	mtu := route.NetworkMtu.Get()
	var mtuInt64 *int64
	if mtu == nil {
		mtuInt64 = nil
	} else {
		i := int64(*mtu)
		mtuInt64 = &i
	}
	return NetworkRouterRouteModel{
		Id:           convert.Int64ToType(route.Id),
		RouterId:     types.Int64Value(routerID),
		Name:         convert.StrToType(route.Name),
		Code:         convert.StrToType(route.Code.Get()),
		Description:  convert.StrToType(route.Description.Get()),
		Mtu:          convert.Int64ToType(mtuInt64),
		RouteType:    convert.StrToType(route.RouteType),
		SourceType:   convert.StrToType(route.SourceType),
		DefaultRoute: convert.BoolToType(route.DefaultRoute),
		Enabled:      convert.BoolToType(route.Enabled),
		ExternalId:   convert.StrToType(route.ExternalId),
		ProviderId:   convert.StrToType(route.ProviderId),
	}
}

func getRouteByID(
	ctx context.Context,
	id int64,
	routerID int64,
	apiClient *sdk.APIClient,
) (*sdk.GetNetworkRouterRoute200ResponseNetworkRoute, error) {
	r, hresp, err := apiClient.NetworksAPI.GetNetworkRouterRoute(
		ctx, id, routerID,
	).Execute()
	if r == nil || err != nil || hresp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"GET failed for network router route %d: %s",
			id, providererrors.ErrMsg(err, hresp),
		)
	}

	route := r.GetNetworkRoute()

	return &route, nil
}

func getRouteByName(
	ctx context.Context,
	name string,
	routerID int64,
	apiClient *sdk.APIClient,
) (*sdk.GetNetworkRouterRoute200ResponseNetworkRoute, error) {
	rs, hresp, err := apiClient.NetworksAPI.GetNetworkRoutersRoutes(
		ctx, routerID,
	).Execute()
	if rs == nil || err != nil || hresp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"GET failed for network router routes with name %s: %s",
			name, providererrors.ErrMsg(err, hresp),
		)
	}

	items := rs.GetNetworkRoutes()
	if len(items) == 0 {
		return nil, errors.New(ErrorNoNetworkRouterRouteFound)
	}

	var matchedIDs []int64

	for i := range items {
		if items[i].GetName() != name {
			continue
		}

		matchedIDs = append(matchedIDs, items[i].GetId())
	}

	if len(matchedIDs) == 0 {
		return nil, errors.New(ErrorNoNetworkRouterRouteFound)
	} else if len(matchedIDs) > 1 {
		return nil, errors.New(ErrorMultipleNetworkRouterRoutes)
	}

	return getRouteByID(ctx, matchedIDs[0], routerID, apiClient)
}

func getRoute(
	ctx context.Context,
	config *NetworkRouterRouteModel,
	apiClient *sdk.APIClient,
) (*sdk.GetNetworkRouterRoute200ResponseNetworkRoute, error) {
	routerID := config.RouterId.ValueInt64()

	if !config.Id.IsNull() {
		return getRouteByID(ctx, config.Id.ValueInt64(), routerID, apiClient)
	} else if !config.Name.IsNull() {
		return getRouteByName(ctx, config.Name.ValueString(), routerID, apiClient)
	}

	return nil, errors.New(ErrorNoValidSearchTerms)
}

// Read refreshes the Terraform state with the latest data.
func (d *DataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	var config NetworkRouterRouteModel

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

	route, err := getRoute(ctx, &config, apiClient)
	if err != nil {
		resp.Diagnostics.AddError(
			summary,
			err.Error(),
		)

		return
	}

	routerID := config.RouterId.ValueInt64()
	state := routeAsState(route, routerID)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
