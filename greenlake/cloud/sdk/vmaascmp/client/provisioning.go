// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package client

import (
	"context"
	"encoding/json"

	consts "github.com/HPE/terraform-provider-hpe/greenlake/cloud/sdk/vmaascmp/common"
	models "github.com/HPE/terraform-provider-hpe/greenlake/cloud/sdk/vmaascmp/models"
)

type ProvisioningAPIService struct {
	Client APIClientHandler
	Cfg    Configuration
}

func (a *ProvisioningAPIService) GetAllProvisioningTypes(ctx context.Context,
	param map[string]string) (models.GetAllProvisioningTypes, error) {
	Response := models.GetAllProvisioningTypes{}

	allProvisionAPI := &api{
		method: "GET",
		path:   consts.ProvisionTypesPath,
		client: a.Client,

		jsonParser: func(body []byte) error {
			return json.Unmarshal(body, &Response)
		},
	}

	err := allProvisionAPI.do(ctx, nil, param)

	return Response, err
}
