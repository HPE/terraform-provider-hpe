// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package sweep

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"

	testsweep "github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/sweep"
)

func init() {
	testsweep.RegisterTypedAPISweeper(
		"hpe_morpheus_policy",
		// List all policy resources.
		func(ctx context.Context, client *sdk.APIClient) (
			[]sdk.ListPolicies200ResponseAllOfPoliciesInner,
			*http.Response,
			error,
		) {
			resp, hresp, err := client.PoliciesAPI.ListPolicies(ctx).Execute()
			if resp == nil {
				return nil, hresp, err
			}

			return resp.GetPolicies(), hresp, err
		},
		// Is this a test policy?
		func(item sdk.ListPolicies200ResponseAllOfPoliciesInner) bool {
			name, ok := item.GetNameOk()
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
			id, ok := item.GetIdOk()
			if !ok || id == nil {
				return nil, fmt.Errorf("could not get ID")
			}

			_, hresp, err := client.PoliciesAPI.RemovePolicies(ctx, *id).Execute()

			return hresp, err
		},
		// Ignore (i.e. just log) "not found" and "forbidden" errors
		testsweep.WithIgnoreListStatuses[sdk.ListPolicies200ResponseAllOfPoliciesInner](
			http.StatusNotFound,
			http.StatusForbidden,
		),
	)
}
