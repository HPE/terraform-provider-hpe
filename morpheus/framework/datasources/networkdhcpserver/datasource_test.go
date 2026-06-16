// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package networkdhcpserver_test

import (
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/datasources/networkdhcpserver"
	dhcpresource "github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/networkdhcpserver"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
)

func TestMain(m *testing.M) {
	code := m.Run()
	testhelpers.WriteMergedResults()
	os.Exit(code)
}

const providerConfigOffline = `
provider "hpe" {
  morpheus {
    url          = ""
    username     = ""
    password     = ""
  }
}
`

// dhcpFixture renders a self-contained DHCP server on the QA NSX-T network
// integration (id 5), labelled hpe_morpheus_network_dhcp_server.example.
//
// QA verify: NSX-T integration id 5 and edge cluster "qa-edge-cluster-01" (the
// resource example default) are the QA appliance values.
func dhcpFixture(t *testing.T, name, serverIP string) string {
	t.Helper()

	cfg, err := dhcpresource.RenderNetworkDhcpServerConfig(t, map[string]string{
		"NetworkIntegrationId": "5",
		"Name":                 name,
		"ServerIpAddress":      serverIP,
	})
	if err != nil {
		t.Fatal(err)
	}

	return cfg
}

func TestAccMorpheusFindNetworkDhcpServerByName(t *testing.T) {
	if capabilities.Missing(t, capabilities.NetworkDHCP) {
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

	// depends_on defers the data source read until the DHCP server exists.
	dataSourceConfig := `
data "hpe_morpheus_network_dhcp_server" "example" {
  name                   = "` + name + `"
  network_integration_id = 5
  depends_on             = [hpe_morpheus_network_dhcp_server.example]
}
`

	checkFn := resource.ComposeAggregateTestCheckFunc(networkDhcpServerChecks()...)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + dhcpFixture(t, name, "192.168.40.1/24") + dataSourceConfig,
				Check:  checkFn,
			},
		},
	})
}

func TestAccMorpheusFindNetworkDhcpServerById(t *testing.T) {
	if capabilities.Missing(t, capabilities.NetworkDHCP) {
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

	// id references the created DHCP server, deferring the read.
	dataSourceConfig, err := networkdhcpserver.RenderNetworkDhcpServerByIdConfig(t, map[string]string{
		"Id":                   "hpe_morpheus_network_dhcp_server.example.id",
		"NetworkIntegrationId": "5",
	})
	if err != nil {
		t.Fatal(err)
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(networkDhcpServerChecks()...)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + dhcpFixture(t, name, "192.168.41.1/24") + dataSourceConfig,
				Check:  checkFn,
			},
		},
	})
}

func TestAccMorpheusFindNetworkDhcpServerNotFound(t *testing.T) {
	if capabilities.Missing(t, capabilities.NetworkDHCP) {
		t.Log("Skipping test due to missing capabilities")

		return
	}
	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()

	// Search the real NSX-T integration (5) for a DHCP server name that does
	// not exist.
	dataSourceConfig, err := networkdhcpserver.RenderNetworkDhcpServerByNameConfig(t,
		map[string]string{
			"Name":                 "______",
			"NetworkIntegrationId": "5",
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	expected := regexp.MustCompile(`no network DHCP server found`)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
		Steps: []resource.TestStep{
			{
				Config:      providerConfig + dataSourceConfig,
				ExpectError: expected,
			},
		},
	})
}

func TestAccMorpheusFindNetworkDhcpServerNoSearchAttrs(t *testing.T) {
	if capabilities.Missing(t, capabilities.NetworkDHCP) {
		t.Log("Skipping test due to missing capabilities")

		return
	}
	defer testhelpers.RecordResult(t)

	t.Parallel()

	config := providerConfigOffline + `
      data "hpe_morpheus_network_dhcp_server" "test" {
        network_integration_id = 1
      }`

	expected := networkdhcpserver.ErrorNoValidSearchTerms

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
		Steps: []resource.TestStep{
			{
				Config:      config,
				ExpectError: regexp.MustCompile(expected),
			},
		},
	})
}

func TestAccMorpheusFindNetworkDhcpServerBothSearchAttrs(t *testing.T) {
	if capabilities.Missing(t, capabilities.NetworkDHCP) {
		t.Log("Skipping test due to missing capabilities")

		return
	}
	defer testhelpers.RecordResult(t)

	t.Parallel()

	config := providerConfigOffline + `
      data "hpe_morpheus_network_dhcp_server" "test" {
        id                     = 1
        name                   = "______"
        network_integration_id = 1
      }`

	expected := networkdhcpserver.ErrorRunningPreApply

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
		Steps: []resource.TestStep{
			{
				Config:      config,
				ExpectError: regexp.MustCompile(expected),
			},
		},
	})
}

func networkDhcpServerChecks() []resource.TestCheckFunc {
	ds := "data.hpe_morpheus_network_dhcp_server.example"

	return []resource.TestCheckFunc{
		resource.TestCheckResourceAttrSet(ds, "id"),
		resource.TestCheckResourceAttrSet(ds, "name"),
		resource.TestCheckResourceAttrSet(ds, "network_integration_id"),
		// lease_time is intentionally not checked: the DHCP server GET response
		// does not return leaseTime, so the data source cannot populate it (the
		// resource preserves it from prior state, but a data source has none).
	}
}
