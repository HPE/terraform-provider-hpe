// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package backup_test

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"

	testsweep "github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/sweep"
)

const sweeperName = "hpe_morpheus_backup"

func init() {
	testsweep.RegisterTypedAPISweeper(
		sweeperName,
		// List backup resources.
		func(ctx context.Context, client *sdk.APIClient) (
			[]sdk.ListBackups200ResponseAllOfBackupsInner,
			*http.Response,
			error,
		) {
			resp, hresp, err := client.BackupsAPI.ListBackups(ctx).Execute()
			if resp == nil {
				return nil, hresp, err
			}

			return resp.GetBackups(), hresp, err
		},
		// Is this a test backup?
		func(item sdk.ListBackups200ResponseAllOfBackupsInner) bool {
			name, ok := item.GetNameOk()
			if !ok || name == nil {
				return false
			}

			return strings.HasPrefix(*name, testsweep.TestResourcePrefix)
		},
		// Delete the test backup.
		func(
			ctx context.Context,
			client *sdk.APIClient,
			item sdk.ListBackups200ResponseAllOfBackupsInner,
		) (*http.Response, error) {
			id, ok := item.GetIdOk()
			if !ok || id == nil {
				return nil, fmt.Errorf("could not get ID")
			}

			_, hresp, err := client.BackupsAPI.RemoveBackups(ctx, *id).Execute()

			return hresp, err
		},
		testsweep.WithIgnoreListStatuses[sdk.ListBackups200ResponseAllOfBackupsInner](
			http.StatusNotFound,
			http.StatusForbidden,
		),
	)
}
