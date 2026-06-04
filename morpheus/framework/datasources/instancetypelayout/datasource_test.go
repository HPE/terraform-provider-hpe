// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package instancetypelayout_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/datasources/instancetypelayout"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/getsafe"
	"github.com/HPE/terraform-provider-hpe/provider"
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

func newProviderWithError() (tfprotov6.ProviderServer, error) {
	providerInstance := provider.New("test", morpheus.New())()

	return providerserver.NewProtocol6WithError(providerInstance)()
}

var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"hpe": newProviderWithError,
}

func TestAccMorpheusFindInstanceTypeLayoutById(t *testing.T) {
	if capabilities.Missing(t, capabilities.All) {
		t.Log("Skipping test due to missing capabilities")

		return
	}
	defer testhelpers.RecordResult(t)
	t.Parallel()

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	layouts, err := testhelpers.CreateInstanceTypeLayout(t, 1)
	if err != nil || len(layouts) == 0 {
		t.Fatal(err)
	}

	layout := layouts[0]

	t.Cleanup(func() {
		testhelpers.DeleteInstanceTypeLayout(t, getsafe.GetSafe(layout.Id))
	})

	layoutID := fmt.Sprintf("%d", getsafe.GetSafe(layout.Id))

	layoutName := getsafe.GetSafe(layout.Name)

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_instance_type_layout.example",
			"name",
			layoutName,
		),
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_instance_type_layout.example",
			"id",
			layoutID,
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)

	providerConfig := testhelpers.ProviderBlock()

	dataSourceConfig, err := testhelpers.RenderExample(t, "example-id.tf.tmpl", "Id", layoutID)
	if err != nil {
		t.Fatal(err)
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + dataSourceConfig,
				Check:  checkFn,
			},
		},
	})
}

func TestAccMorpheusFindInstanceTypeLayoutByName(t *testing.T) {
	if capabilities.Missing(t, capabilities.All) {
		t.Log("Skipping test due to missing capabilities")

		return
	}
	defer testhelpers.RecordResult(t)
	t.Parallel()

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	layouts, err := testhelpers.CreateInstanceTypeLayout(t, 1)
	if err != nil || len(layouts) == 0 {
		t.Fatal(err)
	}

	layout := layouts[0]

	t.Cleanup(func() {
		testhelpers.DeleteInstanceTypeLayout(t, getsafe.GetSafe(layout.Id))
	})

	layoutID := fmt.Sprintf("%d", getsafe.GetSafe(layout.Id))
	layoutName := getsafe.GetSafe(layout.Name)

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_instance_type_layout.example",
			"name",
			layoutName,
		),
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_instance_type_layout.example",
			"id",
			layoutID,
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)

	providerConfig := testhelpers.ProviderBlock()

	dataSourceConfig, err := testhelpers.RenderExample(t, "example-name.tf.tmpl", "Name", layoutName)
	if err != nil {
		t.Fatal(err)
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + dataSourceConfig,
				Check:  checkFn,
			},
		},
	})
}

func TestAccMorpheusFindInstanceTypeLayoutByNameAndVersion(t *testing.T) {
	if capabilities.Missing(t, capabilities.All) {
		t.Log("Skipping test due to missing capabilities")

		return
	}
	defer testhelpers.RecordResult(t)
	t.Parallel()

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	layouts, err := testhelpers.CreateInstanceTypeLayout(t, 1)
	if err != nil || len(layouts) == 0 {
		t.Fatal(err)
	}

	layout := layouts[0]

	t.Cleanup(func() {
		testhelpers.DeleteInstanceTypeLayout(t, getsafe.GetSafe(layout.Id))
	})

	layoutID := fmt.Sprintf("%d", getsafe.GetSafe(layout.Id))
	layoutName := getsafe.GetSafe(layout.Name)
	layoutVersion := getsafe.GetSafe(layout.InstanceVersion)

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_instance_type_layout.example",
			"name",
			layoutName,
		),
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_instance_type_layout.example",
			"version",
			layoutVersion,
		),
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_instance_type_layout.example",
			"id",
			layoutID,
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)

	providerConfig := testhelpers.ProviderBlock()

	dataSourceConfig, err := testhelpers.RenderExample(t, "example-name-version.tf.tmpl",
		"Name", layoutName, "Version", layoutVersion)
	if err != nil {
		t.Fatal(err)
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + dataSourceConfig,
				Check:  checkFn,
			},
		},
	})
}

