// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package network_router_nat_test

import (
	"context"
	"net/http"
	"strings"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"

	testsweep "github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/sweep"
)

const sweeperName = "hpe_morpheus_network_router_nat"

type natSweeperItem struct {
	routerID int64
	nat      sdk.GetNetworkRoutersNats200ResponseNetworkRouterNATsInner
}

func init() {
	testsweep.RegisterTypedAPISweeper(
		sweeperName,
		// List network router NAT resources by iterating routers.
		func(ctx context.Context, client *sdk.APIClient) ([]natSweeperItem, *http.Response, error) {
			routersResp, routersHTTPResp, err := client.NetworksAPI.GetNetworkRouters(ctx).Execute()
			if err != nil || routersResp == nil {
				return nil, routersHTTPResp, err
			}

			items := make([]natSweeperItem, 0)

			for _, router := range routersResp.GetNetworkRouters() {
				routerID, ok := router.GetIdOk()
				if !ok || routerID == nil {
					continue
				}

				natsResp, _, listErr := client.NetworksAPI.GetNetworkRoutersNats(ctx, *routerID).Execute()
				if listErr != nil || natsResp == nil {
					continue
				}

				for _, nat := range natsResp.GetNetworkRouterNATs() {
					items = append(items, natSweeperItem{routerID: *routerID, nat: nat})
				}
			}

			return items, routersHTTPResp, nil
		},
		// Is this a test network router NAT?
		func(item natSweeperItem) bool {
			name, ok := item.nat.GetNameOk()
			if !ok || name == nil {
				return false
			}

			return strings.HasPrefix(*name, testsweep.TestResourcePrefix)
		},
		// Delete the test network router NAT.
		func(ctx context.Context, client *sdk.APIClient, item natSweeperItem) (*http.Response, error) {
			id, ok := item.nat.GetIdOk()
			if !ok || id == nil {
				return &http.Response{StatusCode: http.StatusOK}, nil
			}

			_, hresp, err := client.NetworksAPI.DeleteNetworkRouterNat(ctx, int64(*id), item.routerID).Execute()
			if err != nil && hresp != nil && hresp.StatusCode == http.StatusNotFound {
				return &http.Response{StatusCode: http.StatusOK}, nil
			}

			return hresp, err
		},
	)
}
