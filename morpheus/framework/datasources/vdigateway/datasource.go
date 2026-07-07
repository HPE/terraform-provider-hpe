// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package vdigateway

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
	summary                 = "read vdi gateway data source"
	ErrorNoValidSearchTerms = `no valid search terms - an id or name is required`
	ErrorRunningPreApply    = `Error running pre-apply plan: exit status 1`
	ErrorNoVdiGatewayFound  = `no vdi gateway found`
	ErrorMultipleVdiGWs     = `multiple vdi gateways were returned`
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
	resp.TypeName = req.ProviderTypeName + "_" + "vdi_gateway"
}

// Schema defines the schema for the data source.
func (d *DataSource) Schema(
	ctx context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = VdiGatewayDataSourceSchema(ctx)
}

func vdiGatewayAsState(
	gw *sdk.GetVDIGateways200ResponseVdiGateway,
) VdiGatewayModel {
	state := VdiGatewayModel{
		Id:          convert.Int64ToType(gw.Id),
		Name:        convert.StrToType(gw.Name),
		DateCreated: convert.TimeToType(gw.DateCreated),
		LastUpdated: convert.TimeToType(gw.LastUpdated),
	}

	// Nullable fields
	if gw.Description.IsSet() {
		state.Description = convert.StrToType(gw.Description.Get())
	} else {
		state.Description = types.StringNull()
	}

	if gw.GatewayUrl.IsSet() {
		state.GatewayUrl = convert.StrToType(gw.GatewayUrl.Get())
	} else {
		state.GatewayUrl = types.StringNull()
	}

	return state
}

func getByID(
	ctx context.Context,
	id int64,
	apiClient *sdk.APIClient,
) (*sdk.GetVDIGateways200ResponseVdiGateway, error) {
	r, hresp, err := apiClient.VDIAPI.GetVDIGateways(ctx, id).Execute()
	if r == nil || err != nil || hresp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GET failed for vdi gateway %d: %s", id, providererrors.ErrMsg(err, hresp))
	}

	if r.VdiGateway == nil {
		return nil, fmt.Errorf("GET failed for vdi gateway %d: response missing vdiGateway", id)
	}

	gw := *r.VdiGateway

	return &gw, nil
}

func getByName(
	ctx context.Context,
	name string,
	apiClient *sdk.APIClient,
) (*sdk.GetVDIGateways200ResponseVdiGateway, error) {
	rs, hresp, err := apiClient.VDIAPI.ListVDIGateways(ctx).Name(name).Execute()
	if rs == nil || err != nil || hresp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GET failed for vdi gateway %s: %s", name, providererrors.ErrMsg(err, hresp))
	}

	var matched []sdk.ListVDIGateways200ResponseAllOfVdiGatewaysInner

	for _, o := range rs.VdiGateways {
		if o.Name != nil && *o.Name == name {
			matched = append(matched, o)
		}
	}

	if len(matched) == 0 {
		return nil, errors.New(ErrorNoVdiGatewayFound)
	} else if len(matched) > 1 {
		return nil, errors.New(ErrorMultipleVdiGWs)
	}

	if matched[0].Id == nil {
		return nil, fmt.Errorf("GET failed for vdi gateway %s: response missing id", name)
	}

	return getByID(ctx, *matched[0].Id, apiClient)
}

func getVdiGateway(
	ctx context.Context,
	config *VdiGatewayModel,
	apiClient *sdk.APIClient,
) (*sdk.GetVDIGateways200ResponseVdiGateway, error) {
	if !config.Id.IsNull() {
		return getByID(ctx, config.Id.ValueInt64(), apiClient)
	} else if !config.Name.IsNull() {
		return getByName(ctx, config.Name.ValueString(), apiClient)
	}

	return nil, errors.New(ErrorNoValidSearchTerms)
}

// Read refreshes the Terraform state with the latest data.
func (d *DataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	var config VdiGatewayModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiClient, err := d.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError(summary, "could not create sdk client")

		return
	}

	gw, err := getVdiGateway(ctx, &config, apiClient)
	if err != nil {
		resp.Diagnostics.AddError(summary, err.Error())

		return
	}

	state := vdiGatewayAsState(gw)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
