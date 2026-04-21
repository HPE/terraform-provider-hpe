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
	testhelpers.RegisterTypedAPISweeper(
		"hpe_morpheus_cloud",
		testCloudPrefix,
		func(ctx context.Context, client *sdk.APIClient, _ string) ([]sdk.ListClouds200ResponseAllOfZonesInner, *http.Response, error) {
			resp, hresp, err := client.CloudsAPI.ListClouds(ctx).Execute()
			if resp == nil {
				return nil, hresp, err
			}

			return resp.Zones, hresp, err
		},
		func(item sdk.ListClouds200ResponseAllOfZonesInner) (string, bool) {
			name, ok := item.GetNameOk()
			if !ok || name == nil {
				return "", false
			}

			return *name, true
		},
		func(item sdk.ListClouds200ResponseAllOfZonesInner) (int64, bool) {
			id, ok := item.GetIdOk()
			if !ok || id == nil {
				return 0, false
			}

			return *id, true
		},
		func(
			ctx context.Context,
			client *sdk.APIClient,
			id int64,
			_ sdk.ListClouds200ResponseAllOfZonesInner,
		) (*http.Response, error) {
			_, hresp, err := client.CloudsAPI.RemoveClouds(ctx, id).Force(true).Execute()
			return hresp, err
		},
	)
}
