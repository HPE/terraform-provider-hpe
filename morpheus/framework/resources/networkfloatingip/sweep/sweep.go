// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

//go:build sweep

package sweep

import (
	"context"
	"net/http"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"

	testsweep "github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/sweep"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/getsafe"
)

const sweeperName = "hpe_morpheus_network_floating_ip"

// init registers a deliberately inert sweeper. Floating IPs carry no name to
// match on, and their only distinguishing status ("free") is shared with IPs
// allocated by other users and integrations: sweeping by status would release
// allocations this provider never created. So the "is this a test resource?"
// callback always returns false and nothing is ever released. The sweeper is
// still registered so the resource participates in the sweep framework and a
// safe identification rule can be added later if the API grows one.
func init() {
	testsweep.RegisterTypedAPISweeper(
		sweeperName,
		// List network floating IPs.
		func(ctx context.Context, client *sdk.APIClient) (
			[]sdk.GetAllNetworkFloatingIps200ResponseAllOfNetworkFloatingIpsInner,
			*http.Response,
			error,
		) {
			resp, hresp, err := client.NetworksAPI.GetAllNetworkFloatingIps(ctx).Execute()
			if resp == nil {
				return nil, hresp, err
			}

			return getsafe.Get(&resp.NetworkFloatingIps), hresp, err
		},
		// No floating IP is ever treated as a test resource; see the comment on
		// init for why sweeping by status would be unsafe.
		func(item sdk.GetAllNetworkFloatingIps200ResponseAllOfNetworkFloatingIpsInner) bool {
			return false
		},
		// Unreachable: the check above never returns true, so nothing is
		// released. Present only to satisfy the typed sweeper signature.
		func(
			ctx context.Context,
			client *sdk.APIClient,
			item sdk.GetAllNetworkFloatingIps200ResponseAllOfNetworkFloatingIpsInner,
		) (*http.Response, error) {
			return nil, nil
		},
		testsweep.WithIgnoreListStatuses[sdk.GetAllNetworkFloatingIps200ResponseAllOfNetworkFloatingIpsInner](
			http.StatusNotFound,
			http.StatusForbidden,
		),
	)
}
