// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package usergroup_test

import (
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	sdkv2morpheus "github.com/HPE/terraform-provider-hpe/morpheus/sdkv2"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
	"github.com/HPE/terraform-provider-hpe/provider/adapter"
)

// TestAccMorpheusUserGroupValidationDescriptionTooLong verifies that a
// description exceeding 255 characters is rejected at plan time rather than
// triggering a server 500 (MORPH-9344).
func TestAccMorpheusUserGroupValidationDescriptionTooLong(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)

	providerConfig := testhelpers.ProviderBlock()

	longDescription := strings.Repeat("a", 256)

	config := providerConfig + `
resource "hpe_morpheus_user_group" "test" {
	name        = "TestAccMorpheusUserGroupValidationDescriptionTooLong"
	description = "` + longDescription + `"
}
`

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), sdkv2morpheus.Provider()),
		Steps: []resource.TestStep{
			{
				Config:      config,
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`to be in the range \(0 - 255\)`),
			},
		},
	})
}
