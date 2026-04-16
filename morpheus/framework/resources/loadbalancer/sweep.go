// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package loadbalancer

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

// testLoadBalancerPrefix is truncated to 16 characters to account for the
// 32-character name limit after RandomWithPrefix appends a random suffix.
const testLoadBalancerPrefix = "TestAccMorpheusL"

// loadBalancerSweeper handles cleanup of test load balancers
type loadBalancerSweeper struct {
	client *sdk.APIClient
}

// NewLoadBalancerSweeper creates and registers a load balancer sweeper
func NewLoadBalancerSweeper(client *sdk.APIClient) {
	s := &loadBalancerSweeper{
		client: client,
	}

	resource.AddTestSweepers(
		"hpe_morpheus_load_balancer",
		&resource.Sweeper{
			Name: "hpe_morpheus_load_balancer",
			F:    s.Sweep,
		})
}

// Sweep cleans up test load balancers
func (s *loadBalancerSweeper) Sweep(_ string) error {
	ctx := context.Background()

	if s.client == nil {
		log.Printf("[INFO] No client provided, skipping load balancer sweep")

		return nil
	}

	lbs, hresp, err := s.client.LoadBalancersAPI.ListLoadBalancers(ctx).
		Phrase(testLoadBalancerPrefix).Execute()
	if err != nil || hresp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to list load balancers: %s", errfmt.ErrMsg(err, hresp))
	}

	lbList := lbs.GetLoadBalancers()
	var sweptCount int
	var sweepErrors []string

	for _, lb := range lbList {
		name, ok := lb.GetNameOk()
		if !ok || name == nil {
			continue
		}

		if strings.HasPrefix(*name, testLoadBalancerPrefix) {
			id, ok := lb.GetIdOk()
			if !ok || id == nil {
				log.Printf("[INFO] Skipping load balancer (id): %s", *name)

				continue
			}

			log.Printf("[INFO] Sweeping load balancer: %s (id: %d)", *name, *id)

			_, hresp, err := s.client.LoadBalancersAPI.DeleteLoadBalancer(ctx, *id).Execute()
			if err != nil || hresp.StatusCode != http.StatusOK {
				errMsg := fmt.Sprintf(
					"failed to delete load balancer %s (id: %d): %s",
					*name, *id, errfmt.ErrMsg(err, hresp),
				)
				log.Printf("[ERROR] %s", errMsg)
				sweepErrors = append(sweepErrors, errMsg)

				continue
			}

			sweptCount++
		} else {
			log.Printf("[INFO] Skipping load balancer (name): %s", *name)
		}
	}

	log.Printf(
		"[INFO] Load balancer sweep completed. Load balancers swept: %d, errors: %d",
		sweptCount, len(sweepErrors),
	)

	if len(sweepErrors) > 0 {
		return fmt.Errorf("%s", strings.Join(sweepErrors, "\n"))
	}

	return nil
}
