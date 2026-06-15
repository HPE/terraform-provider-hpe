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
	state.DnsServers = convert.StrSliceToSet(pool.DnsServers)
	state.DnsSuffixList = convert.StrSliceToSet(pool.DnsSuffixList)

	// NullableString fields
	state.DisplayName = convert.StrToType(pool.DisplayName.Get())
	state.Gateway = convert.StrToType(pool.Gateway.Get())
	state.Netmask = convert.StrToType(pool.Netmask.Get())
	state.SubnetAddress = convert.StrToType(pool.SubnetAddress.Get())
	state.DnsDomain = convert.StrToType(pool.DnsDomain.Get())
	state.DnsSearchPath = convert.StrToType(pool.DnsSearchPath.Get())
	state.BootFile = convert.StrToType(pool.BootFile.Get())
	state.Category = convert.StrToType(pool.Category.Get())
	state.Code = convert.StrToType(pool.Code.Get())
	state.InternalId = convert.StrToType(pool.InternalId.Get())
	state.ExternalId = convert.StrToType(pool.ExternalId.Get())
	state.DhcpIp = convert.StrToType(pool.DhcpIp.Get())
	state.HostPrefix = convert.StrToType(pool.HostPrefix.Get())
	state.HttpProxy = convert.StrToType(pool.HttpProxy.Get())
	state.TftpServer = convert.StrToType(pool.TftpServer.Get())

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
