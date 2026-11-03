// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package loadbalancermonitor_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"

	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/loadbalancer"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
	"github.com/HPE/terraform-provider-hpe/provider/adapter"
)

// Removing an optional attribute clears it on the monitor (MORPH-17990): the
// description and the free-text probe settings, all Optional+Computed and
// cleared once an apply has recorded them. The monitor sits on a load
// balancer the test builds on the NSX-T fixture.
func TestAccMorpheusLoadBalancerMonitorRemovedAttributesAreCleared(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.NSXT)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	lbName := acctest.RandomWithPrefix(t.Name())
	lbName = lbName[0:min(32, len(lbName))]

	lbConfig, err := loadbalancer.RenderLoadBalancerNsxtConfig(t, map[string]string{"Name": lbName})
	if err != nil {
		t.Fatalf("failed to render lb config: %s", err)
	}

	name := acctest.RandomWithPrefix(t.Name())
	const rn = "hpe_morpheus_load_balancer_monitor.test"

	config := func(with ...string) string {
		set := map[string]string{
			"description":  `description  = "to be cleared"`,
			"send_data":    `send_data    = "GET / HTTP/1.1"`,
			"receive_code": `receive_code = "200"`,
		}

		lines := make([]string, 0, len(with))
		for _, a := range with {
			lines = append(lines, "  "+set[a])
		}

		return testhelpers.ProviderBlock() + lbConfig + fmt.Sprintf(`
resource "hpe_morpheus_load_balancer_monitor" "test" {
  load_balancer_id    = hpe_morpheus_load_balancer.lb.id
  name                = %q
  monitor_type        = "http"
  monitor_interval    = 5
  monitor_timeout     = 15
  monitor_destination = "/"
  fall_count          = 3
  rise_count          = 3
  send_type           = "GET"
  send_version        = "HTTP_VERSION_1_1"
%s
}
`, name, strings.Join(lines, "\n"))
	}

	all := []string{"description", "send_data", "receive_code"}

	nullChecks := []plancheck.PlanCheck{plancheck.ExpectResourceAction(rn, plancheck.ResourceActionUpdate)}
	var noAttrs []resource.TestCheckFunc
	for _, a := range all {
		nullChecks = append(nullChecks, plancheck.ExpectKnownValue(rn, tfjsonpath.New(a), knownvalue.Null()))
		noAttrs = append(noAttrs, resource.TestCheckNoResourceAttr(rn, a))
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config: config(all...),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(rn, "description", "to be cleared"),
					resource.TestCheckResourceAttr(rn, "send_data", "GET / HTTP/1.1"),
					resource.TestCheckResourceAttr(rn, "receive_code", "200"),
				),
			},
			{
				Config:           config(),
				ConfigPlanChecks: resource.ConfigPlanChecks{PreApply: nullChecks},
				Check:            resource.ComposeAggregateTestCheckFunc(noAttrs...),
			},
			{
				Config:   config(),
				PlanOnly: true,
			},
		},
	})
}
