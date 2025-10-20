// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package sweep

import (
	"context"
	"errors"
	"os"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/sdk"

	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/clientfactory"
)

func NewSweepClient(ctx context.Context) (*sdk.APIClient, error) {
	var username, password string

	url, ok := os.LookupEnv("TF_VAR_testacc_morpheus_url")
	if !ok {
		return nil, errors.New("TF_VAR_testacc_morpheus_url not set")
	}

	token, ok := os.LookupEnv("TF_VAR_testacc_morpheus_access_token")
	if !ok {
		username, ok = os.LookupEnv("TF_VAR_testacc_morpheus_username")
		if !ok {
			return nil, errors.New(
				"one of TF_VAR_testacc_morpheus_access_token or " +
					"TF_VAR_testacc_morpheus_username must be set",
			)
		}

		password, ok = os.LookupEnv("TF_VAR_testacc_morpheus_password")
		if !ok {
			return nil, errors.New(
				"one of TF_VAR_testacc_morpheus_access_token or " +
					"TF_VAR_testacc_morpheus_password must be set",
			)
		}
	}

	// If set to any value, use insecure
	_, insecure := os.LookupEnv("TF_VAR_testacc_morpheus_insecure")
	var opts []clientfactory.ClientOption
	if insecure {
		opts = append(opts, clientfactory.WithInsecureTLS())
	}

	client := clientfactory.NewAPIClient(
		ctx,
		url,
		username,
		password,
		token,
		opts...,
	)

	return client, nil
}
