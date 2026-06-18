// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package form_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus"
	sdkv2morpheus "github.com/HPE/terraform-provider-hpe/morpheus/sdkv2"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
)

// TestAccMorpheusFormRejectsSelfReferentialDependentField verifies the plan-time
// guard that rejects an option type whose dependent_field equals its own
// field_name (a field cannot depend on itself, which would create a circular
// dependsOnCode). This is a pure plan-time validation, so it needs no live
// appliance.
func TestAccMorpheusFormRejectsSelfReferentialDependentField(t *testing.T) {
	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()

	invalidConfig := `
resource "hpe_morpheus_form" "bad_self_ref" {
  name        = "bad-self-ref-test"
  code        = "bad-self-ref-test"
  description = "test"

  option_type {
    name            = "bad option"
    code            = "bad-option"
    type            = "text"
    field_name      = "myField"
    field_label     = "My Field"
    dependent_field = "myField"
  }
}
`
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), sdkv2morpheus.Provider()),
		Steps: []resource.TestStep{
			{
				Config:      providerConfig + invalidConfig,
				ExpectError: regexp.MustCompile(`dependent_field must not equal field_name`),
			},
		},
	})
}
