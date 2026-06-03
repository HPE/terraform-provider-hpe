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

const sweeperName = "hpe_morpheus_load_balancer_virtual_server"

type virtualServerSweepItem struct {
	loadBalancerID int64
	id             int64
	vipName        string
}

func init() {
	testsweep.RegisterTypedAPISweeper(
		sweeperName,
		// List load balancer virtual server resources by iterating load balancers.
		func(ctx context.Context, client *sdk.APIClient) (
			[]virtualServerSweepItem,
			*http.Response,
			error,
		) {
			lbResp, lbHresp, lbErr := client.LoadBalancersAPI.ListLoadBalancers(ctx).Execute()
			if lbErr != nil || lbResp == nil {
				return nil, lbHresp, lbErr
			}

			items := make([]virtualServerSweepItem, 0)

			for _, lb := range lbResp.LoadBalancers {
				lbID, ok := getsafe.GetSafeOk(lb.Id)
				if !ok || lbID == nil {
					continue
				}

				vsResp, _, vsErr := client.LoadBalancersAPI.
					ListLoadBalancerVirtualServers(ctx, *lbID).Execute()
				if vsErr != nil || vsResp == nil {
					continue
				}

				for _, vs := range vsResp.LoadBalancerInstances {
					id, ok := getsafe.GetSafeOk(vs.Id)
					if !ok || id == nil {
						continue
					}

					name, ok := getsafe.GetSafeOk(vs.VipName)
					if !ok || name == nil {
						continue
					}

					items = append(items, virtualServerSweepItem{
						loadBalancerID: *lbID,
						id:             *id,
						vipName:        *name,
					})
				}
			}

			return items, lbHresp, nil
		},
		// Is this a test virtual server?
		func(item virtualServerSweepItem) bool {
			return strings.HasPrefix(item.vipName, testsweep.TestResourcePrefix)
		},
		// Delete the test virtual server.
		func(
			ctx context.Context,
			client *sdk.APIClient,
			item virtualServerSweepItem,
		) (*http.Response, error) {
			if item.id == 0 {
				return nil, fmt.Errorf("could not get ID")
			}

			_, hresp, err := client.LoadBalancersAPI.
				DeleteLoadBalancerVirtualServer(ctx, item.loadBalancerID, item.id).
				Execute()

			return hresp, err
		},
	)
}
