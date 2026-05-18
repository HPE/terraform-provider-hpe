// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkrouterroute_test

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

const sweeperName = "hpe_morpheus_network_router_route"

func init() {
	resource.AddTestSweepers(sweeperName, &resource.Sweeper{
		Name: sweeperName,
		F:    sweepNetworkRouterRoutes,
	})
}

func sweepNetworkRouterRoutes(system string) error {
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

		if err := sweepRoutesForRouter(ctx, client, *routerID); err != nil {
			log.Printf("[ERROR] %s", err)
		}
	}

	return nil
}

func sweepRoutesForRouter(ctx context.Context, client *sdk.APIClient, routerID int64) error {
	resp, hresp, err := client.NetworksAPI.GetNetworkRoutersRoutes(ctx, routerID).Execute()
	if err != nil || hresp == nil || hresp.StatusCode != http.StatusOK {
		return fmt.Errorf("error listing routes for router %d: %s", routerID, err)
	}

	for _, route := range resp.GetNetworkRoutes() {
		if !isTestRoute(route) {
			continue
		}

		routeID, ok := route.GetIdOk()
		if !ok || routeID == nil {
			continue
		}

		log.Printf("[INFO] Sweeping route %d (name=%s) on router %d",
			*routeID, route.GetName(), routerID)

		_, delResp, delErr := client.NetworksAPI.
			DeleteNetworkRouterRoute(ctx, *routeID, routerID).
			Execute()
		if delErr != nil {
			if delResp != nil && delResp.StatusCode == http.StatusNotFound {
				continue
			}

			log.Printf("[ERROR] Error destroying route %d: %s", *routeID, delErr)
		}
	}

	return nil
}

func isTestRoute(
	route sdk.GetNetworkRoutersRoutes200ResponseNetworkRoutesInner,
) bool {
	desc, ok := route.GetDescriptionOk()
	if ok && desc != nil && strings.HasPrefix(*desc, "TestAcc") {
		return true
	}

	return false
}
