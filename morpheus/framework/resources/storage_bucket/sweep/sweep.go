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
)

const sweeperName = "hpe_morpheus_storage_bucket"

func init() {
	testsweep.RegisterTypedAPISweeper(
		sweeperName,
		// List storage bucket resources.
		func(ctx context.Context, client *sdk.APIClient) (
			[]sdk.ListStorageBuckets200ResponseAllOfStorageBucketsInner,
			*http.Response,
			error,
		) {
			resp, hresp, err := client.StorageAPI.ListStorageBuckets(ctx).Execute()
			if resp == nil {
				return nil, hresp, err
			}

			return resp.GetStorageBuckets(), hresp, err
		},
		// Is this a test storage bucket?
		func(item sdk.ListStorageBuckets200ResponseAllOfStorageBucketsInner) bool {
			name, ok := item.GetNameOk()
			if !ok || name == nil {
				return false
			}

			return strings.HasPrefix(*name, testsweep.TestResourcePrefix)
		},
		// Delete the test storage bucket.
		func(
			ctx context.Context,
			client *sdk.APIClient,
			item sdk.ListStorageBuckets200ResponseAllOfStorageBucketsInner,
		) (*http.Response, error) {
			id, ok := item.GetIdOk()
			if !ok || id == nil {
				return nil, fmt.Errorf("could not get ID")
			}

			_, hresp, err := client.StorageAPI.RemoveStorageBuckets(ctx, *id).Execute()

			return hresp, err
		},
		testsweep.WithIgnoreListStatuses[sdk.ListStorageBuckets200ResponseAllOfStorageBucketsInner](
			http.StatusNotFound,
			http.StatusForbidden,
		),
	)
}
