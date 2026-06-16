// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package loadbalancervirtualserver_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/HPE/terraform-provider-hpe/morpheus"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/loadbalancer"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/loadbalancervirtualserver"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
)

func TestMain(m *testing.M) {
	code := testhelpers.TestMain(m)
	testhelpers.WriteMergedResults()
	os.Exit(code)
}

func TestAccMorpheusLoadBalancerVirtualServerResourceNsxtExampleOk(t *testing.T) {
	if capabilities.Missing(t, capabilities.NSXT) {
		t.Log("Skipping test due to missing capabilities")

		return
	}
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

	// Use port 443 with protocol "http" intentionally: this verifies the API
	// stores the configured port as-is rather than normalizing based on protocol.
	vsConfig, err := loadbalancervirtualserver.RenderLoadBalancerVirtualServerNsxtFullConfig(t, map[string]string{
		"LoadBalancerId":     "hpe_morpheus_load_balancer.lb.id",
		"VipName":            vipName,
		"Description":        "test nsxt virtual server",
		"VipAddress":         "10.0.0.5",
		"VipPort":            "443",
		"VipProtocol":        "http",
		"PoolId":             "11",
		"SslCert":            "12",
		"SslServerCert":      "0",
		"ApplicationProfile": "13",
		"Persistence":        "COOKIE",
		"PersistenceProfile": "16",
		"SslClientProfile":   "19",
		"SslServerProfile":   "",
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
		resource.TestCheckResourceAttr(resourceName, "vip_address", "10.0.0.5"),
		resource.TestCheckResourceAttr(resourceName, "vip_port", "443"),
		resource.TestCheckResourceAttr(resourceName, "vip_protocol", "http"),
		resource.TestCheckResourceAttr(resourceName, "pool_id", "11"),
		resource.TestCheckResourceAttr(resourceName, "ssl_cert", "12"),
		resource.TestCheckResourceAttr(resourceName, "ssl_server_cert", "0"),
		resource.TestCheckResourceAttr(resourceName, "config_nsxt.application_profile", "13"),
		resource.TestCheckResourceAttr(resourceName, "config_nsxt.persistence", "COOKIE"),
		resource.TestCheckResourceAttr(resourceName, "config_nsxt.persistence_profile", "16"),
		resource.TestCheckResourceAttr(resourceName, "config_nsxt.ssl_client_profile", "19"),
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
				ImportStateVerifyIgnore: []string{"vip_pool"},
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

func TestAccMorpheusLoadBalancerVirtualServerResourceNsxtUpdateOk(t *testing.T) {
	if capabilities.Missing(t, capabilities.NSXT) {
		t.Log("Skipping test due to missing capabilities")

		return
	}
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
  vip_address      = "10.0.0.201"
  vip_port         = 443
  vip_protocol     = "http"
  pool_id          = 42
  ssl_cert         = 12
  ssl_server_cert  = 0

  config_nsxt = {
    application_profile = 85
    persistence         = "SOURCE_IP"
    persistence_profile = 78
    ssl_client_profile  = 33
  }
}
`

	updatedVipName := vipName + "-upd"
	updateConfig := providerConfig + lbConfig + `
resource "hpe_morpheus_load_balancer_virtual_server" "nsxt_update" {
  load_balancer_id = hpe_morpheus_load_balancer.lb.id
  vip_name         = "` + updatedVipName + `"
  description      = "nsxt update test updated"
  vip_address      = "10.0.0.201"
  vip_port         = 443
  vip_protocol     = "http"
  pool_id          = 42
  ssl_cert         = 12
  ssl_server_cert  = 0

  config_nsxt = {
    application_profile = 85
    persistence         = "COOKIE"
    persistence_profile = 79
    ssl_client_profile  = 33
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
		resource.TestCheckResourceAttr(resourceName, "vip_protocol", "http"),
		resource.TestCheckResourceAttr(resourceName, "pool_id", "42"),
		resource.TestCheckResourceAttr(resourceName, "ssl_cert", "12"),
		resource.TestCheckResourceAttr(resourceName, "config_nsxt.application_profile", "85"),
		resource.TestCheckResourceAttr(resourceName, "config_nsxt.persistence", "SOURCE_IP"),
		resource.TestCheckResourceAttr(resourceName, "config_nsxt.persistence_profile", "78"),
		resource.TestCheckResourceAttr(resourceName, "config_nsxt.ssl_client_profile", "33"),
	)

	updateChecks := resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttrSet(resourceName, "id"),
		resource.TestCheckResourceAttr(resourceName, "vip_name", updatedVipName),
		resource.TestCheckResourceAttr(resourceName, "description", "nsxt update test updated"),
		resource.TestCheckResourceAttr(resourceName, "config_nsxt.application_profile", "85"),
		resource.TestCheckResourceAttr(resourceName, "config_nsxt.persistence", "COOKIE"),
		resource.TestCheckResourceAttr(resourceName, "config_nsxt.persistence_profile", "79"),
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

func TestAccMorpheusLoadBalancerVirtualServerResourceNsxtConfigChangeRequiresReplace(t *testing.T) {
	if capabilities.Missing(t, capabilities.NSXT) {
		t.Log("Skipping test due to missing capabilities")

		return
	}
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
  vip_address      = "10.0.0.202"
  vip_port         = 443
  vip_protocol     = "http"
  pool_id          = 42

  config_nsxt = {
    application_profile = 85
    persistence         = "SOURCE_IP"
    persistence_profile = 78
  }
}
`

	replaceConfig := providerConfig + lbConfig + `
resource "hpe_morpheus_load_balancer_virtual_server" "nsxt_replace" {
  load_balancer_id = hpe_morpheus_load_balancer.lb.id
  vip_name         = "` + vipName + `"
  description      = "nsxt replace test"
  vip_address      = "10.0.0.202"
  vip_port         = 443
  vip_protocol     = "http"
  pool_id          = 42

  config_nsxt = {
    application_profile = 90
    persistence         = ""
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
		resource.TestCheckResourceAttr(resourceName, "config_nsxt.application_profile", "85"),
		resource.TestCheckResourceAttr(resourceName, "config_nsxt.persistence", "SOURCE_IP"),
		resource.TestCheckResourceAttr(resourceName, "config_nsxt.persistence_profile", "78"),
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
		resource.TestCheckResourceAttr(resourceName, "config_nsxt.application_profile", "90"),
		resource.TestCheckResourceAttr(resourceName, "config_nsxt.persistence", ""),
		func(s *terraform.State) error {
			rs, ok := s.RootModule().Resources[resourceName]
			if !ok {
				return fmt.Errorf("resource not found: %s", resourceName)
			}

			if rs.Primary.ID == initialResourceID {
				return fmt.Errorf(
					"expected resource ID to change due to config_nsxt change (RequiresReplace), "+
						"but ID remained the same: %s", rs.Primary.ID,
				)
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

func TestAccMorpheusLoadBalancerVirtualServerResourceNsxtMinimalExampleOk(t *testing.T) {
	if capabilities.Missing(t, capabilities.NSXT) {
		t.Log("Skipping test due to missing capabilities")

		return
	}
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

	// Use port 443 with protocol "http" intentionally: this verifies the API
	// stores the configured port as-is rather than normalizing based on protocol.
	vsConfig, err := loadbalancervirtualserver.RenderLoadBalancerVirtualServerNsxtMinimalConfig(t, map[string]string{
		"LoadBalancerId":     "hpe_morpheus_load_balancer.lb.id",
		"VipName":            vipName,
		"Description":        "test minimal nsxt virtual server",
		"VipAddress":         "10.0.0.4",
		"VipPort":            "443",
		"VipProtocol":        "http",
		"PoolId":             "11",
		"ApplicationProfile": "13",
	})
	if err != nil {
		t.Fatalf("failed to render vs config: %s", err)
	}

	providerConfig := testhelpers.ProviderBlock()
	resourceName := "hpe_morpheus_load_balancer_virtual_server.nsxt_minimal"

	checks := resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttrSet(resourceName, "id"),
		resource.TestCheckResourceAttr(resourceName, "vip_name", vipName),
		resource.TestCheckResourceAttr(resourceName, "description", "test minimal nsxt virtual server"),
		resource.TestCheckResourceAttr(resourceName, "vip_address", "10.0.0.4"),
		resource.TestCheckResourceAttr(resourceName, "vip_port", "443"),
		resource.TestCheckResourceAttr(resourceName, "vip_protocol", "http"),
		resource.TestCheckResourceAttr(resourceName, "pool_id", "11"),
		resource.TestCheckResourceAttr(resourceName, "config_nsxt.application_profile", "13"),
		resource.TestCheckResourceAttr(resourceName, "config_nsxt.persistence", ""),
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
				ImportStateVerifyIgnore: []string{"vip_pool"},
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
