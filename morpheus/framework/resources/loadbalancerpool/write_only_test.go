// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package loadbalancerpool_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/loadbalancer"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
	"github.com/HPE/terraform-provider-hpe/provider/adapter"
)

// TestAccMorpheusLoadBalancerPoolResourceWriteOnlyOk: changes to the
// write-only partition are detected by digest (MORPH-18005): an in-place
// update that sends it; the resource has no marker; both import forms -- on
// the composite "lb/pool" id -- leave a clean plan. Creates an NSX-T load
// balancer: do not run against an appliance where deleting one NSX-T load
// balancer is known to delete others.
func TestAccMorpheusLoadBalancerPoolResourceWriteOnlyOk(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.NSXT)
	testhelpers.SkipUnlessAcc(t)

	name := acctest.RandomWithPrefix(t.Name())
	lbName := name[0:min(32, len(name))]

	lbConfig, err := loadbalancer.RenderLoadBalancerNsxtConfig(t, map[string]string{"Name": lbName})
	if err != nil {
		t.Fatal(err)
	}

	addr := "hpe_morpheus_load_balancer_pool.wo"

	c := testhelpers.WriteOnlyCase{
		Address:   addr,
		Attribute: "partition",
		Marker:    "", // none
		APIPath:   "/pools",
		APIKey:    "loadBalancerPool.partition",
		Config: func(value, _ string, unrelated bool) string {
			desc := "write-only matrix"
			if unrelated {
				desc = "write-only matrix, changed"
			}

			return testhelpers.ProviderBlock() + lbConfig + `
resource "hpe_morpheus_load_balancer_pool" "wo" {
  load_balancer_id = hpe_morpheus_load_balancer.lb.id
  name             = "` + name + `"
  description      = "` + desc + `"
  vip_balance      = "roundrobin"
  min_active       = 1
  partition        = "` + value + `"

  config_nsxt = {
    snat_translation_type   = "LBSnatAutoMap"
    tcp_multiplexing        = false
    tcp_multiplexing_number = 6
  }
}
`
		},
		ImportIDFunc: func(s *terraform.State) (string, error) {
			rs, ok := s.RootModule().Resources[addr]
			if !ok {
				return "", fmt.Errorf("%s not in state", addr)
			}

			return rs.Primary.Attributes["load_balancer_id"] + "/" + rs.Primary.ID, nil
		},
	}

	rec := testhelpers.RecordRequests(t)

	// No marker attribute: the value-only matrix.
	steps := testhelpers.WriteOnlyValueSteps(c, rec, "Common", "Partition2")
	steps = append(steps, testhelpers.WriteOnlyImportSteps(c, "Partition2", "config_nsxt")...)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps:                    steps,
	})
}
