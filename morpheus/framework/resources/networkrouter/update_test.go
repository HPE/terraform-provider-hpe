// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkrouter_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/networkrouter"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
)

func TestAccMorpheusNetworkRouterGenericUpdateOk(t *testing.T) {
	if capabilities.Missing(t, capabilities.All, capabilities.NetworkRouter) {
		t.Log("Skipping test due to missing capabilities")

		return
	}
	t.Parallel()
	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow acceptance test in short mode")
	}

	name := acctest.RandomWithPrefix(t.Name())
	updatedName := name + "-updated"

	createConfig, err := networkrouter.RenderNetworkRouterGenericConfig(t, map[string]string{
		"Name":                 name,
		"TypeId":               "9",
		"GroupId":              "3",
		"NetworkIntegrationId": "5",
	})
	if err != nil {
		t.Fatalf("failed to render create config: %s", err)
	}

	updateConfig, err := networkrouter.RenderNetworkRouterGenericConfig(t, map[string]string{
		"Name":                 updatedName,
		"TypeId":               "9",
		"GroupId":              "3",
		"NetworkIntegrationId": "5",
	})
	if err != nil {
		t.Fatalf("failed to render update config: %s", err)
	}

	resourceName := testResourceName
	providerConfig := testhelpers.ProviderBlock()

	createChecks := resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttr(resourceName, "name", name),
		resource.TestCheckResourceAttr(resourceName, "type_id", "9"),
		resource.TestCheckResourceAttr(resourceName, "group_id", "3"),
		resource.TestCheckResourceAttr(resourceName, "network_integration_id", "5"),
		resource.TestCheckResourceAttr(resourceName, "config.ipManagementType", "dhcpLocal"),
	)

	updateChecks := resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttr(resourceName, "name", updatedName),
		resource.TestCheckResourceAttr(resourceName, "type_id", "9"),
		resource.TestCheckResourceAttr(resourceName, "group_id", "3"),
		resource.TestCheckResourceAttr(resourceName, "network_integration_id", "5"),
		resource.TestCheckResourceAttr(resourceName, "config.ipManagementType", "dhcpLocal"),
	)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + createConfig,
				Check:  createChecks,
			},
			{
				Config:             providerConfig + createConfig,
				Check:              createChecks,
				ExpectNonEmptyPlan: false,
				PlanOnly:           true,
			},
			{
				Config: providerConfig + updateConfig,
				Check:  updateChecks,
			},
			{
				Config:             providerConfig + updateConfig,
				Check:              updateChecks,
				ExpectNonEmptyPlan: false,
				PlanOnly:           true,
			},
		},
	})
}

