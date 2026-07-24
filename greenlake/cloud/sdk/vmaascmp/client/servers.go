// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package client

import (
	"context"
	"encoding/json"
	"fmt"

	consts "github.com/HPE/terraform-provider-hpe/greenlake/cloud/sdk/vmaascmp/common"
	models "github.com/HPE/terraform-provider-hpe/greenlake/cloud/sdk/vmaascmp/models"
)

type ServersAPIService struct {
	Client APIClientHandler
	Cfg    Configuration
}

/*
ServersApiService
Get All Servers
  - @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc.
    Passed from http.Request or context.Background().
  - @param name/phrase optional

@return models.Servers
*/
func (a *ServersAPIService) GetAllServers(
	ctx context.Context,
	queryParams map[string]string,
) (models.ServersResponse, error) {
	ServersResponse := models.ServersResponse{}
	serverAPI := &api{
		method: "GET",
		path:   consts.ServerPath,
		client: a.Client,
		jsonParser: func(body []byte) error {
			return json.Unmarshal(body, &ServersResponse)
		},
	}
	err := serverAPI.do(ctx, nil, queryParams)

	return ServersResponse, err
}

func (a *ServersAPIService) GetSpecificServer(
	ctx context.Context,
	serverID int) (models.GetSpecificServerResponse, error) {
	ServersResponse := models.GetSpecificServerResponse{}
	serverAPI := &api{
		method: "GET",
		path: fmt.Sprintf("%s/%d",
			consts.ServerPath, serverID),
		client: a.Client,
		jsonParser: func(body []byte) error {
			return json.Unmarshal(body, &ServersResponse)
		},
		validations: []validationFunc{
			func() error {
				if serverID < 1 {
					return fmt.Errorf("%s", "server id should be greater than or equal to 1")
				}

				return nil
			},
		},
	}
	err := serverAPI.do(ctx, nil, nil)

	return ServersResponse, err
}
