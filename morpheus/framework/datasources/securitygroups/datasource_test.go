// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package securitygroups_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
)

func TestMain(m *testing.M) {
	code := m.Run()
	testhelpers.WriteMergedResults()
	os.Exit(code)
}

func TestAccMorpheusSecurityGroupsFilterByName(t *testing.T) {
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

	name := acctest.RandomWithPrefix("tf-acc-sgs-ds")

	createConfig := providerConfig + `
resource "hpe_morpheus_security_group" "test" {
  name = "` + name + `"
}`

	dataSourceConfig := createConfig + `
data "hpe_morpheus_security_groups" "example" {
  filter {
    name   = "name"
    values = ["^` + name + `$"]
  }
}`

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_security_groups.example", "security_groups.#", "1",
		),
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_security_groups.example", "security_groups.0.name", name,
		),
		resource.TestCheckResourceAttrSet(
			"data.hpe_morpheus_security_groups.example", "security_groups.0.id",
		),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
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

// TestAccMorpheusSecurityGroupsMultiFilter exercises multiple filter blocks,
// which are ANDed together. A new security group defaults to private visibility.
func TestAccMorpheusSecurityGroupsMultiFilter(t *testing.T) {
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

	name := acctest.RandomWithPrefix("tf-acc-sgs-ds")

	createConfig := providerConfig + `
resource "hpe_morpheus_security_group" "test" {
  name = "` + name + `"
}`

	dataSourceConfig := createConfig + `
data "hpe_morpheus_security_groups" "example" {
  filter {
    name   = "name"
    values = ["^` + name + `$"]
  }

  filter {
    name   = "visibility"
    values = ["private"]
  }
}`

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_security_groups.example", "security_groups.#", "1",
		),
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_security_groups.example", "security_groups.0.visibility", "private",
		),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
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

func TestAccMorpheusSecurityGroupsEmptyResult(t *testing.T) {
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
      data "hpe_morpheus_security_groups" "test" {
        filter {
          name   = "name"
          values = ["this-name-should-not-exist-______"]
        }
      }`

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_security_groups.test", "security_groups.#", "0",
		),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check:  resource.ComposeAggregateTestCheckFunc(checks...),
			},
		},
	})
}
