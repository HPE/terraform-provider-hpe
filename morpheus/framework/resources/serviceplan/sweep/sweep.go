// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

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
		"hpe_morpheus_service_plan",
		// List all service plan resources.
		func(ctx context.Context, client *sdk.APIClient) (
			[]sdk.ListServicePlans200ResponseAllOfServicePlansInner,
			*http.Response,
			error,
		) {
			resp, hresp, err := client.ServicePlansAPI.ListServicePlans(ctx).Execute()
			if resp == nil {
				return nil, hresp, err
			}

			return resp.GetServicePlans(), hresp, err
		},
		// Is this a test service plan?
		func(item sdk.ListServicePlans200ResponseAllOfServicePlansInner) bool {
			name, ok := item.GetNameOk()
			if !ok || name == nil {
				return false
			}

			return strings.HasPrefix(*name, testsweep.TestResourcePrefix)
		},
		// Delete the test service plan.
		func(
			ctx context.Context,
			client *sdk.APIClient,
			item sdk.ListServicePlans200ResponseAllOfServicePlansInner,
		) (*http.Response, error) {
			id, ok := item.GetIdOk()
			if !ok || id == nil {
				return nil, fmt.Errorf("could not get ID")
			}

			_, hresp, err := client.ServicePlansAPI.RemoveServicePlans(ctx, *id).Execute()

			return hresp, err
		},
		testsweep.WithIgnoreListStatuses[sdk.ListServicePlans200ResponseAllOfServicePlansInner](
			http.StatusNotFound,
			http.StatusForbidden,
		),
	)
}
