// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package loadbalancervirtualserver_test

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

// Removing the description or the host name clears it on the virtual server
// (MORPH-17990). Both are Optional-only: Terraform plans the null, and Update
// now sends it. The virtual server, its pool and its application profile are
// built on a load balancer the test creates on the NSX-T fixture.
func TestAccMorpheusLoadBalancerVirtualServerRemovedAttributesAreCleared(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.NSXT)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	lbName := acctest.RandomWithPrefix(t.Name())
	lbName = lbName[0:16] + lbName[len(lbName)-16:]

	lbConfig, err := loadbalancer.RenderLoadBalancerNsxtConfig(t, map[string]string{"Name": lbName})
	if err != nil {
		t.Fatalf("failed to render lb config: %s", err)
	}

	name := acctest.RandomWithPrefix(t.Name())
	const rn = "hpe_morpheus_load_balancer_virtual_server.test"

	config := func(with ...string) string {
		set := map[string]string{
			"description":  `description  = "to be cleared"`,
			"vip_hostname": `vip_hostname = "vip.example.test"`,
		}

		lines := make([]string, 0, len(with))
		for _, a := range with {
			lines = append(lines, "  "+set[a])
		}

		return testhelpers.ProviderBlock() + lbConfig + fmt.Sprintf(`
resource "hpe_morpheus_load_balancer_pool" "test" {
  load_balancer_id = hpe_morpheus_load_balancer.lb.id
  name             = "%[1]s-pool"
  vip_balance      = "ROUND_ROBIN"
  min_active       = 1

  config_nsxt = {
    snat_translation_type   = "LBSnatAutoMap"
    tcp_multiplexing        = true
    tcp_multiplexing_number = 6
  }
}

resource "hpe_morpheus_load_balancer_profile" "test" {
  load_balancer_id = hpe_morpheus_load_balancer.lb.id
  name             = "%[1]s-profile"
  service_type     = "LBFastTcpProfile"

  config_fast_tcp = {
    fast_tcp_idle_timeout    = 1800
    connection_close_timeout = 8
  }
}

resource "hpe_morpheus_load_balancer_virtual_server" "test" {
  load_balancer_id = hpe_morpheus_load_balancer.lb.id
  vip_name         = %[1]q
  vip_address      = "10.0.0.5"
  vip_port         = 8080
  vip_protocol     = "tcp"
  pool_id          = hpe_morpheus_load_balancer_pool.test.id
%[2]s

  config_nsxt = {
    application_profile = hpe_morpheus_load_balancer_profile.test.id
    persistence         = ""
  }
}
`, name, strings.Join(lines, "\n"))
	}

	all := []string{"description", "vip_hostname"}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config: config(all...),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(rn, "description", "to be cleared"),
					resource.TestCheckResourceAttr(rn, "vip_hostname", "vip.example.test"),
				),
			},
			{
				Config: config(),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(rn, plancheck.ResourceActionUpdate),
						plancheck.ExpectKnownValue(rn, tfjsonpath.New("description"), knownvalue.Null()),
						plancheck.ExpectKnownValue(rn, tfjsonpath.New("vip_hostname"), knownvalue.Null()),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckNoResourceAttr(rn, "description"),
					resource.TestCheckNoResourceAttr(rn, "vip_hostname"),
				),
			},
			{
				Config:   config(),
				PlanOnly: true,
			},
		},
	})
}
