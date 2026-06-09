// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

// Package networkrouter implements a data source for network_router
package networkrouter

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
	summary                     = "read network router data source"
	ErrorNoValidSearchTerms     = `no valid search terms - an id or name is required`
	ErrorRunningPreApply        = `Error running pre-apply plan: exit status 1`
	ErrorNoNetworkRouterFound   = `no network router found`
	ErrorMultipleNetworkRouters = `multiple network routers were returned`
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
	resp.TypeName = req.ProviderTypeName + "_" + constants.SubProviderName + "_network_router"
}

// Schema defines the schema for the data source.
func (d *DataSource) Schema(
	ctx context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = NetworkRouterDataSourceSchema(ctx)
}

// networkRouterListItem is a minimal struct used to decode the untyped list response.
type networkRouterListItem struct {
	Id   *int64  `json:"id,omitempty"`
	Name *string `json:"name,omitempty"`
}

func networkRouterAsState(
	ctx context.Context,
	router *sdk.GetNetworkRouter200ResponseNetworkRouter,
) (NetworkRouterModel, error) {
	state := NetworkRouterModel{
		Id:         convert.Int64ToType(router.Id),
		Name:       convert.StrToType(router.Name),
		Code:       convert.StrToType(router.Code),
		Enabled:    convert.BoolToType(router.Enabled),
		EnableBgp:  convert.BoolToType(router.EnableBgp),
		ProviderId: convert.StrToType(router.ProviderId),
	}

	config, err := convert.MapToDynamic(ctx, router.Config)
	if err != nil {
		return NetworkRouterModel{}, fmt.Errorf("network router config mapping failed: %w", err)
	}

	state.Config = config

	if err := setCloudState(ctx, &state, router); err != nil {
		return NetworkRouterModel{}, err
	}

	if err := setGroupState(ctx, &state, router); err != nil {
		return NetworkRouterModel{}, err
	}

	if err := setNetworkIntegrationState(ctx, &state, router); err != nil {
		return NetworkRouterModel{}, err
	}

	if err := setPermissionsState(ctx, &state, router); err != nil {
		return NetworkRouterModel{}, err
	}

	if err := setInterfacesState(ctx, &state, router); err != nil {
		return NetworkRouterModel{}, err
	}

	return state, nil
}

func setCloudState(
	ctx context.Context,
	state *NetworkRouterModel,
	router *sdk.GetNetworkRouter200ResponseNetworkRouter,
) error {
	if router.Zone != nil {
		cloud, diags := NewCloudValue(
			CloudValue{}.AttributeTypes(ctx),
			map[string]attr.Value{
				"id":   convert.Int64ToType(router.Zone.Id),
				"code": convert.StrToType(router.Zone.Code),
				"name": convert.StrToType(router.Zone.Name),
			},
		)
		if diags.HasError() {
			return fmt.Errorf("error creating cloud value")
		}

		state.Cloud = cloud
	} else {
		state.Cloud = NewCloudValueNull()
	}

	return nil
}

func setGroupState(
	ctx context.Context,
	state *NetworkRouterModel,
	router *sdk.GetNetworkRouter200ResponseNetworkRouter,
) error {
	if router.Site != nil {
		group, diags := NewGroupValue(
			GroupValue{}.AttributeTypes(ctx),
			map[string]attr.Value{
				"id":   convert.Int64ToType(router.Site.Id),
				"name": convert.StrToType(router.Site.Name),
			},
		)
		if diags.HasError() {
			return fmt.Errorf("error creating group value")
		}

		state.Group = group
	} else {
		state.Group = NewGroupValueNull()
	}

	return nil
}

func setNetworkIntegrationState(
	ctx context.Context,
	state *NetworkRouterModel,
	router *sdk.GetNetworkRouter200ResponseNetworkRouter,
) error {
	if router.NetworkServer != nil {
		ni, diags := NewNetworkIntegrationValue(
			NetworkIntegrationValue{}.AttributeTypes(ctx),
			map[string]attr.Value{
				"id":   convert.Int64ToType(router.NetworkServer.Id),
				"name": convert.StrToType(router.NetworkServer.Name),
			},
		)
		if diags.HasError() {
			return fmt.Errorf("error creating network integration value")
		}

		state.NetworkIntegration = ni
	} else {
		state.NetworkIntegration = NewNetworkIntegrationValueNull()
	}

	return nil
}