func TestAccMorpheusNetworkRouterNSXTGatewayTier0UpdateOk(t *testing.T) {
	if capabilities.Missing(t, capabilities.NSXT, capabilities.NetworkRouter) {
		t.Log("Skipping test due to missing capabilities")

		return
	}
	t.Parallel()
	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow acceptance test in short mode")
	}

	name := acctest.RandomWithPrefix(t.Name())
	updatedName := name + "-updated"

	createConfig, err := networkrouter.RenderNetworkRouterNSXTGatewayTier0Config(t, map[string]string{
		"Name":                 name,
		"GroupId":              "3",
		"NetworkIntegrationId": "5",
	})
	if err != nil {
		t.Fatalf("failed to render create config: %s", err)
	}

	updateConfig, err := networkrouter.RenderNetworkRouterNSXTGatewayTier0Config(t, map[string]string{
		"Name":                 updatedName,
		"GroupId":              "3",
		"NetworkIntegrationId": "5",
	})
	if err != nil {
		t.Fatalf("failed to render update config: %s", err)
	}

	resourceName := testResourceName
	providerConfig := testhelpers.ProviderBlock()

	createChecks := resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttr(resourceName, "name", name),
		resource.TestCheckResourceAttr(resourceName, "type_id", "8"),
		resource.TestCheckResourceAttr(resourceName, "group_id", "3"),
		resource.TestCheckResourceAttr(resourceName, "network_integration_id", "5"),
		resource.TestCheckResourceAttr(resourceName, "config_nsxt_gateway_tier0.ha_mode", "ACTIVE_ACTIVE"),
		resource.TestCheckResourceAttr(resourceName, "config_nsxt_gateway_tier0.restart_mode", "HELPER_ONLY"),
	)

	updateChecks := resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttr(resourceName, "name", updatedName),
		resource.TestCheckResourceAttr(resourceName, "type_id", "8"),
		resource.TestCheckResourceAttr(resourceName, "group_id", "3"),
		resource.TestCheckResourceAttr(resourceName, "network_integration_id", "5"),
		resource.TestCheckResourceAttr(resourceName, "config_nsxt_gateway_tier0.ha_mode", "ACTIVE_ACTIVE"),
		resource.TestCheckResourceAttr(resourceName, "config_nsxt_gateway_tier0.restart_mode", "HELPER_ONLY"),
	)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + createConfig,
				Check:  createChecks,
			},
			{
				Config:             providerConfig + createConfig,
				Check:              createChecks,
				ExpectNonEmptyPlan: false,
				PlanOnly:           true,
			},
			{
				Config: providerConfig + updateConfig,
				Check:  updateChecks,
			},
			{
				Config:             providerConfig + updateConfig,
				Check:              updateChecks,
				ExpectNonEmptyPlan: false,
				PlanOnly:           true,
			},
		},
	})
}

func TestAccMorpheusNetworkRouterNSXTGatewayTier1UpdateOk(t *testing.T) {
	if capabilities.Missing(t, capabilities.NSXT, capabilities.NetworkRouter) {
		t.Log("Skipping test due to missing capabilities")

		return
	}
	t.Parallel()
	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow acceptance test in short mode")
	}

	name := acctest.RandomWithPrefix(t.Name())
	updatedName := name + "-updated"

	createConfig, err := networkrouter.RenderNetworkRouterNSXTGatewayTier1Config(t, map[string]string{
		"Name":                 name,
		"GroupId":              "3",
		"NetworkIntegrationId": "5",
	})
	if err != nil {
		t.Fatalf("failed to render create config: %s", err)
	}

	updateConfig, err := networkrouter.RenderNetworkRouterNSXTGatewayTier1Config(t, map[string]string{
		"Name":                 updatedName,
		"GroupId":              "3",
		"NetworkIntegrationId": "5",
	})
	if err != nil {
		t.Fatalf("failed to render update config: %s", err)
	}

	resourceName := testResourceName
	providerConfig := testhelpers.ProviderBlock()

	createChecks := resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttr(resourceName, "name", name),
		resource.TestCheckResourceAttr(resourceName, "type_id", "9"),
		resource.TestCheckResourceAttr(resourceName, "group_id", "3"),
		resource.TestCheckResourceAttr(resourceName, "network_integration_id", "5"),
		resource.TestCheckResourceAttr(resourceName, "config_nsxt_gateway_tier1.ip_management_type", "dhcpLocal"),
	)

	updateChecks := resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttr(resourceName, "name", updatedName),
		resource.TestCheckResourceAttr(resourceName, "type_id", "9"),
		resource.TestCheckResourceAttr(resourceName, "group_id", "3"),
		resource.TestCheckResourceAttr(resourceName, "network_integration_id", "5"),
		resource.TestCheckResourceAttr(resourceName, "config_nsxt_gateway_tier1.ip_management_type", "dhcpLocal"),
	)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + createConfig,
				Check:  createChecks,
			},
			{
				Config:             providerConfig + createConfig,
				Check:              createChecks,
				ExpectNonEmptyPlan: false,
				PlanOnly:           true,
			},
			{
				Config: providerConfig + updateConfig,
				Check:  updateChecks,
			},
			{
				Config:             providerConfig + updateConfig,
				Check:              updateChecks,
				ExpectNonEmptyPlan: false,
				PlanOnly:           true,
			},
		},
	})
}
