// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package loadbalancervirtualserver_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/HPE/terraform-provider-hpe/morpheus"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/loadbalancer"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/loadbalancervirtualserver"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
)

func TestAccMorpheusLoadBalancerVirtualServerNsxtExampleOk(t *testing.T) {
	t.Parallel()
	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow acceptance test in short mode")
	}

	lbName := acctest.RandomWithPrefix(t.Name())
	lbName = lbName[0:16] + lbName[len(lbName)-16:]

	lbConfig, err := loadbalancer.RenderLoadBalancerNsxtConfig(t, map[string]string{
		"Name": lbName,
	})
	if err != nil {
		t.Fatalf("failed to render lb config: %s", err)
	}

	vipName := acctest.RandomWithPrefix(t.Name())

	vsConfig, err := loadbalancervirtualserver.RenderLoadBalancerVirtualServerNsxtConfig(t, map[string]string{
		"LoadBalancerId":     "hpe_morpheus_load_balancer.lb.id",
		"VipName":            vipName,
		"Description":        "test nsxt virtual server",
		"VipAddress":         "10.0.0.200",
		"VipPort":            "443",
		"VipProtocol":        "https",
		"ApplicationProfile": "/infra/lb-app-profiles/default-http-lb-app-profile",
	})
	if err != nil {
		t.Fatalf("failed to render vs config: %s", err)
	}

	providerConfig := testhelpers.ProviderBlock()
	resourceName := "hpe_morpheus_load_balancer_virtual_server.nsxt"

	checks := resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttrSet(resourceName, "id"),
		resource.TestCheckResourceAttr(resourceName, "vip_name", vipName),
		resource.TestCheckResourceAttr(resourceName, "description", "test nsxt virtual server"),
		resource.TestCheckResourceAttr(resourceName, "vip_address", "10.0.0.200"),
		resource.TestCheckResourceAttr(resourceName, "vip_port", "443"),
		resource.TestCheckResourceAttr(resourceName, "vip_protocol", "https"),
		resource.TestCheckResourceAttr(resourceName, "config_nsxt.application_profile",
			"/infra/lb-app-profiles/default-http-lb-app-profile"),
		resource.TestCheckResourceAttrPair(resourceName, "load_balancer_id",
			"hpe_morpheus_load_balancer.lb", "id"),
	)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + lbConfig + vsConfig,
				Check:  checks,
			},
			{
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"ssl_server_cert"},
				ResourceName:            resourceName,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources[resourceName]
					if !ok {
						return "", fmt.Errorf("resource not found: %s", resourceName)
					}

					lbID := rs.Primary.Attributes["load_balancer_id"]
					id := rs.Primary.Attributes["id"]

					return lbID + "." + id, nil
				},
			},
		},
	})
}

func TestAccMorpheusLoadBalancerVirtualServerNsxtUpdateOk(t *testing.T) {
	t.Parallel()
	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow acceptance test in short mode")
	}

	lbName := acctest.RandomWithPrefix(t.Name())
	lbName = lbName[0:16] + lbName[len(lbName)-16:]

	lbConfig, err := loadbalancer.RenderLoadBalancerNsxtConfig(t, map[string]string{
		"Name": lbName,
	})
	if err != nil {
		t.Fatalf("failed to render lb config: %s", err)
	}

	vipName := acctest.RandomWithPrefix(t.Name())
	providerConfig := testhelpers.ProviderBlock()
	resourceName := "hpe_morpheus_load_balancer_virtual_server.nsxt_update"

	createConfig := providerConfig + lbConfig + `
resource "hpe_morpheus_load_balancer_virtual_server" "nsxt_update" {
  load_balancer_id = hpe_morpheus_load_balancer.lb.id
  vip_name         = "` + vipName + `"
  description      = "nsxt update test"
  vip_port         = 443
  vip_protocol     = "https"

  config_nsxt = {
    application_profile = "/infra/lb-app-profiles/default-http-lb-app-profile"
  }
}
`

	updatedVipName := vipName + "-upd"
	updateConfig := providerConfig + lbConfig + `
resource "hpe_morpheus_load_balancer_virtual_server" "nsxt_update" {
  load_balancer_id = hpe_morpheus_load_balancer.lb.id
  vip_name         = "` + updatedVipName + `"
  description      = "nsxt update test updated"
  vip_port         = 443
  vip_protocol     = "https"

  config_nsxt = {
    application_profile = "/infra/lb-app-profiles/default-http-lb-app-profile"
  }
}
`

	checkInPlaceUpdate := resource.ConfigPlanChecks{
		PreApply: []plancheck.PlanCheck{
			plancheck.ExpectResourceAction(
				resourceName,
				plancheck.ResourceActionUpdate,
			),
		},
	}

	createChecks := resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttrSet(resourceName, "id"),
		resource.TestCheckResourceAttr(resourceName, "vip_name", vipName),
		resource.TestCheckResourceAttr(resourceName, "description", "nsxt update test"),
		resource.TestCheckResourceAttr(resourceName, "vip_port", "443"),
		resource.TestCheckResourceAttr(resourceName, "vip_protocol", "https"),
		resource.TestCheckResourceAttr(resourceName, "config_nsxt.application_profile",
			"/infra/lb-app-profiles/default-http-lb-app-profile"),
	)

	updateChecks := resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttrSet(resourceName, "id"),
		resource.TestCheckResourceAttr(resourceName, "vip_name", updatedVipName),
		resource.TestCheckResourceAttr(resourceName, "description", "nsxt update test updated"),
		resource.TestCheckResourceAttr(resourceName, "config_nsxt.application_profile",
			"/infra/lb-app-profiles/default-http-lb-app-profile"),
	)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
		Steps: []resource.TestStep{
			{
				Config: createConfig,
				Check:  createChecks,
			},
			{
				Config:           updateConfig,
				ConfigPlanChecks: checkInPlaceUpdate,
				Check:            updateChecks,
			},
		},
	})
}

