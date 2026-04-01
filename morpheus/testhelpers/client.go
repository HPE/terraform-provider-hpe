// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package testhelpers

import (
	"context"
	"os"
	"testing"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/clientfactory"
)

func newClient(ctx context.Context, t *testing.T) *sdk.APIClient {
	t.Helper()

	return clientfactory.NewAPIClient(
		ctx,
		os.Getenv("TF_VAR_testacc_morpheus_url"),
		os.Getenv("TF_VAR_testacc_morpheus_username"),
		os.Getenv("TF_VAR_testacc_morpheus_password"),
		os.Getenv("TF_VAR_testacc_morpheus_subdomain"),
		os.Getenv("TF_VAR_testacc_morpheus_access_token"),
		clientfactory.WithInsecureTLS())
}
