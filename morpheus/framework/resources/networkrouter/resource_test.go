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
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/systemoverride"
)

func TestMain(m *testing.M) {
	systemoverride.ParseFlags()
	code := m.Run()
	testhelpers.WriteMergedResults()
	os.Exit(code)
}

func TestAccMorpheusNetworkRouterNSXGatewayTier0ExampleOk(t *testing.T) {
	t.Parallel()
	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow acceptance test in short mode")
	}

	t.Skip("Skipping due to missing infrastructure in test environment")

	name := acctest.RandomWithPrefix(t.Name())

	config, err := networkrouter.RenderNetworkRouterNSXGatewayTier0Config(t, map[string]string{
		"Name": name,
	})
	if err != nil {
		t.Fatalf("failed to render config: %s", err)
	}

	resourceName := "hpe_morpheus_network_router.example"
	config = testhelpers.ProviderBlock() + config

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "config_nsx_gateway_tier0.ha_mode", "ACTIVE_ACTIVE"),
					resource.TestCheckResourceAttr(resourceName, "config_nsx_gateway_tier0.ip_server_id", "1"),
					resource.TestCheckResourceAttr(resourceName, "config_nsx_gateway_tier0.restart_mode", "HELPER_ONLY"),
				),
			},
		},
	})
}

func TestAccMorpheusNetworkRouterNSXGatewayTier1ExampleOk(t *testing.T) {
	t.Parallel()
	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow acceptance test in short mode")
	}

	t.Skip("Skipping due to missing infrastructure in test environment")

	name := acctest.RandomWithPrefix(t.Name())

	config, err := networkrouter.RenderNetworkRouterNSXGatewayTier1Config(t, map[string]string{
		"Name": name,
	})
	if err != nil {
		t.Fatalf("failed to render config: %s", err)
	}

	resourceName := "hpe_morpheus_network_router.example"
	config = testhelpers.ProviderBlock() + config

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", name),
				),
			},
		},
	})
}
