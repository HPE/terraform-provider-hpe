// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package monitoring_check_test

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"

	testsweep "github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/sweep"
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

			return resp.GetChecks(), hresp, err
		},
		// Is this a test monitoring check?
		func(item sdk.ListChecks200ResponseAllOfChecksInner) bool {
			name, ok := item.GetNameOk()
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
			id, ok := item.GetIdOk()
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
