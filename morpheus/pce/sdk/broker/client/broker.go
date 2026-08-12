// (C) Copyright 2024-2025 Hewlett Packard Enterprise Development LP

package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	consts "github.com/HPE/terraform-provider-hpe/morpheus/pce/sdk/broker/common"
	"github.com/HPE/terraform-provider-hpe/morpheus/pce/sdk/broker/models"
)

// BrokerAPIService is a service that provides methods to interact with the broker API
type BrokerAPIService struct {
	Client APIClientHandler
	Cfg    Configuration
}

// GetMorpheusDetails returns Morpheus details to terraform
func (a *BrokerAPIService) GetMorpheusDetails(ctx context.Context) (models.TFMorpheusDetails, error) {
	// Get the Morpheus token
	CMPDetails := models.CMPDetails{}
	cmpDetailsAPI := &api{
		method:                 http.MethodGet,
		path:                   consts.CMPDetails,
		client:                 a.Client,
		removeVmaasCMPBasePath: true,
		jsonParser: func(body []byte) error {
			return json.Unmarshal(body, &CMPDetails)
		},
	}
	if err := cmpDetailsAPI.do(ctx, nil, a.Cfg.DefaultQueryParams); err != nil {
		return models.TFMorpheusDetails{}, fmt.Errorf("error getting Morpheus connection details: %v", err)
	}
	// build response
	ret := models.TFMorpheusDetails{
		ID:          CMPDetails.ServiceInstanceID,
		AccessToken: CMPDetails.TokenDetails.AccessToken,
		ValidTill:   CMPDetails.TokenDetails.ExpiresIn,
		URL:         strings.TrimSuffix(CMPDetails.URL, "/"),
	}

	return ret, nil
}
