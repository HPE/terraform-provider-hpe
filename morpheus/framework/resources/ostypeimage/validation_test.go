// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package ostypeimage_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
	"github.com/HPE/terraform-provider-hpe/provider/adapter"
)

// TestAccMorpheusOsTypeImageResourceValidationZeroOsTypeId verifies that a
// zero os_type_id is rejected at plan time rather than being sent to the API
// (MORPH-16357).
func TestAccMorpheusOsTypeImageResourceValidationZeroOsTypeId(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)

	providerConfig := testhelpers.ProviderBlock()

	config := providerConfig + `
resource "hpe_morpheus_os_type_image" "test" {
	os_type_id       = 0
	virtual_image_id = 1
}
`

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config:      config,
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`must be at least 1`),
			},
		},
	})
}

// TestAccMorpheusOsTypeImageResourceValidationZeroVirtualImageId verifies that
// a zero virtual_image_id is rejected at plan time (MORPH-16357).
func TestAccMorpheusOsTypeImageResourceValidationZeroVirtualImageId(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)

	providerConfig := testhelpers.ProviderBlock()

	config := providerConfig + `
resource "hpe_morpheus_os_type_image" "test" {
	os_type_id       = 1
	virtual_image_id = 0
}
`

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config:      config,
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`must be at least 1`),
			},
		},
	})
}
