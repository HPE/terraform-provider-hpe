// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package cluster_test

import (
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/datasources/cluster"
	sdkv2morpheus "github.com/HPE/terraform-provider-hpe/morpheus/sdkv2"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/systemoverride"
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
	systemoverride.ParseFlags()
	code := m.Run()
	testhelpers.WriteMergedResults()

	os.Exit(code)
}

func TestAccMorpheusFindClusterByID(t *testing.T) {
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

	checkFn := resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttrSet("data.hpe_morpheus_cluster.example", "id"),
		resource.TestCheckResourceAttrSet("data.hpe_morpheus_cluster.example", "name"),
		resource.TestCheckResourceAttrSet("data.hpe_morpheus_cluster.example", "cloud_id"),
		resource.TestCheckResourceAttrSet("data.hpe_morpheus_cluster.example", "group_id"),
		resource.TestCheckResourceAttrSet("data.hpe_morpheus_cluster.example", "layout_id"),
		resource.TestCheckResourceAttrSet("data.hpe_morpheus_cluster.example", "description"),
		resource.TestCheckResourceAttrSet("data.hpe_morpheus_cluster.example", "service_url"),
		resource.TestCheckResourceAttrSet("data.hpe_morpheus_cluster.example", "uuid"),
	)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), sdkv2morpheus.Provider()),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + dataSourceConfig,
				Check:  checkFn,
			},
		},
	})
}

func TestAccMorpheusFindClusterByName(t *testing.T) {
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

	checkFn := resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttrSet("data.hpe_morpheus_cluster.example", "id"),
		resource.TestCheckResourceAttrSet("data.hpe_morpheus_cluster.example", "name"),
		resource.TestCheckResourceAttrSet("data.hpe_morpheus_cluster.example", "cloud_id"),
		resource.TestCheckResourceAttrSet("data.hpe_morpheus_cluster.example", "group_id"),
		resource.TestCheckResourceAttrSet("data.hpe_morpheus_cluster.example", "layout_id"),
		resource.TestCheckResourceAttrSet("data.hpe_morpheus_cluster.example", "description"),
		resource.TestCheckResourceAttrSet("data.hpe_morpheus_cluster.example", "service_url"),
		resource.TestCheckResourceAttrSet("data.hpe_morpheus_cluster.example", "uuid"),
	)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), sdkv2morpheus.Provider()),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + dataSourceConfig,
				Check:  checkFn,
			},
		},
	})
}

func TestAccMorpheusFindClusterNotFound(t *testing.T) {
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
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), sdkv2morpheus.Provider()),
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
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), sdkv2morpheus.Provider()),
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
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), sdkv2morpheus.Provider()),
		Steps: []resource.TestStep{
			{
				Config:      config,
				Check:       checkFn,
				ExpectError: regexp.MustCompile(expected),
			},
		},
	})
}
