// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package network

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
)

// Networks whose name begins with this string will be eligible for deletion
const testNetworkPrefix = "TestAccMorpheusNetworkResource"

// All of these labels must be present for the network to be deleted
var requiredSweepLabels = []string{
	"terraform",
	"acctest",
	"hpe_morpheus_network",
	"sweepable",
}

func hasRequiredLabels(labels []string) bool {
	if labels == nil {
		return false
	}

	labelMap := make(map[string]bool)
	for _, label := range labels {
		labelMap[label] = true
	}

	for _, requiredLabel := range requiredSweepLabels {
		if !labelMap[requiredLabel] {
			return false
		}
	}

	return true
}

func init() {
	testhelpers.RegisterTypedAPISweeper(
		"hpe_morpheus_network",
		// List candidate networks for sweeping.
		func(ctx context.Context, client *sdk.APIClient, _ string) (
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
		// Match only acceptance-test networks with the required labels.
		func(item sdk.ListNetworks200ResponseAllOfNetworksInner) bool {
			name, ok := item.GetNameOk()
			if !ok || name == nil || !strings.HasPrefix(*name, testNetworkPrefix) {
				return false
			}

			labels, ok := item.GetLabelsOk()
			if !ok || !hasRequiredLabels(labels) {
				return false
			}

			return true
		},
		// Delete a matched network.
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
