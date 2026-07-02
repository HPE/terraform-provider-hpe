// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package template_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	sdkv2morpheus "github.com/HPE/terraform-provider-hpe/morpheus/sdkv2"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/provider/adapter"
)

// TestAccMorpheusSpecTemplateTerraformRejectsMissingRepositoryId is a regression
// test for MORPH-13329: source_type = "repository" without repository_id must
// fail at plan time (the Morpheus UI requires the repository, but the API
// silently accepts the omission). Runs as a unit test - the CustomizeDiff error
// is raised during plan before any API call.
func TestAccMorpheusSpecTemplateTerraformRejectsMissingRepositoryId(t *testing.T) {
	t.Parallel()

	providerConfig := testhelpers.ProviderBlockUnitTest()

	invalidConfig := `
resource "hpe_morpheus_spec_template_terraform" "test" {
  name        = "qatf-terraform-spec-no-repo-id"
  source_type = "repository"
  version_ref = "main"
  spec_path   = "test.tf"
}
`
	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), sdkv2morpheus.Provider()),
		Steps: []resource.TestStep{
			{
				Config:      providerConfig + invalidConfig,
				ExpectError: regexp.MustCompile(`repository_id is required when source_type is "repository"`),
			},
		},
	})
}

// TestAccMorpheusSpecTemplateTerraformRejectsEmptySpecContent is a regression
// test for MORPH-13325: source_type = "local" with an empty spec_content must
// fail at plan time. Runs as a unit test - the CustomizeDiff error is raised
// during plan before any API call.
func TestAccMorpheusSpecTemplateTerraformRejectsEmptySpecContent(t *testing.T) {
	t.Parallel()

	providerConfig := testhelpers.ProviderBlockUnitTest()

	invalidConfig := `
resource "hpe_morpheus_spec_template_terraform" "test" {
  name         = "qatf-terraform-spec-empty-content"
  source_type  = "local"
  spec_content = ""
}
`
	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), sdkv2morpheus.Provider()),
		Steps: []resource.TestStep{
			{
				Config:      providerConfig + invalidConfig,
				ExpectError: regexp.MustCompile(`spec_content is required and must not be empty when source_type is "local"`),
			},
		},
	})
}
