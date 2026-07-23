// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

//go:build sweep

package sweep

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	testsweep "github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/sweep"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/getsafe"
)

const (
	sweeperName                 = "hpe_morpheus_load_balancer"
	orphanedInstanceSweeperName = "hpe_morpheus_load_balancer_orphaned_instances"
)

func init() {
	testsweep.RegisterTypedAPISweeper(
		sweeperName,
		// List load balancer resources.
		func(ctx context.Context, client *sdk.APIClient) (
			[]sdk.ListLoadBalancers200ResponseAllOfLoadBalancersInner,
			*http.Response,
			error,
		) {
			resp, hresp, err := client.LoadBalancersAPI.ListLoadBalancers(ctx).Execute()
			if resp == nil {
				return nil, hresp, err
			}

			return getsafe.Get(&resp.LoadBalancers), hresp, err
		},
		// Is this a test load balancer?
		func(item sdk.ListLoadBalancers200ResponseAllOfLoadBalancersInner) bool {
			name, ok := getsafe.GetOk(item.Name)
			if !ok || name == nil {
				return false
			}

			return strings.HasPrefix(*name, testsweep.TestResourcePrefix)
		},
		// Delete the test load balancer.
		func(
			ctx context.Context,
			client *sdk.APIClient,
			item sdk.ListLoadBalancers200ResponseAllOfLoadBalancersInner,
		) (*http.Response, error) {
			id, ok := getsafe.GetOk(item.Id)
			if !ok || id == nil {
				return nil, fmt.Errorf("could not get ID")
			}

			_, hresp, err := client.LoadBalancersAPI.DeleteLoadBalancer(ctx, *id).Execute()

			return hresp, err
		},
		testsweep.WithDependencies[sdk.ListLoadBalancers200ResponseAllOfLoadBalancersInner](
			"hpe_morpheus_load_balancer_monitor",
			"hpe_morpheus_load_balancer_virtual_server",
		),
	)

	// Register a second sweeper that runs after load balancers are deleted.
	// HAProxy load balancers provision instances on Docker clusters that are not
	// automatically removed when the load balancer is deleted. These managed
	// instances are hidden from the standard ListInstances API, so the Search
	// API is used to discover them.
	resource.AddTestSweepers(orphanedInstanceSweeperName, &resource.Sweeper{
		Name:         orphanedInstanceSweeperName,
		Dependencies: []string{sweeperName},
		F: func(system string) error {
			ctx := context.Background()

			client, err := testhelpers.NewClientForServer(ctx, system)
			if err != nil {
				log.Printf("[WARN] Cannot create sweep client for %q: %v", system, err)

				return nil
			}

			sweepOrphanedInstances(ctx, client)

			return nil
		},
	})
}

// sweepOrphanedInstances removes test-prefixed instances that were provisioned
// by HAProxy load balancers on Docker clusters. These instances represent the
// actual Docker containers and are not automatically cleaned up when the load
// balancer is deleted. The instances are hidden from the standard ListInstances
// API, so the Search API is used to discover them.
func sweepOrphanedInstances(ctx context.Context, client *sdk.APIClient) {
	instances, err := findOrphanedInstances(ctx, client)
	if err != nil {
		log.Printf("[ERROR] %v", err)

		return
	}

	var sweptCnt, errCnt int

	for i := range instances {
		swept, errored := sweepHitIfInstance(ctx, client, instances[i])
		sweptCnt += swept
		errCnt += errored
	}

	log.Printf(
		"[INFO] Instance sweep completed. Instances swept: %d, errors: %d",
		sweptCnt, errCnt,
	)
}

func findOrphanedInstances(
	ctx context.Context,
	client *sdk.APIClient,
) ([]testsweep.SearchHit, error) {
	hits, hresp, err := testsweep.SearchHits(ctx, client, testsweep.TestResourcePrefix)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to search for orphaned instances: %s", errfmt.ErrMsg(err, hresp),
		)
	}

	return hits, nil
}

// sweepHitIfInstance deletes a search hit if it is a test-prefixed Instance.
// Returns (1, 0) on success, (0, 1) on error, and (0, 0) if the hit is not
// an applicable instance.
func sweepHitIfInstance(
	ctx context.Context,
	client *sdk.APIClient,
	hit testsweep.SearchHit,
) (swept, errored int) {
	if hit.Type != "Instance" {
		return 0, 0
	}

	if !strings.HasPrefix(hit.Name, testsweep.TestResourcePrefix) {
		return 0, 0
	}

	id, err := strconv.ParseInt(hit.ID, 10, 64)
	if err != nil {
		log.Printf("[ERROR] Failed to parse instance ID %q: %v", hit.ID, err)

		return 0, 1
	}

	log.Printf("[INFO] Sweeping orphaned instance %q (id=%d)", hit.Name, id)

	_, delResp, delErr := client.InstancesAPI.
		DeleteInstance(ctx, id).
		Force("on").
		Execute()
	if delErr != nil {
		if delResp != nil && delResp.StatusCode == http.StatusNotFound {
			log.Printf("[INFO] Instance %q already removed (404)", hit.Name)

			return 0, 0
		}

		log.Printf("[ERROR] Failed to delete instance %q (id=%d): %s",
			hit.Name, id, errfmt.ErrMsg(delErr, delResp))

		return 0, 1
	}

	return 1, 0
}
