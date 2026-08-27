// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

//go:build sweep

package sweep

import (
	"context"
	"net/http"
	"strings"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"

	testsweep "github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/sweep"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/getsafe"
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

			for _, router := range routersResp.NetworkRouters {
				routerID, ok := getsafe.GetOk(router.Id)
				if !ok || routerID == nil {
					continue
				}

				natsResp, _, listErr := client.NetworksAPI.GetNetworkRoutersNats(ctx, *routerID).Execute()
				if listErr != nil || natsResp == nil {
					continue
				}

				for _, nat := range natsResp.NetworkRouterNATs {
					items = append(items, natSweeperItem{routerID: *routerID, nat: nat})
				}
			}

			return items, routersHTTPResp, nil
		},
		// Is this a test network router NAT?
		func(item natSweeperItem) bool {
			name, ok := getsafe.GetOk(item.nat.Name)
			if !ok || name == nil {
				return false
			}

			return strings.HasPrefix(*name, testsweep.TestResourcePrefix)
		},
		// Delete the test network router NAT.
		func(ctx context.Context, client *sdk.APIClient, item natSweeperItem) (*http.Response, error) {
			id, ok := getsafe.GetOk(item.nat.Id)
			if !ok || id == nil {
				return &http.Response{StatusCode: http.StatusOK}, nil
			}

			_, hresp, err := client.NetworksAPI.DeleteNetworkRouterNat(ctx, int64(*id), item.routerID).Execute()
			if err != nil && hresp != nil && hresp.StatusCode == http.StatusNotFound {
				return &http.Response{StatusCode: http.StatusOK}, nil
			}

			return hresp, err
		},
		// The router endpoints this sweeper lists from are not present on every
		// appliance version, and a sweep account without network permissions
		// gets 403. Either way there is nothing to sweep, but without this the
		// sweeper reports a failure, and a failing sweeper leaves its own test
		// resources behind while making the whole sweep look broken.
		testsweep.WithIgnoreListStatuses[natSweeperItem](
			http.StatusNotFound,
			http.StatusForbidden,
		),
	)
}
