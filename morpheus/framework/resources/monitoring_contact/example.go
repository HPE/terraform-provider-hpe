// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package monitoring_contact

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
)

//go:generate ../../../../bin/render -out examples/resources/morpheus_monitoring_contact/example.tf example.tf.tmpl Name "Ops Team" EmailAddress "ops-team@example.com" SmsAddress "+15551234567"

func RenderMonitoringContactConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Name":         "Ops Team",
		"EmailAddress": "ops-team@example.com",
		"SmsAddress":   "+15551234567",
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
