// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package network

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
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

// networkSweeper handles cleanup of test networks
type networkSweeper struct {
	client *sdk.APIClient
}

// NewNetworkSweeper creates and registers a network sweeper
func NewNetworkSweeper(client *sdk.APIClient) {
	s := &networkSweeper{
		client: client,
	}

	resource.AddTestSweepers(
		"hpe_morpheus_network",
		&resource.Sweeper{
			Name: "hpe_morpheus_network",
			F:    s.Sweep,
		})
}

// Sweep cleans up test networks
func (s *networkSweeper) Sweep(_ string) error {
	ctx := context.Background()

	if s.client == nil {
		log.Printf("[INFO] No client provided, skipping network sweep")

		return nil
	}

	networks, hresp, err := s.client.NetworksAPI.ListNetworks(ctx).
		Phrase(testNetworkPrefix).Execute()
	if err != nil || hresp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to list networks: %s", errfmt.ErrMsg(err, hresp))
	}

	networkList := networks.GetNetworks()
	var sweptCount int
	var sweepErrors []string

	for _, network := range networkList {
		name, ok := network.GetNameOk()
		if !ok || name == nil {
			continue
		}

		if strings.HasPrefix(*name, testNetworkPrefix) {
			labels, ok := network.GetLabelsOk()
			if !ok || !hasRequiredLabels(labels) {
				log.Printf("[INFO] Skipping network (labels): %s", *name)

				continue
			}

			id, ok := network.GetIdOk()
			if !ok || id == nil {
				log.Printf("[INFO] Skipping network (id): %s", *name)

				continue
			}

			log.Printf("[INFO] Sweeping network: %s (id: %d)", *name, *id)

			// Delete the network
			_, hresp, err := s.client.NetworksAPI.DeleteNetwork(ctx, *id).Execute()
			if err != nil || hresp.StatusCode != http.StatusOK {
				errMsg := fmt.Sprintf(
					"failed to delete network %s (id: %d): %s",
					*name, *id, errfmt.ErrMsg(err, hresp),
				)
				log.Printf("[ERROR] %s", errMsg)
				sweepErrors = append(sweepErrors, errMsg)

				continue
			}

			sweptCount++
		} else {
			log.Printf("[INFO] Skipping network (name): %s", *name)
		}
	}

	log.Printf(
		"[INFO] Network sweep completed. Networks swept: %d, errors: %d",
		sweptCount, len(sweepErrors),
	)

	if len(sweepErrors) > 0 {
		return fmt.Errorf("%s", strings.Join(sweepErrors, "\n"))
	}

	return nil
}
