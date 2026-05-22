// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkrouterbgpneighbor_test

import (
	"context"
	"net/http"
	"strings"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"

	testsweep "github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/sweep"
)

const sweeperName = "hpe_morpheus_network_router_bgp_neighbor"

type bgpNeighborSweeperItem struct {
	routerID int64
	neighbor sdk.GetNetworkRoutersBgpNeighbors200ResponseNetworkRouterBgpNeighborsInner
}

func init() {
	testsweep.RegisterTypedAPISweeper(
		sweeperName,
		// List network router BGP neighbor resources by iterating routers.
		func(ctx context.Context, client *sdk.APIClient) ([]bgpNeighborSweeperItem, *http.Response, error) {
			routersResp, routersHTTPResp, err := client.NetworksAPI.GetNetworkRouters(ctx).Execute()
			if err != nil || routersResp == nil {
				return nil, routersHTTPResp, err
			}

			items := make([]bgpNeighborSweeperItem, 0)

			for _, router := range routersResp.GetNetworkRouters() {
				routerID, ok := router.GetIdOk()
				if !ok || routerID == nil {
					continue
				}

				neighborsResp, neighborsHTTPResp, listErr := client.NetworksAPI.GetNetworkRoutersBgpNeighbors(ctx, *routerID).Execute()
				if listErr != nil || neighborsResp == nil {
					return nil, neighborsHTTPResp, listErr
				}

				for _, neighbor := range neighborsResp.GetNetworkRouterBgpNeighbors() {
					items = append(items, bgpNeighborSweeperItem{routerID: *routerID, neighbor: neighbor})
				}
			}

			return items, routersHTTPResp, nil
		},
		// Is this a test network router BGP neighbor?
		func(item bgpNeighborSweeperItem) bool {
			return isTestBgpNeighbor(item.neighbor)
		},
		// Delete the test network router BGP neighbor.
		func(ctx context.Context, client *sdk.APIClient, item bgpNeighborSweeperItem) (*http.Response, error) {
			neighborID, ok := item.neighbor.GetIdOk()
			if !ok || neighborID == nil {
				return &http.Response{StatusCode: http.StatusOK}, nil
			}

			_, delResp, delErr := client.NetworksAPI.DeleteNetworkRouterBgpNeighbor(ctx, *neighborID, item.routerID).Execute()
			if delErr != nil && delResp != nil && delResp.StatusCode == http.StatusNotFound {
				return &http.Response{StatusCode: http.StatusOK}, nil
			}

			return delResp, delErr
		},
	)
}

func isTestBgpNeighbor(
	neighbor sdk.GetNetworkRoutersBgpNeighbors200ResponseNetworkRouterBgpNeighborsInner,
) bool {
	desc, ok := neighbor.GetDescriptionOk()
	if ok && desc != nil && strings.HasPrefix(*desc, "TestAcc") {
		return true
	}

	return false
}
