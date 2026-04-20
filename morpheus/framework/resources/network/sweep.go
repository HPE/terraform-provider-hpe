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
	testhelpers.RegisterAPISweeper(
		"hpe_morpheus_network",
		testNetworkPrefix,
		func(ctx context.Context, client *sdk.APIClient, prefix string) (any, *http.Response, error) {
			return client.NetworksAPI.ListNetworks(ctx).Phrase(prefix).Execute()
		},
		"GetNetworks",
		func(ctx context.Context, client *sdk.APIClient, id int64, _ any) (*http.Response, error) {
			_, hresp, err := client.NetworksAPI.DeleteNetwork(ctx, id).Execute()
			return hresp, err
		},
		testhelpers.WithFilter(func(_ context.Context, _ *sdk.APIClient, item any) (bool, string, error) {
			network := item.(sdk.Network)
			labels, ok := network.GetLabelsOk()
			if !ok || !hasRequiredLabels(labels) {
				return false, "labels", nil
			}

			return true, "", nil
		}),
	)
}
