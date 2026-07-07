// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

// Package networkrouternat implements a data source for network_router_nat
package networkrouternat

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"

	"github.com/HPE/terraform-provider-hpe/morpheus/configure"
	providererrors "github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/utils/convert"
)

const (
	summary                        = "read network router nat data source"
	ErrorNoValidSearchTerms        = `no valid search terms - an id or name is required`
	ErrorNoNetworkRouterNatFound   = `no network router nat found`
	ErrorMultipleNetworkRouterNats = `multiple network router nats were returned`
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
	resp.TypeName = req.ProviderTypeName + "_" + "network_router_nat"
}

// Schema defines the schema for the data source.
func (d *DataSource) Schema(
	ctx context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = NetworkRouterNatDataSourceSchema(ctx)
}

func natAsState(
	nat *sdk.GetNetworkRouterNat200ResponseNetworkRouterNAT,
	routerID int64,
) NetworkRouterNatModel {
	// Id is *int32 in SDK — convert to *int64
	var idPtr *int64
	if nat.Id != nil {
		i := int64(*nat.Id)
		idPtr = &i
	}

	// Priority is *int32 in SDK — convert to *int64
	var priorityPtr *int64
	if nat.Priority != nil {
		p := int64(*nat.Priority)
		priorityPtr = &p
	}

	return NetworkRouterNatModel{
		Id:                 convert.Int64ToType(idPtr),
		RouterId:           types.Int64Value(routerID),
		Name:               convert.StrToType(nat.Name),
		Action:             convert.StrToType(nat.Action),
		Description:        convert.StrToType(nat.Description),
		Enabled:            convert.BoolToType(nat.Enabled),
		SourceNetwork:      convert.StrToType(nat.SourceNetwork),
		DestinationNetwork: convert.StrToType(nat.DestinationNetwork.Get()),
		TranslatedNetwork:  convert.StrToType(nat.TranslatedNetwork),
		SourcePorts:        convert.StrToType(nat.SourcePorts.Get()),
		DestinationPorts:   convert.StrToType(nat.DestinationPorts.Get()),
		TranslatedPorts:    convert.StrToType(nat.TranslatedPorts.Get()),
		Priority:           convert.Int64ToType(priorityPtr),
		Protocol:           convert.StrToType(nat.Protocol.Get()),
	}
}

func getNatByID(
	ctx context.Context,
	id int64,
	routerID int64,
	apiClient *sdk.APIClient,
) (*sdk.GetNetworkRouterNat200ResponseNetworkRouterNAT, error) {
	r, hresp, err := apiClient.NetworksAPI.GetNetworkRouterNat(
		ctx, id, routerID,
	).Execute()
	if r == nil || err != nil || hresp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"GET failed for network router nat %d: %s",
			id, providererrors.ErrMsg(err, hresp),
		)
	}

	return r.NetworkRouterNAT, nil
}

func getNatByName(
	ctx context.Context,
	name string,
	routerID int64,
	apiClient *sdk.APIClient,
) (*sdk.GetNetworkRouterNat200ResponseNetworkRouterNAT, error) {
	rs, hresp, err := apiClient.NetworksAPI.GetNetworkRoutersNats(
		ctx, routerID,
	).Execute()
	if rs == nil || err != nil || hresp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"GET failed for network router nats with name %s: %s",
			name, providererrors.ErrMsg(err, hresp),
		)
	}

	items := rs.NetworkRouterNATs
	if len(items) == 0 {
		return nil, errors.New(ErrorNoNetworkRouterNatFound)
	}

	var matchedIDs []int64

	for i := range items {
		if items[i].Name == nil || *items[i].Name != name {
			continue
		}
		if items[i].Id == nil {
			continue
		}

		matchedIDs = append(matchedIDs, int64(*items[i].Id))
	}

	if len(matchedIDs) == 0 {
		return nil, errors.New(ErrorNoNetworkRouterNatFound)
	} else if len(matchedIDs) > 1 {
		return nil, errors.New(ErrorMultipleNetworkRouterNats)
	}

	return getNatByID(ctx, matchedIDs[0], routerID, apiClient)
}

func getNat(
	ctx context.Context,
	config *NetworkRouterNatModel,
	apiClient *sdk.APIClient,
) (*sdk.GetNetworkRouterNat200ResponseNetworkRouterNAT, error) {
	routerID := config.RouterId.ValueInt64()

	if !config.Id.IsNull() {
		return getNatByID(ctx, config.Id.ValueInt64(), routerID, apiClient)
	} else if !config.Name.IsNull() {
		return getNatByName(ctx, config.Name.ValueString(), routerID, apiClient)
	}

	return nil, errors.New(ErrorNoValidSearchTerms)
}

// Read refreshes the Terraform state with the latest data.
func (d *DataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	var config NetworkRouterNatModel

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

	nat, err := getNat(ctx, &config, apiClient)
	if err != nil {
		resp.Diagnostics.AddError(
			summary,
			err.Error(),
		)

		return
	}

	routerID := config.RouterId.ValueInt64()
	state := natAsState(nat, routerID)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
