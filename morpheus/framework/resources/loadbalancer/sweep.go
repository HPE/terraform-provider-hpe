// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package loadbalancer

import (
	"context"
	"net/http"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
)

// testLoadBalancerPrefix is truncated to 16 characters to account for the
// 32-character name limit after RandomWithPrefix appends a random suffix.
const testLoadBalancerPrefix = "TestAccMorpheusL"

func init() {
	testhelpers.RegisterTypedAPISweeper(
		"hpe_morpheus_load_balancer",
		testLoadBalancerPrefix,
		func(ctx context.Context, client *sdk.APIClient, prefix string) ([]sdk.ListLoadBalancers200ResponseAllOfLoadBalancersInner, *http.Response, error) {
			resp, hresp, err := client.LoadBalancersAPI.ListLoadBalancers(ctx).Phrase(prefix).Execute()
			if resp == nil {
				return nil, hresp, err
			}

			return resp.GetLoadBalancers(), hresp, err
		},
		func(item sdk.ListLoadBalancers200ResponseAllOfLoadBalancersInner) (string, bool) {
			name, ok := item.GetNameOk()
			if !ok || name == nil {
				return "", false
			}

			return *name, true
		},
		func(item sdk.ListLoadBalancers200ResponseAllOfLoadBalancersInner) (int64, bool) {
			id, ok := item.GetIdOk()
			if !ok || id == nil {
				return 0, false
			}

			return *id, true
		},
		func(
			ctx context.Context,
			client *sdk.APIClient,
			id int64,
			_ sdk.ListLoadBalancers200ResponseAllOfLoadBalancersInner,
		) (*http.Response, error) {
			_, hresp, err := client.LoadBalancersAPI.DeleteLoadBalancer(ctx, id).Execute()
			return hresp, err
		},
	)
}
