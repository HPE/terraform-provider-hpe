// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package sweep

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"

	testsweep "github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/sweep"
)

// Networks whose name begins with this string will be eligible for deletion
const testResourcePrefix = "TestAccMorpheus"

func init() {
	testsweep.RegisterTypedAPISweeper(
		"hpe_morpheus_network",
		// List all network resources.
		func(ctx context.Context, client *sdk.APIClient) (
			[]sdk.ListNetworks200ResponseAllOfNetworksInner,
			*http.Response,
			error,
		) {
			resp, hresp, err := client.NetworksAPI.ListNetworks(ctx).Execute()
			if resp == nil {
				return nil, hresp, err
			}

			return resp.GetNetworks(), hresp, err
		},
		// Is this a test network?
		func(item sdk.ListNetworks200ResponseAllOfNetworksInner) bool {
			name, ok := item.GetNameOk()
			if !ok || name == nil {
				return false
			}

			return strings.HasPrefix(*name, testResourcePrefix)
		},
		// Delete the test network.
		func(
			ctx context.Context,
			client *sdk.APIClient,
			item sdk.ListNetworks200ResponseAllOfNetworksInner,
		) (*http.Response, error) {
			id, ok := item.GetIdOk()
			if !ok || id == nil {
				return nil, fmt.Errorf("could not get ID")
			}

			_, hresp, err := client.NetworksAPI.DeleteNetwork(ctx, *id).Execute()

			return hresp, err
		},
	)
}
