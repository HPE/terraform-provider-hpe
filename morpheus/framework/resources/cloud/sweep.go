// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package cloud

import (
	"context"
	"net/http"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
)

const testCloudPrefix = "TestAccMorpheusCloud"

func init() {
	testhelpers.RegisterAPISweeper(
		"hpe_morpheus_cloud",
		testCloudPrefix,
		func(ctx context.Context, client *sdk.APIClient, _ string) (any, *http.Response, error) {
			return client.CloudsAPI.ListClouds(ctx).Execute()
		},
		"Zones",
		func(ctx context.Context, client *sdk.APIClient, id int64, _ any) (*http.Response, error) {
			_, hresp, err := client.CloudsAPI.RemoveClouds(ctx, id).Force(true).Execute()
			return hresp, err
		},
	)
}
