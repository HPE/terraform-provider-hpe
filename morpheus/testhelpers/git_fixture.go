// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package testhelpers

import "testing"

// The Git fixture: a Git integration the appliance has, one of the code
// repositories Morpheus recorded from it, and the branch that repository
// serves -- named by environment variables and resolved by the
// hpe_morpheus_integration_git data source inside each test configuration,
// which carries the integration's id and its repository_ids map. Nothing here
// is an id: ids are allocated per appliance, a name resolves wherever the
// object exists.
//
// A blueprint sourced from a repository is verified when it is saved: Morpheus
// clones the repository at the branch and parses the files under the working
// path (TerraformProvisionService, ArmProvisionService), so the repository has
// to exist and the branch has to be one it has. There are no defaults: a
// guessed name turns a test that should skip into one that fails against the
// API.
const (
	// EnvGitIntegrationName names the appliance's Git integration.
	EnvGitIntegrationName = "TF_VAR_testacc_morpheus_git_integration_name"
	// EnvGitRepositoryName names a code repository Morpheus recorded from that
	// integration, as it appears in the integration's repository_ids map --
	// the repository's own name, without the owner.
	EnvGitRepositoryName = "TF_VAR_testacc_morpheus_git_repository_name"
	// EnvGitRepositoryBranch names a branch of that repository, the one a
	// blueprint's version_ref points at. Morpheus falls back to the
	// repository's default branch when none is sent, so this is normally that.
	EnvGitRepositoryBranch = "TF_VAR_testacc_morpheus_git_repository_branch"
)

// GitIntegrationName returns the name of the appliance's Git integration,
// skipping the test when EnvGitIntegrationName is unset.
func GitIntegrationName(t *testing.T) string {
	t.Helper()

	return envOrSkip(t, EnvGitIntegrationName, "a Git integration on the appliance")
}

// GitRepositoryName returns the name of a code repository of that
// integration, skipping the test when EnvGitRepositoryName is unset.
func GitRepositoryName(t *testing.T) string {
	t.Helper()

	return envOrSkip(t, EnvGitRepositoryName, "a code repository of the Git integration")
}

// GitRepositoryBranch returns the branch of that repository a blueprint's
// version_ref should name, skipping the test when EnvGitRepositoryBranch is
// unset.
func GitRepositoryBranch(t *testing.T) string {
	t.Helper()

	return envOrSkip(t, EnvGitRepositoryBranch, "the branch the code repository serves")
}

// GitIntegrationBlock renders a data source that resolves the appliance's Git
// integration by name, as data.hpe_morpheus_integration_git.git; its id is
// GitIntegrationIDRef and a repository's id is GitRepositoryIDRef. The name is
// quoted for HCL (HCLQuote), so one with a quote, a backslash or a `${` in it
// is looked up as written.
func GitIntegrationBlock(name string) string {
	return "data \"hpe_morpheus_integration_git\" \"git\" {\n  name = " + HCLQuote(name) + "\n}\n"
}

// GitIntegrationIDRef is the HCL expression for the id of the integration
// GitIntegrationBlock resolves.
const GitIntegrationIDRef = "data.hpe_morpheus_integration_git.git.id"

// GitRepositoryIDRef is the HCL expression for the id Morpheus recorded for
// the named repository of that integration; the name is quoted for HCL.
func GitRepositoryIDRef(repository string) string {
	return "data.hpe_morpheus_integration_git.git.repository_ids[" + HCLQuote(repository) + "]"
}
