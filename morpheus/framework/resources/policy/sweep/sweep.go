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

const sweeperName = "hpe_morpheus_policy"

func init() {
	testsweep.RegisterTypedAPISweeper(
		sweeperName,
		// List policy resources.
		func(ctx context.Context, client *sdk.APIClient) (
			[]sdk.ListPolicies200ResponseAllOfPoliciesInner,
			*http.Response,
			error,
		) {
			resp, hresp, err := client.PoliciesAPI.ListPolicies(ctx).Execute()
			if resp == nil {
				return nil, hresp, err
			}

			return getsafe.GetSafe(&resp.Policies), hresp, err
		},
		// Is this a test policy?
		func(item sdk.ListPolicies200ResponseAllOfPoliciesInner) bool {
			name, ok := getsafe.GetSafeOk(item.Name)
			if !ok || name == nil {
				return false
			}

			return strings.HasPrefix(*name, testsweep.TestResourcePrefix)
		},
		// Delete the test policy.
		func(
			ctx context.Context,
			client *sdk.APIClient,
			item sdk.ListPolicies200ResponseAllOfPoliciesInner,
		) (*http.Response, error) {
			id, ok := getsafe.GetSafeOk(item.Id)
			if !ok || id == nil {
				return nil, fmt.Errorf("could not get ID")
			}

			_, hresp, err := client.PoliciesAPI.RemovePolicies(ctx, *id).Execute()

			return hresp, err
		},
		// Ignore (i.e. just log) "not found" and "forbidden" errors.
		testsweep.WithIgnoreListStatuses[sdk.ListPolicies200ResponseAllOfPoliciesInner](
			http.StatusNotFound,
			http.StatusForbidden,
		),
	)
}
