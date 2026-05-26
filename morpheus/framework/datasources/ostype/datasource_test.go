// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package ostype_test

import (
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/datasources/ostype"
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

func TestAccMorpheusFindOsTypeByName(t *testing.T) {
	if capabilities.Missing(t, capabilities.All) {
		t.Log("Skipping test due to missing capabilities")

		return
	}
	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()

	dataSourceConfig, err := testhelpers.RenderExample(t,
		"example-name.tf.tmpl",
		"Name", "debian 12 64-bit",
	)
	if err != nil {
		t.Fatal(err)
	}

	checks := osTypeChecks()

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)

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

func TestAccMorpheusFindOsTypeById(t *testing.T) {
	if capabilities.Missing(t, capabilities.All) {
		t.Log("Skipping test due to missing capabilities")

		return
	}
	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()

	dataSourceConfig, err := testhelpers.RenderExample(t,
		"example-id.tf.tmpl",
		"Id", "1",
	)
	if err != nil {
		t.Fatal(err)
	}

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_os_type.example",
			"id",
			"1",
		),
		resource.TestCheckResourceAttrSet(
			"data.hpe_morpheus_os_type.example",
			"name",
		),
		resource.TestCheckResourceAttrSet(
			"data.hpe_morpheus_os_type.example",
			"code",
		),
		resource.TestCheckResourceAttrSet(
			"data.hpe_morpheus_os_type.example",
			"platform",
		),
		resource.TestCheckResourceAttrSet(
			"data.hpe_morpheus_os_type.example",
			"bit_count",
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)

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

func TestAccMorpheusFindOsTypeNotFound(t *testing.T) {
	if capabilities.Missing(t, capabilities.All) {
		t.Log("Skipping test due to missing capabilities")

		return
	}
	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()
	config := providerConfig + `
      data "hpe_morpheus_os_type" "test" {
        name = "______"
      }`

	checks := []resource.TestCheckFunc{
		resource.TestCheckNoResourceAttr(
			"data.hpe_morpheus_os_type.test",
			"id",
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)

	expected := regexp.MustCompile(`no os type found`)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
		Steps: []resource.TestStep{
			{
				Config:      config,
				Check:       checkFn,
				ExpectError: expected,
			},
		},
	})
}

func TestAccMorpheusFindOsTypeNoSearchAttrs(t *testing.T) {
	if capabilities.Missing(t, capabilities.All) {
		t.Log("Skipping test due to missing capabilities")

		return
	}
	defer testhelpers.RecordResult(t)

	t.Parallel()

	config := providerConfigOffline + `
      data "hpe_morpheus_os_type" "test" {
      }`

	checks := []resource.TestCheckFunc{
		resource.TestCheckNoResourceAttr(
			"data.hpe_morpheus_os_type.test",
			"id",
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)

	expected := ostype.ErrorNoValidSearchTerms

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

func TestAccMorpheusFindOsTypeBothSearchAttrs(t *testing.T) {
	if capabilities.Missing(t, capabilities.All) {
		t.Log("Skipping test due to missing capabilities")

		return
	}
	defer testhelpers.RecordResult(t)

	t.Parallel()

	config := providerConfigOffline + `
      data "hpe_morpheus_os_type" "test" {
        id = 1
        name = "______"
      }`

	checks := []resource.TestCheckFunc{
		resource.TestCheckNoResourceAttr(
			"data.hpe_morpheus_os_type.test",
			"id",
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)

	expected := ostype.ErrorRunningPreApply

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

func osTypeChecks() []resource.TestCheckFunc {
	ds := "data.hpe_morpheus_os_type.example"

	return []resource.TestCheckFunc{
		resource.TestCheckResourceAttrSet(ds, "id"),
		resource.TestCheckResourceAttr(ds, "name", "debian 12 64-bit"),
		resource.TestCheckResourceAttr(ds, "code", "debian.12.64"),
		resource.TestCheckResourceAttr(ds, "platform", "linux"),
		resource.TestCheckResourceAttr(ds, "category", "debian"),
		resource.TestCheckResourceAttr(ds, "vendor", "debian"),
		resource.TestCheckResourceAttr(ds, "os_name", "debian"),
		resource.TestCheckResourceAttr(ds, "os_version", "12"),
		resource.TestCheckResourceAttr(ds, "os_family", "debian"),
		resource.TestCheckResourceAttr(ds, "bit_count", "64"),
		resource.TestCheckResourceAttr(ds, "install_agent", "true"),
		resource.TestCheckResourceAttr(ds, "owner", "System"),
		resource.TestCheckNoResourceAttr(ds, "description"),
		resource.TestCheckNoResourceAttr(ds, "cloud_init_version"),
		resource.TestCheckNoResourceAttr(ds, "os_codename"),
		resource.TestCheckResourceAttrSet(ds, "images.#"),
	}
}
