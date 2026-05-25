// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkrouter_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/networkrouter"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
)

const testResourceName = "hpe_morpheus_network_router.example"

func TestMain(m *testing.M) {
	code := m.Run()
	testhelpers.WriteMergedResults()
	os.Exit(code)
}

func TestAccMorpheusNetworkRouterResourceGenericExampleOk(t *testing.T) {
	if capabilities.Missing(t, capabilities.NetworkRouter) {
		t.Log("Skipping test due to missing capabilities")

		return
	}
	t.Parallel()
	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow acceptance test in short mode")
	}

	name := acctest.RandomWithPrefix(t.Name())

	config, err := networkrouter.RenderNetworkRouterGenericConfig(t, map[string]string{
		"Name":                 name,
		"TypeId":               "9", // tier 1 gateway
		"GroupId":              "3",
		"NetworkIntegrationId": "5",
	})
	if err != nil {
		t.Fatalf("failed to render config: %s", err)
	}

	resourceName := testResourceName
	config = testhelpers.ProviderBlock() + config

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "type_id", "9"), // tier 1
					resource.TestCheckResourceAttr(resourceName, "group_id", "3"),
					resource.TestCheckResourceAttr(resourceName, "network_integration_id", "5"),
					resource.TestCheckResourceAttr(resourceName, "config.ipManagementType", "dhcpLocal"),
				),
			},
		},
	})
}

func TestAccMorpheusNetworkRouterResourceNSXTGatewayTier0ExampleOk(t *testing.T) {
	if capabilities.Missing(t, capabilities.NSXT) {
		t.Log("Skipping test due to missing capabilities")

		return
	}
	t.Parallel()
	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow acceptance test in short mode")
	}

	name := acctest.RandomWithPrefix(t.Name())

	config, err := networkrouter.RenderNetworkRouterNSXTGatewayTier0Config(t, map[string]string{
		"Name":                 name,
		"GroupId":              "3",
		"NetworkIntegrationId": "5",
	})
	if err != nil {
		t.Fatalf("failed to render config: %s", err)
	}

	resourceName := testResourceName
	config = testhelpers.ProviderBlock() + config

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "type_id", "8"), // tier 0
					resource.TestCheckResourceAttr(resourceName, "group_id", "3"),
					resource.TestCheckResourceAttr(resourceName, "network_integration_id", "5"),
					resource.TestCheckResourceAttr(resourceName, "config_nsxt_gateway_tier0.ha_mode", "ACTIVE_ACTIVE"),
					resource.TestCheckResourceAttr(resourceName, "config_nsxt_gateway_tier0.restart_mode", "HELPER_ONLY"),
				),
			},
		},
	})
}

func TestAccMorpheusNetworkRouterResourceNSXTGatewayTier1ExampleOk(t *testing.T) {
	if capabilities.Missing(t, capabilities.NSXT) {
		t.Log("Skipping test due to missing capabilities")

		return
	}
	t.Parallel()
	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow acceptance test in short mode")
	}

	name := acctest.RandomWithPrefix(t.Name())

	config, err := networkrouter.RenderNetworkRouterNSXTGatewayTier1Config(t, map[string]string{
		"Name":                 name,
		"GroupId":              "3",
		"NetworkIntegrationId": "5",
	})
	if err != nil {
		t.Fatalf("failed to render config: %s", err)
	}

	resourceName := testResourceName
	config = testhelpers.ProviderBlock() + config

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "type_id", "9"), // tier 1
					resource.TestCheckResourceAttr(resourceName, "group_id", "3"),
					resource.TestCheckResourceAttr(resourceName, "network_integration_id", "5"),
					resource.TestCheckResourceAttr(resourceName, "config_nsxt_gateway_tier1.ip_management_type", "dhcpLocal"),
				),
			},
		},
	})
}
