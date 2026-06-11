// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkpool_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/datasources/networkpool"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
)

func TestMain(m *testing.M) {
	code := m.Run()
	testhelpers.WriteMergedResults()
	os.Exit(code)
}

func TestAccMorpheusNetworkPoolDataSourceByNameExampleOk(t *testing.T) {
	if capabilities.Missing(t, capabilities.NetworkPool) {
		t.Log("Skipping test due to missing capabilities")

		return
	}

	t.Parallel()
	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()

	dataSourceConfig, err := networkpool.RenderNetworkPoolDataSourceByNameConfig(t, nil)
	if err != nil {
		t.Fatal(err)
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + dataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.hpe_morpheus_network_pool.example", "id"),
					resource.TestCheckResourceAttr("data.hpe_morpheus_network_pool.example", "name", "Morpheus"),
				),
			},
		},
	})
}

func TestAccMorpheusNetworkPoolDataSourceByIdExampleOk(t *testing.T) {
	if capabilities.Missing(t, capabilities.NetworkPool) {
		t.Log("Skipping test due to missing capabilities")

		return
	}

	t.Parallel()
	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()

	dataSourceConfig, err := networkpool.RenderNetworkPoolDataSourceByIDConfig(t, nil)
	if err != nil {
		t.Fatal(err)
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + dataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.hpe_morpheus_network_pool.example", "name"),
					resource.TestCheckResourceAttr("data.hpe_morpheus_network_pool.example", "id", "1"),
				),
			},
		},
	})
}
