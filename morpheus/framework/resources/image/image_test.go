// (C) Copyright 2025-2026 Hewlett Packard Enterprise Development LP

package image_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/systemoverride"
	"github.com/HPE/terraform-provider-hpe/provider"
)

func TestMain(m *testing.M) {
	systemoverride.ParseFlags()
	code := m.Run()
	testhelpers.WriteMergedResults()
	os.Exit(code)
}

func newProviderWithError() (tfprotov6.ProviderServer, error) {
	providerInstance := provider.New("test", morpheus.New())()

	return providerserver.NewProtocol6WithError(providerInstance)()
}

var testAccProtoV6ProviderFactories = map[string]func() (
	tfprotov6.ProviderServer, error,
){
	"hpe": newProviderWithError,
}

// Tests that our example file template used for docs is a valid config
func TestAccMorpheusImageExampleOk(t *testing.T) {
	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	testSystem := systemoverride.GetPreferred(t, "feature")
	providerConfig := testhelpers.ProviderBlock(testSystem)

	name := acctest.RandomWithPrefix(t.Name())

	resourceConfig, err := testhelpers.RenderExample(t, "example.tf.tmpl",
		"Name", name,
		"StorageProviderId", "196",
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
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:             providerConfig + resourceConfig,
				ExpectNonEmptyPlan: false,
				Check:              checkFn,
				PlanOnly:           false,
			},
			{
				// Check that a post-apply plan detects no changes
				Config:             providerConfig + resourceConfig,
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
