// (C) Copyright 2026 Hewlett Packard Enterprise Development LP


//go:build sweep

package sweep

import (
	"context"
	"net/http"
	"strings"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"

	testsweep "github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/sweep"
)

const sweeperName = "hpe_morpheus_network_router_route"

type routeSweeperItem struct {
	routerID int64
	route    sdk.GetNetworkRoutersRoutes200ResponseNetworkRoutesInner
}

func init() {
	testsweep.RegisterTypedAPISweeper(
		sweeperName,
		// List network router route resources by iterating routers.
		func(ctx context.Context, client *sdk.APIClient) ([]routeSweeperItem, *http.Response, error) {
			routersResp, routersHTTPResp, err := client.NetworksAPI.GetNetworkRouters(ctx).Execute()
			if err != nil || routersResp == nil {
				return nil, routersHTTPResp, err
			}

			items := make([]routeSweeperItem, 0)

			for _, router := range routersResp.GetNetworkRouters() {
				routerID, ok := router.GetIdOk()
				if !ok || routerID == nil {
					continue
				}

				routesResp, _, listErr := client.NetworksAPI.GetNetworkRoutersRoutes(ctx, *routerID).Execute()
				if listErr != nil || routesResp == nil {
					continue
				}

				for _, route := range routesResp.GetNetworkRoutes() {
					items = append(items, routeSweeperItem{routerID: *routerID, route: route})
				}
			}

			return items, routersHTTPResp, nil
		},
		// Is this a test network router route?
		func(item routeSweeperItem) bool {
			name, ok := item.route.GetNameOk()
			if ok && name != nil && strings.HasPrefix(*name, testsweep.TestResourcePrefix) {
				return true
			}

			desc, ok := item.route.GetDescriptionOk()
			if ok && desc != nil && strings.HasPrefix(*desc, testsweep.TestResourcePrefix) {
				return true
			}

			return false
		},
		// Delete the test network router route.
		func(ctx context.Context, client *sdk.APIClient, item routeSweeperItem) (*http.Response, error) {
			id, ok := item.route.GetIdOk()
			if !ok || id == nil {
				return &http.Response{StatusCode: http.StatusOK}, nil
			}

			_, hresp, err := client.NetworksAPI.DeleteNetworkRouterRoute(ctx, *id, item.routerID).Execute()
			if err != nil && hresp != nil && hresp.StatusCode == http.StatusNotFound {
				return &http.Response{StatusCode: http.StatusOK}, nil
			}

			return hresp, err
		},
	)
}
