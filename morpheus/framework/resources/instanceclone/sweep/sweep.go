// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

//go:build sweep

package sweep

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	testsweep "github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/sweep"
)

const sweeperName = "hpe_morpheus_instance_clone"

func init() {
	resource.AddTestSweepers(
		sweeperName,
		&resource.Sweeper{
			Name: sweeperName,
			F:    sweepClones,
		},
	)
}

func sweepClones(system string) error {
	ctx := context.Background()

	client, err := testhelpers.NewClientForServer(ctx, system)
	if err != nil {
		log.Printf("[WARN] Cannot create sweep client for %q: %v", system, err)

		return nil
	}

	instancesResp, hresp, err := client.InstancesAPI.ListInstances(ctx).Execute()
	if err != nil {
		if hresp != nil && (hresp.StatusCode == http.StatusNotFound ||
			hresp.StatusCode == http.StatusForbidden) {
			log.Printf("[INFO] No instances accessible for clone sweep")

			return nil
		}

		return fmt.Errorf("failed to list instances for clone sweep: %v", err)
	}

	if instancesResp == nil {
		return nil
	}

	var sweepErr error
	sweptCount := 0

	for _, inst := range instancesResp.Instances {
		if inst.Id == nil || inst.Name == nil {
			continue
		}

		if !strings.HasPrefix(*inst.Name, testsweep.TestResourcePrefix) {
			continue
		}

		instanceID := *inst.Id
		log.Printf("[INFO] Sweeping clone instance %d (%s)", instanceID, *inst.Name)

		_, hresp, err := client.InstancesAPI.DeleteInstance(ctx, instanceID).Execute()
		if err != nil {
			if hresp != nil && hresp.StatusCode == http.StatusNotFound {
				continue
			}
			log.Printf("[ERROR] Failed to delete clone instance %d: %v", instanceID, err)
			sweepErr = fmt.Errorf("failed to delete clone instance %d: %v", instanceID, err)

			continue
		}

		sweptCount++
	}

	log.Printf("[INFO] Instance clone sweep completed. Resources swept: %d", sweptCount)

	return sweepErr
}
