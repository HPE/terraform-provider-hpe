// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkdomain_test

import (
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus/framework/datasources/networkdomain"
	sdkv2morpheus "github.com/HPE/terraform-provider-hpe/morpheus/sdkv2"
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

// networkDomainResourceConfig returns HCL for a private network domain to use
// as a per-run dependency.
//
// MORPH-16523: the data source tests previously looked up a static
// master-tenant domain (id 193 / "acceptance test.localdomain") that does not
// exist as a sub-tenant. Create the domain per-run instead. visibility MUST be
// "private" (public triggers MORPH-16419 churn) and no tenant_id is set
// (a hardcoded tenant_id is MORPH-16416's defect).
func networkDomainResourceConfig(name string) string {
	return `
resource "hpe_morpheus_network_domain" "example" {
  name        = "` + name + `"
  description = "Terraform acceptance test network domain"
  public_zone = false
  visibility  = "private"
  active      = true
}
`
}

func TestAccMorpheusFindNetworkDomainByName(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()

	// The server lower-cases the domain name; use a lower-case fixture name so
	// there is no perpetual plan diff (avoids the MORPH-16419 churn).
	name := strings.ToLower(acctest.RandomWithPrefix(t.Name()))
	dependenciesConfig := networkDomainResourceConfig(name)

	dataSourceConfig, err := networkdomain.RenderNetworkDomainByNameConfig(t,
		map[string]string{
			"Name": "hpe_morpheus_network_domain.example.name",
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	checks := networkDomainChecks()
	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), sdkv2morpheus.Provider()),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + dependenciesConfig + dataSourceConfig,
				Check:  checkFn,
			},
		},
	})
}

func TestAccMorpheusFindNetworkDomainById(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()

	// The server lower-cases the domain name; use a lower-case fixture name so
	// there is no perpetual plan diff (avoids the MORPH-16419 churn).
	name := strings.ToLower(acctest.RandomWithPrefix(t.Name()))
	dependenciesConfig := networkDomainResourceConfig(name)

	dataSourceConfig, err := networkdomain.RenderNetworkDomainByIdConfig(t,
		map[string]string{
			"Id": "hpe_morpheus_network_domain.example.id",
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttrPair(
			"data.hpe_morpheus_network_domain.example",
			"id",
			"hpe_morpheus_network_domain.example",
			"id",
		),
		resource.TestCheckResourceAttrPair(
			"data.hpe_morpheus_network_domain.example",
			"name",
			"hpe_morpheus_network_domain.example",
			"name",
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), sdkv2morpheus.Provider()),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + dependenciesConfig + dataSourceConfig,
				Check:  checkFn,
			},
		},
	})
}

func TestAccMorpheusFindNetworkDomainNotFound(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()

	dataSourceConfig, err := networkdomain.RenderNetworkDomainByNameConfig(t,
		map[string]string{
			"Name": "\"______\"",
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	expected := regexp.MustCompile(`no network domain found`)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config:      providerConfig + dataSourceConfig,
				ExpectError: expected,
			},
		},
	})
}

func TestAccMorpheusFindNetworkDomainNoSearchAttrs(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)

	t.Parallel()

	// A real connection is used so the data source Read runs and returns the
	// "no valid search terms" error; with an unconfigured provider the mux
	// provider fails earlier with a connection error and the validation path is
	// never reached.
	config := testhelpers.ProviderBlock() + `
      data "hpe_morpheus_network_domain" "test" {
      }`

	expected := networkdomain.ErrorNoValidSearchTerms

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

func TestAccMorpheusFindNetworkDomainBothSearchAttrs(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)

	t.Parallel()

	config := providerConfigOffline + `
      data "hpe_morpheus_network_domain" "test" {
        id   = 1
        name = "______"
      }`

	expected := networkdomain.ErrorRunningPreApply

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

func networkDomainChecks() []resource.TestCheckFunc {
	ds := "data.hpe_morpheus_network_domain.example"

	return []resource.TestCheckFunc{
		resource.TestCheckResourceAttrSet(ds, "id"),
		// The server normalises the domain name (e.g. lower-cases it), so
		// compare against the resource's stored name rather than the config.
		resource.TestCheckResourceAttrPair(
			ds, "name",
			"hpe_morpheus_network_domain.example", "name",
		),
		resource.TestCheckResourceAttrSet(ds, "visibility"),
	}
}
