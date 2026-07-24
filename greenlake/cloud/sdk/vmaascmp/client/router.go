// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package client

import (
	"context"
	"encoding/json"
	"fmt"

	consts "github.com/HPE/terraform-provider-hpe/greenlake/cloud/sdk/vmaascmp/common"
	"github.com/HPE/terraform-provider-hpe/greenlake/cloud/sdk/vmaascmp/models"
)

const routerCompatibleVersion = "5.2.10"

type RouterAPIService struct {
	Client APIClientHandler
	Cfg    Configuration
}

func (r *RouterAPIService) GetAllRouter(
	ctx context.Context,
	queryParams map[string]string,
) (models.GetAllNetworkRouter, error) {
	RouterResp := models.GetAllNetworkRouter{}
	routerAPI := &api{
		compatibleVersion: routerCompatibleVersion,
		method:            "GET",
		path:              fmt.Sprintf("%s/%s", consts.NetworksPath, consts.NetworkRouterPath),
		client:            r.Client,
		jsonParser: func(body []byte) error {
			return json.Unmarshal(body, &RouterResp)
		},
	}
	err := routerAPI.do(ctx, nil, queryParams)

	return RouterResp, err
}

func (r *RouterAPIService) GetSpecificRouter(
	ctx context.Context,
	routerID int,
) (models.GetSpecificRouterResp, error) {
	RouterResp := models.GetSpecificRouterResp{}
	routerAPI := &api{
		compatibleVersion: routerCompatibleVersion,
		method:            "GET",
		path: fmt.Sprintf("%s/%s/%d",
			consts.NetworksPath, consts.NetworkRouterPath, routerID),
		client: r.Client,
		jsonParser: func(body []byte) error {
			return json.Unmarshal(body, &RouterResp)
		},
	}
	err := routerAPI.do(ctx, nil, nil)

	return RouterResp, err
}

func (r *RouterAPIService) CreateRouter(
	ctx context.Context,
	request models.CreateRouterRequest,
) (models.CreateRouterResp, error) {
	RouterResp := models.CreateRouterResp{}
	routerAPI := &api{
		compatibleVersion: routerCompatibleVersion,
		method:            "POST",
		path:              fmt.Sprintf("%s/%s", consts.NetworksPath, consts.NetworkRouterPath),
		client:            r.Client,
		jsonParser: func(body []byte) error {
			return json.Unmarshal(body, &RouterResp)
		},
	}
	err := routerAPI.do(ctx, request, nil)

	return RouterResp, err
}

func (r *RouterAPIService) UpdateRouter(
	ctx context.Context,
	routerID int,
	request models.CreateRouterRequest,
) (models.SuccessOrErrorMessage, error) {
	RouterResp := models.SuccessOrErrorMessage{}
	routerAPI := &api{
		compatibleVersion: routerCompatibleVersion,
		method:            "PUT",
		path: fmt.Sprintf("%s/%s/%d",
			consts.NetworksPath, consts.NetworkRouterPath, routerID),
		client: r.Client,
		jsonParser: func(body []byte) error {
			return json.Unmarshal(body, &RouterResp)
		},
	}
	err := routerAPI.do(ctx, request, nil)

	return RouterResp, err
}

func (r *RouterAPIService) DeleteRouter(
	ctx context.Context,
	routerID int,
) (models.SuccessOrErrorMessage, error) {
	RouterResp := models.SuccessOrErrorMessage{}
	routerAPI := &api{
		compatibleVersion: routerCompatibleVersion,
		method:            "DELETE",
		path: fmt.Sprintf("%s/%s/%d",
			consts.NetworksPath, consts.NetworkRouterPath, routerID),
		client: r.Client,
		jsonParser: func(body []byte) error {
			return json.Unmarshal(body, &RouterResp)
		},
	}
	err := routerAPI.do(ctx, nil, nil)

	return RouterResp, err
}

func (r *RouterAPIService) GetRouterTypes(
	ctx context.Context,
	queryParams map[string]string,
) (models.GetNetworlRouterTypes, error) {
	RouterResp := models.GetNetworlRouterTypes{}
	routerAPI := &api{
		compatibleVersion: routerCompatibleVersion,
		method:            "GET",
		path:              consts.NetworkRouterTypePath,
		client:            r.Client,
		jsonParser: func(body []byte) error {
			return json.Unmarshal(body, &RouterResp)
		},
	}
	err := routerAPI.do(ctx, nil, queryParams)

	return RouterResp, err
}

