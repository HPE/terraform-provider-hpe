// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package clusterlayout_test

import (
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus/framework/datasources/clusterlayout"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
	"github.com/HPE/terraform-provider-hpe/provider/adapter"
)

func TestMain(m *testing.M) {
	code := testhelpers.TestMain(m)
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

func TestAccMorpheusFindClusterLayoutByName(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()

	dataSourceConfig, err := clusterlayout.RenderClusterLayoutDataSourceByNameConfig(
		t,
		map[string]string{"Name": "HVM 1.3 Cluster on HVM/Ubuntu 24.04"},
	)
	if err != nil {
		t.Fatal(err)
	}

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttrSet(
			"data.hpe_morpheus_cluster_layout.example",
			"id",
		),
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_cluster_layout.example",
			"name",
			"HVM 1.3 Cluster on HVM/Ubuntu 24.04",
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + dataSourceConfig,
				Check:  checkFn,
			},
		},
	})
}

// TestAccMorpheusFindClusterLayoutById reads a cluster layout by id. Layout ids
// differ from one appliance to the next, so the id is not hard-coded: the
// hpe_morpheus_cluster_layouts data source resolves an HVM cluster layout
// (code "mvm-..."; every one of them carries a cluster_version), the singular
// data source reads that layout back by its id, and the id it returns is
// checked against the id that was requested, so the lookup itself is proved
// rather than just a read of some HVM layout.
func TestAccMorpheusFindClusterLayoutById(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()

	const pluralConfig = `
data "hpe_morpheus_cluster_layouts" "hvm" {
  filter {
    name   = "code"
    values = ["^mvm-"]
  }
}
`

	dataSourceConfig, err := clusterlayout.RenderClusterLayoutDataSourceByIDConfig(
		t,
		map[string]string{
			"Id": "tolist(data.hpe_morpheus_cluster_layouts.hvm.cluster_layouts)[0].id",
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	checks := []resource.TestCheckFunc{
		// The layout read back is the one whose id was requested.
		resource.TestCheckResourceAttrPair(
			"data.hpe_morpheus_cluster_layout.example", "id",
			"data.hpe_morpheus_cluster_layouts.hvm", "cluster_layouts.0.id",
		),
		resource.TestCheckResourceAttrPair(
			"data.hpe_morpheus_cluster_layout.example", "code",
			"data.hpe_morpheus_cluster_layouts.hvm", "cluster_layouts.0.code",
		),
		resource.TestCheckResourceAttrSet(
			"data.hpe_morpheus_cluster_layout.example",
			"name",
		),
		resource.TestMatchResourceAttr(
			"data.hpe_morpheus_cluster_layout.example",
			"code",
			regexp.MustCompile(`^mvm-`),
		),
		resource.TestCheckResourceAttrSet(
			"data.hpe_morpheus_cluster_layout.example",
			"cluster_version",
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + pluralConfig + dataSourceConfig,
				Check:  checkFn,
			},
		},
	})
}

func TestAccMorpheusFindClusterLayoutNotFound(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()
	config := providerConfig + `
      data "hpe_morpheus_cluster_layout" "test" {
        name = "______"
      }`

	checks := []resource.TestCheckFunc{
		resource.TestCheckNoResourceAttr(
			"data.hpe_morpheus_cluster_layout.test",
			"id",
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)

	expected := regexp.MustCompile(`no cluster layout found`)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config:      config,
				Check:       checkFn,
				ExpectError: expected,
			},
		},
	})
}

func TestAccMorpheusFindClusterLayoutNoSearchAttrs(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)

	t.Parallel()

	// A real connection is used so the data source Read runs and returns the
	// "no valid search terms" error; with an unconfigured provider the mux
	// provider fails earlier with a connection error and the validation path is
	// never reached.
	config := testhelpers.ProviderBlock() + `
      data "hpe_morpheus_cluster_layout" "test" {
      }`

	checks := []resource.TestCheckFunc{
		resource.TestCheckNoResourceAttr(
			"data.hpe_morpheus_cluster_layout.test",
			"id",
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)

	expected := clusterlayout.ErrorNoValidSearchTerms

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config:      config,
				Check:       checkFn,
				ExpectError: regexp.MustCompile(expected),
			},
		},
	})
}

func TestAccMorpheusFindClusterLayoutBothSearchAttrs(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)

	t.Parallel()

	config := providerConfigOffline + `
      data "hpe_morpheus_cluster_layout" "test" {
        id = 7
        name = "______"
      }`

	checks := []resource.TestCheckFunc{
		resource.TestCheckNoResourceAttr(
			"data.hpe_morpheus_cluster_layout.test",
			"id",
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)

	expected := clusterlayout.ErrorRunningPreApply

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config:      config,
				Check:       checkFn,
				ExpectError: regexp.MustCompile(expected),
			},
		},
	})
}
