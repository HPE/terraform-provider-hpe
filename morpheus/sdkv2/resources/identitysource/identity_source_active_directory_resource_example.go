// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package identitysource

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
)

//go:generate ../../../../bin/render -out examples/resources/morpheus_identity_source_active_directory/resource.tf identity_source_active_directory_resource.tf.tmpl TenantId 1 Name "addemo" Description "TF example AD identity source" AdServer "\"dc01.contoso.com\"" Domain "\"contoso.com\"" UseSsl false BindingUsername "\"administrator\"" BindingPassword "\"Password123\"" RequiredGroup "administrators" SearchMemberGroups true DefaultAccountRoleId 7 RoleMappingRoleId 2 RoleMappingRoleName "developers" RoleMappingGroupName "\"developers\"" RoleMappingGroupFqn "\"CN=developers,CN=Users,DC=contoso,DC=com\""

// RenderIdentitySourceActiveDirectoryConfig generates a Terraform
// configuration for the Active Directory identity source resource using
// default values that can be overridden via the overrides map. The defaults
// name a placeholder directory, which is fine for the rendered docs example
// but will not apply: when the identity source is saved, Morpheus binds to
// the server as the binding user, looks the required group and every mapped
// group up in the directory, and rejects the save if any of that fails. So
// the acceptance test overrides AdServer, Domain, BindingUsername,
// BindingPassword and RoleMappingGroupName with a reachable directory's
// details. Morpheus also rewrites a mapping's group name and FQN to the
// group's sAMAccountName and distinguished name as the directory reports
// them; an empty RoleMappingGroupFqn leaves the (Computed) FQN out of the
// configuration so that Morpheus fills it in.
//
// The values the test takes from the environment — AdServer, Domain,
// BindingUsername, BindingPassword, RoleMappingGroupName and
// RoleMappingGroupFqn — are rendered as HCL expressions, like the two role
// ids: their defaults carry their own quotes, and the test passes a Terraform
// variable reference (var.testacc_active_directory_server, …) instead of the
// value, so a DOMAIN\user binding name, a password with a quote, backslash or
// ${, or a distinguished name with an escaped comma reaches the provider
// intact without ever being spliced into HCL text.
func RenderIdentitySourceActiveDirectoryConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"TenantId":             "1",
		"Name":                 "Example",
		"Description":          "TF example AD identity source",
		"AdServer":             `"dc01.contoso.com"`,
		"Domain":               `"contoso.com"`,
		"UseSsl":               "false",
		"BindingUsername":      `"administrator"`,
		"BindingPassword":      `"Password123"`,
		"RequiredGroup":        "administrators",
		"SearchMemberGroups":   "true",
		"DefaultAccountRoleId": "7",
		"RoleMappingRoleId":    "2",
		"RoleMappingRoleName":  "developers",
		"RoleMappingGroupName": `"developers"`,
		"RoleMappingGroupFqn":  `"CN=developers,CN=Users,DC=contoso,DC=com"`,
	}

	for key, value := range overrides {
		defaults[key] = value
	}

	var args []string
	for key, value := range defaults {
		args = append(args, key, value)
	}

	// Get the directory where this source file is located
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("unable to get current file path")
	}
	dir := filepath.Dir(filename)
	templatePath := filepath.Join(dir, "identity_source_active_directory_resource.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath,
		args...,
	)
}
