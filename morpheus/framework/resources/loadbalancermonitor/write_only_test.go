// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package loadbalancermonitor_test

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

// TestAccMorpheusLoadBalancerMonitorResourceWriteOnlyOk: changes to
// monitor_password_wo are detected by digest (MORPH-18005); the deprecated
// marker still works (and no longer replaces the monitor); both import forms
// -- on this resource's composite "lb.monitor" id -- leave a clean plan.
func TestAccMorpheusLoadBalancerMonitorResourceWriteOnlyOk(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.NSXT)
	testhelpers.SkipUnlessAcc(t)

	name := acctest.RandomWithPrefix(t.Name())
	lbName := name[0:min(32, len(name))]

	lbConfig, err := loadbalancer.RenderLoadBalancerNsxtConfig(t, map[string]string{"Name": lbName})
	if err != nil {
		t.Fatal(err)
	}

	addr := "hpe_morpheus_load_balancer_monitor.wo"

	c := testhelpers.WriteOnlyCase{
		Address:   addr,
		Attribute: "monitor_password_wo",
		Marker:    "monitor_password_wo_version",
		APIPath:   "/monitors",
		APIKey:    "loadBalancerMonitor.monitorPassword",
		Config: func(value, extra string, unrelated bool) string {
			interval := "5"
			if unrelated {
				interval = "10"
			}

			return testhelpers.ProviderBlock() + lbConfig + `
resource "hpe_morpheus_load_balancer_monitor" "wo" {
  load_balancer_id    = hpe_morpheus_load_balancer.lb.id
  name                = "` + name + `"
  monitor_type        = "http"
  monitor_interval    = ` + interval + `
  monitor_timeout     = 3
  monitor_destination = "/health"
  monitor_username    = "probe"
  monitor_password_wo = "` + value + `"` + extra + `
}
`
		},
		ImportIDFunc: func(s *terraform.State) (string, error) {
			rs, ok := s.RootModule().Resources[addr]
			if !ok {
				return "", fmt.Errorf("%s not in state", addr)
			}

			return rs.Primary.Attributes["load_balancer_id"] + "." + rs.Primary.ID, nil
		},
	}

	rec := testhelpers.RecordRequests(t)

	steps := testhelpers.WriteOnlyMatrixSteps(c, rec, "Secret123!", "Secret456!", "Secret789!")
	steps = append(steps, testhelpers.WriteOnlyImportSteps(c, "Secret789!", "monitor_type", "extra_config", "config")...)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps:                    steps,
	})
}
