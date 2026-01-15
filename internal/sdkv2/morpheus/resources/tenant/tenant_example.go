// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package tenant

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/testhelpers"
)

//go:generate go run ../../../../../cmd/render -out examples/resources/morpheus_tenant/resource.tf tenant_resource.tf.tmpl RoleName "Tenant Admin" Name tftenant Description "Terraform example tenant" Enabled true Subdomain tfexample Currency USD BaseRoleId data.hpe_morpheus_role.example.id AccountNumber 12345 AccountName "tenant 12345" CustomerNumber 12345

// RenderResourceConfig generates a Terraform configuration for the tenant resource.
// It accepts optional overrides for field values. Default values are used if not overridden.
func RenderTenantConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Name":           "tftenant",
		"Description":    "Terraform example tenant",
		"Enabled":        "true",
		"Subdomain":      "tfexample",
		"Currency":       "USD",
		"BaseRoleId":     "data.hpe_morpheus_role.example.id",
		"AccountNumber":  "12345",
		"AccountName":    "tenant 12345",
		"CustomerNumber": "12345",
	}

	// Apply overrides to defaults
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
	templatePath := filepath.Join(dir, "tenant_resource.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath,
		args...,
	)
}
