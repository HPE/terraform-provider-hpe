// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package loadbalancerprofile_test

import (
	"fmt"
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

// Removing the description clears it on the profile (MORPH-17990). It is
// Optional-only: Terraform plans the null, and Update now sends it.
func TestAccMorpheusLoadBalancerProfileRemovedAttributesAreCleared(t *testing.T) {
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
	const rn = "hpe_morpheus_load_balancer_profile.test"

	config := func(extra string) string {
		return testhelpers.ProviderBlock() + lbConfig + fmt.Sprintf(`
resource "hpe_morpheus_load_balancer_profile" "test" {
  load_balancer_id = hpe_morpheus_load_balancer.lb.id
  name             = %q
  service_type     = "LBFastTcpProfile"
%s

  config_fast_tcp = {
    fast_tcp_idle_timeout    = 1800
    connection_close_timeout = 8
  }
}
`, name, extra)
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config: config(`  description      = "to be cleared"`),
				Check:  resource.TestCheckResourceAttr(rn, "description", "to be cleared"),
			},
			{
				Config: config(""),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(rn, plancheck.ResourceActionUpdate),
						plancheck.ExpectKnownValue(rn, tfjsonpath.New("description"), knownvalue.Null()),
					},
				},
				Check: resource.TestCheckNoResourceAttr(rn, "description"),
			},
			{
				Config:   config(""),
				PlanOnly: true,
			},
		},
	})
}
