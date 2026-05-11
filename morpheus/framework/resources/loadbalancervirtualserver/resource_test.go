// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package loadbalancervirtualserver_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/HPE/terraform-provider-hpe/morpheus"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/loadbalancer"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/loadbalancervirtualserver"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
)

func TestMain(m *testing.M) {
	code := m.Run()
	testhelpers.WriteMergedResults()
	os.Exit(code)
}

func TestAccMorpheusLoadBalancerVirtualServerExampleOk(t *testing.T) {
	t.Parallel()
	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow acceptance test in short mode")
	}

	lbName := acctest.RandomWithPrefix(t.Name())
	// Load balancer names are limited to 32 characters.
	lbName = lbName[0:16] + lbName[len(lbName)-16:]

	lbConfig, err := loadbalancer.RenderLoadBalancerHAProxyConfig(t, map[string]string{
		"Name": lbName,
	})
	if err != nil {
		t.Fatalf("failed to render lb config: %s", err)
	}

	vipName := acctest.RandomWithPrefix(t.Name())

	vsConfig, err := loadbalancervirtualserver.RenderLoadBalancerVirtualServerConfig(t, map[string]string{
		"LoadBalancerId": "hpe_morpheus_load_balancer.haproxy.id",
		"VipName":        vipName,
		"Description":    "test virtual server",
		"VipAddress":     "10.0.0.100",
		"VipPort":        "80",
		"VipProtocol":    "http",
	})
	if err != nil {
		t.Fatalf("failed to render vs config: %s", err)
	}

	providerConfig := testhelpers.ProviderBlock()

	resourceName := "hpe_morpheus_load_balancer_virtual_server.example"

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttrSet(resourceName, "id"),
		resource.TestCheckResourceAttr(resourceName, "vip_name", vipName),
		resource.TestCheckResourceAttr(resourceName, "description", "test virtual server"),
		resource.TestCheckResourceAttr(resourceName, "vip_address", "10.0.0.100"),
		resource.TestCheckResourceAttr(resourceName, "vip_port", "80"),
		resource.TestCheckResourceAttr(resourceName, "vip_protocol", "http"),
		resource.TestCheckResourceAttrPair(
			resourceName, "load_balancer_id",
			"hpe_morpheus_load_balancer.haproxy", "id",
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + lbConfig + vsConfig,
				Check:  checkFn,
			},
			{
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"config", "ssl_server_cert"},
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

func TestAccMorpheusLoadBalancerVirtualServerUpdateOk(t *testing.T) {
	t.Parallel()
	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow acceptance test in short mode")
	}

	lbName := acctest.RandomWithPrefix(t.Name())
	lbName = lbName[0:16] + lbName[len(lbName)-16:]

	lbConfig, err := loadbalancer.RenderLoadBalancerHAProxyConfig(t, map[string]string{
		"Name": lbName,
	})
	if err != nil {
		t.Fatalf("failed to render lb config: %s", err)
	}

	vipName := acctest.RandomWithPrefix(t.Name())
	providerConfig := testhelpers.ProviderBlock()
	resourceName := "hpe_morpheus_load_balancer_virtual_server.test"

	createConfig := providerConfig + lbConfig + `
resource "hpe_morpheus_load_balancer_virtual_server" "test" {
  load_balancer_id = hpe_morpheus_load_balancer.haproxy.id
  vip_name         = "` + vipName + `"
  description      = "initial description"
  vip_port         = 80
  vip_protocol     = "http"
}
`

	updatedVipName := vipName + "-updated"
	updateConfig := providerConfig + lbConfig + `
resource "hpe_morpheus_load_balancer_virtual_server" "test" {
  load_balancer_id = hpe_morpheus_load_balancer.haproxy.id
  vip_name         = "` + updatedVipName + `"
  description      = "updated description"
  vip_port         = 443
  vip_protocol     = "https"
}
`

	createChecks := resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttrSet(resourceName, "id"),
		resource.TestCheckResourceAttr(resourceName, "vip_name", vipName),
		resource.TestCheckResourceAttr(resourceName, "description", "initial description"),
		resource.TestCheckResourceAttr(resourceName, "vip_port", "80"),
		resource.TestCheckResourceAttr(resourceName, "vip_protocol", "http"),
	)

	updateChecks := resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttrSet(resourceName, "id"),
		resource.TestCheckResourceAttr(resourceName, "vip_name", updatedVipName),
		resource.TestCheckResourceAttr(resourceName, "description", "updated description"),
		resource.TestCheckResourceAttr(resourceName, "vip_port", "443"),
		resource.TestCheckResourceAttr(resourceName, "vip_protocol", "https"),
	)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
		Steps: []resource.TestStep{
			{
				Config: createConfig,
				Check:  createChecks,
			},
			{
				Config: updateConfig,
				Check:  updateChecks,
			},
		},
	})
}
