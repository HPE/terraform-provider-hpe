//  (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package client

import (
	"context"
	"encoding/json"
	"fmt"

	consts "github.com/HPE/terraform-provider-hpe/greenlake/cloud/sdk/vmaascmp/common"
	"github.com/HPE/terraform-provider-hpe/greenlake/cloud/sdk/vmaascmp/models"
)

const networkCompatibleVersion = "5.2.10"

type NetworksAPIService struct {
	Client APIClientHandler
	Cfg    Configuration
}

func (n *NetworksAPIService) GetAllNetworks(
	ctx context.Context,
	param map[string]string,
) (models.ListNetworksBody, error) {
	var NetworksResp models.ListNetworksBody
	networkAPI := &api{
		method: "GET",
		path:   consts.NetworksPath,
		client: n.Client,

		jsonParser: func(body []byte) error {
			return json.Unmarshal(body, &NetworksResp)
		},
	}
	err := networkAPI.do(ctx, nil, param)

	return NetworksResp, err
}

func (n *NetworksAPIService) GetSpecificNetwork(
	ctx context.Context,
	networkID int,
) (models.GetSpecificNetworkBody, error) {
	var NetworksResp models.GetSpecificNetworkBody
	networkAPI := &api{
		method: "GET",
		path:   fmt.Sprintf("%s/%d", consts.NetworksPath, networkID),
		client: n.Client,

		jsonParser: func(body []byte) error {
			return json.Unmarshal(body, &NetworksResp)
		},
	}
	err := networkAPI.do(ctx, nil, nil)

	return NetworksResp, err
}

func (n *NetworksAPIService) CreateNetwork(
	ctx context.Context,
	networkReq models.CreateNetworkRequest,
) (models.CreateNetworkResponse, error) {
	var NetworksResp models.CreateNetworkResponse
	if v, _ := parseVersion("5.4.4"); v <= n.Client.getVersion() {
		// network Pool is not required for DHCP
		if networkReq.Network.NetworkPool != nil {
			networkReq.Network.NetworkPool.Pool = networkReq.Network.PoolID
			networkReq.Network.PoolID = 0
		}
	}
	networkAPI := &api{
		compatibleVersion: "5.2.13",
		method:            "POST",
		path:              consts.NetworksPath,
		client:            n.Client,

		jsonParser: func(body []byte) error {
			return json.Unmarshal(body, &NetworksResp)
		},
	}
	err := networkAPI.do(ctx, networkReq, nil)

	return NetworksResp, err
}

func (n *NetworksAPIService) DeleteNetwork(
	ctx context.Context,
	networkID int,
) (models.SuccessOrErrorMessage, error) {
	var Output models.SuccessOrErrorMessage
	networkAPI := &api{
		compatibleVersion: networkCompatibleVersion,
		method:            "DELETE",
		path:              fmt.Sprintf("%s/%d", consts.NetworksPath, networkID),
		client:            n.Client,

		jsonParser: func(body []byte) error {
			return json.Unmarshal(body, &Output)
		},
	}
	err := networkAPI.do(ctx, nil, nil)

	return Output, err
}

func (n *NetworksAPIService) GetNetworkType(
	ctx context.Context,
	params map[string]string,
) (models.GetNetworkTypesResponse, error) {
	var Resp models.GetNetworkTypesResponse
	networkAPI := &api{
		compatibleVersion: networkCompatibleVersion,
		method:            "GET",
		path:              consts.NetworkTypePath,
		client:            n.Client,

		jsonParser: func(body []byte) error {
			return json.Unmarshal(body, &Resp)
		},
	}
	err := networkAPI.do(ctx, nil, params)

	return Resp, err
}

func (n *NetworksAPIService) GetNetworkPool(
	ctx context.Context,
	params map[string]string,
) (models.GetNetworkPoolsResp, error) {
	var Resp models.GetNetworkPoolsResp
	networkAPI := &api{
		method: "GET",
		path:   fmt.Sprintf("%s/%s", consts.NetworksPath, consts.NetworkPoolPath),
		client: n.Client,

		jsonParser: func(body []byte) error {
			return json.Unmarshal(body, &Resp)
		},
	}
	err := networkAPI.do(ctx, nil, params)

	return Resp, err
}

func (n *NetworksAPIService) GetSpecificNetworkPool(
	ctx context.Context,
	networkPoolID int,
) (models.GetSpecificNetworkPool, error) {
	var Resp models.GetSpecificNetworkPool
	networkAPI := &api{
		method: "GET",
		path: fmt.Sprintf("%s/%s/%d",
			consts.NetworksPath, consts.NetworkPoolPath, networkPoolID),
		client: n.Client,

		jsonParser: func(body []byte) error {
			return json.Unmarshal(body, &Resp)
		},
	}
	err := networkAPI.do(ctx, nil, nil)

	return Resp, err
}

func (n *NetworksAPIService) UpdateNetwork(
	ctx context.Context,
	networkID int,
	request models.CreateNetworkRequest,
) (models.SuccessOrErrorMessage, error) {
	var Output models.SuccessOrErrorMessage
	networkAPI := &api{
		compatibleVersion: "5.2.13",
		method:            "PUT",
		path:              fmt.Sprintf("%s/%d", consts.NetworksPath, networkID),
		client:            n.Client,

		jsonParser: func(body []byte) error {
			return json.Unmarshal(body, &Output)
		},
	}
	err := networkAPI.do(ctx, request, nil)

	return Output, err
}

func (n *NetworksAPIService) GetNetworkProxy(
	ctx context.Context,
	params map[string]string,
) (models.GetAllNetworkProxies, error) {
	var ProxyResp models.GetAllNetworkProxies
	networkAPI := &api{
		method: "GET",
		path:   fmt.Sprintf("%s/%s", consts.NetworksPath, consts.NetworkProxyPath),
		client: n.Client,

		jsonParser: func(body []byte) error {
			return json.Unmarshal(body, &ProxyResp)
		},
	}
	err := networkAPI.do(ctx, nil, params)

	return ProxyResp, err
}

func (n *NetworksAPIService) GetSpecificNetworkProxy(
	ctx context.Context,
	proxyID int,
) (models.GetSpecificNetworkProxy, error) {
	var ProxyResp models.GetSpecificNetworkProxy
	networkAPI := &api{
		method: "GET",
		path:   fmt.Sprintf("%s/%s/%d", consts.NetworksPath, consts.NetworkProxyPath, proxyID),
		client: n.Client,

		jsonParser: func(body []byte) error {
			return json.Unmarshal(body, &ProxyResp)
		},
	}
	err := networkAPI.do(ctx, nil, nil)

	return ProxyResp, err
}

func (n *NetworksAPIService) GetSpecificNetworkType(
	ctx context.Context,
	typeID int,
) (models.GetaNetworkType, error) {
	var Resp models.GetaNetworkType
	networkAPI := &api{
		compatibleVersion: networkCompatibleVersion,
		method:            "GET",
		path:              fmt.Sprintf("%s/%d", consts.NetworkTypePath, typeID),
		client:            n.Client,

		jsonParser: func(body []byte) error {
			return json.Unmarshal(body, &Resp)
		},
	}
	err := networkAPI.do(ctx, nil, nil)

	return Resp, err
}
