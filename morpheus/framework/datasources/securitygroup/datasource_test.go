// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package securitygroup_test

import (
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/datasources/securitygroup"
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

func TestAccMorpheusFindSecurityGroupByName(t *testing.T) {
	if capabilities.Missing(t, capabilities.NetworkFirewall) {
		t.Log("Skipping test due to missing capabilities")

		return
	}

	t.Parallel()
	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()

	dataSourceConfig, err := securitygroup.RenderSecurityGroupDataSourceByNameConfig(
		t, map[string]string{"Name": "MorpheusUbuntu-nsg"},
	)
	if err != nil {
		t.Fatal(err)
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + dataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.hpe_morpheus_security_group.example", "id"),
					resource.TestCheckResourceAttr("data.hpe_morpheus_security_group.example", "name", "MorpheusUbuntu-nsg"),
				),
			},
		},
	})
}

func TestAccMorpheusFindSecurityGroupById(t *testing.T) {
	if capabilities.Missing(t, capabilities.NetworkFirewall) {
		t.Log("Skipping test due to missing capabilities")

		return
	}

	t.Parallel()
	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()

	dataSourceConfig, err := securitygroup.RenderSecurityGroupDataSourceByIDConfig(t, map[string]string{"Id": "31"})
	if err != nil {
		t.Fatal(err)
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + dataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.hpe_morpheus_security_group.example", "name"),
					resource.TestCheckResourceAttr("data.hpe_morpheus_security_group.example", "id", "31"),
				),
			},
		},
	})
}

func TestAccMorpheusFindSecurityGroupNotFound(t *testing.T) {
	if capabilities.Missing(t, capabilities.NetworkFirewall) {
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
data "hpe_morpheus_security_group" "test" {
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

func TestAccMorpheusFindSecurityGroupNoSearchAttrs(t *testing.T) {
	if capabilities.Missing(t, capabilities.NetworkFirewall) {
		t.Log("Skipping test due to missing capabilities")

		return
	}

	t.Parallel()
	defer testhelpers.RecordResult(t)

	config := providerConfigOffline + `
data "hpe_morpheus_security_group" "test" {
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

func TestAccMorpheusFindSecurityGroupBothSearchAttrs(t *testing.T) {
	if capabilities.Missing(t, capabilities.NetworkFirewall) {
		t.Log("Skipping test due to missing capabilities")

		return
	}

	t.Parallel()
	defer testhelpers.RecordResult(t)

	config := providerConfigOffline + `
data "hpe_morpheus_security_group" "test" {
  id   = 1
  name = "MorpheusUbuntu-nsg"
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
