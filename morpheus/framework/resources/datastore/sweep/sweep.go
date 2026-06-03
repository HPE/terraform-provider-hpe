// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

//go:build sweep

package sweep

import (
	"context"
	"net/http"
	"strings"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"

	testsweep "github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/sweep"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/getsafe"
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

			return getsafe.GetSafe(&resp.Datastores), hresp, err
		},
		// Is this a test datastore?
		func(item sdk.ListDatastores200ResponseAllOfDatastoresInner) bool {
			return strings.HasPrefix(item.Name, testsweep.TestResourcePrefix)
		},
		// Delete the test datastore.
		func(
			ctx context.Context,
			client *sdk.APIClient,
			item sdk.ListDatastores200ResponseAllOfDatastoresInner,
		) (*http.Response, error) {
			_, hresp, err := client.DatastoresAPI.DeleteDatastores(ctx, item.Id).Execute()

			return hresp, err
		},
		testsweep.WithIgnoreListStatuses[sdk.ListDatastores200ResponseAllOfDatastoresInner](
			http.StatusNotFound,
			http.StatusForbidden,
		),
	)
}
