// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package setting

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/testhelpers"
)

//go:generate sh -c "go run ../../../../../cmd/render -out examples/resources/morpheus_setting_appliance/resource.tf setting_appliance_resource.tf.tmpl ApiAllowedOrigins 'demo' ApplianceUrl 'https://morpheus.test.local' CurrencyProvider 'fixer' CurrencyProviderApiKey '5a4b220e-6f9f-43da-a572-390c8f6afed8' DefaultRoleId '5' DefaultUserRoleId '10' DockerPrivilegedMode 'false' InternalApplianceUrl 'https://pxemorpheus.test.local' ProxyDomain 'test.local' ProxyHost '10.0.0.100' ProxyPassword 'Password123456' ProxyPort '3128' ProxyUser 'jsmith' ProxyWorkstation 'work123' RegistrationEnabled 'true' SmtpFromAddress 'testemail@test.local' SmtpPassword 'Password12' SmtpPort '465' SmtpServer 'smtp01.test.local' SmtpUseSsl 'true' SmtpUseTls 'true' SmtpUsername 'jsmith'"

// RenderSettingApplianceConfig generates a Terraform configuration for the tenant resource.
// It accepts optional overrides for field values. Default values are used if not overridden.
func RenderSettingApplianceConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"ApiAllowedOrigins":      "demo",
		"ApplianceUrl":           "https://morpheus.test.local",
		"CurrencyProvider":       "fixer",
		"CurrencyProviderApiKey": "5a4b220e-6f9f-43da-a572-390c8f6afed8",
		"DefaultRoleId":          "1",
		"DefaultUserRoleId":      "1",
		"DockerPrivilegedMode":   "false",
		"InternalApplianceUrl":   "https://pxemorpheus.test.local",
		"ProxyDomain":            "test.local",
		"ProxyHost":              "10.0.0.100",
		"ProxyPassword":          "Password123456",
		"ProxyPort":              "3128",
		"ProxyUser":              "jsmith",
		"ProxyWorkstation":       "work123",
		"RegistrationEnabled":    "true",
		"SmtpFromAddress":        "testemail@test.local",
		"SmtpPassword":           "Password12",
		"SmtpPort":               "465",
		"SmtpServer":             "smtp01.test.local",
		"SmtpUseSsl":             "true",
		"SmtpUseTls":             "true",
		"SmtpUsername":           "jsmith",
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
	templatePath := filepath.Join(dir, "setting_appliance_resource.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath,
		args...,
	)
}
