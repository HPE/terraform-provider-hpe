// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package plugin_test

import (
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	sdkv2morpheus "github.com/HPE/terraform-provider-hpe/morpheus/sdkv2"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
	"github.com/HPE/terraform-provider-hpe/provider/adapter"
)

// A prepackaged plugin present on the target appliance (Morpheus 9.1.0.2).
const (
	testPluginCode = "morpheus-infoblox-plugin"
)

func TestMain(m *testing.M) {
	code := testhelpers.TestMain(m)
	testhelpers.WriteMergedResults()
	os.Exit(code)
}

func TestAccMorpheusPluginDataSourceByCodeOk(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.Plugins)

	t.Parallel()

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()

	config := providerConfig + `
data "hpe_morpheus_plugin" "by_code" {
  code = "` + testPluginCode + `"
}
`

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttrSet(
			"data.hpe_morpheus_plugin.by_code", "id"),
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_plugin.by_code", "code", testPluginCode),
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_plugin.by_code", "status", "loaded"),
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_plugin.by_code", "is_prepackaged", "true"),
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_plugin.by_code", "enabled", "true"),
		resource.TestCheckResourceAttrSet(
			"data.hpe_morpheus_plugin.by_code", "name"),
		resource.TestCheckResourceAttrSet(
			"data.hpe_morpheus_plugin.by_code", "version"),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), sdkv2morpheus.Provider()),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check:  resource.ComposeAggregateTestCheckFunc(checks...),
			},
		},
	})
}

func TestAccMorpheusPluginDataSourceByIdOk(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.Plugins)

	t.Parallel()

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()

	config := providerConfig + `
data "hpe_morpheus_plugin" "by_code" {
  code = "` + testPluginCode + `"
}

data "hpe_morpheus_plugin" "by_id" {
  id = data.hpe_morpheus_plugin.by_code.id
}
`

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttrSet(
			"data.hpe_morpheus_plugin.by_id", "id"),
		resource.TestCheckResourceAttrPair(
			"data.hpe_morpheus_plugin.by_id", "code",
			"data.hpe_morpheus_plugin.by_code", "code"),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), sdkv2morpheus.Provider()),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check:  resource.ComposeAggregateTestCheckFunc(checks...),
			},
		},
	})
}

func TestAccMorpheusPluginDataSourceNotFoundErr(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.Plugins)

	t.Parallel()

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	config := testhelpers.ProviderBlock() + `
data "hpe_morpheus_plugin" "missing" {
  code = "does-not-exist-tf-acc"
}
`

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), sdkv2morpheus.Provider()),
		Steps: []resource.TestStep{
			{
				Config:      config,
				ExpectError: regexp.MustCompile("no plugin found"),
			},
		},
	})
}

// TestAccMorpheusPluginDataSourceConflict verifies that setting both id and
// code triggers the schema-level ConflictsWith validation error.
func TestAccMorpheusPluginDataSourceConflict(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.Plugins)

	t.Parallel()

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	config := testhelpers.ProviderBlock() + `
data "hpe_morpheus_plugin" "conflict" {
  id   = 1
  code = "` + testPluginCode + `"
}
`

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), sdkv2morpheus.Provider()),
		Steps: []resource.TestStep{
			{
				Config:      config,
				ExpectError: regexp.MustCompile("(?s)(Invalid Attribute Combination|cannot be specified when)"),
			},
		},
	})
}

// TestAccMorpheusPluginDataSourceNoSearchAttrs verifies that providing neither
// id nor code produces the "no valid search terms" error. A real connection is
// used so the data source Read runs and returns the validation error.
func TestAccMorpheusPluginDataSourceNoSearchAttrs(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.Plugins)

	t.Parallel()

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	config := testhelpers.ProviderBlock() + `
data "hpe_morpheus_plugin" "empty" {
}
`

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), sdkv2morpheus.Provider()),
		Steps: []resource.TestStep{
			{
				Config:      config,
				ExpectError: regexp.MustCompile("no valid search terms"),
			},
		},
	})
}
