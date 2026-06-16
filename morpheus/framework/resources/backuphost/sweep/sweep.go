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
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/getsafe"
)

const sweeperName = "hpe_morpheus_backup_host"

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

			return getsafe.Get(&resp.Backups), hresp, err
		},
		// Is this a test host backup?
		func(item sdk.ListBackups200ResponseAllOfBackupsInner) bool {
			if item.LocationType == nil || *item.LocationType != "server" {
				return false
			}

			name, ok := getsafe.GetOk(item.Name)
			if !ok || name == nil {
				return false
			}

			return strings.HasPrefix(*name, testsweep.TestResourcePrefix)
		},
		// Delete the test host backup.
		func(
			ctx context.Context,
			client *sdk.APIClient,
			item sdk.ListBackups200ResponseAllOfBackupsInner,
		) (*http.Response, error) {
			id, ok := getsafe.GetOk(item.Id)
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
