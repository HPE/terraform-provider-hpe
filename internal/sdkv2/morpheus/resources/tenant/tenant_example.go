// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package tenant

import (
	"testing"

	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/testhelpers"
)

//go:generate go run ../../../../../cmd/render -out examples/resources/morpheus_tenant/resource.tf tenant_resource.tf.tmpl RoleName "Tenant Admin" Name tftenant Description "Terraform example tenant" Enabled true Subdomain tfexample Currency USD AccountNumber 12345 AccountName "tenant 12345" CustomerNumber 12345

// RenderResourceConfig generates a Terraform configuration for the tenant resource.
// It accepts optional overrides for field values. Default values are used if not overridden.
func RenderTenantConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"RoleName":       "Tenant Admin",
		"Name":           "tftenant",
		"Description":    "Terraform example tenant",
		"Enabled":        "true",
		"Subdomain":      "tfexample",
		"Currency":       "USD",
		"AccountNumber":  "12345",
		"AccountName":    "tenant 12345",
		"CustomerNumber": "12345",
	}

	// Apply overrides to defaults
	for key, value := range overrides {
		defaults[key] = value
	}

	return testhelpers.RenderExample(
		t,
		"tenant_resource.tf.tmpl",
		"RoleName",
		defaults["RoleName"],
		"Name",
		defaults["Name"],
		"Description",
		defaults["Description"],
		"Enabled",
		defaults["Enabled"],
		"Subdomain",
		defaults["Subdomain"],
		"Currency",
		defaults["Currency"],
		"AccountNumber",
		defaults["AccountNumber"],
		"AccountName",
		defaults["AccountName"],
		"CustomerNumber",
		defaults["CustomerNumber"],
	)
}
