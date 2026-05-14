// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkrouterbgpneighbor_test

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
)

const sweeperName = "hpe_morpheus_network_router_bgp_neighbor"

func init() {
	resource.AddTestSweepers(sweeperName, &resource.Sweeper{
		Name: sweeperName,
		F:    sweepNetworkRouterBgpNeighbors,
	})
}

func sweepNetworkRouterBgpNeighbors(system string) error {
	ctx := context.Background()

	client, err := testhelpers.NewClientForServer(ctx, system)
	if err != nil {
		return fmt.Errorf("error getting client: %s", err)
	}

	resp, hresp, err := client.NetworksAPI.GetNetworkRouters(ctx).Execute()
	if err != nil || hresp == nil || hresp.StatusCode != http.StatusOK {
		return fmt.Errorf("error listing routers: %s", err)
	}

	for _, router := range resp.GetNetworkRouters() {
		routerID, ok := router.GetIdOk()
		if !ok || routerID == nil {
			continue
		}

		if err := sweepBgpNeighborsForRouter(ctx, client, *routerID); err != nil {
			log.Printf("[ERROR] %s", err)
		}
	}

	return nil
}

func sweepBgpNeighborsForRouter(ctx context.Context, client *sdk.APIClient, routerID int64) error {
	resp, hresp, err := client.NetworksAPI.GetNetworkRoutersBgpNeighbors(ctx, routerID).Execute()
	if err != nil || hresp == nil || hresp.StatusCode != http.StatusOK {
		return fmt.Errorf("error listing BGP neighbors for router %d: %s", routerID, err)
	}

	for _, neighbor := range resp.GetNetworkRouterBgpNeighbors() {
		if !isTestBgpNeighbor(neighbor) {
			continue
		}

		neighborID, ok := neighbor.GetIdOk()
		if !ok || neighborID == nil {
			continue
		}

		log.Printf("[INFO] Sweeping BGP neighbor %d (ip=%s) on router %d",
			*neighborID, neighbor.GetIpAddress(), routerID)

		_, delResp, delErr := client.NetworksAPI.
			DeleteNetworkRouterBgpNeighbor(ctx, *neighborID, routerID).
			Execute()
		if delErr != nil {
			if delResp != nil && delResp.StatusCode == http.StatusNotFound {
				continue
			}

			log.Printf("[ERROR] Error destroying BGP neighbor %d: %s", *neighborID, delErr)
		}
	}

	return nil
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
