// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package setting

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/testhelpers"
)

//go:generate go run ../../../../../cmd/render -out examples/resources/morpheus_setting_provisioning/resource.tf morpheus_setting_provisioning_resource.tf.tmpl AllowZoneSelection false AllowHostSelection false RequireEnvironments false ShowPricing true HideDatastoreStats true CrossTenantNamingPolicies false CloudinitUsername cloudinit CloudinitPassword Pa55w0rd! WindowsPassword Pa55w0rd! PxeRootPassword Pa55w0rd!

func RenderSettingProvisioningConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"AllowZoneSelection":        "false",
		"AllowHostSelection":        "false",
		"RequireEnvironments":       "false",
		"ShowPricing":               "true",
		"HideDatastoreStats":        "true",
		"CrossTenantNamingPolicies": "false",
		"CloudinitUsername":         "cloudinit",
		"CloudinitPassword":         "Pa55w0rd!",
		"WindowsPassword":           "Pa55w0rd!",
		"PxeRootPassword":           "Pa55w0rd!",
	}

	for key, value := range overrides {
		defaults[key] = value
	}

	// Get the directory where this source file is located
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("unable to get current file path")
	}
	dir := filepath.Dir(filename)
	templatePath := filepath.Join(dir, "morpheus_setting_provisioning_resource.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath,
		"AllowZoneSelection", defaults["AllowZoneSelection"],
		"AllowHostSelection", defaults["AllowHostSelection"],
		"RequireEnvironments", defaults["RequireEnvironments"],
		"ShowPricing", defaults["ShowPricing"],
		"HideDatastoreStats", defaults["HideDatastoreStats"],
		"CrossTenantNamingPolicies", defaults["CrossTenantNamingPolicies"],
		"CloudinitUsername", defaults["CloudinitUsername"],
		"CloudinitPassword", defaults["CloudinitPassword"],
		"WindowsPassword", defaults["WindowsPassword"],
		"PxeRootPassword", defaults["PxeRootPassword"],
	)
}

