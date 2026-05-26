// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package monitoring_alert_test

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"

	testsweep "github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/sweep"
)

const sweeperName = "hpe_morpheus_monitoring_alert"

func init() {
	testsweep.RegisterTypedAPISweeper(
		sweeperName,
		// List monitoring alert resources.
		func(ctx context.Context, client *sdk.APIClient) (
			[]sdk.ListAlerts200ResponseAllOfAlertsInner,
			*http.Response,
			error,
		) {
			resp, hresp, err := client.AlertsAPI.ListAlerts(ctx).Execute()
			if resp == nil {
				return nil, hresp, err
			}

			return resp.GetAlerts(), hresp, err
		},
		// Is this a test monitoring alert?
		func(item sdk.ListAlerts200ResponseAllOfAlertsInner) bool {
			name, ok := item.GetNameOk()
			if !ok || name == nil {
				return false
			}

			return strings.HasPrefix(*name, testsweep.TestResourcePrefix)
		},
		// Delete the test monitoring alert.
		func(
			ctx context.Context,
			client *sdk.APIClient,
			item sdk.ListAlerts200ResponseAllOfAlertsInner,
		) (*http.Response, error) {
			id, ok := item.GetIdOk()
			if !ok || id == nil {
				return nil, fmt.Errorf("could not get ID")
			}

			_, hresp, err := client.AlertsAPI.DeleteAlerts(ctx, *id).Execute()

			return hresp, err
		},
		testsweep.WithIgnoreListStatuses[sdk.ListAlerts200ResponseAllOfAlertsInner](
			http.StatusNotFound,
			http.StatusForbidden,
		),
	)
}
