// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkpool

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/utils/convert"
)

func getNetworkPool(
	ctx context.Context,
	config NetworkPoolModel,
	client *sdk.APIClient,
) (*NetworkPoolModel, error) {
	if !config.Id.IsNull() {
		return getNetworkPoolByID(ctx, config.Id.ValueInt64(), client)
	}

	if !config.Name.IsNull() {
		return getNetworkPoolByName(ctx, config.Name.ValueString(), client)
	}

	return nil, fmt.Errorf("either id or name must be specified")
}

func getNetworkPoolByID(
	ctx context.Context,
	id int64,
	client *sdk.APIClient,
) (*NetworkPoolModel, error) {
	response, hresp, err := client.NetworksAPI.GetNetworkPool(ctx, id).Execute()
	if err != nil || hresp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("network pool %d GET failed: %s", id, errfmt.ErrMsg(err, hresp))
	}

	if response.NetworkPool == nil {
		return nil, fmt.Errorf("network pool %d is nil", id)
	}

	pool := response.NetworkPool
	state := &NetworkPoolModel{}

	state.Id = types.Int64Value(id)
	state.Name = convert.StrToType(pool.Name)
	state.IpCount = convert.Int64ToType(pool.IpCount)
	state.FreeCount = convert.Int64ToType(pool.FreeCount)
	state.PoolEnabled = convert.BoolToType(pool.PoolEnabled)
	state.DhcpServer = convert.BoolToType(pool.DhcpServer)

	if pool.DisplayName.IsSet() {
		state.DisplayName = types.StringValue(*pool.DisplayName.Get())
	} else {
		state.DisplayName = types.StringNull()
	}

	if pool.Gateway.IsSet() {
		state.Gateway = types.StringValue(*pool.Gateway.Get())
	} else {
		state.Gateway = types.StringNull()
	}

	if pool.Netmask.IsSet() {
		state.Netmask = types.StringValue(*pool.Netmask.Get())
	} else {
		state.Netmask = types.StringNull()
	}

	if pool.SubnetAddress.IsSet() {
		state.SubnetAddress = types.StringValue(*pool.SubnetAddress.Get())
	} else {
		state.SubnetAddress = types.StringNull()
	}

	if pool.DnsDomain.IsSet() {
		state.DnsDomain = types.StringValue(*pool.DnsDomain.Get())
	} else {
		state.DnsDomain = types.StringNull()
	}

	if pool.DnsSearchPath.IsSet() {
		state.DnsSearchPath = types.StringValue(*pool.DnsSearchPath.Get())
	} else {
		state.DnsSearchPath = types.StringNull()
	}

	if pool.BootFile.IsSet() {
		state.BootFile = types.StringValue(*pool.BootFile.Get())
	} else {
		state.BootFile = types.StringNull()
	}

	state.DnsServers = convert.StrSliceToSet(pool.DnsServers)

	if pool.Type != nil && pool.Type.Code != nil {
		state.TypeCode = types.StringValue(*pool.Type.Code)
	} else {
		state.TypeCode = types.StringNull()
	}

	return state, nil
}

func getNetworkPoolByName(
	ctx context.Context,
	name string,
	client *sdk.APIClient,
) (*NetworkPoolModel, error) {
	response, hresp, err := client.NetworksAPI.GetNetworkPools(ctx).Name(name).Execute()
	if err != nil || hresp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("network pool %s list failed: %s", name, errfmt.ErrMsg(err, hresp))
	}

	// NetworkPools is typed as interface{} in the SDK because the OpenAPI spec
	// defines it as a free-form object. Type assert to []interface{} to extract
	// individual pool entries.
	poolsRaw := response.NetworkPools
	pools, ok := poolsRaw.([]interface{})
	if !ok || len(pools) == 0 {
		return nil, fmt.Errorf("network pool %s not found", name)
	}

	var matchID int64
	var matchCount int

	for _, p := range pools {
		poolMap, ok := p.(map[string]interface{})
		if !ok {
			continue
		}

		poolName, _ := poolMap["name"].(string)
		if poolName == name {
			matchCount++

			if id, ok := poolMap["id"].(float64); ok {
				matchID = int64(id)
			}
		}
	}

	if matchCount == 0 {
		return nil, fmt.Errorf("network pool %s not found", name)
	}

	if matchCount > 1 {
		return nil, fmt.Errorf(
			"multiple network pools found with name %s, please specify an ID instead",
			name,
		)
	}

	return getNetworkPoolByID(ctx, matchID, client)
}

func (d *DataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	var config NetworkPoolModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client, err := d.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"read network pool data source",
			fmt.Sprintf("failed to create client: %s", err.Error()),
		)

		return
	}

	state, err := getNetworkPool(ctx, config, client)
	if err != nil {
		resp.Diagnostics.AddError(
			"read network pool data source",
			err.Error(),
		)

		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}
