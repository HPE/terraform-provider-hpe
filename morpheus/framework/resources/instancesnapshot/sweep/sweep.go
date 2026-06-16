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

const sweeperName = "hpe_morpheus_instance_snapshot"

func init() {
	resource.AddTestSweepers(
		sweeperName,
		&resource.Sweeper{
			Name: sweeperName,
			F:    sweepSnapshots,
		},
	)
}

func sweepSnapshots(system string) error {
	ctx := context.Background()

	client, err := testhelpers.NewClientForServer(ctx, system)
	if err != nil {
		log.Printf("[WARN] Cannot create sweep client for %q: %v", system, err)

		return nil
	}

	// List all instances, then for each instance list snapshots,
	// and delete any that match the test prefix.
	instancesResp, hresp, err := client.InstancesAPI.ListInstances(ctx).Execute()
	if err != nil {
		if hresp != nil && (hresp.StatusCode == http.StatusNotFound ||
			hresp.StatusCode == http.StatusForbidden) {
			log.Printf("[INFO] No instances accessible for snapshot sweep")

			return nil
		}

		return fmt.Errorf("failed to list instances for snapshot sweep: %v", err)
	}

	if instancesResp == nil {
		return nil
	}

	var sweepErr error
	sweptCount := 0

	for _, inst := range instancesResp.Instances {
		if inst.Id == nil {
			continue
		}

		instanceID := *inst.Id

		snapsResp, hresp, err := client.InstancesAPI.SnapshotsInstance(ctx, instanceID).Execute()
		if err != nil {
			if hresp != nil && (hresp.StatusCode == http.StatusNotFound ||
				hresp.StatusCode == http.StatusForbidden) {
				continue
			}
			log.Printf("[WARN] Failed to list snapshots for instance %d: %v", instanceID, err)

			continue
		}

		if snapsResp == nil {
			continue
		}

		for _, snap := range snapsResp.Snapshots {
			if snap.Id == nil || snap.Name == nil {
				continue
			}

			if !strings.HasPrefix(*snap.Name, testsweep.TestResourcePrefix) {
				continue
			}

			snapshotID := *snap.Id
			log.Printf("[INFO] Sweeping snapshot %d (%s) on instance %d",
				snapshotID, *snap.Name, instanceID)

			_, hresp, err := client.InstancesAPI.DeleteSnapshotInstance(ctx, snapshotID).Execute()
			if err != nil {
				if hresp != nil && hresp.StatusCode == http.StatusNotFound {
					continue
				}
				log.Printf("[ERROR] Failed to delete snapshot %d: %v", snapshotID, err)
				sweepErr = fmt.Errorf("failed to delete snapshot %d: %v", snapshotID, err)

				continue
			}

			sweptCount++
		}
	}

	log.Printf("[INFO] Snapshot sweep completed. Resources swept: %d", sweptCount)

	return sweepErr
}
