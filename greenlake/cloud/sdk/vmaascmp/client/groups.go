// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package client

import (
	"context"
	"encoding/json"
	"fmt"

	consts "github.com/HPE/terraform-provider-hpe/greenlake/cloud/sdk/vmaascmp/common"
	"github.com/HPE/terraform-provider-hpe/greenlake/cloud/sdk/vmaascmp/models"
)

type GroupsAPIService struct {
	Client APIClientHandler
	Cfg    Configuration
}

/*
GroupsAPIService
Get a Specific Group
  - @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc.
    Passed from http.Request or context.Background().
  - @param serviceInstanceID
  - @param groupID The Group ID
*/
func (a *GroupsAPIService) GetASpecificGroup(ctx context.Context,
	groupID int) (models.GroupResp, error) {
	SpecificGrpResp := models.GroupResp{}

	specificGrpRespAPI := &api{
		method: "GET",
		path: fmt.Sprintf("%s/%d",
			consts.GroupsPath, groupID),
		client: a.Client,

		jsonParser: func(body []byte) error {
			return json.Unmarshal(body, &SpecificGrpResp)
		},
		validations: []validationFunc{
			func() error {
				if groupID < 1 {
					return fmt.Errorf("%s", "group id should be greater than or equal to 1")
				}

				return nil
			},
		},
	}
	err := specificGrpRespAPI.do(ctx, nil, nil)

	return SpecificGrpResp, err
}

/*
GroupsAPIService
This endpoint retrieves all groups and a list of zones associated with the group by id.
  - @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc.
    Passed from http.Request or context.Background().
  - @param serviceInstanceID
*/
func (a *GroupsAPIService) GetAllGroups(ctx context.Context,
	queryParams map[string]string) (models.Groups, error) {
	AllGrpResp := models.Groups{}

	allGrpRespAPI := &api{
		method: "GET",
		path:   consts.GroupsPath,
		client: a.Client,

		jsonParser: func(body []byte) error {
			return json.Unmarshal(body, &AllGrpResp)
		},
	}
	err := allGrpRespAPI.do(ctx, nil, queryParams)

	return AllGrpResp, err
}
