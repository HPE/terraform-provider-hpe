// (C) Copyright 2025-2026 Hewlett Packard Enterprise Development LP

package storage_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	sdkv2morpheus "github.com/HPE/terraform-provider-hpe/morpheus/sdkv2"
	dsstorage "github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/datasources/storage"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
	"github.com/HPE/terraform-provider-hpe/utils/adapter"
)

func TestAccMorpheusDataSourceStorageBucketExampleOk(t *testing.T) {
	defer testhelpers.RecordResult(t)

	if capabilities.Missing(t, capabilities.All) {
		t.Log("Skipping test due to missing capabilities")
	}
	t.Parallel()

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()

	var datasourceConfig string

	// There should be a storage bucket with ID 1 for local storage.
	// We can use this to test the data source without creating a storage bucket.
	// Note that we lack a storage bucket resource.
	if current, err := dsstorage.RenderStorageBucketIdConfig(t, map[string]string{
		"Id": "1",
	}); err != nil {
		t.Fatal(err)
	} else {
		datasourceConfig += current
	}

	// now test the example config used for docs
	if current, err := dsstorage.RenderStorageBucketConfig(t, map[string]string{
		"Name": "data.hpe_morpheus_storage_bucket.example_id.name",
	}); err != nil {
		t.Fatal(err)
	} else {
		datasourceConfig += current
	}

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttrSet(
			"data.hpe_morpheus_storage_bucket.example_id",
			"name",
		),
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_storage_bucket.example_id",
			"id",
			"1",
		),
		// don't care what the name is for the example with name config,
		// just that it has been set correctly by provider
		resource.TestCheckResourceAttrPair(
			"data.hpe_morpheus_storage_bucket.example_id",
			"name",
			"data.hpe_morpheus_storage_bucket.example",
			"name",
		),
		// should have same id as the example_id config
		resource.TestCheckResourceAttrPair(
			"data.hpe_morpheus_storage_bucket.example_id",
			"id",
			"data.hpe_morpheus_storage_bucket.example",
			"id",
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewAdaptedMorpheus(), sdkv2morpheus.Provider()),
		Steps: []resource.TestStep{
			{
				Config:             providerConfig + datasourceConfig,
				ExpectNonEmptyPlan: false,
				Check:              checkFn,
			},
		},
	})
}
