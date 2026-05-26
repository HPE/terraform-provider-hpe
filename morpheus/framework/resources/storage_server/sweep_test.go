// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package storage_server_test

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"

	testsweep "github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/sweep"
)

const sweeperName = "hpe_morpheus_storage_server"

func init() {
	testsweep.RegisterTypedAPISweeper(
		sweeperName,
		// List storage server resources.
		func(ctx context.Context, client *sdk.APIClient) (
			[]sdk.ListStorageServers200ResponseAllOfStorageServersInner,
			*http.Response,
			error,
		) {
			resp, hresp, err := client.StorageAPI.ListStorageServers(ctx).Execute()
			if resp == nil {
				return nil, hresp, err
			}

			return resp.GetStorageServers(), hresp, err
		},
		// Is this a test storage server?
		func(item sdk.ListStorageServers200ResponseAllOfStorageServersInner) bool {
			name, ok := item.GetNameOk()
			if !ok || name == nil {
				return false
			}

			return strings.HasPrefix(*name, testsweep.TestResourcePrefix)
		},
		// Delete the test storage server.
		func(
			ctx context.Context,
			client *sdk.APIClient,
			item sdk.ListStorageServers200ResponseAllOfStorageServersInner,
		) (*http.Response, error) {
			id, ok := item.GetIdOk()
			if !ok || id == nil {
				return nil, fmt.Errorf("could not get ID")
			}

			_, hresp, err := client.StorageAPI.RemoveStorageServers(ctx, *id).Execute()

			return hresp, err
		},
		testsweep.WithIgnoreListStatuses[sdk.ListStorageServers200ResponseAllOfStorageServersInner](
			http.StatusNotFound,
			http.StatusForbidden,
		),
	)
}
