// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

//go:build sweep

package sweep

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"strings"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"

	testsweep "github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/sweep"
)

const sweeperName = "hpe_morpheus_storage_volume"

func init() {
	testsweep.RegisterTypedAPISweeper(
		sweeperName,
		// list
		func(ctx context.Context, client *sdk.APIClient) ([]testsweep.SearchHit, *http.Response, error) {
			return testsweep.SearchHits(ctx, client, testsweep.TestResourcePrefix)
		},
		// is this
		func(item testsweep.SearchHit) bool {
			return strings.HasPrefix(item.Name, testsweep.TestResourcePrefix)
		},
		// delete
		func(ctx context.Context, client *sdk.APIClient, item testsweep.SearchHit) (*http.Response, error) {
			if item.ID == "" {
				return &http.Response{StatusCode: http.StatusOK}, nil
			}

			id, err := strconv.ParseInt(item.ID, 10, 64)
			if err != nil {
				log.Printf("[ERROR] Failed to parse storage volume ID %q: %v", item.ID, err)

				return &http.Response{StatusCode: http.StatusOK}, nil
			}

			idParam := sdk.UpdateStorageVolumesIdParameter{Int64: &id}

			_, hresp, delErr := client.StorageAPI.RemoveStorageVolumes(ctx, idParam).Execute()
			if delErr != nil && hresp != nil && hresp.StatusCode == http.StatusNotFound {
				return &http.Response{StatusCode: http.StatusOK}, nil
			}

			return hresp, delErr
		},
		testsweep.WithIgnoreListStatuses[testsweep.SearchHit](
			http.StatusNotFound,
			http.StatusForbidden,
		),
	)
}
