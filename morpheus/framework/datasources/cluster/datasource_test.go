// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package cluster_test

import (
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/datasources/cluster"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
)

const providerConfigOffline = `
provider "hpe" {
  morpheus {
    url          = ""
    username     = ""
    password     = ""
  }
}
`

func TestMain(m *testing.M) {
	code := m.Run()
	testhelpers.WriteMergedResults()

	os.Exit(code)
}

func TestAccMorpheusFindClusterById(t *testing.T) {
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

	dataSourceConfig, err := testhelpers.RenderExample(t, "example-id.tf.tmpl", "Id", "1")
	if err != nil {
		t.Fatal(err)
	}

	// We can't test all the possible schema attributes with this type of cluster.
	// We'll just check that what was required to create it was read correctly,
	// as well as some config properties.
	checkFn := resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttr("data.hpe_morpheus_cluster.example", "id", "1"),
		resource.TestCheckResourceAttrSet("data.hpe_morpheus_cluster.example", "name"),
		resource.TestCheckResourceAttrSet("data.hpe_morpheus_cluster.example", "cloud.id"),
		resource.TestCheckResourceAttrSet("data.hpe_morpheus_cluster.example", "cloud.name"),
		resource.TestCheckResourceAttrSet("data.hpe_morpheus_cluster.example", "cloud.cloud_type.id"),
		resource.TestCheckResourceAttrSet("data.hpe_morpheus_cluster.example", "group.id"),
		resource.TestCheckResourceAttrSet("data.hpe_morpheus_cluster.example", "group.name"),
		resource.TestCheckResourceAttrSet("data.hpe_morpheus_cluster.example", "layout.id"),
		resource.TestCheckResourceAttrSet("data.hpe_morpheus_cluster.example", "layout.name"),
		resource.TestCheckResourceAttrSet("data.hpe_morpheus_cluster.example", "layout.provision_type_code"),
		resource.TestCheckResourceAttrSet("data.hpe_morpheus_cluster.example", "config.cpuArch"),
		resource.TestCheckResourceAttrSet("data.hpe_morpheus_cluster.example", "config.cpuModel"),
		resource.TestCheckResourceAttrSet("data.hpe_morpheus_cluster.example", "config.dynamicPlacementMode"),
		resource.TestCheckResourceAttrSet("data.hpe_morpheus_cluster.example", "config.powerPolicy"),
	)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + dataSourceConfig,
				Check:  checkFn,
			},
		},
	})
}

func TestAccMorpheusFindClusterByName(t *testing.T) {
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

	dataSourceConfig, err := testhelpers.RenderExample(t, "example-name.tf.tmpl", "Name", `"Duck"`)
	if err != nil {
		t.Fatal(err)
	}

	// We can't test all the possible schema attributes with this type of cluster.
	// We'll just check that what was required to create it was read correctly,
	// as well as some config properties.
	checkFn := resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttrSet("data.hpe_morpheus_cluster.example", "id"),
		resource.TestCheckResourceAttr("data.hpe_morpheus_cluster.example", "name", "Duck"),
		resource.TestCheckResourceAttrSet("data.hpe_morpheus_cluster.example", "cloud.id"),
		resource.TestCheckResourceAttrSet("data.hpe_morpheus_cluster.example", "cloud.name"),
		resource.TestCheckResourceAttrSet("data.hpe_morpheus_cluster.example", "cloud.cloud_type.id"),
		resource.TestCheckResourceAttrSet("data.hpe_morpheus_cluster.example", "group.id"),
		resource.TestCheckResourceAttrSet("data.hpe_morpheus_cluster.example", "group.name"),
		resource.TestCheckResourceAttrSet("data.hpe_morpheus_cluster.example", "layout.id"),
		resource.TestCheckResourceAttrSet("data.hpe_morpheus_cluster.example", "layout.name"),
		resource.TestCheckResourceAttrSet("data.hpe_morpheus_cluster.example", "layout.provision_type_code"),
		resource.TestCheckResourceAttrSet("data.hpe_morpheus_cluster.example", "config.cpuArch"),
		resource.TestCheckResourceAttrSet("data.hpe_morpheus_cluster.example", "config.cpuModel"),
		resource.TestCheckResourceAttrSet("data.hpe_morpheus_cluster.example", "config.dynamicPlacementMode"),
		resource.TestCheckResourceAttrSet("data.hpe_morpheus_cluster.example", "config.powerPolicy"),
	)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + dataSourceConfig,
				Check:  checkFn,
			},
		},
	})
}

func TestAccMorpheusFindClusterNotFound(t *testing.T) {
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
      data "hpe_morpheus_cluster" "test" {
        name = "______"
      }`

	checkFn := resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckNoResourceAttr("data.hpe_morpheus_cluster.test", "id"),
	)

	expected := cluster.ErrorNoClusterFound

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
		Steps: []resource.TestStep{
			{
				Config:      config,
				Check:       checkFn,
				ExpectError: regexp.MustCompile(expected),
			},
		},
	})
}

func TestAccMorpheusFindClusterNoSearchAttrs(t *testing.T) {
	if capabilities.Missing(t, capabilities.All) {
		t.Log("Skipping test due to missing capabilities")

		return
	}
	t.Parallel()
	defer testhelpers.RecordResult(t)

	config := providerConfigOffline + `
      data "hpe_morpheus_cluster" "test" {
      }`

	checkFn := resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckNoResourceAttr("data.hpe_morpheus_cluster.test", "id"),
	)

	expected := cluster.ErrorNoValidSearchTerms

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
		Steps: []resource.TestStep{
			{
				Config:      config,
				Check:       checkFn,
				ExpectError: regexp.MustCompile(expected),
			},
		},
	})
}

func TestAccMorpheusFindClusterBothSearchAttrs(t *testing.T) {
	if capabilities.Missing(t, capabilities.All) {
		t.Log("Skipping test due to missing capabilities")

		return
	}
	t.Parallel()
	defer testhelpers.RecordResult(t)

	config := providerConfigOffline + `
      data "hpe_morpheus_cluster" "test" {
        id = 1
        name = "______"
      }`

	checkFn := resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckNoResourceAttr("data.hpe_morpheus_cluster.test", "id"),
	)

	expected := cluster.ErrorRunningPreApply

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
		Steps: []resource.TestStep{
			{
				Config:      config,
				Check:       checkFn,
				ExpectError: regexp.MustCompile(expected),
			},
		},
	})
}
