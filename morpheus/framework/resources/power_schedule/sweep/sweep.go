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

const sweeperName = "hpe_morpheus_power_schedule"

func init() {
	testsweep.RegisterTypedAPISweeper(
		sweeperName,
		// List power schedule resources.
		func(ctx context.Context, client *sdk.APIClient) (
			[]sdk.ListPowerSchedules200ResponseAllOfSchedulesInner,
			*http.Response,
			error,
		) {
			resp, hresp, err := client.AutomationAPI.ListPowerSchedules(ctx).Execute()
			if resp == nil {
				return nil, hresp, err
			}

			return getsafe.GetSafe(&resp.Schedules), hresp, err
		},
		// Is this a test power schedule?
		func(item sdk.ListPowerSchedules200ResponseAllOfSchedulesInner) bool {
			name, ok := getsafe.GetSafeOk(item.Name)
			if !ok || name == nil {
				return false
			}

			return strings.HasPrefix(*name, testsweep.TestResourcePrefix)
		},
		// Delete the test power schedule.
		func(
			ctx context.Context,
			client *sdk.APIClient,
			item sdk.ListPowerSchedules200ResponseAllOfSchedulesInner,
		) (*http.Response, error) {
			id, ok := getsafe.GetSafeOk(item.Id)
			if !ok || id == nil {
				return nil, fmt.Errorf("could not get ID")
			}

			_, hresp, err := client.AutomationAPI.RemovePowerSchedules(ctx, *id).Execute()

			return hresp, err
		},
		testsweep.WithIgnoreListStatuses[sdk.ListPowerSchedules200ResponseAllOfSchedulesInner](
			http.StatusNotFound,
			http.StatusForbidden,
		),
	)
}
