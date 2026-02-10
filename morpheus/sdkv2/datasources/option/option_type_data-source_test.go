// (C) Copyright 2025-2026 Hewlett Packard Enterprise Development LP

package option_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus"
	sdkv2morpheus "github.com/HPE/terraform-provider-hpe/morpheus/sdkv2"
	dsoption "github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/datasources/option"
	"github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/resources/optiontype"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/systemoverride"
)

func TestAccMorpheusDataSourceOptionTypeExampleOk(t *testing.T) {
	t.Parallel()

	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	testSystem := systemoverride.GetPreferred(t, "feature")
	providerConfig := testhelpers.ProviderBlock(testSystem)

	name := acctest.RandomWithPrefix(t.Name())

	var dependenciesConfig string

	if currentDependency, err := optiontype.RenderOptionTypeNumberConfig(t, map[string]string{
		"Name": name,
	}); err != nil {
		t.Fatal(err)
	} else {
		dependenciesConfig += currentDependency
	}

	datasourceConfig, err := dsoption.RenderOptionTypeConfig(t, map[string]string{
		"Name": "resource.hpe_morpheus_option_type_number.example.name",
	})
	if err != nil {
		t.Fatal(err)
	}

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_option_type.example",
			"name",
			name,
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), sdkv2morpheus.Provider()),
		Steps: []resource.TestStep{
			{
				Config:             providerConfig + "\n" + dependenciesConfig + "\n" + datasourceConfig,
				ExpectNonEmptyPlan: false,
				Check:              checkFn,
			},
		},
	})
}
