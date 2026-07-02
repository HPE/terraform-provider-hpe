// (C) Copyright 2025-2026 Hewlett Packard Enterprise Development LP

package provisiontype_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	sdkv2morpheus "github.com/HPE/terraform-provider-hpe/morpheus/sdkv2"
	dsprovisiontype "github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/datasources/provisiontype"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
	"github.com/HPE/terraform-provider-hpe/provider/adapter"
)

func TestMain(m *testing.M) {
	code := m.Run()

	testhelpers.WriteMergedResults()

	os.Exit(code)
}

func TestAccMorpheusDataSourceProvisionTypeExampleOk(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)

	t.Parallel()

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()

	// We have no resource for provion type so we'll search for a system provision type
	datasourceConfig, err := dsprovisiontype.RenderProvisionTypeConfig(t, map[string]string{
		"Name": "\"KVM\"",
	})
	if err != nil {
		t.Fatal(err)
	}

	t.Log(providerConfig + datasourceConfig)

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_provision_type.example",
			"name",
			"KVM",
		),
		resource.TestCheckResourceAttrSet(
			"data.hpe_morpheus_provision_type.example",
			"id",
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), sdkv2morpheus.Provider()),
		Steps: []resource.TestStep{
			{
				Config:             providerConfig + datasourceConfig,
				ExpectNonEmptyPlan: false,
				Check:              checkFn,
			},
		},
	})
}
