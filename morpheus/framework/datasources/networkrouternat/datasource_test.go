// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkrouternat_test

import (
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus/framework/datasources/networkrouternat"
	natresource "github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/networkrouternat"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/nsxt"
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

// NSX-T NAT rules require a realized tier-1 gateway; nsxt.Tier1Config renders
// one by resolving the pre-provisioned tier-0 (and its group, network
// integration and edge cluster) by name.

func TestAccMorpheusFindNetworkRouterNatByName(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.NetworkRouter)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()
	name := acctest.RandomWithPrefix(t.Name())

	routerConfig := nsxt.Tier1Config(name)

	resourceConfig, err := natresource.RenderNetworkRouterNatConfig(t, map[string]string{
		"RouterId": nsxt.Tier1RouterIDRef,
		"Name":     name,
	})
	if err != nil {
		t.Fatal(err)
	}

	dataSourceConfig := `
data "hpe_morpheus_network_router_nat" "example" {
  name       = "` + name + `"
  router_id  = ` + nsxt.Tier1RouterIDRef + `
  depends_on = [hpe_morpheus_network_router_nat.example]
}
`

	checkFn := resource.ComposeAggregateTestCheckFunc(natChecks()...)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + routerConfig + resourceConfig + dataSourceConfig,
				Check:  checkFn,
			},
		},
	})
}

func TestAccMorpheusFindNetworkRouterNatById(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.NetworkRouter)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()
	name := acctest.RandomWithPrefix(t.Name())

	routerConfig := nsxt.Tier1Config(name)

	resourceConfig, err := natresource.RenderNetworkRouterNatConfig(t, map[string]string{
		"RouterId": nsxt.Tier1RouterIDRef,
		"Name":     name,
	})
	if err != nil {
		t.Fatal(err)
	}

	dataSourceConfig, err := networkrouternat.RenderNatByIdConfig(t, map[string]string{
		"Id":       "hpe_morpheus_network_router_nat.example.id",
		"RouterId": nsxt.Tier1RouterIDRef,
	})
	if err != nil {
		t.Fatal(err)
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(natChecks()...)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + routerConfig + resourceConfig + dataSourceConfig,
				Check:  checkFn,
			},
		},
	})
}

func TestAccMorpheusFindNetworkRouterNatNotFound(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.NetworkRouter)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()
	name := acctest.RandomWithPrefix(t.Name())

	routerConfig := nsxt.Tier1Config(name)

	dataSourceConfig := `
data "hpe_morpheus_network_router_nat" "example" {
  name      = "nonexistent-nat-name-that-should-not-exist"
  router_id = ` + nsxt.Tier1RouterIDRef + `
}
`

	expected := regexp.MustCompile(`no network router nat found`)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config:      providerConfig + routerConfig + dataSourceConfig,
				ExpectError: expected,
			},
		},
	})
}

func TestAccMorpheusFindNetworkRouterNatNoSearchAttrs(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.NetworkRouter)

	t.Parallel()

	config := providerConfigOffline + `
      data "hpe_morpheus_network_router_nat" "test" {
        router_id = 1
      }`

	expected := networkrouternat.ErrorNoValidSearchTerms

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

func natChecks() []resource.TestCheckFunc {
	ds := "data.hpe_morpheus_network_router_nat.example"

	return []resource.TestCheckFunc{
		resource.TestCheckResourceAttrSet(ds, "id"),
		resource.TestCheckResourceAttrSet(ds, "router_id"),
		resource.TestCheckResourceAttrSet(ds, "name"),
		// action is a create-only input that the NAT read API does not return,
		// so the data source cannot populate it. Follow-up: revisit whether
		// action/firewall should remain in the data source schema.
		resource.TestCheckResourceAttrSet(ds, "enabled"),
	}
}
