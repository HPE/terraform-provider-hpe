// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package plan

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/testhelpers"
)

//go:generate sh -c "go run ../../../../../cmd/render -out examples/resources/morpheus_price_set/resource.tf price_set_resource.tf.tmpl Name 'terraform-test' Code 'terraform-test' RegionCode 'us-west-2' PoolId '1' PriceUnit 'minute' Type 'fixed' PriceIds '[1]'"

// RenderPriceSetConfig generates a Terraform configuration for the price_set resource.
// It accepts optional overrides for field values. Default values are used if not overridden.
func RenderPriceSetConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Name":       "terraform-test",
		"Code":       "terraform-test",
		"RegionCode": "us-west-2",
		"PriceUnit":  "minute",
		"Type":       "fixed",
		"PriceIds":   "[1]",
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
	templatePath := filepath.Join(dir, "price_set_resource.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath,
		args...,
	)
}
