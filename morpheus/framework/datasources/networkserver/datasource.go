// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkserver

import (
	"context"
	"errors"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"

	"github.com/HPE/terraform-provider-hpe/morpheus/configure"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/networkserverlookup"
	"github.com/HPE/terraform-provider-hpe/utils/convert"
)

const (
	summary                 = "read network server data source"
	ErrorNoValidSearchTerms = `no valid search terms - an id or name is required`
	ErrorRunningPreApply    = `Error running pre-apply plan: exit status 1`
	// The lookup's messages (networkserverlookup.ErrNotFound / ErrMultiple),
	// under the names tests have always used; a unit test pins them equal.
	ErrorNoNetworkServerFound   = "no network server found"
	ErrorMultipleNetworkServers = "multiple network servers were returned"
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
	resp.TypeName = req.ProviderTypeName + "_" + "network_server"
}

// Schema defines the schema for the data source.
func (d *DataSource) Schema(
	ctx context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = NetworkServerDataSourceSchema(ctx)
}

func int32PtrToInt64Ptr(v *int32) *int64 {
	if v == nil {
		return nil
	}

	i := int64(*v)

	return &i
}

func networkServerAsState(
	ctx context.Context,
	ns *sdk.GetNetworkServer200ResponseNetworkServer,
) (NetworkServerModel, error) {
	state := NetworkServerModel{
		Id:              convert.Int64ToType(ns.Id),
		Name:            convert.StrToType(ns.Name),
		Description:     convert.StrToType(ns.Description.Get()),
		Enabled:         convert.BoolToType(ns.Enabled),
		Visibility:      convert.StrToType(ns.Visibility),
		Visible:         convert.BoolToType(ns.Visible),
		ExternalId:      convert.StrToType(ns.ExternalId.Get()),
		InternalId:      convert.StrToType(ns.InternalId.Get()),
		ServiceUrl:      convert.StrToType(ns.ServiceUrl.Get()),
		ServiceHost:     convert.StrToType(ns.ServiceHost.Get()),
		ServicePort:     convert.Int64ToType(int32PtrToInt64Ptr(ns.ServicePort.Get())),
		ServiceMode:     convert.StrToType(ns.ServiceMode.Get()),
		ServicePath:     convert.StrToType(ns.ServicePath.Get()),
		ServiceUsername: convert.StrToType(ns.ServiceUsername.Get()),
		ServicePassword: convert.StrToType(ns.ServicePassword.Get()),
		ServiceToken:    convert.StrToType(ns.ServiceToken.Get()),
		ApiPort:         convert.Int64ToType(int32PtrToInt64Ptr(ns.ApiPort.Get())),
		AdminPort:       convert.Int64ToType(int32PtrToInt64Ptr(ns.AdminPort.Get())),
		NetworkFilter:   convert.StrToType(ns.NetworkFilter.Get()),
		TenantMatch:     convert.StrToType(ns.TenantMatch.Get()),
		ZoneId:          convert.Int64ToType(ns.ZoneId.Get()),
	}

	if ns.Config != nil {
		dyn, err := convert.MapToDynamic(ctx, ns.Config)
		if err == nil {
			state.Config = dyn
		} else {
			state.Config = types.DynamicNull()
		}
	} else {
		state.Config = types.DynamicNull()
	}

	return state, nil
}

func getNetworkServer(
	ctx context.Context,
	config *NetworkServerModel,
	apiClient *sdk.APIClient,
) (*sdk.GetNetworkServer200ResponseNetworkServer, error) {
	// Hidden servers — a cluster's own — are found by either lookup; see
	// networkserverlookup.
	if !config.Id.IsNull() {
		return networkserverlookup.ByID(ctx, apiClient, config.Id.ValueInt64())
	} else if !config.Name.IsNull() {
		return networkserverlookup.ByName(ctx, apiClient, config.Name.ValueString())
	}

	return nil, errors.New(ErrorNoValidSearchTerms)
}

// Read refreshes the Terraform state with the latest data.
func (d *DataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	var config NetworkServerModel

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

	ns, err := getNetworkServer(ctx, &config, apiClient)
	if err != nil {
		resp.Diagnostics.AddError(
			summary,
			err.Error(),
		)

		return
	}

	state, err := networkServerAsState(ctx, ns)
	if err != nil {
		resp.Diagnostics.AddError(
			summary,
			err.Error(),
		)

		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
