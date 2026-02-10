// (C) Copyright 2025-2026 Hewlett Packard Enterprise Development LP

package provisiontype_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus"
	sdkv2morpheus "github.com/HPE/terraform-provider-hpe/morpheus/sdkv2"
	dsprovisiontype "github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/datasources/provisiontype"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/systemoverride"
)

func TestMain(m *testing.M) {
	systemoverride.ParseFlags()
	code := m.Run()

	testhelpers.WriteMergedResults()

	os.Exit(code)
}

func TestAccMorpheusDataSourceProvisionTypeExampleOk(t *testing.T) {
	t.Parallel()

	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	testSystem := systemoverride.GetPreferred(t, "feature")
	providerConfig := testhelpers.ProviderBlock(testSystem)

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
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), sdkv2morpheus.Provider()),
		Steps: []resource.TestStep{
			{
				Config:             providerConfig + datasourceConfig,
				ExpectNonEmptyPlan: false,
				Check:              checkFn,
			},
		},
	})
}
