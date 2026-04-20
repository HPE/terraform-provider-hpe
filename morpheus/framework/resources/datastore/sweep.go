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
	testhelpers.RegisterAPISweeper(
		"hpe_morpheus_datastore",
		testDatastorePrefix,
		func(ctx context.Context, client *sdk.APIClient, prefix string) (any, *http.Response, error) {
			return client.DatastoresAPI.ListDatastores(ctx).Phrase(prefix).Execute()
		},
		"GetDatastores",
		func(ctx context.Context, client *sdk.APIClient, id int64, _ any) (*http.Response, error) {
			_, hresp, err := client.DatastoresAPI.DeleteDatastores(ctx, id).Execute()
			return hresp, err
		},
		testhelpers.WithIgnoreListStatuses(http.StatusNotFound, http.StatusForbidden),
		testhelpers.WithFilter(func(_ context.Context, _ *sdk.APIClient, _ any) (bool, string, error) {
			if enableDatastoreDelete {
				return true, "", nil
			}

			return false, "delete disabled", nil
		}),
	)
}
