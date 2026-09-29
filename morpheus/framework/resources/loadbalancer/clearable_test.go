// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package loadbalancer_test

import (
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

// Removing the description clears it on the load balancer (MORPH-17990); it
// is Optional+Computed and cleared once an apply has recorded it. The load
// balancer is an NSX-T one on the fixture's tier-0.
func TestAccMorpheusLoadBalancerRemovedAttributesAreCleared(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.NSXT)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	name := acctest.RandomWithPrefix(t.Name())
	name = name[0:min(32, len(name))]
	const rn = "hpe_morpheus_load_balancer.lb"

	// The rendered load balancer has no description; the first step adds one
	// at the top of its block, the second renders the block without it.
	withoutDescription, err := loadbalancer.RenderLoadBalancerNsxtConfig(t, map[string]string{"Name": name})
	if err != nil {
		t.Fatalf("failed to render lb config: %s", err)
	}

	const open = `resource "hpe_morpheus_load_balancer" "lb" {` + "\n"
	if !strings.Contains(withoutDescription, open) {
		t.Fatalf("rendered config has no %q block", open)
	}

	described := strings.Replace(withoutDescription, open, open+`  description = "to be cleared"`+"\n", 1)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config: testhelpers.ProviderBlock() + described,
				Check:  resource.TestCheckResourceAttr(rn, "description", "to be cleared"),
			},
			{
				Config: testhelpers.ProviderBlock() + withoutDescription,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(rn, plancheck.ResourceActionUpdate),
						plancheck.ExpectKnownValue(rn, tfjsonpath.New("description"), knownvalue.Null()),
					},
				},
				Check: resource.TestCheckNoResourceAttr(rn, "description"),
			},
			{
				Config:   testhelpers.ProviderBlock() + withoutDescription,
				PlanOnly: true,
			},
		},
	})
}