func (r *RouterAPIService) GetNetworkServices(
	ctx context.Context,
	queryParams map[string]string,
) (models.GetNetworkServicesResp, error) {
	RouterResp := models.GetNetworkServicesResp{}
	routerAPI := &api{
		compatibleVersion: routerCompatibleVersion,
		method:            "GET",
		path: fmt.Sprintf("%s/%s", consts.NetworksPath,
			consts.NetworkServicePath),
		client: r.Client,
		jsonParser: func(body []byte) error {
			return json.Unmarshal(body, &RouterResp)
		},
	}
	err := routerAPI.do(ctx, nil, queryParams)

	return RouterResp, err
}

func (r *RouterAPIService) RefreshNetworkServices(
	ctx context.Context,
	serverID int,
	queryParams map[string]string,
) (models.SuccessOrErrorMessage, error) {
	ServerResp := models.SuccessOrErrorMessage{}
	serverAPI := &api{
		compatibleVersion: consts.CMPSixZeroFiveVersion,
		method:            "POST",
		path: fmt.Sprintf("%s/%s/%d/%s", consts.NetworksPath,
			consts.ServerPath, serverID, consts.RefreshPath),
		client: r.Client,
		jsonParser: func(body []byte) error {
			return json.Unmarshal(body, &ServerResp)
		},
	}
	err := serverAPI.do(ctx, nil, queryParams)

	return ServerResp, err
}

func (r *RouterAPIService) CreateRouterNat(
	ctx context.Context,
	routerID int,
	request models.CreateRouterNatRequest,
) (models.SuccessOrErrorMessage, error) {
	NatResp := models.SuccessOrErrorMessage{}
	routerAPI := &api{
		compatibleVersion: routerCompatibleVersion,
		method:            "POST",
		path: fmt.Sprintf("%s/%s/%d/%s", consts.NetworksPath,
			consts.NetworkRouterPath, routerID, consts.RoutersNatPath),
		client: r.Client,
		jsonParser: func(body []byte) error {
			return json.Unmarshal(body, &NatResp)
		},
	}
	err := routerAPI.do(ctx, request, nil)

	return NatResp, err
}

func (r *RouterAPIService) GetSpecificRouterNat(
	ctx context.Context,
	routerID, natID int,
) (models.GetSpecificRouterNatResponse, error) {
	NatResp := models.GetSpecificRouterNatResponse{}
	routerAPI := &api{
		compatibleVersion: routerCompatibleVersion,
		method:            "GET",
		path: fmt.Sprintf("%s/%s/%d/%s/%d", consts.NetworksPath,
			consts.NetworkRouterPath, routerID, consts.RoutersNatPath, natID),
		client: r.Client,
		jsonParser: func(body []byte) error {
			return json.Unmarshal(body, &NatResp)
		},
	}
	err := routerAPI.do(ctx, nil, nil)

	return NatResp, err
}

func (r *RouterAPIService) UpdateRouterNat(
	ctx context.Context,
	routerID, natID int,
	req models.CreateRouterNatRequest,
) (models.SuccessOrErrorMessage, error) {
	NatResp := models.SuccessOrErrorMessage{}
	routerAPI := &api{
		compatibleVersion: routerCompatibleVersion,
		method:            "PUT",
		path: fmt.Sprintf("%s/%s/%d/%s/%d", consts.NetworksPath,
			consts.NetworkRouterPath, routerID, consts.RoutersNatPath, natID),
		client: r.Client,
		jsonParser: func(body []byte) error {
			return json.Unmarshal(body, &NatResp)
		},
	}
	err := routerAPI.do(ctx, req, nil)

	return NatResp, err
}

func (r *RouterAPIService) DeleteRouterNat(
	ctx context.Context,
	routerID, natID int,
) (models.SuccessOrErrorMessage, error) {
	NatResp := models.SuccessOrErrorMessage{}
	routerAPI := &api{
		compatibleVersion: routerCompatibleVersion,
		method:            "DELETE",
		path: fmt.Sprintf("%s/%s/%d/%s/%d", consts.NetworksPath,
			consts.NetworkRouterPath, routerID, consts.RoutersNatPath, natID),
		client: r.Client,
		jsonParser: func(body []byte) error {
			return json.Unmarshal(body, &NatResp)
		},
	}
	err := routerAPI.do(ctx, nil, nil)

	return NatResp, err
}

