// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networktype

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/utils/convert"
)

func getNetworkType(
	ctx context.Context,
	config NetworkTypeModel,
	client *sdk.APIClient,
) (*NetworkTypeModel, error) {
	if !config.Id.IsNull() {
		return getNetworkTypeByID(ctx, config.Id.ValueInt64(), client)
	}

	if !config.Name.IsNull() {
		return getNetworkTypeByName(ctx, config.Name.ValueString(), client)
	}

	return nil, fmt.Errorf("either id or name must be specified")
}

func getNetworkTypeByID(
	ctx context.Context,
	id int64,
	client *sdk.APIClient,
) (*NetworkTypeModel, error) {
	response, hresp, err := client.NetworksAPI.GetNetworkType(ctx, id).Execute()
	if err != nil || hresp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("network type %d GET failed: %s", id, errfmt.ErrMsg(err, hresp))
	}

	if response.NetworkType == nil {
		return nil, fmt.Errorf("network type %d is nil", id)
	}

	nt := response.NetworkType
	state := &NetworkTypeModel{}

	state.Id = types.Int64Value(id)
	state.Name = convert.StrToType(nt.Name)
	state.Code = convert.StrToType(nt.Code)
	state.Creatable = convert.BoolToType(nt.Creatable)
	state.Overlay = convert.BoolToType(nt.Overlay)
	state.NameEditable = convert.BoolToType(nt.NameEditable)
	state.CidrEditable = convert.BoolToType(nt.CidrEditable)
	state.CidrRequired = convert.BoolToType(nt.CidrRequired)
	state.CanAssignPool = convert.BoolToType(nt.CanAssignPool)
	state.HasCidr = convert.BoolToType(nt.HasCidr)
	state.VlanIdEditable = convert.BoolToType(nt.VlanIdEditable)
	state.Deletable = convert.BoolToType(nt.Deletable)
	state.DhcpServerEditable = convert.BoolToType(nt.DhcpServerEditable)
	state.DnsEditable = convert.BoolToType(nt.DnsEditable)
	state.GatewayEditable = convert.BoolToType(nt.GatewayEditable)
	state.StaticOverrideEditable = convert.BoolToType(nt.StaticOverrideEditable)
	state.NetworkDomainEditable = convert.BoolToType(nt.NetworkDomainEditable)
	state.HasNetworkServer = convert.BoolToType(nt.HasNetworkServer)
	state.HasStaticRoutes = convert.BoolToType(nt.HasStaticRoutes)
	state.HasFloatingIps = convert.BoolToType(nt.HasFloatingIps)

	// NullableString fields
	state.Description = convert.StrToType(nt.Description.Get())
	state.Category = convert.StrToType(nt.Category.Get())
	state.ExternalType = convert.StrToType(nt.ExternalType.Get())

	return state, nil
}

func getNetworkTypeByName(
	ctx context.Context,
	name string,
	client *sdk.APIClient,
) (*NetworkTypeModel, error) {
	response, hresp, err := client.NetworksAPI.ListNetworkTypes(ctx).Name(name).Execute()
	if err != nil || hresp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("network type %s list failed: %s", name, errfmt.ErrMsg(err, hresp))
	}

	var matchingNetworkTypes []sdk.ListNetworkTypes200ResponseAllOfNetworkTypesInner

	for _, networkType := range response.NetworkTypes {
		if networkType.Name != nil && *networkType.Name == name {
			matchingNetworkTypes = append(matchingNetworkTypes, networkType)
		}
	}

	if len(matchingNetworkTypes) == 0 {
		return nil, fmt.Errorf("network type %s not found", name)
	}

	if len(matchingNetworkTypes) > 1 {
		var networkTypeIDs []string

		for _, n := range matchingNetworkTypes {
			if n.Id != nil {
				networkTypeIDs = append(networkTypeIDs, fmt.Sprintf("%d", *n.Id))
			}
		}

		return nil, fmt.Errorf(
			"multiple network types found with name %s. Network type IDs: %s. "+
				"Please specify an ID instead",
			name,
			strings.Join(networkTypeIDs, ", "),
		)
	}

	id := matchingNetworkTypes[0].Id
	if id == nil {
		return nil, fmt.Errorf("network type %s has missing ID", name)
	}

	return getNetworkTypeByID(ctx, *id, client)
}

func (d *DataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	var config NetworkTypeModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client, err := d.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"read network type data source",
			fmt.Sprintf("failed to create client: %s", err.Error()),
		)

		return
	}

	state, err := getNetworkType(ctx, config, client)
	if err != nil {
		resp.Diagnostics.AddError(
			"read network type data source",
			err.Error(),
		)

		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}
