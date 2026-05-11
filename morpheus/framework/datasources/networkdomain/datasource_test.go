// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkdomain_test

import (
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/datasources/networkdomain"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/systemoverride"
)

func TestMain(m *testing.M) {
	systemoverride.ParseFlags()
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

func TestAccMorpheusFindNetworkDomainByName(t *testing.T) {
	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	testSystem := systemoverride.GetPreferred(t, "zodiac")
	providerConfig := testhelpers.ProviderBlockForServer(testSystem)

	dataSourceConfig, err := networkdomain.RenderNetworkDomainByNameConfig(t, nil)
	if err != nil {
		t.Fatal(err)
	}

	checks := networkDomainChecks()
	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + dataSourceConfig,
				Check:  checkFn,
			},
		},
	})
}

func TestAccMorpheusFindNetworkDomainById(t *testing.T) {
	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	testSystem := systemoverride.GetPreferred(t, "zodiac")
	providerConfig := testhelpers.ProviderBlockForServer(testSystem)

	dataSourceConfig, err := networkdomain.RenderNetworkDomainByIdConfig(t, nil)
	if err != nil {
		t.Fatal(err)
	}

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttrSet(
			"data.hpe_morpheus_network_domain.example",
			"id",
		),
		resource.TestCheckResourceAttrSet(
			"data.hpe_morpheus_network_domain.example",
			"name",
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + dataSourceConfig,
				Check:  checkFn,
			},
		},
	})
}

func TestAccMorpheusFindNetworkDomainNotFound(t *testing.T) {
	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	testSystem := systemoverride.GetPreferred(t, "zodiac")
	providerConfig := testhelpers.ProviderBlockForServer(testSystem)

	dataSourceConfig, err := networkdomain.RenderNetworkDomainByNameConfig(t,
		map[string]string{
			"Name": "______",
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	expected := regexp.MustCompile(`no network domain found`)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
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

	t.Parallel()

	config := providerConfigOffline + `
      data "hpe_morpheus_network_domain" "test" {
      }`

	expected := networkdomain.ErrorNoValidSearchTerms

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
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

	t.Parallel()

	config := providerConfigOffline + `
      data "hpe_morpheus_network_domain" "test" {
        id   = 1
        name = "______"
      }`

	expected := networkdomain.ErrorRunningPreApply

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
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
		resource.TestCheckResourceAttr(ds, "name", "acceptance test.localdomain"),
		resource.TestCheckResourceAttrSet(ds, "visibility"),
	}
}
