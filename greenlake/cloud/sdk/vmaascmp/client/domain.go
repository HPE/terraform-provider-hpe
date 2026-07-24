//  (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package client

import (
	"context"
	"encoding/json"
	"fmt"

	consts "github.com/HPE/terraform-provider-hpe/greenlake/cloud/sdk/vmaascmp/common"
	"github.com/HPE/terraform-provider-hpe/greenlake/cloud/sdk/vmaascmp/models"
)

type DomainAPIService struct {
	Client APIClientHandler
	Cfg    Configuration
}

func (d *DomainAPIService) GetAllDomains(
	ctx context.Context,
	param map[string]string,
) (models.GetAllDomains, error) {
	var DomainResp models.GetAllDomains
	networkAPI := &api{
		method: "GET",
		path: fmt.Sprintf("%s/%s",
			consts.NetworksPath, consts.DomainPath),
		client: d.Client,

		jsonParser: func(body []byte) error {
			return json.Unmarshal(body, &DomainResp)
		},
	}
	err := networkAPI.do(ctx, nil, param)

	return DomainResp, err
}

func (d *DomainAPIService) GetSpecificDomain(
	ctx context.Context,
	domainID int,
) (models.GetSpecificDomain, error) {
	var DomainResp models.GetSpecificDomain
	networkAPI := &api{
		method: "GET",
		path: fmt.Sprintf("%s/%s/%d",
			consts.NetworksPath, consts.DomainPath, domainID),
		client: d.Client,

		jsonParser: func(body []byte) error {
			return json.Unmarshal(body, &DomainResp)
		},
	}
	err := networkAPI.do(ctx, nil, nil)

	return DomainResp, err
}