func (r *RouterAPIService) CreateRouterFirewallRuleGroup(
	ctx context.Context,
	routerID int,
	request models.CreateRouterFirewallRuleGroupRequest,
) (models.SuccessOrErrorMessage, error) {
	FirewallGroupResp := models.SuccessOrErrorMessage{}
	routerAPI := &api{
		compatibleVersion: routerCompatibleVersion,
		method:            "POST",
		path: fmt.Sprintf("%s/%s/%d/%s", consts.NetworksPath,
			consts.NetworkRouterPath, routerID, consts.RoutersFirewallRuleGroupPath),
		client: r.Client,
		jsonParser: func(body []byte) error {
			return json.Unmarshal(body, &FirewallGroupResp)
		},
	}
	err := routerAPI.do(ctx, request, nil)

	return FirewallGroupResp, err
}

func (r *RouterAPIService) GetSpecificRouterFirewallRuleGroup(
	ctx context.Context,
	routerID, firewallGroupID int,
) (models.GetSpecificRouterFirewallRuleGroupResponse, error) {
	FirewallGroupResp := models.GetSpecificRouterFirewallRuleGroupResponse{}
	routerAPI := &api{
		compatibleVersion: routerCompatibleVersion,
		method:            "GET",
		path: fmt.Sprintf("%s/%s/%d/%s/%d", consts.NetworksPath,
			consts.NetworkRouterPath, routerID, consts.RoutersFirewallRuleGroupPath, firewallGroupID),
		client: r.Client,
		jsonParser: func(body []byte) error {
			return json.Unmarshal(body, &FirewallGroupResp)
		},
	}
	err := routerAPI.do(ctx, nil, nil)

	return FirewallGroupResp, err
}

func (r *RouterAPIService) DeleteRouterFirewallRuleGroup(
	ctx context.Context,
	routerID, firewallGroupID int,
) (models.SuccessOrErrorMessage, error) {
	FirewallGroupResp := models.SuccessOrErrorMessage{}
	routerAPI := &api{
		compatibleVersion: routerCompatibleVersion,
		method:            "DELETE",
		path: fmt.Sprintf("%s/%s/%d/%s/%d", consts.NetworksPath,
			consts.NetworkRouterPath, routerID, consts.RoutersFirewallRuleGroupPath, firewallGroupID),
		client: r.Client,
		jsonParser: func(body []byte) error {
			return json.Unmarshal(body, &FirewallGroupResp)
		},
	}
	err := routerAPI.do(ctx, nil, nil)

	return FirewallGroupResp, err
}

func (r *RouterAPIService) CreateRouterRoute(
	ctx context.Context,
	routerID int,
	req models.CreateRouterRoute,
) (models.SuccessOrErrorMessage, error) {
	Resp := models.SuccessOrErrorMessage{}
	routerAPI := &api{
		compatibleVersion: "5.2.12",
		method:            "POST",
		path: fmt.Sprintf("%s/%s/%d/%s", consts.NetworksPath,
			consts.NetworkRouterPath, routerID, consts.RouterRoutePath),
		client: r.Client,
		jsonParser: func(body []byte) error {
			return json.Unmarshal(body, &Resp)
		},
	}
	err := routerAPI.do(ctx, req, nil)

	return Resp, err
}

func (r *RouterAPIService) GetSpecificRouterRoute(
	ctx context.Context,
	routerID, routeID int,
) (models.GetSpecificRouterRoute, error) {
	Resp := models.GetSpecificRouterRoute{}
	routerAPI := &api{
		compatibleVersion: "5.2.12",
		method:            "GET",
		path: fmt.Sprintf("%s/%s/%d/%s/%d", consts.NetworksPath,
			consts.NetworkRouterPath, routerID, consts.RouterRoutePath, routeID),
		client: r.Client,
		jsonParser: func(body []byte) error {
			return json.Unmarshal(body, &Resp)
		},
	}
	err := routerAPI.do(ctx, nil, nil)

	return Resp, err
}

func (r *RouterAPIService) DeleteRouterRoute(
	ctx context.Context,
	routerID, routeID int,
) (models.SuccessOrErrorMessage, error) {
	Resp := models.SuccessOrErrorMessage{}
	routerAPI := &api{
		compatibleVersion: "5.2.12",
		method:            "DELETE",
		path: fmt.Sprintf("%s/%s/%d/%s/%d", consts.NetworksPath,
			consts.NetworkRouterPath, routerID, consts.RouterRoutePath, routeID),
		client: r.Client,
		jsonParser: func(body []byte) error {
			return json.Unmarshal(body, &Resp)
		},
	}
	err := routerAPI.do(ctx, nil, nil)

	return Resp, err
}

