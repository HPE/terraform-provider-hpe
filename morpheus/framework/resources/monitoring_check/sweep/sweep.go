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

const sweeperName = "hpe_morpheus_monitoring_check"

func init() {
	testsweep.RegisterTypedAPISweeper(
		sweeperName,
		// List monitoring check resources.
		func(ctx context.Context, client *sdk.APIClient) (
			[]sdk.ListChecks200ResponseAllOfChecksInner,
			*http.Response,
			error,
		) {
			resp, hresp, err := client.ChecksAPI.ListChecks(ctx).Execute()
			if resp == nil {
				return nil, hresp, err
			}

			return getsafe.GetSafe(&resp.Checks), hresp, err
		},
		// Is this a test monitoring check?
		func(item sdk.ListChecks200ResponseAllOfChecksInner) bool {
			name, ok := getsafe.GetSafeOk(item.Name)
			if !ok || name == nil {
				return false
			}

			return strings.HasPrefix(*name, testsweep.TestResourcePrefix)
		},
		// Delete the test monitoring check.
		func(
			ctx context.Context,
			client *sdk.APIClient,
			item sdk.ListChecks200ResponseAllOfChecksInner,
		) (*http.Response, error) {
			id, ok := getsafe.GetSafeOk(item.Id)
			if !ok || id == nil {
				return nil, fmt.Errorf("could not get ID")
			}

			_, hresp, err := client.ChecksAPI.DeleteChecks(ctx, *id).Execute()

			return hresp, err
		},
		testsweep.WithIgnoreListStatuses[sdk.ListChecks200ResponseAllOfChecksInner](
			http.StatusNotFound,
			http.StatusForbidden,
		),
	)
}
