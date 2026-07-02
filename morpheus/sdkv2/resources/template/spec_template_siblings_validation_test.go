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

// The arm, cloud_formation, helm and kubernetes spec templates share the same
// source schema as the terraform spec template and the same
// validateSpecTemplateSource CustomizeDiff (MORPH-13329 / MORPH-13325). These
// unit tests confirm the validator is wired into each sibling; the validator's
// full behaviour is covered by the terraform tests. The two validation paths
// (missing repository_id / empty spec_content) are alternated across the
// siblings. Errors are raised at plan time, so no live environment is needed.

func expectSpecTemplateValidationError(t *testing.T, config string, errRe string) {
	t.Helper()

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), sdkv2morpheus.Provider()),
		Steps: []resource.TestStep{
			{
				Config:      testhelpers.ProviderBlockUnitTest() + config,
				ExpectError: regexp.MustCompile(errRe),
			},
		},
	})
}

func TestAccMorpheusSpecTemplateArmRejectsMissingRepositoryId(t *testing.T) {
	t.Parallel()
	expectSpecTemplateValidationError(t, `
resource "hpe_morpheus_spec_template_arm" "test" {
  name        = "qatf-arm-spec-no-repo-id"
  source_type = "repository"
  version_ref = "main"
  spec_path   = "test.json"
}
`, `repository_id is required when source_type is "repository"`)
}

func TestAccMorpheusSpecTemplateCloudFormationRejectsEmptySpecContent(t *testing.T) {
	t.Parallel()
	expectSpecTemplateValidationError(t, `
resource "hpe_morpheus_spec_template_cloud_formation" "test" {
  name         = "qatf-cf-spec-empty-content"
  source_type  = "local"
  spec_content = ""
}
`, `spec_content is required and must not be empty when source_type is "local"`)
}

func TestAccMorpheusSpecTemplateHelmRejectsMissingRepositoryId(t *testing.T) {
	t.Parallel()
	expectSpecTemplateValidationError(t, `
resource "hpe_morpheus_spec_template_helm" "test" {
  name        = "qatf-helm-spec-no-repo-id"
  source_type = "repository"
  version_ref = "main"
  spec_path   = "chart"
}
`, `repository_id is required when source_type is "repository"`)
}

func TestAccMorpheusSpecTemplateKubernetesRejectsEmptySpecContent(t *testing.T) {
	t.Parallel()
	expectSpecTemplateValidationError(t, `
resource "hpe_morpheus_spec_template_kubernetes" "test" {
  name         = "qatf-k8s-spec-empty-content"
  source_type  = "local"
  spec_content = ""
}
`, `spec_content is required and must not be empty when source_type is "local"`)
}
