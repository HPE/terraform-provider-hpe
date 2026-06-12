// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkpool_test

import (
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/datasources/networkpool"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
)

const providerConfigOffline = `
provider "hpe" {
  morpheus {
    url      = ""
    username = ""
    password = ""
  }
}
`

func TestMain(m *testing.M) {
	code := m.Run()
	testhelpers.WriteMergedResults()
	os.Exit(code)
}

func TestAccMorpheusFindNetworkPoolByName(t *testing.T) {
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

	dataSourceConfig, err := networkpool.RenderNetworkPoolDataSourceByNameConfig(t, map[string]string{"Name": "CAN"})
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
					resource.TestCheckResourceAttr("data.hpe_morpheus_network_pool.example", "name", "CAN"),
				),
			},
		},
	})
}

func TestAccMorpheusFindNetworkPoolById(t *testing.T) {
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

func TestAccMorpheusFindNetworkPoolNotFound(t *testing.T) {
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

	config := providerConfig + `
data "hpe_morpheus_network_pool" "test" {
  name = "____nonexistent____"
}`

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
		Steps: []resource.TestStep{
			{
				Config:      config,
				ExpectError: regexp.MustCompile(`not found`),
			},
		},
	})
}

func TestAccMorpheusFindNetworkPoolNoSearchAttrs(t *testing.T) {
	if capabilities.Missing(t, capabilities.NetworkPool) {
		t.Log("Skipping test due to missing capabilities")

		return
	}

	t.Parallel()
	defer testhelpers.RecordResult(t)

	config := providerConfigOffline + `
data "hpe_morpheus_network_pool" "test" {
}`

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
		Steps: []resource.TestStep{
			{
				Config:      config,
				ExpectError: regexp.MustCompile(`Error running pre-apply plan|at least one`),
			},
		},
	})
}

func TestAccMorpheusFindNetworkPoolBothSearchAttrs(t *testing.T) {
	if capabilities.Missing(t, capabilities.NetworkPool) {
		t.Log("Skipping test due to missing capabilities")

		return
	}

	t.Parallel()
	defer testhelpers.RecordResult(t)

	config := providerConfigOffline + `
data "hpe_morpheus_network_pool" "test" {
  id   = 1
  name = "CAN"
}`

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
		Steps: []resource.TestStep{
			{
				Config:      config,
				ExpectError: regexp.MustCompile(`Conflicting configuration arguments|conflicts with`),
			},
		},
	})
}
