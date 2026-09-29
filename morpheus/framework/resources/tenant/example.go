// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package tenant

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
)

//go:generate ../../../../bin/render -out examples/resources/morpheus_tenant/example.tf example.tf.tmpl ResourceLabel "example" RoleName "Tenant Admin" Name "Example Tenant" Description "Terraform example tenant" Enabled "true" Subdomain "tfexample" BaseRoleId "data.hpe_morpheus_role.example.id" Currency "USD" AccountNumber "12345" AccountName "tenant 12345" CustomerNumber "12345"
//go:generate ../../../../bin/render -out examples/resources/morpheus_tenant/example_subtenant.tf example_subtenant.tf.tmpl RoleName "Tenant Admin" ParentName "Parent Tenant" Name "Example Subtenant" Description "Terraform example subtenant" BaseRoleId "data.hpe_morpheus_role.example.id" Subdomain "subtenant" Currency "USD"

// RenderTenantConfig renders example.tf.tmpl. The defaults carry no RoleName,
// so no role data source is rendered and BaseRoleId is used as given; the
// published example (see go:generate above) sets RoleName to resolve the id
// from a seeded role by name, while the acceptance tests reference a role
// they create.
func RenderTenantConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"ResourceLabel":  "example",
		"Name":           "Example Tenant",
		"Description":    "Terraform example tenant",
		"Enabled":        "true",
		"Subdomain":      "tfexample",
		"BaseRoleId":     "3",
		"Currency":       "USD",
		"AccountNumber":  "12345",
		"AccountName":    "tenant 12345",
		"CustomerNumber": "12345",
	}

	for key, value := range overrides {
		defaults[key] = value
	}

	var args []string
	for key, value := range defaults {
		args = append(args, key, value)
	}

	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("unable to get current file path")
	}
	dir := filepath.Dir(filename)
	templatePath := filepath.Join(dir, "example.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath,
		args...,
	)
}
