// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package client

import (
	"context"
	"encoding/json"
	"fmt"

	consts "github.com/HPE/terraform-provider-hpe/greenlake/cloud/sdk/vmaascmp/common"
	"github.com/HPE/terraform-provider-hpe/greenlake/cloud/sdk/vmaascmp/models"
)

type LibraryAPIService struct {
	Client APIClientHandler
	Cfg    Configuration
}

/*
LibrariesAPIService
Get All layouts
  - @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc.
    Passed from http.Request or context.Background().
  - @param serviceInstanceID
*/
func (a *LibraryAPIService) GetAllLayouts(ctx context.Context,
	param map[string]string) (models.LayoutsResp, error) {
	Response := models.LayoutsResp{}

	allLayoutsAPI := &api{
		method: "GET",
		path:   consts.LibraryLayoutPath,
		client: a.Client,

		jsonParser: func(body []byte) error {
			return json.Unmarshal(body, &Response)
		},
	}
	err := allLayoutsAPI.do(ctx, nil, param)

	return Response, err
}

func (a *LibraryAPIService) GetAllInstanceTypes(ctx context.Context,
	param map[string]string) (models.InstanceTypesResp, error) {
	Response := models.InstanceTypesResp{}

	allInstTypeAPI := &api{
		method: "GET",
		path:   consts.LibraryInstanceTypesPath,
		client: a.Client,

		jsonParser: func(body []byte) error {
			return json.Unmarshal(body, &Response)
		},
	}
	err := allInstTypeAPI.do(ctx, nil, param)

	return Response, err
}

func (a *LibraryAPIService) GetSpecificLayout(
	ctx context.Context,
	layoutID int,
) (models.GetSpecificLayout, error) {
	Response := models.GetSpecificLayout{}

	allLayoutsAPI := &api{
		method: "GET",
		path:   fmt.Sprintf("%s/%d", consts.LibraryLayoutPath, layoutID),
		client: a.Client,

		jsonParser: func(body []byte) error {
			return json.Unmarshal(body, &Response)
		},
	}
	err := allLayoutsAPI.do(ctx, nil, nil)

	return Response, err
}
