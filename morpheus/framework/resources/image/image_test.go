// (C) Copyright 2025-2026 Hewlett Packard Enterprise Development LP

package image_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus"
	sdkv2morpheus "github.com/HPE/terraform-provider-hpe/morpheus/sdkv2"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
)

func TestMain(m *testing.M) {
	code := m.Run()
	testhelpers.WriteMergedResults()
	os.Exit(code)
}

// Tests that our example file template used for docs is a valid config
func TestAccMorpheusImageExampleOk(t *testing.T) {
	if capabilities.Missing(t, capabilities.Alletra) {
		t.Log("Skipping test due to missing capabilities")

		return
	}
	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()

	name := acctest.RandomWithPrefix(t.Name())

	datasourceConfig := `
data "hpe_morpheus_os_type" "test" {
	name = "linux"
}

data "hpe_morpheus_storage_bucket" "test" {
	name = "Local Storage"
}
`

	resourceConfig, err := testhelpers.RenderExample(
		t, "example.tf.tmpl",
		"Name", name,
		"StorageProviderId", "data.hpe_morpheus_storage_bucket.test.id",
		"OsTypeId", "data.hpe_morpheus_os_type.test.id",
	)
	if err != nil {
		t.Fatal(err)
	}

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"hpe_morpheus_image.example_image",
			"name",
			name,
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_image.example_image",
			"image_type",
			"qcow2",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_image.example_image",
			"user_data",
			`#!/bin/sh
apk add --no-cache bash`,
		),
		resource.TestCheckNoResourceAttr(
			"hpe_morpheus_image.example_image",
			"ssh_password_wo",
		),
		resource.TestCheckResourceAttrPair(
			"hpe_morpheus_image.example_image",
			"os_type_id",
			"data.hpe_morpheus_os_type.test",
			"id",
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), sdkv2morpheus.Provider()),
		Steps: []resource.TestStep{
			{
				Config:             providerConfig + datasourceConfig + resourceConfig,
				ExpectNonEmptyPlan: false,
				Check:              checkFn,
				PlanOnly:           false,
			},
			{
				// Check that a post-apply plan detects no changes
				Config:             providerConfig + datasourceConfig + resourceConfig,
				ExpectNonEmptyPlan: false,
				Check:              checkFn,
				PlanOnly:           true,
			},
			{
				ImportState:       true,
				ImportStateVerify: true, // Check state post import
				ResourceName:      "hpe_morpheus_image.example_image",
				// ignore these fields as they are not available from the API
				ImportStateVerifyIgnore: []string{"url", "ssh_password_wo_version", "ssh_key_wo_version"},
				Check:                   checkFn,
			},
		},
	})
}
