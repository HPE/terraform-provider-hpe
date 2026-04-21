// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package loadbalancer

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
)

// testLoadBalancerPrefix is truncated to 16 characters to account for the
// 32-character name limit after RandomWithPrefix appends a random suffix.
const testLoadBalancerPrefix = "TestAccMorpheusL"

func init() {
	testhelpers.RegisterTypedAPISweeper(
		"hpe_morpheus_load_balancer",
		// List candidate load balancers for sweeping.
		func(ctx context.Context, client *sdk.APIClient, _ string) (
			[]sdk.ListLoadBalancers200ResponseAllOfLoadBalancersInner,
			*http.Response,
			error,
		) {
			resp, hresp, err := client.LoadBalancersAPI.ListLoadBalancers(ctx).Execute()
			if resp == nil {
				return nil, hresp, err
			}

			return resp.GetLoadBalancers(), hresp, err
		},
		// Match only acceptance-test load balancers.
		func(item sdk.ListLoadBalancers200ResponseAllOfLoadBalancersInner) bool {
			name, ok := item.GetNameOk()
			if !ok || name == nil {
				return false
			}

			return strings.HasPrefix(*name, testLoadBalancerPrefix)
		},
		// Delete a matched load balancer.
		func(
			ctx context.Context,
			client *sdk.APIClient,
			item sdk.ListLoadBalancers200ResponseAllOfLoadBalancersInner,
		) (*http.Response, error) {
			id, ok := item.GetIdOk()
			if !ok || id == nil {
				return nil, fmt.Errorf("could not get ID")
			}

			_, hresp, err := client.LoadBalancersAPI.DeleteLoadBalancer(ctx, *id).Execute()

			return hresp, err
		},
	)
}
