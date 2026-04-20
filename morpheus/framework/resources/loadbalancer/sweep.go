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
	testhelpers.RegisterAPISweeper(
		"hpe_morpheus_load_balancer",
		testLoadBalancerPrefix,
		func(ctx context.Context, client *sdk.APIClient, prefix string) (any, *http.Response, error) {
			return client.LoadBalancersAPI.ListLoadBalancers(ctx).Phrase(prefix).Execute()
		},
		"GetLoadBalancers",
		func(ctx context.Context, client *sdk.APIClient, id int64, _ any) (*http.Response, error) {
			_, hresp, err := client.LoadBalancersAPI.DeleteLoadBalancer(ctx, id).Execute()
			return hresp, err
		},
	)
}
