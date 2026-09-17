// (C) Copyright 2025-2026 Hewlett Packard Enterprise Development LP

package image_test

import (
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	sdkv2morpheus "github.com/HPE/terraform-provider-hpe/morpheus/sdkv2"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
	"github.com/HPE/terraform-provider-hpe/provider/adapter"
)

func TestMain(m *testing.M) {
	code := testhelpers.TestMain(m)
	testhelpers.WriteMergedResults()
	os.Exit(code)
}

// Tests that our example file template used for docs is a valid config
func TestAccMorpheusImageResourceExampleOk(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()

	name := acctest.RandomWithPrefix(t.Name())
	visibility := testhelpers.TenantVisibility(t)

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
		"Visibility", visibility,
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
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), sdkv2morpheus.Provider()),
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

// TestAccMorpheusImageVisibilityPublicRequiresMasterTenant_MORPH16419 verifies
// that a sub-tenant caller setting visibility = "public" is rejected at plan
// time with a clear message. It skips on the master tenant.
func TestAccMorpheusImageVisibilityPublicRequiresMasterTenant_MORPH16419(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	if testhelpers.IsMasterTenant(t) {
		t.Skip("visibility = \"public\" is allowed for the master tenant; " +
			"this negative test only applies to sub-tenant callers")
	}

	providerConfig := testhelpers.ProviderBlock()
	name := acctest.RandomWithPrefix(t.Name())

	// Literal placeholder ids: the step is PlanOnly and the master-tenant
	// guard rejects the plan before anything is created or resolved, so the
	// ids never need to exist. Data-source lookups (os type, storage bucket)
	// must be avoided here -- they are read at plan time and a sub-tenant
	// cannot see the master's "Local Storage" bucket, which would fail the
	// plan with an unrelated error before the guard's message can match.
	resourceConfig, err := testhelpers.RenderExample(
		t, "example.tf.tmpl",
		"Name", name,
		"StorageProviderId", "1",
		"OsTypeId", "1",
		"Visibility", "public",
	)
	if err != nil {
		t.Fatal(err)
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), sdkv2morpheus.Provider()),
		Steps: []resource.TestStep{
			{
				Config:      providerConfig + resourceConfig,
				PlanOnly:    true,
				ExpectError: regexp.MustCompile("master tenant"),
			},
		},
	})
}
