// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package whitelabel_settings

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
)

//go:generate ../../../../bin/render -out examples/resources/morpheus_whitelabel_settings/example.tf example.tf.tmpl Enabled "true" ApplianceName "Acme Cloud Platform" PrimaryColor "#1a73e8" SecondaryColor "#ffffff"

func RenderWhitelabelSettingsConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Enabled":        "true",
		"ApplianceName":  "Acme Cloud Platform",
		"PrimaryColor":   "#1a73e8",
		"SecondaryColor": "#ffffff",
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
