// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

//go:build sweep

package sweep

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"

	testsweep "github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/sweep"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/getsafe"
)

const sweeperName = "hpe_morpheus_cloud"

func init() {
	testsweep.RegisterTypedAPISweeper(
		sweeperName,
		// List cloud resources.
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
			name, ok := getsafe.GetOk(item.Name)
			if !ok || name == nil {
				return false
			}

			return strings.HasPrefix(*name, testsweep.TestResourcePrefix)
		},
		// Delete the test cloud.
		func(
			ctx context.Context,
			client *sdk.APIClient,
			item sdk.ListClouds200ResponseAllOfZonesInner,
		) (*http.Response, error) {
			id, ok := getsafe.GetOk(item.Id)
			if !ok || id == nil {
				return nil, fmt.Errorf("could not get ID")
			}

			_, hresp, err := client.CloudsAPI.RemoveClouds(ctx, *id).Force(true).Execute()

			return hresp, err
		},
	)
}
