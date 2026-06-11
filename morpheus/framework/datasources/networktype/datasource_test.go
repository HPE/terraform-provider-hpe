// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networktype_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/datasources/networktype"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
)

func TestMain(m *testing.M) {
	code := m.Run()
	testhelpers.WriteMergedResults()
	os.Exit(code)
}

func TestAccMorpheusNetworkTypeDataSourceByNameExampleOk(t *testing.T) {
	if capabilities.Missing(t, capabilities.Network) {
		t.Log("Skipping test due to missing capabilities")

		return
	}

	t.Parallel()
	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()

	dataSourceConfig, err := networktype.RenderNetworkTypeDataSourceByNameConfig(t, nil)
	if err != nil {
		t.Fatal(err)
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + dataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.hpe_morpheus_network_type.example", "id"),
					resource.TestCheckResourceAttr("data.hpe_morpheus_network_type.example", "name", "Host Network"),
				),
			},
		},
	})
}

func TestAccMorpheusNetworkTypeDataSourceByIdExampleOk(t *testing.T) {
	if capabilities.Missing(t, capabilities.Network) {
		t.Log("Skipping test due to missing capabilities")

		return
	}

	t.Parallel()
	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()

	dataSourceConfig, err := networktype.RenderNetworkTypeDataSourceByIDConfig(t, nil)
	if err != nil {
		t.Fatal(err)
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + dataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.hpe_morpheus_network_type.example", "name"),
					resource.TestCheckResourceAttr("data.hpe_morpheus_network_type.example", "id", "1"),
				),
			},
		},
	})
}
