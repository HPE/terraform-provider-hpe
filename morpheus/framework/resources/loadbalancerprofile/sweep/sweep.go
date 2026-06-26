// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

//go:build sweep

package sweep

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"

	testsweep "github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/sweep"
)

const sweeperName = "hpe_morpheus_load_balancer_profile"

// loadBalancerProfileSweepItem pairs a profile with its parent load balancer ID.
type loadBalancerProfileSweepItem struct {
	loadBalancerID int64
	id             int64
	name           string
}

func init() {
	testsweep.RegisterTypedAPISweeper(
		sweeperName,
		// List load balancer profile resources.
		func(ctx context.Context, client *sdk.APIClient) (
			[]loadBalancerProfileSweepItem,
			*http.Response,
			error,
		) {
			lbResp, hresp, err := client.LoadBalancersAPI.ListLoadBalancers(ctx).Execute()
			if err != nil {
				return nil, hresp, err
			}

			var items []loadBalancerProfileSweepItem

			for _, lb := range lbResp.LoadBalancers {
				lbID := lb.Id
				if lbID == nil {
					continue
				}

				profileResp, _, err := client.LoadBalancersAPI.
					ListLoadBalancerProfiles(ctx, *lbID).Execute()
				if err != nil || profileResp == nil {
					continue
				}

				for _, profile := range profileResp.LoadBalancerProfiles {
					id := profile.Id
					if id == nil {
						continue
					}

					name := profile.Name
					if name == nil {
						continue
					}

					items = append(items, loadBalancerProfileSweepItem{
						loadBalancerID: *lbID,
						id:             *id,
						name:           *name,
					})
				}
			}

			return items, hresp, nil
		},
		// Is this a test load balancer profile?
		func(item loadBalancerProfileSweepItem) bool {
			return strings.HasPrefix(item.name, testsweep.TestResourcePrefix)
		},
		// Delete the test load balancer profile.
		func(
			ctx context.Context,
			client *sdk.APIClient,
			item loadBalancerProfileSweepItem,
		) (*http.Response, error) {
			if item.id == 0 {
				return nil, fmt.Errorf("could not get ID")
			}

			_, hresp, err := client.LoadBalancersAPI.
				DeleteLoadBalancerProfile(ctx, item.loadBalancerID, item.id).Execute()

			return hresp, err
		},
		testsweep.WithIgnoreListStatuses[loadBalancerProfileSweepItem](
			http.StatusNotFound,
			http.StatusForbidden,
		),
	)
}
