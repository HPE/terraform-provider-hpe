// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkinterfacetype_test

import (
	"os"
	"regexp"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
	"github.com/HPE/terraform-provider-hpe/provider/adapter"
)

func TestMain(m *testing.M) {
	code := testhelpers.TestMain(m)
	testhelpers.WriteMergedResults()
	os.Exit(code)
}

// TestAccMorpheusFindNetworkInterfaceType discovers a real (cloud, provision
// type, NIC type) triple on the appliance and asserts the data source resolves
// the same id. The fixture is discovered rather than configured, so the test is
// self-contained and needs no environment variables.
func TestAccMorpheusFindNetworkInterfaceType(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	fx := testhelpers.DiscoverNetworkInterfaceType(t)

	providerConfig := testhelpers.ProviderBlock()

	dataSourceConfig, err := testhelpers.RenderExample(t,
		"example.tf.tmpl",
		"Name", fx.Name,
		"CloudId", strconv.FormatInt(fx.CloudID, 10),
		"ProvisionTypeCode", fx.ProvisionTypeCode,
	)
	if err != nil {
		t.Fatal(err)
	}

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_network_interface_type.example",
			"id",
			strconv.FormatInt(fx.ID, 10),
		),
		resource.TestCheckResourceAttrSet(
			"data.hpe_morpheus_network_interface_type.example",
			"code",
		),
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_network_interface_type.example",
			"name",
			fx.Name,
		),
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_network_interface_type.example",
			"provision_type_code",
			fx.ProvisionTypeCode,
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

// TestAccMorpheusFindNetworkInterfaceTypeNotFound uses a discovered cloud and
// provision type (so the lookup reaches the option source) with a name that
// does not exist, and asserts the clear not-found error.
func TestAccMorpheusFindNetworkInterfaceTypeNotFound(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	fx := testhelpers.DiscoverNetworkInterfaceType(t)

	providerConfig := testhelpers.ProviderBlock()

	dataSourceConfig, err := testhelpers.RenderExample(t,
		"example.tf.tmpl",
		"Name", "__nonexistent_nic_type__",
		"CloudId", strconv.FormatInt(fx.CloudID, 10),
		"ProvisionTypeCode", fx.ProvisionTypeCode,
	)
	if err != nil {
		t.Fatal(err)
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config:      providerConfig + dataSourceConfig,
				ExpectError: regexp.MustCompile(`no network interface type found`),
			},
		},
	})
}
