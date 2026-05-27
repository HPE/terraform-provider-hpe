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

const sweeperName = "hpe_morpheus_budget"

func init() {
	testsweep.RegisterTypedAPISweeper(
		sweeperName,
		// List budget resources.
		func(ctx context.Context, client *sdk.APIClient) (
			[]sdk.ListBudgets200ResponseAllOfBudgetsInner,
			*http.Response,
			error,
		) {
			resp, hresp, err := client.BudgetsAPI.ListBudgets(ctx).Execute()
			if resp == nil {
				return nil, hresp, err
			}

			return resp.GetBudgets(), hresp, err
		},
		// Is this a test budget?
		func(item sdk.ListBudgets200ResponseAllOfBudgetsInner) bool {
			name, ok := item.GetNameOk()
			if !ok || name == nil {
				return false
			}

			return strings.HasPrefix(*name, testsweep.TestResourcePrefix)
		},
		// Delete the test budget.
		func(
			ctx context.Context,
			client *sdk.APIClient,
			item sdk.ListBudgets200ResponseAllOfBudgetsInner,
		) (*http.Response, error) {
			id, ok := item.GetIdOk()
			if !ok || id == nil {
				return nil, fmt.Errorf("could not get ID")
			}

			_, hresp, err := client.BudgetsAPI.RemoveBudgets(ctx, *id).Execute()

			return hresp, err
		},
		testsweep.WithIgnoreListStatuses[sdk.ListBudgets200ResponseAllOfBudgetsInner](
			http.StatusNotFound,
			http.StatusForbidden,
		),
	)
}
