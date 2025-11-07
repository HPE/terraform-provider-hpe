// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

// Package sweep allows deletion of dangling test resources
package sweep

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/errors"
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

func Networks() {
	resource.AddTestSweepers(
		"hpe_morpheus_network",
		&resource.Sweeper{
			Name: "hpe_morpheus_network",
			F:    testSweepMorpheusNetworks,
		})
}

func testSweepMorpheusNetworks(_ string) error {
	ctx := context.Background()

	client, err := NewSweepClient(ctx)
	if err != nil {
		return err
	}

	networks, hresp, err := client.NetworksAPI.ListNetworks(ctx).
		Phrase(testNetworkPrefix).Execute()
	if err != nil || hresp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to list networks: %s", errors.ErrMsg(err, hresp))
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
			_, hresp, err := client.NetworksAPI.DeleteNetwork(ctx, *id).Execute()
			if err != nil || hresp.StatusCode != http.StatusOK {
				errMsg := fmt.Sprintf(
					"failed to delete network %s (id: %d): %s",
					*name, *id, errors.ErrMsg(err, hresp),
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