func TestAccMorpheusLoadBalancerVirtualServerNsxtConfigChangeRequiresReplace(t *testing.T) {
	t.Parallel()
	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow acceptance test in short mode")
	}

	lbName := acctest.RandomWithPrefix(t.Name())
	lbName = lbName[0:16] + lbName[len(lbName)-16:]

	lbConfig, err := loadbalancer.RenderLoadBalancerNsxtConfig(t, map[string]string{
		"Name": lbName,
	})
	if err != nil {
		t.Fatalf("failed to render lb config: %s", err)
	}

	vipName := acctest.RandomWithPrefix(t.Name())
	providerConfig := testhelpers.ProviderBlock()
	resourceName := "hpe_morpheus_load_balancer_virtual_server.nsxt_replace"

	createConfig := providerConfig + lbConfig + `
resource "hpe_morpheus_load_balancer_virtual_server" "nsxt_replace" {
  load_balancer_id = hpe_morpheus_load_balancer.lb.id
  vip_name         = "` + vipName + `"
  description      = "nsxt replace test"
  vip_port         = 443
  vip_protocol     = "https"

  config_nsxt = {
    application_profile = "/infra/lb-app-profiles/default-http-lb-app-profile"
  }
}
`

	replaceConfig := providerConfig + lbConfig + `
resource "hpe_morpheus_load_balancer_virtual_server" "nsxt_replace" {
  load_balancer_id = hpe_morpheus_load_balancer.lb.id
  vip_name         = "` + vipName + `"
  description      = "nsxt replace test"
  vip_port         = 443
  vip_protocol     = "https"

  config_nsxt = {
    application_profile = "/infra/lb-app-profiles/default-tcp-lb-app-profile"
  }
}
`

	checkReplace := resource.ConfigPlanChecks{
		PreApply: []plancheck.PlanCheck{
			plancheck.ExpectResourceAction(
				resourceName,
				plancheck.ResourceActionReplace,
			),
		},
	}

	var initialResourceID string

	createChecks := resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttrSet(resourceName, "id"),
		resource.TestCheckResourceAttr(resourceName, "vip_name", vipName),
		resource.TestCheckResourceAttr(resourceName, "config_nsxt.application_profile",
			"/infra/lb-app-profiles/default-http-lb-app-profile"),
		func(s *terraform.State) error {
			rs, ok := s.RootModule().Resources[resourceName]
			if !ok {
				return fmt.Errorf("resource not found: %s", resourceName)
			}

			initialResourceID = rs.Primary.ID

			return nil
		},
	)

	replaceChecks := resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttrSet(resourceName, "id"),
		resource.TestCheckResourceAttr(resourceName, "vip_name", vipName),
		resource.TestCheckResourceAttr(resourceName, "config_nsxt.application_profile",
			"/infra/lb-app-profiles/default-tcp-lb-app-profile"),
		func(s *terraform.State) error {
			rs, ok := s.RootModule().Resources[resourceName]
			if !ok {
				return fmt.Errorf("resource not found: %s", resourceName)
			}

			if rs.Primary.ID == initialResourceID {
				return fmt.Errorf(
					"expected resource ID to change due to config_nsxt change (RequiresReplace), "+
						"but ID remained the same: %s", rs.Primary.ID)
			}

			return nil
		},
	)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
		Steps: []resource.TestStep{
			{
				Config: createConfig,
				Check:  createChecks,
			},
			{
				Config:           replaceConfig,
				ConfigPlanChecks: checkReplace,
				Check:            replaceChecks,
			},
		},
	})
}
