// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package networkdhcpserver_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/networkdhcpserver"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
)

func TestMain(m *testing.M) {
	code := m.Run()
	testhelpers.WriteMergedResults()
	os.Exit(code)
}

func TestAccMorpheusNetworkDhcpServerExampleOk(t *testing.T) {
	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()

	name := acctest.RandomWithPrefix(t.Name())

	resourceConfig, err := networkdhcpserver.RenderNetworkDhcpServerConfig(t,
		map[string]string{
			"Name": name,
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"hpe_morpheus_network_dhcp_server.example",
			"name",
			name,
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_network_dhcp_server.example",
			"server_ip_address",
			"192.168.1.1/24",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_network_dhcp_server.example",
			"lease_time",
			"86400",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_network_dhcp_server.example",
			"config_nsx.edge_cluster",
			"qa-edge-cluster-01",
		),
		resource.TestCheckResourceAttrSet(
			"hpe_morpheus_network_dhcp_server.example",
			"id",
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(
			t, morpheus.New(), nil,
		),
		Steps: []resource.TestStep{
			{
				Config:             providerConfig + resourceConfig,
				ExpectNonEmptyPlan: false,
				Check:              checkFn,
				PlanOnly:           false,
			},
		},
	})
}

func TestAccMorpheusNetworkDhcpServerUpdateOk(t *testing.T) {
	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()

	name := acctest.RandomWithPrefix(t.Name())

	createConfig, err := networkdhcpserver.RenderNetworkDhcpServerConfig(t,
		map[string]string{
			"Name":            name,
			"ServerIpAddress": "192.168.1.1/24",
			"LeaseTime":       "86400",
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	updatedName := name + "-updated"

	updateConfig, err := networkdhcpserver.RenderNetworkDhcpServerConfig(t,
		map[string]string{
			"Name":            updatedName,
			"ServerIpAddress": "192.168.1.2/24",
			"LeaseTime":       "43200",
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	createChecks := resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttr(
			"hpe_morpheus_network_dhcp_server.example",
			"name",
			name,
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_network_dhcp_server.example",
			"server_ip_address",
			"192.168.1.1/24",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_network_dhcp_server.example",
			"lease_time",
			"86400",
		),
	)

	updateChecks := resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttr(
			"hpe_morpheus_network_dhcp_server.example",
			"name",
			updatedName,
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_network_dhcp_server.example",
			"server_ip_address",
			"192.168.1.2/24",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_network_dhcp_server.example",
			"lease_time",
			"43200",
		),
	)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(
			t, morpheus.New(), nil,
		),
		Steps: []resource.TestStep{
			{
				Config:   providerConfig + createConfig,
				Check:    createChecks,
				PlanOnly: false,
			},
			{
				Config:             providerConfig + createConfig,
				Check:              createChecks,
				ExpectNonEmptyPlan: false,
				PlanOnly:           true,
			},
			{
				Config:   providerConfig + updateConfig,
				Check:    updateChecks,
				PlanOnly: false,
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
