// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package identitysource

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
)

//go:generate ../../../../bin/render -out examples/resources/morpheus_identity_source_saml/resource.tf identity_source_saml_resource.tf.tmpl TenantId 1 Name "samldemo" Description "TF example SAML identity source" LoginRedirectUrl "https://tfexamplesaml.test.local:8443/realms/master/protocol/saml" LogoutRedirectUrl "https://tfexamplesaml.test.local:8443/realms/master/protocol/saml" IncludeSamlRequestParameter true SamlRequest "SelfSigned" ValidateAssertionSignature false GivenNameAttribute "givenName" SurnameAttribute "surname" EmailAttribute "email" DefaultAccountRoleId 4 RoleAttributeName "memberOf" RequiredRoleAttributeValue "test" EnableRoleMappingPermission false RoleMappingRoleId 4 RoleMappingRoleName "Demo" RoleMappingAssertionAttribute "developers" RoleMapping2RoleId 5 RoleMapping2RoleName "tf-example-user-role" RoleMapping2AssertionAttribute "operators"

// RenderIdentitySourceSAMLConfig generates a Terraform configuration for the
// SAML identity source resource using default values that can be overridden
// via the overrides map. The defaults name a placeholder identity provider:
// Morpheus validates a SAML identity source field by field (SamlUserService.
// validateUserSource) without contacting the IdP, so the example applies
// without one.
func RenderIdentitySourceSAMLConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"TenantId":                       "1",
		"Name":                           "Example",
		"Description":                    "TF example SAML identity source",
		"LoginRedirectUrl":               "https://tfexamplesaml.test.local:8443/realms/master/protocol/saml",
		"LogoutRedirectUrl":              "https://tfexamplesaml.test.local:8443/realms/master/protocol/saml",
		"IncludeSamlRequestParameter":    "true",
		"SamlRequest":                    "SelfSigned",
		"ValidateAssertionSignature":     "false",
		"GivenNameAttribute":             "givenName",
		"SurnameAttribute":               "surname",
		"EmailAttribute":                 "email",
		"DefaultAccountRoleId":           "4",
		"RoleAttributeName":              "memberOf",
		"RequiredRoleAttributeValue":     "test",
		"EnableRoleMappingPermission":    "false",
		"RoleMappingRoleId":              "4",
		"RoleMappingRoleName":            "Demo",
		"RoleMappingAssertionAttribute":  "developers",
		"RoleMapping2RoleId":             "5",
		"RoleMapping2RoleName":           "tf-example-user-role",
		"RoleMapping2AssertionAttribute": "operators",
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
	templatePath := filepath.Join(dir, "identity_source_saml_resource.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath,
		args...,
	)
}
