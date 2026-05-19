// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

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
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
)

func init() {
	resource.AddTestSweepers(
		"hpe_morpheus_load_balancer_monitor",
		&resource.Sweeper{
			Name: "hpe_morpheus_load_balancer_monitor",
			F:    sweepLoadBalancerMonitors,
		},
	)
}

func sweepLoadBalancerMonitors(system string) error {
	ctx := context.Background()

	client, err := testhelpers.NewClientForServer(ctx, system)
	if err != nil {
		log.Printf("[WARN] Cannot create sweep client for %q: %v", system, err)

		return nil
	}

	// List all load balancers to iterate their monitors.
	lbResp, hresp, err := client.LoadBalancersAPI.ListLoadBalancers(ctx).Execute()
	if err != nil || hresp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to list load balancers: %s", errfmt.ErrMsg(err, hresp))
	}

	for _, lb := range lbResp.GetLoadBalancers() {
		if lb.Id == nil {
			continue
		}

		lbID := *lb.Id

		monResp, hresp, err := client.LoadBalancersAPI.
			ListLoadBalancerMonitors(ctx, lbID).Execute()
		if err != nil || hresp.StatusCode != http.StatusOK {
			log.Printf("[WARN] Failed to list monitors for LB %d: %s",
				*lb.Id, errfmt.ErrMsg(err, hresp))

			continue
		}

		for _, mon := range monResp.GetLoadBalancerMonitors() {
			name, ok := mon.GetNameOk()
			if !ok || name == nil {
				continue
			}

			if !strings.HasPrefix(*name, testsweep.TestResourcePrefix) {
				continue
			}

			id, ok := mon.GetIdOk()
			if !ok || id == nil {
				continue
			}

			log.Printf("[INFO] Sweeping load balancer monitor %q (ID %d, LB %d)",
				*name, *id, *lb.Id)

			_, hresp, err := client.LoadBalancersAPI.
				DeleteLoadBalancerMonitor(ctx, lbID, *id).Execute()
			if err != nil {
				log.Printf("[ERROR] Failed to delete monitor %d: %s",
					*id, errfmt.ErrMsg(err, hresp))
			}
		}
	}

	return nil
}
