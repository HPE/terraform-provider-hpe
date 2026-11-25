// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package integration_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/resources/integration"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
)

// gitRepositoryFromEnv is the repository the git-backed integration tests
// have Morpheus clone: Morpheus clones a repository when the integration is
// saved and rejects one it cannot reach, so the documented example's public
// repository cannot be relied on from every appliance. The tests take
//
//	TF_VAR_testacc_morpheus_git_repository_url     a repository Morpheus can clone without credentials
//	TF_VAR_testacc_morpheus_git_repository_branch  a branch of it
//
// and skip when either is unset. The values returned are what the tests
// expect to find in state; the configuration itself reads the same
// environment through the Terraform variables gitRepositoryVariables
// declares (gitRepositoryURLExpr, gitRepositoryBranchExpr), so neither value
// is ever spliced into HCL text.
func gitRepositoryFromEnv(t *testing.T) (url, branch string) {
	t.Helper()
	url = os.Getenv("TF_VAR_testacc_morpheus_git_repository_url")
	branch = os.Getenv("TF_VAR_testacc_morpheus_git_repository_branch")
	if url == "" || branch == "" {
		t.Skip("Skipping: set TF_VAR_testacc_morpheus_git_repository_url and " +
			"TF_VAR_testacc_morpheus_git_repository_branch to a repository Morpheus can clone to run this test")
	}

	return url, branch
}

// gitRepositoryVariables declares the Terraform variables the git-backed
// integration configurations take the repository from. Terraform reads their
// values from the TF_VAR_testacc_morpheus_git_repository_* environment itself
// — the convention the TF_VAR_ prefix signals (testhelpers/capabilities/
// README.md) — which is why the templates render url and default_branch as
// expressions and the tests pass these references rather than the values.
const (
	gitRepositoryVariables = `
variable "testacc_morpheus_git_repository_url" {}

variable "testacc_morpheus_git_repository_branch" {}
`
	gitRepositoryURLExpr    = "var.testacc_morpheus_git_repository_url"
	gitRepositoryBranchExpr = "var.testacc_morpheus_git_repository_branch"
)

func TestAccMorpheusIntegrationGitExampleOk(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.Git)

	t.Parallel()

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	repositoryURL, repositoryBranch := gitRepositoryFromEnv(t)

	providerConfig := testhelpers.ProviderBlock()

	name := acctest.RandomWithPrefix(t.Name())

	resourceConfig, err := integration.RenderIntegrationGitConfig(t, map[string]string{
		"Name":          name,
		"Url":           gitRepositoryURLExpr,
		"DefaultBranch": gitRepositoryBranchExpr,
	})
	if err != nil {
		t.Fatal(err)
	}
	config := providerConfig + gitRepositoryVariables + resourceConfig

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"hpe_morpheus_integration_git.tf_example_git_integration",
			"name",
			name,
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_integration_git.tf_example_git_integration",
			"enabled",
			"true",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_integration_git.tf_example_git_integration",
			"url",
			repositoryURL,
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_integration_git.tf_example_git_integration",
			"default_branch",
			repositoryBranch,
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_integration_git.tf_example_git_integration",
			"enable_git_caching",
			"true",
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Plan
			{
				Config:             config,
				ExpectNonEmptyPlan: true,
				Check:              checkFn,
				PlanOnly:           true,
			},
			// Apply
			{
				Config: config,
				Check:  checkFn,
			},
			// Plan after apply
			{
				Config:             config,
				ExpectNonEmptyPlan: false,
				Check:              checkFn,
				PlanOnly:           true,
			},
		},
	})
}
