// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package storage_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus"
	sdkv2morpheus "github.com/HPE/terraform-provider-hpe/morpheus/sdkv2"
	dsstorage "github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/datasources/storage"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/systemoverride"
)

func TestAccMorpheusDataSourceStorageVolumeTypeExampleOk(t *testing.T) {
	t.Parallel()

	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	testSystem := systemoverride.GetPreferred(t, "feature")
	providerConfig := testhelpers.ProviderBlock(testSystem)

	// use a system storage volume type
	datasourceConfig, err := dsstorage.RenderStorageVolumeTypeConfig(t, map[string]string{
		"Name": "\"Kubernetes Volume\"",
	})
	if err != nil {
		t.Fatal(err)
	}

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_storage_volume_type.example",
			"name",
			"Kubernetes Volume",
		),
		resource.TestCheckResourceAttrSet(
			"data.hpe_morpheus_storage_volume_type.example",
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
