package network_pool_server_test

import (
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"

	"github.com/HPE/terraform-provider-hpe/morpheus"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/network_pool_server"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
)

func TestMain(m *testing.M) {
	code := testhelpers.TestMain(m)
	testhelpers.WriteMergedResults()
	os.Exit(code)
}

func TestAccMorpheusNetworkPoolServerResourceExampleOk(t *testing.T) {
	if capabilities.Missing(t, capabilities.NetworkPool) {
		t.Log("Skipping test due to missing capabilities")

		return
	}
	defer testhelpers.RecordResult(t)
	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()
	name := acctest.RandomWithPrefix(t.Name())
	resourceName := "hpe_morpheus_network_pool_server.example"

	resourceConfig, err := network_pool_server.RenderNetworkPoolServerInfobloxConfig(t, map[string]string{
		"Name": name,
	})
	if err != nil {
		t.Fatal(err)
	}

	resourceConfig = strings.Replace(
		resourceConfig,
		`resource "hpe_morpheus_network_pool_server" "infoblox" {`,
		`resource "hpe_morpheus_network_pool_server" "example" {`,
		1,
	)

	checks := resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttr(resourceName, "name", name),
		resource.TestCheckResourceAttr(resourceName, "type_id", "1"),
		resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
		resource.TestCheckResourceAttr(resourceName, "service_url", "https://infoblox.example.com/wapi/v2.12"),
		resource.TestCheckResourceAttr(resourceName, "service_username", "admin"),
		resource.TestCheckResourceAttr(resourceName, "ignore_ssl", "true"),
		resource.TestCheckResourceAttr(resourceName, "network_filter", "10.0.0.0/8"),
		resource.TestCheckResourceAttr(resourceName, "zone_filter", "example.com"),
		resource.TestCheckResourceAttr(resourceName, "tenant_match", ".*"),
		resource.TestCheckResourceAttr(resourceName, "service_mode", "static"),
		resource.TestCheckResourceAttr(resourceName, "service_throttle_rate", "0"),
		resource.TestCheckResourceAttrSet(resourceName, "id"),
	)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + resourceConfig,
				Check:  checks,
			},
			{
				Config:             providerConfig + resourceConfig,
				ExpectNonEmptyPlan: false,
				PlanOnly:           true,
			},
			{
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"service_username", "service_password_wo", "service_password_wo_version"},
				ResourceName:            "hpe_morpheus_network_pool_server.example",
			},
		},
	})
}

func TestAccMorpheusNetworkPoolServerResourceUpdateOk(t *testing.T) {
	if capabilities.Missing(t, capabilities.NetworkPool) {
		t.Log("Skipping test due to missing capabilities")

		return
	}
	defer testhelpers.RecordResult(t)
	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()
	name := acctest.RandomWithPrefix(t.Name())
	resourceName := "hpe_morpheus_network_pool_server.example"

	createConfig, err := network_pool_server.RenderNetworkPoolServerInfobloxConfig(t, map[string]string{
		"Name": name,
	})
	if err != nil {
		t.Fatal(err)
	}

	createConfig = strings.Replace(
		createConfig,
		`resource "hpe_morpheus_network_pool_server" "infoblox" {`,
		`resource "hpe_morpheus_network_pool_server" "example" {`,
		1,
	)

	updateConfig := `
resource "hpe_morpheus_network_pool_server" "example" {
  name                        = "` + name + `"
  type_id                     = 1
  enabled                     = false
  service_url                 = "https://infoblox.example.com/wapi/v2.12"
  service_username            = "admin"
  service_password_wo         = "changeme"
  service_password_wo_version = 1
  ignore_ssl                  = true
  network_filter              = "192.168.0.0/16"
  zone_filter                 = "example.com"
  tenant_match                = ".*"
  service_mode                = "static"
  service_throttle_rate       = 25
}
`

	createChecks := resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttr(resourceName, "name", name),
		resource.TestCheckResourceAttr(resourceName, "type_id", "1"),
		resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
		resource.TestCheckResourceAttr(resourceName, "service_url", "https://infoblox.example.com/wapi/v2.12"),
		resource.TestCheckResourceAttr(resourceName, "service_username", "admin"),
		resource.TestCheckResourceAttr(resourceName, "ignore_ssl", "true"),
		resource.TestCheckResourceAttr(resourceName, "network_filter", "10.0.0.0/8"),
		resource.TestCheckResourceAttr(resourceName, "zone_filter", "example.com"),
		resource.TestCheckResourceAttr(resourceName, "tenant_match", ".*"),
		resource.TestCheckResourceAttr(resourceName, "service_mode", "static"),
		resource.TestCheckResourceAttr(resourceName, "service_throttle_rate", "0"),
	)

	updateChecks := resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttr(resourceName, "name", name),
		resource.TestCheckResourceAttr(resourceName, "type_id", "1"),
		resource.TestCheckResourceAttr(resourceName, "enabled", "false"),
		resource.TestCheckResourceAttr(resourceName, "service_url", "https://infoblox.example.com/wapi/v2.12"),
		resource.TestCheckResourceAttr(resourceName, "service_username", "admin"),
		resource.TestCheckResourceAttr(resourceName, "ignore_ssl", "true"),
		resource.TestCheckResourceAttr(resourceName, "network_filter", "192.168.0.0/16"),
		resource.TestCheckResourceAttr(resourceName, "zone_filter", "example.com"),
		resource.TestCheckResourceAttr(resourceName, "tenant_match", ".*"),
		resource.TestCheckResourceAttr(resourceName, "service_mode", "static"),
		resource.TestCheckResourceAttr(resourceName, "service_throttle_rate", "25"),
	)

	checkInPlaceUpdate := resource.ConfigPlanChecks{
		PreApply: []plancheck.PlanCheck{
			plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionUpdate),
		},
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
		Steps: []resource.TestStep{
			{Config: providerConfig + createConfig, Check: createChecks},
			{Config: providerConfig + updateConfig, Check: updateChecks, ConfigPlanChecks: checkInPlaceUpdate},
			{Config: providerConfig + updateConfig, ExpectNonEmptyPlan: false, PlanOnly: true},
		},
	})
}
