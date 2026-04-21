// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package datastore

import (
	"context"
	"net/http"

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
		testDatastorePrefix,
		func(ctx context.Context, client *sdk.APIClient, prefix string) ([]sdk.ListDatastores200ResponseAllOfDatastoresInner, *http.Response, error) {
			resp, hresp, err := client.DatastoresAPI.ListDatastores(ctx).Phrase(prefix).Execute()
			if resp == nil {
				return nil, hresp, err
			}

			return resp.GetDatastores(), hresp, err
		},
		func(item sdk.ListDatastores200ResponseAllOfDatastoresInner) (string, bool) {
			name, ok := item.GetNameOk()
			if !ok || name == nil {
				return "", false
			}

			return *name, true
		},
		func(item sdk.ListDatastores200ResponseAllOfDatastoresInner) (int64, bool) {
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
			_ sdk.ListDatastores200ResponseAllOfDatastoresInner,
		) (*http.Response, error) {
			_, hresp, err := client.DatastoresAPI.DeleteDatastores(ctx, id).Execute()
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
