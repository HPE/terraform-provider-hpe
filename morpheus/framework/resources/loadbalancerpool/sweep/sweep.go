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

const sweeperName = "hpe_morpheus_load_balancer_pool"

// loadBalancerPoolSweepItem pairs a pool with its parent load balancer ID.
type loadBalancerPoolSweepItem struct {
	loadBalancerID int64
	id             int64
	name           string
}

func init() {
	testsweep.RegisterTypedAPISweeper(
		sweeperName,
		// List load balancer pool resources.
		func(ctx context.Context, client *sdk.APIClient) (
			[]loadBalancerPoolSweepItem,
			*http.Response,
			error,
		) {
			lbResp, hresp, err := client.LoadBalancersAPI.ListLoadBalancers(ctx).Execute()
			if err != nil {
				return nil, hresp, err
			}

			var items []loadBalancerPoolSweepItem

			for _, lb := range lbResp.GetLoadBalancers() {
				lbID, ok := lb.GetIdOk()
				if !ok || lbID == nil {
					continue
				}

				poolResp, _, err := client.LoadBalancersAPI.
					ListLoadBalancerPools(ctx, *lbID).Execute()
				if err != nil || poolResp == nil {
					continue
				}

				for _, pool := range poolResp.GetLoadBalancerPools() {
					id, ok := pool.GetIdOk()
					if !ok || id == nil {
						continue
					}

					name, ok := pool.GetNameOk()
					if !ok || name == nil {
						continue
					}

					items = append(items, loadBalancerPoolSweepItem{
						loadBalancerID: *lbID,
						id:             *id,
						name:           *name,
					})
				}
			}

			return items, hresp, nil
		},
		// Is this a test load balancer pool?
		func(item loadBalancerPoolSweepItem) bool {
			return strings.HasPrefix(item.name, testsweep.TestResourcePrefix)
		},
		// Delete the test load balancer pool.
		func(
			ctx context.Context,
			client *sdk.APIClient,
			item loadBalancerPoolSweepItem,
		) (*http.Response, error) {
			if item.id == 0 {
				return nil, fmt.Errorf("could not get ID")
			}

			_, hresp, err := client.LoadBalancersAPI.
				DeleteLoadBalancerPool(ctx, item.loadBalancerID, item.id).Execute()

			return hresp, err
		},
		testsweep.WithIgnoreListStatuses[loadBalancerPoolSweepItem](
			http.StatusNotFound,
			http.StatusForbidden,
		),
	)
}
