// (C) Copyright 2025-2026 Hewlett Packard Enterprise Development LP

package trust_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus"
	sdkv2morpheus "github.com/HPE/terraform-provider-hpe/morpheus/sdkv2"
	dstrust "github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/datasources/trust"
	"github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/resources/trust"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
)

func TestMain(m *testing.M) {
	code := m.Run()

	testhelpers.WriteMergedResults()

	os.Exit(code)
}

func TestAccMorpheusDataSourceKeyPairExampleOk(t *testing.T) {
	if capabilities.Missing(t, capabilities.All) {
		t.Log("Skipping test due to missing capabilities")

		return
	}
	t.Parallel()

	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Skip("Skipping due to bug in terraform code")

	providerConfig := testhelpers.ProviderBlock()

	name := acctest.RandomWithPrefix(t.Name())

	var dependenciesConfig string

	if currentDependency, err := trust.RenderKeyPairConfig(t, map[string]string{
		"Name": name,
	}); err != nil {
		t.Fatal(err)
	} else {
		dependenciesConfig += currentDependency
	}

	datasourceConfig, err := dstrust.RenderKeyPairConfig(t, map[string]string{
		"Name": "resource.hpe_morpheus_key_pair.example.name",
	})
	if err != nil {
		t.Fatal(err)
	}

	t.Log(providerConfig + dependenciesConfig + datasourceConfig)

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_key_pair.example",
			"name",
			name,
		),
		resource.TestCheckResourceAttrPair(
			"data.hpe_morpheus_key_pair.example",
			"id",
			"hpe_morpheus_key_pair.example",
			"id",
		),
		resource.TestCheckResourceAttrPair(
			"data.hpe_morpheus_key_pair.example",
			"publickey",
			"hpe_morpheus_key_pair.example",
			"public_key",
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), sdkv2morpheus.Provider()),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + dependenciesConfig + datasourceConfig,
				// We need to set this to true for test to pass,
				// so we'll skip until the id plan after apply bug is fixed
				// in key_pair resource.
				ExpectNonEmptyPlan: false,
				Check:              checkFn,
			},
		},
	})
}
