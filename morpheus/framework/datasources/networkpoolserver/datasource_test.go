// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkpoolserver_test

import (
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus/framework/datasources/networkpoolserver"
	poolresource "github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/networkpoolserver"
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

// poolServerFixture renders a minimal infoblox network pool server resource
// labelled hpe_morpheus_network_pool_server.example.
func poolServerFixture(t *testing.T, name string) string {
	t.Helper()

	cfg, err := poolresource.RenderNetworkPoolServerInfobloxConfig(t, map[string]string{
		"Name": name,
	})
	if err != nil {
		t.Fatal(err)
	}

	// The resource template uses "infoblox" as the resource label; rename to "example".
	cfg = strings.Replace(
		cfg,
		`resource "hpe_morpheus_network_pool_server" "infoblox" {`,
		`resource "hpe_morpheus_network_pool_server" "example" {`,
		1,
	)

	return cfg
}

func TestAccMorpheusFindNetworkPoolServerByName(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.NetworkPool)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()
	name := acctest.RandomWithPrefix(t.Name())

	dataSourceConfig := `
data "hpe_morpheus_network_pool_server" "example" {
  name       = "` + name + `"
  depends_on = [hpe_morpheus_network_pool_server.example]
}
`

	checkFn := resource.ComposeAggregateTestCheckFunc(poolServerChecks()...)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + poolServerFixture(t, name) + dataSourceConfig,
				Check:  checkFn,
			},
		},
	})
}

func TestAccMorpheusFindNetworkPoolServerById(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.NetworkPool)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()
	name := acctest.RandomWithPrefix(t.Name())

	dataSourceConfig := `
data "hpe_morpheus_network_pool_server" "example" {
  id         = hpe_morpheus_network_pool_server.example.id
  depends_on = [hpe_morpheus_network_pool_server.example]
}
`

	checkFn := resource.ComposeAggregateTestCheckFunc(poolServerChecks()...)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + poolServerFixture(t, name) + dataSourceConfig,
				Check:  checkFn,
			},
		},
	})
}

func TestAccMorpheusFindNetworkPoolServerNotFound(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.NetworkPool)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()
	config := providerConfig + `
      data "hpe_morpheus_network_pool_server" "test" {
        name = "______"
      }`

	expected := regexp.MustCompile(`no network pool server found`)

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

func TestAccMorpheusFindNetworkPoolServerNoSearchAttrs(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.NetworkPool)

	t.Parallel()

	config := providerConfigOffline + `
      data "hpe_morpheus_network_pool_server" "test" {
      }`

	expected := networkpoolserver.ErrorNoValidSearchTerms

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

func poolServerChecks() []resource.TestCheckFunc {
	ds := "data.hpe_morpheus_network_pool_server.example"

	return []resource.TestCheckFunc{
		resource.TestCheckResourceAttrSet(ds, "id"),
		resource.TestCheckResourceAttrSet(ds, "name"),
		resource.TestCheckResourceAttrSet(ds, "enabled"),
		resource.TestCheckResourceAttrSet(ds, "service_url"),
	}
}
