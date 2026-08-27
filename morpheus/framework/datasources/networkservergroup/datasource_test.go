// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkservergroup_test

import (
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus/framework/datasources/networkservergroup"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
	"github.com/HPE/terraform-provider-hpe/provider/adapter"
)

func TestMain(m *testing.M) {
	code := testhelpers.TestMain(m)
	testhelpers.WriteMergedResults()
	os.Exit(code)
}

func TestAccMorpheusFindNetworkServerGroupByName(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.NSXT)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()

	dataSourceConfig, err := testhelpers.RenderExample(t,
		"example-name.tf.tmpl",
		"Name", "Edge_NSGroup",
	)
	if err != nil {
		t.Fatal(err)
	}

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttrSet(
			"data.hpe_morpheus_network_server_group.example",
			"id",
		),
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_network_server_group.example",
			"name",
			"Edge_NSGroup",
		),
		resource.TestCheckResourceAttrSet(
			"data.hpe_morpheus_network_server_group.example",
			"network_server_id",
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

func TestAccMorpheusFindNetworkServerGroupNotFound(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.NSXT)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()
	config := providerConfig + `
      data "hpe_morpheus_network_server_group" "test" {
        name = "______nonexistent______"
      }`

	expected := regexp.MustCompile(networkservergroup.ErrorNoNetworkServerGroupFound)

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

func TestUnitMorpheusFindNetworkServerGroupNoName(t *testing.T) {
	defer testhelpers.RecordResult(t)

	t.Parallel()

	config := testhelpers.ProviderBlockUnitTest() + `
      data "hpe_morpheus_network_server_group" "test" {
      }`

	expected := regexp.MustCompile(networkservergroup.ErrorNoValidSearchTerms)

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config:      config,
				ExpectError: expected,
			},
		},
	})
}
