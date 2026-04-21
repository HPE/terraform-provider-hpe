// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package datastore

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
)

// Datastores whose name begins with this string will be eligible for deletion
const (
	testDatastorePrefix   = "TestAccMorpheusDatastore"
	enableDatastoreDelete = false
)

func init() {
	testhelpers.RegisterTypedAPISweeper(
		"hpe_morpheus_datastore",
		// List all datastore resources.
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
			if !enableDatastoreDelete {
				return false
			}

			name, ok := item.GetNameOk()
			if !ok || name == nil || !strings.HasPrefix(*name, testDatastorePrefix) {
				return false
			}

			return true
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
		testhelpers.WithIgnoreListStatuses[sdk.ListDatastores200ResponseAllOfDatastoresInner](
			http.StatusNotFound,
			http.StatusForbidden,
		),
		testhelpers.WithFilter(func(
			_ context.Context,
			_ *sdk.APIClient,
			_ sdk.ListDatastores200ResponseAllOfDatastoresInner,
		) (bool, string, error) {
			if enableDatastoreDelete {
				return true, "", nil
			}

			return false, "delete disabled", nil
		}),
	)
}
