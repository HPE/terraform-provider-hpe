// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package client

import (
	"context"
	"encoding/json"

	consts "github.com/HPE/terraform-provider-hpe/greenlake/cloud/sdk/vmaascmp/common"
	"github.com/HPE/terraform-provider-hpe/greenlake/cloud/sdk/vmaascmp/models"
)

type CmpStatus struct {
	Client APIClientHandler
	Cfg    Configuration
}

func (a *CmpStatus) GetCmpVersion(ctx context.Context) (models.CmpVersionModel, error) {
	CheckResp := models.CmpVersionModel{}

	statusAPI := &api{
		method: "GET",
		path:   consts.WhoamiPath,
		client: a.Client,

		jsonParser: func(body []byte) error {
			return json.Unmarshal(body, &CheckResp)
		},
	}
	err := statusAPI.do(ctx, nil, nil)

	return CheckResp, err
}
