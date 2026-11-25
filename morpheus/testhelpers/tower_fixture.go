// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package testhelpers

import "testing"

// The Ansible Tower fixture: the Tower (or AWX, or a stand-in that serves its
// API) the `ansible_tower` tests talk to, and the integration Morpheus already
// has for it, named by environment variables and resolved by Terraform data
// sources inside each test configuration. Nothing here is an id: ids are
// allocated per appliance, a name resolves wherever the object exists.
//
// Morpheus verifies a Tower integration when it is created -- a TCP check,
// GET /api/v2/ping/ and GET /api/v2/me/ -- and refuses the save otherwise, so
// the resource test needs a reachable Tower, not just a URL; the data-source
// and task tests need an integration that has synced an inventory and a job
// template. There are no defaults: a guessed URL turns a test that should skip
// into one that fails against the API.
const (
	// EnvAnsibleTowerIntegrationName names the Ansible Tower integration
	// Morpheus already has, for the inventory, job template and task tests.
	EnvAnsibleTowerIntegrationName = "TF_VAR_testacc_morpheus_ansible_tower_integration_name"
	// EnvAnsibleTowerInventoryName names an inventory that integration has
	// synced from the Tower, as GET /api/options/ansibleTowerInventory lists it.
	EnvAnsibleTowerInventoryName = "TF_VAR_testacc_morpheus_ansible_tower_inventory_name"
	// EnvAnsibleTowerJobTemplateName names a job template it has synced.
	EnvAnsibleTowerJobTemplateName = "TF_VAR_testacc_morpheus_ansible_tower_job_template_name"
	// EnvAnsibleTowerURL, EnvAnsibleTowerUsername and EnvAnsibleTowerPassword
	// are the Tower the integration resource test creates an integration
	// for: its API base URL and an account Morpheus can log in with.
	EnvAnsibleTowerURL      = "TF_VAR_testacc_morpheus_ansible_tower_url"
	EnvAnsibleTowerUsername = "TF_VAR_testacc_morpheus_ansible_tower_username"
	EnvAnsibleTowerPassword = "TF_VAR_testacc_morpheus_ansible_tower_password" //nolint:gosec // G101: a variable name
)

// AnsibleTowerIntegrationName returns the name of the appliance's Ansible
// Tower integration, skipping the test when EnvAnsibleTowerIntegrationName is
// unset.
func AnsibleTowerIntegrationName(t *testing.T) string {
	t.Helper()

	return envOrSkip(t, EnvAnsibleTowerIntegrationName, "an Ansible Tower integration on the appliance")
}

// AnsibleTowerInventoryName returns the name of an inventory the integration
// has synced, skipping the test when EnvAnsibleTowerInventoryName is unset.
func AnsibleTowerInventoryName(t *testing.T) string {
	t.Helper()

	return envOrSkip(t, EnvAnsibleTowerInventoryName, "an inventory synced from the Tower")
}

// AnsibleTowerJobTemplateName returns the name of a job template the
// integration has synced, skipping the test when EnvAnsibleTowerJobTemplateName
// is unset.
func AnsibleTowerJobTemplateName(t *testing.T) string {
	t.Helper()

	return envOrSkip(t, EnvAnsibleTowerJobTemplateName, "a job template synced from the Tower")
}

// AnsibleTowerCredentials returns the URL, username and password of a Tower
// that Morpheus can verify an integration against, skipping the test when any
// of EnvAnsibleTowerURL, EnvAnsibleTowerUsername or EnvAnsibleTowerPassword is
// unset.
func AnsibleTowerCredentials(t *testing.T) (url, username, password string) {
	t.Helper()

	url = envOrSkip(t, EnvAnsibleTowerURL, "a reachable Ansible Tower")
	username = envOrSkip(t, EnvAnsibleTowerUsername, "an Ansible Tower account")
	password = envOrSkip(t, EnvAnsibleTowerPassword, "an Ansible Tower account")

	return url, username, password
}

// AnsibleTowerIntegrationBlock renders a data source that resolves the
// appliance's Ansible Tower integration by name, as
// data.hpe_morpheus_integration.ansible_tower; AnsibleTowerIntegrationIDRef
// is the HCL expression for its id. The name is quoted for HCL (HCLQuote), so
// one with a quote, a backslash or a `${` in it is looked up as written.
func AnsibleTowerIntegrationBlock(name string) string {
	return "data \"hpe_morpheus_integration\" \"ansible_tower\" {\n  name = " + HCLQuote(name) + "\n}\n"
}

// AnsibleTowerIntegrationIDRef is the HCL expression for the id of the
// integration AnsibleTowerIntegrationBlock resolves.
const AnsibleTowerIntegrationIDRef = "data.hpe_morpheus_integration.ansible_tower.id"
