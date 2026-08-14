// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkproxy_test

import (
	"os"
	"regexp"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus/framework/datasources/networkproxy"
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

// TestAccMorpheusFindNetworkProxyByName creates its own network proxy fixture
// (there is no hpe_morpheus_network_proxy resource) and looks it up by name,
// so the test is self-contained and needs no environment configuration.
func TestAccMorpheusFindNetworkProxyByName(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	proxyID, proxyName := testhelpers.CreateNetworkProxy(t)
	t.Cleanup(func() { testhelpers.DeleteNetworkProxy(t, proxyID) })

	providerConfig := testhelpers.ProviderBlock()

	dataSourceConfig, err := testhelpers.RenderExample(t,
		"example-name.tf.tmpl",
		"Name", proxyName,
	)
	if err != nil {
		t.Fatal(err)
	}

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_network_proxy.example",
			"id",
			strconv.FormatInt(proxyID, 10),
		),
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_network_proxy.example",
			"name",
			proxyName,
		),
		resource.TestCheckResourceAttrSet(
			"data.hpe_morpheus_network_proxy.example",
			"proxy_host",
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

// TestAccMorpheusFindNetworkProxyById creates its own fixture and looks it up by
// id.
func TestAccMorpheusFindNetworkProxyById(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	proxyID, proxyName := testhelpers.CreateNetworkProxy(t)
	t.Cleanup(func() { testhelpers.DeleteNetworkProxy(t, proxyID) })

	providerConfig := testhelpers.ProviderBlock()

	dataSourceConfig, err := testhelpers.RenderExample(t,
		"example-id.tf.tmpl",
		"Id", strconv.FormatInt(proxyID, 10),
	)
	if err != nil {
		t.Fatal(err)
	}

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_network_proxy.example",
			"id",
			strconv.FormatInt(proxyID, 10),
		),
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_network_proxy.example",
			"name",
			proxyName,
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

func TestAccMorpheusFindNetworkProxyNotFound(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()
	config := providerConfig + `
      data "hpe_morpheus_network_proxy" "test" {
        name = "______"
      }`

	checks := []resource.TestCheckFunc{
		resource.TestCheckNoResourceAttr(
			"data.hpe_morpheus_network_proxy.test",
			"id",
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)

	expected := regexp.MustCompile(`no network proxy found`)

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

func TestAccMorpheusFindNetworkProxyNoSearchAttrs(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)

	t.Parallel()

	// A real connection is used so the data source Read runs and returns the
	// "no valid search terms" error; with an unconfigured provider the mux
	// provider fails earlier with a connection error and the validation path is
	// never reached.
	config := testhelpers.ProviderBlock() + `
      data "hpe_morpheus_network_proxy" "test" {
      }`

	checks := []resource.TestCheckFunc{
		resource.TestCheckNoResourceAttr(
			"data.hpe_morpheus_network_proxy.test",
			"id",
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)

	expected := networkproxy.ErrorNoValidSearchTerms

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

func TestAccMorpheusFindNetworkProxyBothSearchAttrs(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)

	t.Parallel()

	config := providerConfigOffline + `
      data "hpe_morpheus_network_proxy" "test" {
        id = 1
        name = "______"
      }`

	checks := []resource.TestCheckFunc{
		resource.TestCheckNoResourceAttr(
			"data.hpe_morpheus_network_proxy.test",
			"id",
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)

	expected := networkproxy.ErrorRunningPreApply

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
