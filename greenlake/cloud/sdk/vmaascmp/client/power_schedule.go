// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package client

import (
	"context"
	"encoding/json"
	"fmt"

	consts "github.com/HPE/terraform-provider-hpe/greenlake/cloud/sdk/vmaascmp/common"
	models "github.com/HPE/terraform-provider-hpe/greenlake/cloud/sdk/vmaascmp/models"
)

type PowerSchedulesAPIService struct {
	Client APIClientHandler
	Cfg    Configuration
}

/*
VirtualImageApiService
Get All Virtual images
  - @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc.
    Passed from http.Request or context.Background().
  - @param serviceInstanceID
  - @param name/phrase optional

@return models.VirtualImages
*/
func (a *PowerSchedulesAPIService) GetAllPowerSchedules(ctx context.Context,
	param map[string]string) (models.GetAllPowerSchedules, error) {
	Response := models.GetAllPowerSchedules{}

	allPowerScheduleAPI := &api{
		method: "GET",
		path:   consts.PowerSchedulePath,
		client: a.Client,

		jsonParser: func(body []byte) error {
			return json.Unmarshal(body, &Response)
		},
	}
	err := allPowerScheduleAPI.do(ctx, nil, param)

	return Response, err
}

func (a *PowerSchedulesAPIService) GetSpecificPowerSchedule(
	ctx context.Context,
	powerID int,
) (models.GetSpecificPowerSchedule, error) {
	Response := models.GetSpecificPowerSchedule{}

	allPowerScheduleAPI := &api{
		method: "GET",
		path:   fmt.Sprintf("%s/%d", consts.PowerSchedulePath, powerID),
		client: a.Client,

		jsonParser: func(body []byte) error {
			return json.Unmarshal(body, &Response)
		},
	}
	err := allPowerScheduleAPI.do(ctx, nil, nil)

	return Response, err
}
