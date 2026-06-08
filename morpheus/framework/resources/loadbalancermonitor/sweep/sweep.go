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

const sweeperName = "hpe_morpheus_load_balancer_monitor"

// loadBalancerMonitorSweepItem pairs a monitor with its parent load balancer ID.
type loadBalancerMonitorSweepItem struct {
	loadBalancerID int64
	id             int64
	name           string
}

func init() {
	testsweep.RegisterTypedAPISweeper(
		sweeperName,
		// List load balancer monitor resources.
		func(ctx context.Context, client *sdk.APIClient) (
			[]loadBalancerMonitorSweepItem,
			*http.Response,
			error,
		) {
			lbResp, hresp, err := client.LoadBalancersAPI.ListLoadBalancers(ctx).Execute()
			if err != nil {
				return nil, hresp, err
			}

			var items []loadBalancerMonitorSweepItem

			for _, lb := range lbResp.LoadBalancers {
				lbID, ok := getsafe.GetOk(lb.Id)
				if !ok || lbID == nil {
					continue
				}

				monResp, _, err := client.LoadBalancersAPI.
					ListLoadBalancerMonitors(ctx, *lbID).Execute()
				if err != nil || monResp == nil {
					continue
				}

				for _, mon := range monResp.LoadBalancerMonitors {
					id, ok := getsafe.GetOk(mon.Id)
					if !ok || id == nil {
						continue
					}

					name, ok := getsafe.GetOk(mon.Name)
					if !ok || name == nil {
						continue
					}

					items = append(items, loadBalancerMonitorSweepItem{
						loadBalancerID: *lbID,
						id:             *id,
						name:           *name,
					})
				}
			}

			return items, hresp, nil
		},
		// Is this a test load balancer monitor?
		func(item loadBalancerMonitorSweepItem) bool {
			return strings.HasPrefix(item.name, testsweep.TestResourcePrefix)
		},
		// Delete the test load balancer monitor.
		func(
			ctx context.Context,
			client *sdk.APIClient,
			item loadBalancerMonitorSweepItem,
		) (*http.Response, error) {
			if item.id == 0 {
				return nil, fmt.Errorf("could not get ID")
			}

			_, hresp, err := client.LoadBalancersAPI.
				DeleteLoadBalancerMonitor(ctx, item.loadBalancerID, item.id).Execute()

			return hresp, err
		},
		testsweep.WithIgnoreListStatuses[loadBalancerMonitorSweepItem](
			http.StatusNotFound,
			http.StatusForbidden,
		),
	)
}