func (r *RouterAPIService) CreateRouterBgpNeighbor(
	ctx context.Context,
	routerID int,
	req models.CreateNetworkRouterBgpNeighborRequest,
) (models.SuccessOrErrorMessage, error) {
	Resp := models.SuccessOrErrorMessage{}
	routerAPI := &api{
		compatibleVersion: "5.2.12",
		method:            "POST",
		path: fmt.Sprintf("%s/%s/%d/%s", consts.NetworksPath,
			consts.NetworkRouterPath, routerID, consts.RouterBgpNeighborPath),
		client: r.Client,
		jsonParser: func(body []byte) error {
			return json.Unmarshal(body, &Resp)
		},
	}
	err := routerAPI.do(ctx, req, nil)

	return Resp, err
}

func (r *RouterAPIService) GetSpecificRouterBgpNeighbor(
	ctx context.Context,
	routerID, bgpNeighborID int,
) (models.GetSpecificNetworkRouterBgpNeighbor, error) {
	Resp := models.GetSpecificNetworkRouterBgpNeighbor{}
	routerAPI := &api{
		compatibleVersion: "5.2.12",
		method:            "GET",
		path: fmt.Sprintf("%s/%s/%d/%s/%d", consts.NetworksPath,
			consts.NetworkRouterPath, routerID, consts.RouterBgpNeighborPath, bgpNeighborID),
		client: r.Client,
		jsonParser: func(body []byte) error {
			return json.Unmarshal(body, &Resp)
		},
	}
	err := routerAPI.do(ctx, nil, nil)

	return Resp, err
}

func (r *RouterAPIService) UpdateRouterBgpNeighbor(
	ctx context.Context,
	routerID, bgpNeighborID int,
	req models.CreateNetworkRouterBgpNeighborRequest,
) (models.SuccessOrErrorMessage, error) {
	Resp := models.SuccessOrErrorMessage{}
	routerAPI := &api{
		compatibleVersion: "5.2.12",
		method:            "PUT",
		path: fmt.Sprintf("%s/%s/%d/%s/%d", consts.NetworksPath,
			consts.NetworkRouterPath, routerID, consts.RouterBgpNeighborPath, bgpNeighborID),
		client: r.Client,
		jsonParser: func(body []byte) error {
			return json.Unmarshal(body, &Resp)
		},
	}
	err := routerAPI.do(ctx, req, nil)

	return Resp, err
}

func (r *RouterAPIService) DeleteRouterBgpNeighbor(
	ctx context.Context,
	routerID, bgpNeighborID int,
) (models.SuccessOrErrorMessage, error) {
	Resp := models.SuccessOrErrorMessage{}
	routerAPI := &api{
		compatibleVersion: "5.2.12",
		method:            "DELETE",
		path: fmt.Sprintf("%s/%s/%d/%s/%d", consts.NetworksPath,
			consts.NetworkRouterPath, routerID, consts.RouterBgpNeighborPath, bgpNeighborID),
		client: r.Client,
		jsonParser: func(body []byte) error {
			return json.Unmarshal(body, &Resp)
		},
	}
	err := routerAPI.do(ctx, nil, nil)

	return Resp, err
}

func (r *RouterAPIService) GetTransportZones(
	ctx context.Context,
	serviceID int,
	transportName string,
) (models.NetworkScope, error) {
	Resp := models.TransportZonesResp{}
	routerAPI := &api{
		compatibleVersion: "5.2.13",
		method:            "GET",
		path: fmt.Sprintf("%s/%s/%d/%s",
			consts.NetworksPath, consts.ServerPath, serviceID, consts.NetworkScopePath),
		client: r.Client,
		jsonParser: func(body []byte) error {
			return json.Unmarshal(body, &Resp)
		},
	}

	if err := routerAPI.do(ctx, nil, nil); err != nil {
		return models.NetworkScope{}, err
	}

	for _, t := range Resp.NetworkScopes {
		if t.Name == transportName {
			return t, nil
		}
	}

	return models.NetworkScope{}, nil
}

func (r *RouterAPIService) GetEdgeCluster(
	ctx context.Context,
	serviceID int,
	edgeClusterName string,
) (models.NetworkEdgeClusters, error) {
	Resp := models.NetworkEdgeClustersResp{}
	routerAPI := &api{
		compatibleVersion: "5.2.13",
		method:            "GET",
		path: fmt.Sprintf("%s/%s/%d/%s",
			consts.NetworksPath, consts.ServerPath, serviceID, consts.NetworkEdgeClusterPath),
		client: r.Client,
		jsonParser: func(body []byte) error {
			return json.Unmarshal(body, &Resp)
		},
	}

	if err := routerAPI.do(ctx, nil, nil); err != nil {
		return models.NetworkEdgeClusters{}, err
	}

	for _, t := range Resp.NetworkEdgeClusters {
		if t.Name == edgeClusterName {
			return t, nil
		}
	}

	return models.NetworkEdgeClusters{}, nil
}
