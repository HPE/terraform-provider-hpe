// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package client

import (
	"context"
	"encoding/json"
	"fmt"

	consts "github.com/HPE/terraform-provider-hpe/greenlake/cloud/sdk/vmaascmp/common"
	models "github.com/HPE/terraform-provider-hpe/greenlake/cloud/sdk/vmaascmp/models"
)

type EnvironmentAPIService struct {
	Client APIClientHandler
	Cfg    Configuration
}

func (e *EnvironmentAPIService) GetAllEnvironment(ctx context.Context,
	param map[string]string) (models.GetAllEnvironment, error) {
	Response := models.GetAllEnvironment{}

	allEnvAPI := &api{
		method: "GET",
		path:   consts.EnvironmentPath,
		client: e.Client,

		jsonParser: func(body []byte) error {
			return json.Unmarshal(body, &Response)
		},
	}

	err := allEnvAPI.do(ctx, nil, param)

	return Response, err
}

func (e *EnvironmentAPIService) GetSpecificEnvironment(
	ctx context.Context,
	envID int,
) (models.GetSpecificEnvironment, error) {
	Response := models.GetSpecificEnvironment{}

	allEnvAPI := &api{
		method: "GET",
		path:   fmt.Sprintf("%s/%d", consts.EnvironmentPath, envID),
		client: e.Client,
		jsonParser: func(body []byte) error {
			return json.Unmarshal(body, &Response)
		},
	}

	err := allEnvAPI.do(ctx, nil, nil)

	return Response, err
}
