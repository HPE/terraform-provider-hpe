// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package settingappliance

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
)

//go:generate sh -c "../../../../bin/render -out examples/resources/morpheus_setting_appliance/resource.tf example.tf.tmpl ApiAllowedOrigins 'demo' ApplianceUrl 'https://morpheus.test.local' CurrencyProvider 'fixer' CurrencyProviderApiKey 'example' DefaultRoleId '5' DefaultUserRoleId '10' DockerPrivilegedMode 'false' InternalApplianceUrl 'https://pxemorpheus.test.local' ProxyDomain 'test.local' ProxyHost '10.0.0.100' ProxyPassword 'Password123456' ProxyPort '3128' ProxyUser 'jsmith' ProxyWorkstation 'work123' NoProxy '10.0.0.0/8,localhost,127.0.0.1' RegistrationEnabled 'true' SmtpFromAddress 'testemail@test.local' SmtpPassword 'Password12' SmtpPort '465' SmtpServer 'smtp01.test.local' SmtpUseSsl 'true' SmtpUseTls 'true' SmtpUsername 'jsmith'"

// ExampleValues returns the values RenderSettingApplianceConfig renders unless
// overridden, so a test can assert what it rendered.
func ExampleValues() map[string]string {
	//nolint:gosec // example values, not credentials
	return map[string]string{
		"ApiAllowedOrigins":      "demo",
		"ApplianceUrl":           "https://morpheus.test.local",
		"CurrencyProvider":       "fixer",
		"CurrencyProviderApiKey": "example",
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
		"NoProxy":                "10.0.0.0/8,localhost,127.0.0.1",
		"RegistrationEnabled":    "true",
		"SmtpFromAddress":        "testemail@test.local",
		"SmtpPassword":           "Password12",
		"SmtpPort":               "465",
		"SmtpServer":             "smtp01.test.local",
		"SmtpUseSsl":             "true",
		"SmtpUseTls":             "true",
		"SmtpUsername":           "jsmith",
	}
}

// RenderSettingApplianceConfig renders the example configuration, with
// ExampleValues overridden by overrides.
func RenderSettingApplianceConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	values := ExampleValues()
	for key, value := range overrides {
		values[key] = value
	}

	var args []string
	for key, value := range values {
		args = append(args, key, value)
	}

	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("unable to get current file path")
	}

	templatePath := filepath.Join(filepath.Dir(filename), "example.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath,
		args...,
	)
}
