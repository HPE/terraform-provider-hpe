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

	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/errors"
)

// Clouds whose name begins with this string will be eligible for deletion
const testCloudPrefix = "TestAccMorpheusCloud"

// All of these labels must be present for the cloud to be deleted
var requiredCloudSweepLabels = []string{
	"terraform",
	"acctest",
	"hpe_morpheus_cloud",
	"sweepable",
}

func hasRequiredCloudLabels(labels []string) bool {
	if labels == nil {
		return false
	}

	labelMap := make(map[string]bool)
	for _, label := range labels {
		labelMap[label] = true
	}

	for _, requiredLabel := range requiredCloudSweepLabels {
		if !labelMap[requiredLabel] {
			return false
		}
	}

	return true
}

func Clouds() {
	resource.AddTestSweepers(
		"hpe_morpheus_cloud",
		&resource.Sweeper{
			Name: "hpe_morpheus_cloud",
			F:    testSweepMorpheusClouds,
		})
}

func testSweepMorpheusClouds(_ string) error {
	ctx := context.Background()

	client, err := NewSweepClient(ctx)
	if err != nil {
		return err
	}

	clouds, hresp, err := client.CloudsAPI.ListClouds(ctx).
		Phrase(testCloudPrefix).Execute()
	if err != nil {
		// Handle 404 and 403 as "no matches found" rather than an error
		if hresp != nil && (hresp.StatusCode == http.StatusNotFound ||
			hresp.StatusCode == http.StatusForbidden) {
			log.Printf(
				"[INFO] No clouds found matching prefix (status %d): %s",
				hresp.StatusCode, testCloudPrefix,
			)

			return nil
		}

		return fmt.Errorf("failed to list clouds: %s", errors.ErrMsg(err, hresp))
	}
	if hresp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to list clouds: %s", errors.ErrMsg(err, hresp))
	}

	cloudList := clouds.GetZones()
	var sweptCount int
	var sweepErrors []string

	for _, cloud := range cloudList {
		name, ok := cloud.GetNameOk()
		if !ok || name == nil {
			continue
		}

		if strings.HasPrefix(*name, testCloudPrefix) {
			labels, ok := cloud.GetLabelsOk()
			if !ok || !hasRequiredCloudLabels(labels) {
				log.Printf("[INFO] Skipping cloud (labels): %s", *name)

				continue
			}

			id, ok := cloud.GetIdOk()
			if !ok || id == nil {
				log.Printf("[INFO] Skipping cloud (id): %s", *name)

				continue
			}

			log.Printf("[INFO] Sweeping cloud: %s (id: %d)", *name, *id)

			// Delete the cloud
			_, hresp, err := client.CloudsAPI.RemoveClouds(ctx, *id).
				Force(true).Execute()
			if err != nil || hresp.StatusCode != http.StatusOK {
				errMsg := fmt.Sprintf(
					"failed to delete cloud %s (id: %d): %s",
					*name, *id, errors.ErrMsg(err, hresp),
				)
				log.Printf("[ERROR] %s", errMsg)
				sweepErrors = append(sweepErrors, errMsg)

				continue
			}

			sweptCount++
		} else {
			log.Printf("[INFO] Skipping cloud (name): %s", *name)
		}
	}

	log.Printf(
		"[INFO] Cloud sweep completed. Clouds swept: %d, errors: %d",
		sweptCount, len(sweepErrors),
	)

	if len(sweepErrors) > 0 {
		return fmt.Errorf("%s", strings.Join(sweepErrors, "\n"))
	}

	return nil
}
