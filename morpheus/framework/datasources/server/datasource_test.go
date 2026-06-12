// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package server_test

import (
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/datasources/server"
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

func TestAccMorpheusFindServerByName(t *testing.T) {
	if capabilities.Missing(t, capabilities.All) {
		t.Log("Skipping test due to missing capabilities")

		return
	}

	t.Parallel()
	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()

	dataSourceConfig, err := server.RenderServerDataSourceByNameConfig(t, map[string]string{"Name": "edge01"})
	if err != nil {
		t.Fatal(err)
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + dataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.hpe_morpheus_server.example", "id"),
					resource.TestCheckResourceAttr("data.hpe_morpheus_server.example", "name", "edge01"),
				),
			},
		},
	})
}

func TestAccMorpheusFindServerById(t *testing.T) {
	if capabilities.Missing(t, capabilities.All) {
		t.Log("Skipping test due to missing capabilities")

		return
	}

	t.Parallel()
	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()

	dataSourceConfig, err := server.RenderServerDataSourceByIDConfig(t, map[string]string{"Id": "508"})
	if err != nil {
		t.Fatal(err)
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + dataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.hpe_morpheus_server.example", "name"),
					resource.TestCheckResourceAttr("data.hpe_morpheus_server.example", "id", "508"),
				),
			},
		},
	})
}

func TestAccMorpheusFindServerNotFound(t *testing.T) {
	if capabilities.Missing(t, capabilities.All) {
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
data "hpe_morpheus_server" "test" {
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

func TestAccMorpheusFindServerNoSearchAttrs(t *testing.T) {
	if capabilities.Missing(t, capabilities.All) {
		t.Log("Skipping test due to missing capabilities")

		return
	}

	t.Parallel()
	defer testhelpers.RecordResult(t)

	config := providerConfigOffline + `
data "hpe_morpheus_server" "test" {
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

func TestAccMorpheusFindServerBothSearchAttrs(t *testing.T) {
	if capabilities.Missing(t, capabilities.All) {
		t.Log("Skipping test due to missing capabilities")

		return
	}

	t.Parallel()
	defer testhelpers.RecordResult(t)

	config := providerConfigOffline + `
data "hpe_morpheus_server" "test" {
  id   = 1
  name = "edge01"
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
