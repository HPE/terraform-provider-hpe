// (C) Copyright 2022 Hewlett Packard Enterprise Development LP

package client

import (
	"context"
	"encoding/json"
	"fmt"

	consts "github.com/HPE/terraform-provider-hpe/greenlake/cloud/sdk/vmaascmp/common"
	"github.com/HPE/terraform-provider-hpe/greenlake/cloud/sdk/vmaascmp/models"
)

const dhcpCompatibleVersion = "5.2.13"

type DhcpServerAPIService struct {
	Client APIClientHandler
	Cfg    Configuration
}

func (lb *DhcpServerAPIService) CreateDhcpServer(
	ctx context.Context,
	serverID int,
	request models.CreateNetworkDhcpServerRequest,
) (models.CreateNetworkDhcpServerResp, error) {
	DhcpServerResp := models.CreateNetworkDhcpServerResp{}
	dhcpServerAPI := &api{
		compatibleVersion: dhcpCompatibleVersion,
		method:            "POST",
		path: fmt.Sprintf("%s/%s/%d/%s",
			consts.NetworksPath, consts.ServerPath,
			serverID, consts.DhcpServerPath),
		client: lb.Client,
		jsonParser: func(body []byte) error {
			return json.Unmarshal(body, &DhcpServerResp)
		},
	}
	err := dhcpServerAPI.do(ctx, request, nil)

	return DhcpServerResp, err
}

func (lb *DhcpServerAPIService) UpdateDhcpServer(
	ctx context.Context,
	serverID int,
	dhcpID int,
	request models.CreateNetworkDhcpServerRequest,
) (models.CreateNetworkDhcpServerResp, error) {
	DhcpServerResp := models.CreateNetworkDhcpServerResp{}
	dhcpServerAPI := &api{
		compatibleVersion: dhcpCompatibleVersion,
		method:            "PUT",
		path: fmt.Sprintf("%s/%s/%d/%s/%d",
			consts.NetworksPath, consts.ServerPath,
			serverID, consts.DhcpServerPath, dhcpID),
		client: lb.Client,
		jsonParser: func(body []byte) error {
			return json.Unmarshal(body, &DhcpServerResp)
		},
	}
	err := dhcpServerAPI.do(ctx, request, nil)

	return DhcpServerResp, err
}

func (lb *DhcpServerAPIService) DeleteDhcpServer(
	ctx context.Context,
	serverID int,
	dhcpID int,
) (models.SuccessOrErrorMessage, error) {
	DhcpServerResp := models.SuccessOrErrorMessage{}
	dhcpServerAPI := &api{
		compatibleVersion: dhcpCompatibleVersion,
		method:            "DELETE",
		path: fmt.Sprintf("%s/%s/%d/%s/%d",
			consts.NetworksPath, consts.ServerPath,
			serverID, consts.DhcpServerPath, dhcpID),
		client: lb.Client,
		jsonParser: func(body []byte) error {
			return json.Unmarshal(body, &DhcpServerResp)
		},
	}
	err := dhcpServerAPI.do(ctx, nil, nil)

	return DhcpServerResp, err
}

func (lb *DhcpServerAPIService) GetDhcpServers(
	ctx context.Context,
	serverID int,
) (models.GetNetworkDhcpServers, error) {
	DhcpServerResp := models.GetNetworkDhcpServers{}
	dhcpServerAPI := &api{
		compatibleVersion: dhcpCompatibleVersion,
		method:            "GET",
		path: fmt.Sprintf("%s/%s/%d/%s",
			consts.NetworksPath, consts.ServerPath,
			serverID, consts.DhcpServerPath),
		client: lb.Client,
		jsonParser: func(body []byte) error {
			return json.Unmarshal(body, &DhcpServerResp)
		},
	}
	err := dhcpServerAPI.do(ctx, nil, nil)

	return DhcpServerResp, err
}

func (lb *DhcpServerAPIService) GetSpecificDhcpServer(
	ctx context.Context,
	serverID int,
	dhcpID int,
) (models.GetSpecificNetworkDhcpServer, error) {
	DhcpServerResp := models.GetSpecificNetworkDhcpServer{}
	dhcpServerAPI := &api{
		compatibleVersion: dhcpCompatibleVersion,
		method:            "GET",
		path: fmt.Sprintf("%s/%s/%d/%s/%d",
			consts.NetworksPath, consts.ServerPath,
			serverID, consts.DhcpServerPath, dhcpID),
		client: lb.Client,
		jsonParser: func(body []byte) error {
			return json.Unmarshal(body, &DhcpServerResp)
		},
	}
	err := dhcpServerAPI.do(ctx, nil, nil)

	return DhcpServerResp, err
}
