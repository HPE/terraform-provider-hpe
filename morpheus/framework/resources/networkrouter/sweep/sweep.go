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

const sweeperName = "hpe_morpheus_network_router"

func init() {
	testsweep.RegisterTypedAPISweeper(
		sweeperName,
		// List network router resources.
		func(ctx context.Context, client *sdk.APIClient) (
			[]sdk.GetNetworkRouters200ResponseNetworkRoutersInner,
			*http.Response,
			error,
		) {
			resp, hresp, err := client.NetworksAPI.GetNetworkRouters(ctx).Execute()
			if resp == nil {
				return nil, hresp, err
			}

			return getsafe.Get(&resp.NetworkRouters), hresp, err
		},
		// Is this a test network router?
		func(item sdk.GetNetworkRouters200ResponseNetworkRoutersInner) bool {
			name, ok := getsafe.GetOk(item.Name)
			if !ok || name == nil {
				return false
			}

			return strings.HasPrefix(*name, testsweep.TestResourcePrefix)
		},
		// Delete the test network router.
		func(
			ctx context.Context,
			client *sdk.APIClient,
			item sdk.GetNetworkRouters200ResponseNetworkRoutersInner,
		) (*http.Response, error) {
			id, ok := getsafe.GetOk(item.Id)
			if !ok || id == nil {
				return nil, fmt.Errorf("could not get ID")
			}

			_, hresp, err := client.NetworksAPI.DeleteNetworkRouter(ctx, *id).Execute()

			return hresp, err
		},
		testsweep.WithDependencies[sdk.GetNetworkRouters200ResponseNetworkRoutersInner](
			"hpe_morpheus_network_router_route",
			"hpe_morpheus_network_router_firewall_rule",
			"hpe_morpheus_network_router_bgp_neighbor",
			"hpe_morpheus_network_router_nat",
		),
	)
}
