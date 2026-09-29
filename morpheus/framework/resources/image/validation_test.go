// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package image_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
	"github.com/HPE/terraform-provider-hpe/provider/adapter"
)

// TestAccMorpheusImageResourceValidationNegativeMinRam verifies that a negative
// min_ram is rejected at plan time. The API otherwise accepts it and renders the
// image's memory as "N/A" (MORPH-11809).
func TestAccMorpheusImageResourceValidationNegativeMinRam(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)

	providerConfig := testhelpers.ProviderBlock()

	config := providerConfig + `
resource "hpe_morpheus_image" "test" {
	name       = "TestAccMorpheusImageResourceValidationNegativeMinRam"
	image_type = "qcow2"
	min_ram    = -6
}
`

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config:      config,
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`must be at least 0`),
			},
		},
	})
}

// TestAccMorpheusImageResourceValidationNegativeMinDisk verifies that a negative
// min_disk is rejected at plan time, mirroring the min_ram guard (MORPH-11809).
func TestAccMorpheusImageResourceValidationNegativeMinDisk(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)

	providerConfig := testhelpers.ProviderBlock()

	config := providerConfig + `
resource "hpe_morpheus_image" "test" {
	name       = "TestAccMorpheusImageResourceValidationNegativeMinDisk"
	image_type = "qcow2"
	min_disk   = -6
}
`

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config:      config,
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`must be at least 0`),
			},
		},
	})
}
