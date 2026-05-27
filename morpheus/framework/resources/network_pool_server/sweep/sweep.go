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
)

const sweeperName = "hpe_morpheus_network_pool_server"

func init() {
	testsweep.RegisterTypedAPISweeper(
		sweeperName,
		// List network pool server resources.
		func(ctx context.Context, client *sdk.APIClient) (
			[]sdk.ListNetworkPoolServers200ResponseAllOfNetworkPoolServersInner,
			*http.Response,
			error,
		) {
			resp, hresp, err := client.NetworksAPI.ListNetworkPoolServers(ctx).Execute()
			if resp == nil {
				return nil, hresp, err
			}

			return resp.GetNetworkPoolServers(), hresp, err
		},
		// Is this a test network pool server?
		func(item sdk.ListNetworkPoolServers200ResponseAllOfNetworkPoolServersInner) bool {
			name, ok := item.GetNameOk()
			if !ok || name == nil {
				return false
			}

			return strings.HasPrefix(*name, testsweep.TestResourcePrefix)
		},
		// Delete the test network pool server.
		func(
			ctx context.Context,
			client *sdk.APIClient,
			item sdk.ListNetworkPoolServers200ResponseAllOfNetworkPoolServersInner,
		) (*http.Response, error) {
			id, ok := item.GetIdOk()
			if !ok || id == nil {
				return nil, fmt.Errorf("could not get ID")
			}

			_, hresp, err := client.NetworksAPI.DeleteNetworkPoolServer(ctx, *id).Execute()

			return hresp, err
		},
		testsweep.WithIgnoreListStatuses[sdk.ListNetworkPoolServers200ResponseAllOfNetworkPoolServersInner](
			http.StatusNotFound,
			http.StatusForbidden,
		),
	)
}
