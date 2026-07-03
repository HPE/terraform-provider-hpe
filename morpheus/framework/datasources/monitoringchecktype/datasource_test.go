// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package monitoringchecktype_test

import (
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus/framework/datasources/monitoringchecktype"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
	"github.com/HPE/terraform-provider-hpe/provider/adapter"
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

func TestAccMorpheusFindMonitoringCheckTypeByName(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()

	dataSourceConfig, err := testhelpers.RenderExample(
		t,
		"example-name.tf.tmpl",
		"Name", "Web Check",
	)
	if err != nil {
		t.Fatal(err)
	}

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttrSet(
			"data.hpe_morpheus_monitoring_check_type.example",
			"id",
		),
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_monitoring_check_type.example",
			"name",
			"Web Check",
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

func TestAccMorpheusFindMonitoringCheckTypeById(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()

	dataSourceConfig, err := testhelpers.RenderExample(
		t,
		"example-id.tf.tmpl",
		"Id", "1",
	)
	if err != nil {
		t.Fatal(err)
	}

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_monitoring_check_type.example",
			"id",
			"1",
		),
		resource.TestCheckResourceAttrSet(
			"data.hpe_morpheus_monitoring_check_type.example",
			"name",
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

func TestAccMorpheusFindMonitoringCheckTypeNotFound(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()
	config := providerConfig + `
      data "hpe_morpheus_monitoring_check_type" "test" {
        name = "______"
      }`

	checks := []resource.TestCheckFunc{
		resource.TestCheckNoResourceAttr(
			"data.hpe_morpheus_monitoring_check_type.test",
			"id",
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)

	expected := regexp.MustCompile(`no monitoring check type found`)

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

func TestAccMorpheusFindMonitoringCheckTypeNoSearchAttrs(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)

	t.Parallel()

	config := providerConfigOffline + `
      data "hpe_morpheus_monitoring_check_type" "test" {
      }`

	checks := []resource.TestCheckFunc{
		resource.TestCheckNoResourceAttr(
			"data.hpe_morpheus_monitoring_check_type.test",
			"id",
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)

	expected := monitoringchecktype.ErrorNoValidSearchTerms

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

func TestAccMorpheusFindMonitoringCheckTypeBothSearchAttrs(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)

	t.Parallel()

	config := providerConfigOffline + `
      data "hpe_morpheus_monitoring_check_type" "test" {
        id = 1
        name = "______"
      }`

	checks := []resource.TestCheckFunc{
		resource.TestCheckNoResourceAttr(
			"data.hpe_morpheus_monitoring_check_type.test",
			"id",
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)

	expected := monitoringchecktype.ErrorRunningPreApply

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