func setPermissionsState(
	ctx context.Context,
	state *NetworkRouterModel,
	router *sdk.GetNetworkRouter200ResponseNetworkRouter,
) error {
	if router.Permissions == nil {
		state.Permissions = NewPermissionsValueNull()

		return nil
	}

	var tenantPermsVal types.Object
	if router.Permissions.TenantPermissions != nil {
		accounts := convert.Int64SliceToSet(router.Permissions.TenantPermissions.Accounts)

		tp, diags := NewTenantPermissionsValue(
			TenantPermissionsValue{}.AttributeTypes(ctx),
			map[string]attr.Value{
				"accounts": accounts,
			},
		)
		if diags.HasError() {
			return fmt.Errorf("error creating tenant permissions value")
		}

		tpObj, diags := tp.ToObjectValue(ctx)
		if diags.HasError() {
			return fmt.Errorf("error converting tenant permissions to object")
		}

		tenantPermsVal = tpObj
	} else {
		tenantPermsVal = types.ObjectNull(TenantPermissionsValue{}.AttributeTypes(ctx))
	}

	perms, diags := NewPermissionsValue(
		PermissionsValue{}.AttributeTypes(ctx),
		map[string]attr.Value{
			"visibility":         convert.StrToType(router.Permissions.Visibility),
			"tenant_permissions": tenantPermsVal,
		},
	)
	if diags.HasError() {
		return fmt.Errorf("error creating permissions value")
	}

	state.Permissions = perms

	return nil
}

func setInterfacesState(
	ctx context.Context,
	state *NetworkRouterModel,
	router *sdk.GetNetworkRouter200ResponseNetworkRouter,
) error {
	if len(router.Interfaces) == 0 {
		state.Interfaces = types.SetNull(InterfacesValue{}.Type(ctx))

		return nil
	}

	ifaceValues := make([]attr.Value, 0, len(router.Interfaces))

	for _, iface := range router.Interfaces {
		iv, diags := NewInterfacesValue(
			InterfacesValue{}.AttributeTypes(ctx),
			map[string]attr.Value{
				"id":         types.Int64PointerValue(iface.Id),
				"cidr":       types.StringPointerValue(iface.Cidr),
				"ip_address": types.StringPointerValue(iface.IpAddress),
			},
		)
		if diags.HasError() {
			return fmt.Errorf("error creating interface value")
		}

		objVal, diags := iv.ToObjectValue(ctx)
		if diags.HasError() {
			return fmt.Errorf("error converting interface to object")
		}

		ifaceValues = append(ifaceValues, objVal)
	}

	setVal, diags := types.SetValue(InterfacesValue{}.Type(ctx), ifaceValues)
	if diags.HasError() {
		return fmt.Errorf("error creating interfaces set")
	}

	state.Interfaces = setVal

	return nil
}

func getNetworkRouterByID(
	ctx context.Context,
	id int64,
	apiClient *sdk.APIClient,
) (*sdk.GetNetworkRouter200ResponseNetworkRouter, error) {
	r, hresp, err := apiClient.NetworksAPI.GetNetworkRouter(ctx, id).Execute()
	if r == nil || err != nil || hresp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"GET failed for network router %d: %s",
			id, providererrors.ErrMsg(err, hresp),
		)
	}

	if r.NetworkRouter == nil {
		return nil, fmt.Errorf("GET failed for network router %d: missing networkRouter payload", id)
	}

	router := *r.NetworkRouter

	return &router, nil
}

func getNetworkRouterByName(
	ctx context.Context,
	name string,
	apiClient *sdk.APIClient,
) (*sdk.GetNetworkRouter200ResponseNetworkRouter, error) {
	rs, hresp, err := apiClient.NetworksAPI.GetNetworkRouters(ctx).Execute()
	if rs == nil || err != nil || hresp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"GET failed for network router %s: %s",
			name, providererrors.ErrMsg(err, hresp),
		)
	}

	raw, err := json.Marshal(rs.NetworkRouters)
	if err != nil {
		return nil, fmt.Errorf("error marshalling network routers list: %w", err)
	}

	var items []networkRouterListItem
	if err := json.Unmarshal(raw, &items); err != nil {
		return nil, fmt.Errorf("error decoding network routers list: %w", err)
	}

	var matchedIDs []int64

	for _, item := range items {
		if item.Name != nil && *item.Name == name && item.Id != nil {
			matchedIDs = append(matchedIDs, *item.Id)
		}
	}

	if len(matchedIDs) == 0 {
		return nil, errors.New(ErrorNoNetworkRouterFound)
	} else if len(matchedIDs) > 1 {
		return nil, errors.New(ErrorMultipleNetworkRouters)
	}

	return getNetworkRouterByID(ctx, matchedIDs[0], apiClient)
}

func getNetworkRouter(
	ctx context.Context,
	config *NetworkRouterModel,
	apiClient *sdk.APIClient,
) (*sdk.GetNetworkRouter200ResponseNetworkRouter, error) {
	if !config.Id.IsNull() {
		return getNetworkRouterByID(ctx, config.Id.ValueInt64(), apiClient)
	} else if !config.Name.IsNull() {
		return getNetworkRouterByName(ctx, config.Name.ValueString(), apiClient)
	}

	return nil, errors.New(ErrorNoValidSearchTerms)
}

// Read refreshes the Terraform state with the latest data.
func (d *DataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	var config NetworkRouterModel

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

	router, err := getNetworkRouter(ctx, &config, apiClient)
	if err != nil {
		resp.Diagnostics.AddError(
			summary,
			err.Error(),
		)

		return
	}

	state, err := networkRouterAsState(ctx, router)
	if err != nil {
		resp.Diagnostics.AddError(
			summary,
			err.Error(),
		)

		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
