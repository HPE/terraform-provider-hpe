// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package cloud

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"

	sweep "github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/sweep"
)

const testCloudPrefix = "TestAccMorpheusCloud"

func init() {
	sweep.RegisterTypedAPISweeper(
		"hpe_morpheus_cloud",
		// List all cloud resources.
		func(ctx context.Context, client *sdk.APIClient) (
			[]sdk.ListClouds200ResponseAllOfZonesInner,
			*http.Response,
			error,
		) {
			resp, hresp, err := client.CloudsAPI.ListClouds(ctx).Execute()
			if resp == nil {
				return nil, hresp, err
			}

			return resp.Zones, hresp, err
		},
		// Is this a test cloud?
		func(item sdk.ListClouds200ResponseAllOfZonesInner) bool {
			name, ok := item.GetNameOk()
			if !ok || name == nil {
				return false
			}

			return strings.HasPrefix(*name, testCloudPrefix)
		},
		// Delete the test cloud.
		func(
			ctx context.Context,
			client *sdk.APIClient,
			item sdk.ListClouds200ResponseAllOfZonesInner,
		) (*http.Response, error) {
			id, ok := item.GetIdOk()
			if !ok || id == nil {
				return nil, fmt.Errorf("could not get ID")
			}

			_, hresp, err := client.CloudsAPI.RemoveClouds(ctx, *id).Force(true).Execute()

			return hresp, err
		},
	)
}
