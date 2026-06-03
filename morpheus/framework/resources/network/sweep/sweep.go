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

const sweeperName = "hpe_morpheus_network"

func init() {
	testsweep.RegisterTypedAPISweeper(
		sweeperName,
		// List network resources.
		func(ctx context.Context, client *sdk.APIClient) (
			[]sdk.ListNetworks200ResponseAllOfNetworksInner,
			*http.Response,
			error,
		) {
			resp, hresp, err := client.NetworksAPI.ListNetworks(ctx).Execute()
			if resp == nil {
				return nil, hresp, err
			}

			return getsafe.GetSafe(&resp.Networks), hresp, err
		},
		// Is this a test network?
		func(item sdk.ListNetworks200ResponseAllOfNetworksInner) bool {
			name, ok := getsafe.GetSafeOk(item.Name)
			if !ok || name == nil {
				return false
			}

			return strings.HasPrefix(*name, testsweep.TestResourcePrefix)
		},
		// Delete the test network.
		func(
			ctx context.Context,
			client *sdk.APIClient,
			item sdk.ListNetworks200ResponseAllOfNetworksInner,
		) (*http.Response, error) {
			id, ok := getsafe.GetSafeOk(item.Id)
			if !ok || id == nil {
				return nil, fmt.Errorf("could not get ID")
			}

			_, hresp, err := client.NetworksAPI.DeleteNetwork(ctx, *id).Execute()

			return hresp, err
		},
	)
}
