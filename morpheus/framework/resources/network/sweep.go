// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package network

import (
	"context"
	"net/http"

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
		testNetworkPrefix,
		func(ctx context.Context, client *sdk.APIClient, prefix string) ([]sdk.ListNetworks200ResponseAllOfNetworksInner, *http.Response, error) {
			resp, hresp, err := client.NetworksAPI.ListNetworks(ctx).Phrase(prefix).Execute()
			if resp == nil {
				return nil, hresp, err
			}

			return resp.GetNetworks(), hresp, err
		},
		func(item sdk.ListNetworks200ResponseAllOfNetworksInner) (string, bool) {
			name, ok := item.GetNameOk()
			if !ok || name == nil {
				return "", false
			}

			return *name, true
		},
		func(item sdk.ListNetworks200ResponseAllOfNetworksInner) (int64, bool) {
			id, ok := item.GetIdOk()
			if !ok || id == nil {
				return 0, false
			}

			return *id, true
		},
		func(ctx context.Context, client *sdk.APIClient, id int64, _ sdk.ListNetworks200ResponseAllOfNetworksInner) (*http.Response, error) {
			_, hresp, err := client.NetworksAPI.DeleteNetwork(ctx, id).Execute()
			return hresp, err
		},
		testhelpers.WithFilter(func(
			_ context.Context,
			_ *sdk.APIClient,
			network sdk.ListNetworks200ResponseAllOfNetworksInner,
		) (bool, string, error) {
			labels, ok := network.GetLabelsOk()
			if !ok || !hasRequiredLabels(labels) {
				return false, "labels", nil
			}

			return true, "", nil
		}),
	)
}
