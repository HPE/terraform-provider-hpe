// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package client

import (
	"context"
	"encoding/json"
	"fmt"

	consts "github.com/HPE/terraform-provider-hpe/greenlake/cloud/sdk/vmaascmp/common"
	models "github.com/HPE/terraform-provider-hpe/greenlake/cloud/sdk/vmaascmp/models"
)

type VirtualImagesAPIService struct {
	Client APIClientHandler
	Cfg    Configuration
}

/*
VirtualImageAPIService
Get All Virtual images
  - @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc.
    Passed from http.Request or context.Background().
  - @param serviceInstanceID
  - @param name/phrase optional

@return models.VirtualImages
*/
func (a *VirtualImagesAPIService) GetAllVirtualImages(ctx context.Context,
	param map[string]string) (models.VirtualImages, error) {
	Response := models.VirtualImages{}

	allVirtualImagesAPI := &api{
		method: "GET",
		path:   consts.VirtualImagePath,
		client: a.Client,

		jsonParser: func(body []byte) error {
			return json.Unmarshal(body, &Response)
		},
	}

	err := allVirtualImagesAPI.do(ctx, nil, param)

	return Response, err
}

func (a *VirtualImagesAPIService) GetSpecificVirtualImage(
	ctx context.Context,
	virtualImageID int,
) (models.GetSpecificVirtualImage, error) {
	Response := models.GetSpecificVirtualImage{}

	api := &api{
		method: "GET",
		path:   fmt.Sprintf("%s/%d", consts.VirtualImagePath, virtualImageID),
		client: a.Client,

		jsonParser: func(body []byte) error {
			return json.Unmarshal(body, &Response)
		},
	}

	err := api.do(ctx, nil, nil)

	return Response, err
}
