// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package securitygroup_test

import (
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus/framework/datasources/securitygroup"
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

func TestAccMorpheusFindSecurityGroupByName(t *testing.T) {
	capabilities.MustHaveOrSkip(t, capabilities.All)

	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()

	name := acctest.RandomWithPrefix("tf-acc-sg-ds")

	createConfig := providerConfig + `
resource "hpe_morpheus_security_group" "test" {
  name = "` + name + `"
}`

	dataSourceConfig := createConfig + `
data "hpe_morpheus_security_group" "example" {
  name = hpe_morpheus_security_group.test.name
}`

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttrPair(
			"data.hpe_morpheus_security_group.example", "id",
			"hpe_morpheus_security_group.test", "id",
		),
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_security_group.example", "name", name,
		),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config: createConfig,
			},
			{
				Config: dataSourceConfig,
				Check:  resource.ComposeAggregateTestCheckFunc(checks...),
			},
		},
	})
}

func TestAccMorpheusFindSecurityGroupById(t *testing.T) {
	capabilities.MustHaveOrSkip(t, capabilities.All)

	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()

	name := acctest.RandomWithPrefix("tf-acc-sg-ds")

	createConfig := providerConfig + `
resource "hpe_morpheus_security_group" "test" {
  name = "` + name + `"
}`

	dataSourceConfig := createConfig + `
data "hpe_morpheus_security_group" "example" {
  id = hpe_morpheus_security_group.test.id
}`

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttrPair(
			"data.hpe_morpheus_security_group.example", "id",
			"hpe_morpheus_security_group.test", "id",
		),
		resource.TestCheckResourceAttrSet(
			"data.hpe_morpheus_security_group.example", "name",
		),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config: createConfig,
			},
			{
				Config: dataSourceConfig,
				Check:  resource.ComposeAggregateTestCheckFunc(checks...),
			},
		},
	})
}

func TestAccMorpheusFindSecurityGroupNotFound(t *testing.T) {
	capabilities.MustHaveOrSkip(t, capabilities.All)

	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()
	config := providerConfig + `
      data "hpe_morpheus_security_group" "test" {
        name = "______"
      }`

	expected := regexp.MustCompile(securitygroup.ErrorNoSecurityGroupFound)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config:      config,
				ExpectError: expected,
			},
		},
	})
}

func TestAccMorpheusFindSecurityGroupNoSearchAttrs(t *testing.T) {
	capabilities.MustHaveOrSkip(t, capabilities.All)

	defer testhelpers.RecordResult(t)

	t.Parallel()

	// A real connection is used so the data source Read runs and returns the
	// "no valid search terms" error; with an unconfigured provider the mux
	// provider fails earlier with a connection error and the validation path is
	// never reached.
	config := testhelpers.ProviderBlock() + `
      data "hpe_morpheus_security_group" "test" {
      }`

	expected := securitygroup.ErrorNoValidSearchTerms

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config:      config,
				ExpectError: regexp.MustCompile(expected),
			},
		},
	})
}

func TestAccMorpheusFindSecurityGroupBothSearchAttrs(t *testing.T) {
	capabilities.MustHaveOrSkip(t, capabilities.All)

	defer testhelpers.RecordResult(t)

	t.Parallel()

	config := providerConfigOffline + `
      data "hpe_morpheus_security_group" "test" {
        id = 1
        name = "______"
      }`

	expected := securitygroup.ErrorRunningPreApply

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config:      config,
				ExpectError: regexp.MustCompile(expected),
			},
		},
	})
}
