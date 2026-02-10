// (C) Copyright 2025-2026 Hewlett Packard Enterprise Development LP

package catalog_test

import (
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus"
	sdkv2morpheus "github.com/HPE/terraform-provider-hpe/morpheus/sdkv2"
	dscatalog "github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/datasources/catalog"
	"github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/resources/catalogitem"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/systemoverride"
)

func TestMain(m *testing.M) {
	systemoverride.ParseFlags()
	code := m.Run()

	testhelpers.WriteMergedResults()

	os.Exit(code)
}

func TestAccMorpheusDataSourceCatalogItemTypeExampleOk(t *testing.T) {
	t.Parallel()

	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	testSystem := systemoverride.GetPreferred(t, "feature")
	providerConfig := testhelpers.ProviderBlock(testSystem)

	name := acctest.RandomWithPrefix(t.Name())

	var dependenciesConfig string

	if currentDependency, err := catalogitem.RenderCatalogItemInstanceConfig(t, map[string]string{
		"Name": name,
		"Code": strings.ToLower(name),
	}); err != nil {
		t.Fatal(err)
	} else {
		dependenciesConfig += currentDependency
	}

	datasourceConfig, err := dscatalog.RenderCatalogItemTypeConfig(t, map[string]string{
		"Name": "resource.hpe_morpheus_catalog_item_instance.tf_example_instance_catalog_item.name",
	})
	if err != nil {
		t.Fatal(err)
	}

	t.Log(providerConfig + dependenciesConfig + datasourceConfig)

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_catalog_item_type.example",
			"name",
			name,
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), sdkv2morpheus.Provider()),
		Steps: []resource.TestStep{
			{
				Config:             providerConfig + dependenciesConfig + datasourceConfig,
				ExpectNonEmptyPlan: false,
				Check:              checkFn,
			},
		},
	})
}