func TestAccMorpheusFindInstanceTypeLayoutSortOrder(t *testing.T) {
	if capabilities.Missing(t, capabilities.All) {
		t.Log("Skipping test due to missing capabilities")

		return
	}
	defer testhelpers.RecordResult(t)
	t.Parallel()

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	duplicateCount := 3

	layouts, err := testhelpers.CreateInstanceTypeLayout(t, int64(duplicateCount))
	if err != nil || len(layouts) == 0 {
		t.Fatal(err)
	}

	for _, layout := range layouts {
		t.Cleanup(func() {
			testhelpers.DeleteInstanceTypeLayout(t, getsafe.GetSafe(layout.Id))
		})
	}

	layoutID := fmt.Sprintf("%d", getsafe.GetSafe(layouts[len(layouts)-1].Id))

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_instance_type_layout.test",
			"sort_order",
			"2",
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)

	providerConfig := testhelpers.ProviderBlock()

	dataSourceConfig := providerConfig + `
      data "hpe_morpheus_instance_type_layout" "test" {
        id = ` + layoutID + `
      }`

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: dataSourceConfig,
				Check:  checkFn,
			},
		},
	})
}

func TestAccMorpheusFindInstanceLayoutNotFound(t *testing.T) {
	if capabilities.Missing(t, capabilities.All) {
		t.Log("Skipping test due to missing capabilities")

		return
	}
	defer testhelpers.RecordResult(t)
	t.Parallel()

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	checks := []resource.TestCheckFunc{
		resource.TestCheckNoResourceAttr(
			"data.hpe_morpheus_instance_type_layout.test",
			"id",
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)

	providerConfig := testhelpers.ProviderBlock()

	dataSourceConfig := providerConfig + `
      data "hpe_morpheus_instance_type_layout" "test" {
        name = "______" 
      }`

	expected := instancetypelayout.ErrorNoInstanceTypeLayoutFound

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      dataSourceConfig,
				Check:       checkFn,
				ExpectError: regexp.MustCompile(expected),
			},
		},
	})
}

func TestAccMorpheusFindInstanceTypeLayoutByNameAndVersionNoSearchAttrs(t *testing.T) {
	if capabilities.Missing(t, capabilities.All) {
		t.Log("Skipping test due to missing capabilities")

		return
	}
	defer testhelpers.RecordResult(t)
	t.Parallel()

	checks := []resource.TestCheckFunc{
		resource.TestCheckNoResourceAttr(
			"data.hpe_morpheus_instance_type_layout.test",
			"id",
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)

	dataSourceConfig := providerConfigOffline + `
      data "hpe_morpheus_instance_type_layout" "test" {
      }`

	expected := instancetypelayout.ErrorNoValidSearchTerms

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      dataSourceConfig,
				Check:       checkFn,
				ExpectError: regexp.MustCompile(expected),
			},
		},
	})
}

func TestAccMorpheusFindInstanceLayoutWithIdAndName(t *testing.T) {
	if capabilities.Missing(t, capabilities.All) {
		t.Log("Skipping test due to missing capabilities")

		return
	}
	defer testhelpers.RecordResult(t)
	t.Parallel()

	checks := []resource.TestCheckFunc{
		resource.TestCheckNoResourceAttr(
			"data.hpe_morpheus_instance_type_layout.test",
			"id",
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)

	dataSourceConfig := providerConfigOffline + `
      data "hpe_morpheus_instance_type_layout" "test" {
        id = 1
        name = "______" 
      }`

	expected := instancetypelayout.ErrorRunningPreApply

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      dataSourceConfig,
				Check:       checkFn,
				ExpectError: regexp.MustCompile(expected),
			},
		},
	})
}

func TestAccMorpheusFindInstanceLayoutWithIdAndVersion(t *testing.T) {
	if capabilities.Missing(t, capabilities.All) {
		t.Log("Skipping test due to missing capabilities")

		return
	}
	defer testhelpers.RecordResult(t)
	t.Parallel()

	checks := []resource.TestCheckFunc{
		resource.TestCheckNoResourceAttr(
			"data.hpe_morpheus_instance_type_layout.test",
			"id",
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)

	dataSourceConfig := providerConfigOffline + `
      data "hpe_morpheus_instance_type_layout" "test" {
        id = 1
        version = "123" 
      }`

	expected := instancetypelayout.ErrorRunningPreApply

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      dataSourceConfig,
				Check:       checkFn,
				ExpectError: regexp.MustCompile(expected),
			},
		},
	})
}
