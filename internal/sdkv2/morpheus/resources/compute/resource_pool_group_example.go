// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package compute

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/testhelpers"
)

//go:generate go run ../../../../../cmd/render -out examples/resources/morpheus_resource_pool_group/resource.tf resource_pool_group_resource.tf.tmpl Name "TFExample Resource Pool Group" Description "TFExample Resource Pool Group" Mode roundrobin ResourcePoolIds "[1, 2, 3]" AllGroupAccess true GroupAccessGroupId 2 GroupAccessDefault true Visibility public TenantIds "[1, 2]"

// RenderResourcePoolGroupConfig generates a Terraform configuration for the resource pool group resource.
// It accepts optional overrides for field values. Default values are used if not overridden.
func RenderResourcePoolGroupConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Name":               "TFExample Resource Pool Group",
		"Description":        "TFExample Resource Pool Group",
		"Mode":               "roundrobin",
		"ResourcePoolIds":    "[1, 2, 3]",
		"AllGroupAccess":     "true",
		"GroupAccessGroupId": "2",
		"GroupAccessDefault": "true",
		"Visibility":         "public",
		"TenantIds":          "[1, 2]",
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
	templatePath := filepath.Join(dir, "resource_pool_group_resource.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath,
		args...,
	)
}
