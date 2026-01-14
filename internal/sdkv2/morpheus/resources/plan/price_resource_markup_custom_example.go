// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package plan

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/testhelpers"
)

//go:generate sh -c "go run ../../../../../cmd/render -out examples/resources/morpheus_price/resource_markup_custom.tf price_resource_markup_custom.tf.tmpl Code 'terraform-test' Cost '38.00' Currency 'USD' CustomPrice '45.00' IncurCharges 'always' MarkupType 'custom' Name 'terraform-test' PriceType 'fixed' PriceUnit 'minute' TenantId '1'"

// RenderPriceMarkupCustomConfig generates a Terraform configuration for the tenant resource.
// It accepts optional overrides for field values. Default values are used if not overridden.
func RenderPriceMarkupCustomConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Code":         "terraform-test",
		"Cost":         "38.00",
		"Currency":     "USD",
		"CustomPrice":  "45.00",
		"IncurCharges": "always",
		"MarkupType":   "custom",
		"Name":         "terraform-test",
		"PriceType":    "fixed",
		"PriceUnit":    "minute",
		"TenantId":     "1",
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
	templatePath := filepath.Join(dir, "price_resource_markup_custom.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath,
		args...,
	)
}
