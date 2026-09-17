// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package storage_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	sdkv2morpheus "github.com/HPE/terraform-provider-hpe/morpheus/sdkv2"
	dsstorage "github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/datasources/storage"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
	"github.com/HPE/terraform-provider-hpe/provider/adapter"
)

func TestAccMorpheusDataSourceStorageVolumeTypeExampleOk(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)

	t.Parallel()

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()

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
		// Exact-value check rather than TestCheckResourceAttrSet:
		// AttrSet treats an empty string as "not set", so it is unsafe
		// for description (system types may have an empty description).
		// "Kubernetes Volume" has a known, non-empty description.
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_storage_volume_type.example",
			"description",
			"Kubernetes Volume",
		),
		resource.TestCheckResourceAttrSet(
			"data.hpe_morpheus_storage_volume_type.example",
			"id",
		),
		resource.TestCheckResourceAttrSet(
			"data.hpe_morpheus_storage_volume_type.example",
			"code",
		),
		resource.TestCheckResourceAttrSet(
			"data.hpe_morpheus_storage_volume_type.example",
			"category",
		),
		resource.TestCheckResourceAttrSet(
			"data.hpe_morpheus_storage_volume_type.example",
			"enabled",
		),
		resource.TestCheckResourceAttrSet(
			"data.hpe_morpheus_storage_volume_type.example",
			"default_type",
		),
		resource.TestCheckResourceAttrSet(
			"data.hpe_morpheus_storage_volume_type.example",
			"has_datastore",
		),
		resource.TestCheckResourceAttrSet(
			"data.hpe_morpheus_storage_volume_type.example",
			"configurable_iops",
		),
		resource.TestCheckResourceAttrSet(
			"data.hpe_morpheus_storage_volume_type.example",
			"custom_size",
		),
		resource.TestCheckResourceAttrSet(
			"data.hpe_morpheus_storage_volume_type.example",
			"custom_label",
		),
		resource.TestCheckResourceAttrSet(
			"data.hpe_morpheus_storage_volume_type.example",
			"display_order",
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
