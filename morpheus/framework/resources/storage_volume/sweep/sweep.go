// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

//go:build sweep

package sweep

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"

	testsweep "github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/sweep"
)

const sweeperName = "hpe_morpheus_storage_volume"

func init() {
	testsweep.RegisterTypedAPISweeper(
		sweeperName,
		// list
		func(ctx context.Context, client *sdk.APIClient) ([]sdk.Search200ResponseHitsInner, *http.Response, error) {
			resp, hresp, err := client.SearchAPI.Search(ctx).
				Phrase(testsweep.TestResourcePrefix).
				Max(1000).
				Execute()
			if resp == nil {
				return nil, hresp, err
			}

			return resp.Hits, hresp, err
		},
		// is this
		func(item sdk.Search200ResponseHitsInner) bool {
			if item.Name == nil {
				return false
			}

			return strings.HasPrefix(*item.Name, testsweep.TestResourcePrefix)
		},
		// delete
		func(ctx context.Context, client *sdk.APIClient, item sdk.Search200ResponseHitsInner) (*http.Response, error) {
			if item.Id == nil {
				return &http.Response{StatusCode: http.StatusOK}, nil
			}

			id, err := strconv.ParseInt(*item.Id, 10, 64)
			if err != nil {
				return &http.Response{StatusCode: http.StatusOK}, nil
			}

			idParam := sdk.UpdateStorageVolumesIdParameter{Int64: &id}

			_, hresp, delErr := client.StorageAPI.RemoveStorageVolumes(ctx, idParam).Execute()
			if delErr != nil && hresp != nil && hresp.StatusCode == http.StatusNotFound {
				return &http.Response{StatusCode: http.StatusOK}, nil
			}

			return hresp, delErr
		},
	)
}
