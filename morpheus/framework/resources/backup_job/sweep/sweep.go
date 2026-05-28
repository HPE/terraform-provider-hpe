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

const sweeperName = "hpe_morpheus_backup_job"

func init() {
	testsweep.RegisterTypedAPISweeper(
		sweeperName,
		// List backup job resources.
		func(ctx context.Context, client *sdk.APIClient) (
			[]sdk.ListBackupJobs200ResponseAllOfJobsInner,
			*http.Response,
			error,
		) {
			resp, hresp, err := client.BackupsAPI.ListBackupJobs(ctx).Execute()
			if resp == nil {
				return nil, hresp, err
			}

			return resp.GetJobs(), hresp, err
		},
		// Is this a test backup job?
		func(item sdk.ListBackupJobs200ResponseAllOfJobsInner) bool {
			name, ok := item.GetNameOk()
			if !ok || name == nil {
				return false
			}

			return strings.HasPrefix(*name, testsweep.TestResourcePrefix)
		},
		// Delete the test backup job.
		func(
			ctx context.Context,
			client *sdk.APIClient,
			item sdk.ListBackupJobs200ResponseAllOfJobsInner,
		) (*http.Response, error) {
			id, ok := item.GetIdOk()
			if !ok || id == nil {
				return nil, fmt.Errorf("could not get ID")
			}

			_, hresp, err := client.BackupsAPI.RemoveBackupJobs(ctx, *id).Execute()

			return hresp, err
		},
		testsweep.WithIgnoreListStatuses[sdk.ListBackupJobs200ResponseAllOfJobsInner](
			http.StatusNotFound,
			http.StatusForbidden,
		),
	)
}
