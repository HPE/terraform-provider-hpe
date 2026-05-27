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

const sweeperName = "hpe_morpheus_deployment"

func init() {
	testsweep.RegisterTypedAPISweeper(
		sweeperName,
		// List deployment resources.
		func(ctx context.Context, client *sdk.APIClient) (
			[]sdk.ListDeployments200ResponseAllOfDeploymentsInner,
			*http.Response,
			error,
		) {
			resp, hresp, err := client.DeploymentsAPI.ListDeployments(ctx).Execute()
			if resp == nil {
				return nil, hresp, err
			}

			return resp.GetDeployments(), hresp, err
		},
		// Is this a test deployment?
		func(item sdk.ListDeployments200ResponseAllOfDeploymentsInner) bool {
			name, ok := item.GetNameOk()
			if !ok || name == nil {
				return false
			}

			return strings.HasPrefix(*name, testsweep.TestResourcePrefix)
		},
		// Delete the test deployment.
		func(
			ctx context.Context,
			client *sdk.APIClient,
			item sdk.ListDeployments200ResponseAllOfDeploymentsInner,
		) (*http.Response, error) {
			id, ok := item.GetIdOk()
			if !ok || id == nil {
				return nil, fmt.Errorf("could not get ID")
			}

			_, hresp, err := client.DeploymentsAPI.DeleteDeployment(ctx, *id).Execute()

			return hresp, err
		},
		testsweep.WithIgnoreListStatuses[sdk.ListDeployments200ResponseAllOfDeploymentsInner](
			http.StatusNotFound,
			http.StatusForbidden,
		),
	)
}
