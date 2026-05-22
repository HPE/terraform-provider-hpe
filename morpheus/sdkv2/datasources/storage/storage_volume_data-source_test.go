// (C) Copyright 2025-2026 Hewlett Packard Enterprise Development LP

package storage_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus"
	sdkv2morpheus "github.com/HPE/terraform-provider-hpe/morpheus/sdkv2"
	dsstorage "github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/datasources/storage"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
)

func TestAccMorpheusDataSourceStorageVolumeExampleOk(t *testing.T) {
	if capabilities.Missing(t, capabilities.Alletra) {
		t.Log("Skipping test due to missing capabilities")

		return
	}
	t.Parallel()

	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()

	var dependenciesConfig string

	// The data source only accepts IDs and computes the rest
	// We don't have a storage volume resource, so we'll use the first one created on the system
	datasourceConfig, err := dsstorage.RenderStorageVolumeConfig(t, map[string]string{
		"Id": "1",
	})
	if err != nil {
		t.Fatal(err)
	}

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_storage_volume.example",
			"id",
			"1",
		),
		resource.TestCheckResourceAttrSet(
			"data.hpe_morpheus_storage_volume.example",
			"name",
		),
		resource.TestCheckResourceAttrSet(
			"data.hpe_morpheus_storage_volume.example",
			"active",
		),
		resource.TestCheckResourceAttrSet(
			"data.hpe_morpheus_storage_volume.example",
			"category",
		),
		resource.TestCheckResourceAttrSet(
			"data.hpe_morpheus_storage_volume.example",
			"status",
		),
		resource.TestCheckResourceAttrSet(
			"data.hpe_morpheus_storage_volume.example",
			"type",
		),
		resource.TestCheckResourceAttrSet(
			"data.hpe_morpheus_storage_volume.example",
			"type_id",
		),
		resource.TestCheckResourceAttrSet(
			"data.hpe_morpheus_storage_volume.example",
			"uuid",
		),
		// not testing: cloud, datastore, source properties as they may not apply for
		// the type of storage volume with ID 1
		// i.e. we're missing infrastructure / provider features to test these.
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
