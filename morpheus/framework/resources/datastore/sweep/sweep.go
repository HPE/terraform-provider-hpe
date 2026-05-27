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

const sweeperName = "hpe_morpheus_datastore"

func init() {
	testsweep.RegisterTypedAPISweeper(
		sweeperName,
		// List datastore resources.
		func(ctx context.Context, client *sdk.APIClient) (
			[]sdk.ListDatastores200ResponseAllOfDatastoresInner,
			*http.Response,
			error,
		) {
			resp, hresp, err := client.DatastoresAPI.ListDatastores(ctx).Execute()
			if resp == nil {
				return nil, hresp, err
			}

			return resp.GetDatastores(), hresp, err
		},
		// Is this a test datastore?
		func(item sdk.ListDatastores200ResponseAllOfDatastoresInner) bool {
			name, ok := item.GetNameOk()
			if !ok || name == nil {
				return false
			}

			return strings.HasPrefix(*name, testsweep.TestResourcePrefix)
		},
		// Delete the test datastore.
		func(
			ctx context.Context,
			client *sdk.APIClient,
			item sdk.ListDatastores200ResponseAllOfDatastoresInner,
		) (*http.Response, error) {
			id, ok := item.GetIdOk()
			if !ok || id == nil {
				return nil, fmt.Errorf("could not get ID")
			}

			_, hresp, err := client.DatastoresAPI.DeleteDatastores(ctx, *id).Execute()

			return hresp, err
		},
		testsweep.WithIgnoreListStatuses[sdk.ListDatastores200ResponseAllOfDatastoresInner](
			http.StatusNotFound,
			http.StatusForbidden,
		),
	)
}
