package network_pool_server_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
)

func TestMain(m *testing.M) {
	code := testhelpers.TestMain(m)
	testhelpers.WriteMergedResults()
	os.Exit(code)
}

func TestAccMorpheusNetworkPoolServerResourceBasic(t *testing.T) {
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

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + fmt.Sprintf(`
resource "hpe_morpheus_network_pool_server" "test" {
  name                        = %q
  type_id                     = 4
  enabled                     = true
  service_url                 = "http://localhost:8080"
  service_username            = "admin"
  service_password_wo         = "password123"
  service_password_wo_version = 1
}
`, name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hpe_morpheus_network_pool_server.test", "name", name),
					resource.TestCheckResourceAttr("hpe_morpheus_network_pool_server.test", "type_id", "4"),
					resource.TestCheckResourceAttr("hpe_morpheus_network_pool_server.test", "enabled", "true"),
					resource.TestCheckResourceAttrSet("hpe_morpheus_network_pool_server.test", "id"),
				),
			},
			{
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"service_username", "service_password_wo", "service_password_wo_version"},
				ResourceName:            "hpe_morpheus_network_pool_server.test",
			},
		},
	})
}

func TestAccMorpheusNetworkPoolServerResourceUpdate(t *testing.T) {
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

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + fmt.Sprintf(`
resource "hpe_morpheus_network_pool_server" "test" {
  name                        = %q
  type_id                     = 4
  enabled                     = true
  service_url                 = "http://localhost:8080"
  service_username            = "admin"
  service_password_wo         = "password123"
  service_password_wo_version = 1
}
`, name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hpe_morpheus_network_pool_server.test", "enabled", "true"),
				),
			},
			{
				Config: providerConfig + fmt.Sprintf(`
resource "hpe_morpheus_network_pool_server" "test" {
  name                        = %q
  type_id                     = 4
  enabled                     = false
  service_url                 = "http://localhost:8080"
  service_username            = "admin"
  service_password_wo         = "password123"
  service_password_wo_version = 1
}
`, name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hpe_morpheus_network_pool_server.test", "enabled", "false"),
				),
			},
		},
	})
}

func TestAccMorpheusNetworkPoolServerCredential(t *testing.T) {
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

	name := acctest.RandomWithPrefix("tf-acc-pool-srv")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + fmt.Sprintf(`
resource "hpe_morpheus_network_pool_server" "test" {
  name          = %q
  type_id       = 4
  enabled       = true
  service_url   = "http://localhost:8080"
  credential_id = 1
}
`, name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hpe_morpheus_network_pool_server.test", "name", name),
					resource.TestCheckResourceAttr("hpe_morpheus_network_pool_server.test", "credential_id", "1"),
					resource.TestCheckResourceAttrSet("hpe_morpheus_network_pool_server.test", "id"),
				),
			},
			{
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"credential_id"},
				ResourceName:            "hpe_morpheus_network_pool_server.test",
			},
		},
	})
}

func TestAccMorpheusNetworkPoolServerWithFilters(t *testing.T) {
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

	name := acctest.RandomWithPrefix("tf-acc-pool-srv")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + fmt.Sprintf(`
resource "hpe_morpheus_network_pool_server" "test" {
  name                        = %q
  type_id                     = 4
  enabled                     = true
  service_url                 = "http://localhost:8080"
  service_username            = "admin"
  service_password_wo         = "password123"
  service_password_wo_version = 1
  network_filter              = "10.0.0.0/8"
  zone_filter                 = "example.com"
  tenant_match                = ".*"
  service_mode                = "static"
  service_throttle_rate       = 100
}
`, name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hpe_morpheus_network_pool_server.test", "name", name),
					resource.TestCheckResourceAttr("hpe_morpheus_network_pool_server.test", "network_filter", "10.0.0.0/8"),
					resource.TestCheckResourceAttr("hpe_morpheus_network_pool_server.test", "zone_filter", "example.com"),
					resource.TestCheckResourceAttr("hpe_morpheus_network_pool_server.test", "tenant_match", ".*"),
					resource.TestCheckResourceAttr("hpe_morpheus_network_pool_server.test", "service_mode", "static"),
					resource.TestCheckResourceAttr("hpe_morpheus_network_pool_server.test", "service_throttle_rate", "100"),
				),
			},
			// Update filters
			{
				Config: providerConfig + fmt.Sprintf(`
resource "hpe_morpheus_network_pool_server" "test" {
  name                        = %q
  type_id                     = 4
  enabled                     = true
  service_url                 = "http://localhost:8080"
  service_username            = "admin"
  service_password_wo         = "password123"
  service_password_wo_version = 1
  network_filter              = "192.168.0.0/16"
  zone_filter                 = "updated.example.com"
  tenant_match                = "tenant-.*"
  service_mode                = "static"
  service_throttle_rate       = 200
}
`, name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hpe_morpheus_network_pool_server.test", "network_filter", "192.168.0.0/16"),
					resource.TestCheckResourceAttr("hpe_morpheus_network_pool_server.test", "zone_filter", "updated.example.com"),
					resource.TestCheckResourceAttr("hpe_morpheus_network_pool_server.test", "tenant_match", "tenant-.*"),
					resource.TestCheckResourceAttr("hpe_morpheus_network_pool_server.test", "service_throttle_rate", "200"),
				),
			},
			{
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"service_username", "service_password_wo", "service_password_wo_version"},
				ResourceName:            "hpe_morpheus_network_pool_server.test",
			},
		},
	})
}
