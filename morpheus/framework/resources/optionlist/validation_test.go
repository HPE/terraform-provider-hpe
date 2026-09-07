// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package optionlist_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
	"github.com/HPE/terraform-provider-hpe/provider/adapter"
)

// TestAccMorpheusOptionListResourceValidationRestMissingSourceUrl verifies that
// omitting source_url when type is unset (the server defaults it to "rest") is
// rejected at plan time (MORPH-14615).
func TestAccMorpheusOptionListResourceValidationRestMissingSourceUrl(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)

	providerConfig := testhelpers.ProviderBlock()

	// type omitted -> server defaults to "rest" -> source_url required.
	config := providerConfig + `
resource "hpe_morpheus_option_list" "test" {
	name = "TestAccMorpheusOptionListValidationRest"
}
`

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config:      config,
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`Missing required attribute`),
			},
		},
	})
}

// TestAccMorpheusOptionListResourceValidationApiMissingApiType verifies that
// omitting api_type when type is "api" is rejected at plan time (MORPH-14615).
func TestAccMorpheusOptionListResourceValidationApiMissingApiType(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)

	providerConfig := testhelpers.ProviderBlock()

	config := providerConfig + `
resource "hpe_morpheus_option_list" "test" {
	name       = "TestAccMorpheusOptionListValidationApi"
	type       = "api"
	source_url = "https://example.com/list"
}
`

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config:      config,
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`Missing required attribute`),
			},
		},
	})
}
