// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package setting

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/testhelpers"
)

//go:generate sh -c " ../../../../../bin/render -out examples/resources/morpheus_setting_monitoring/resource.tf setting_monitoring_resource.tf.tmpl MorpheusAutoCreateChecks 'true' MorpheusAvailabilityPrecision '4' MorpheusAvailabilityTimeFrame '30' MorpheusDefaultCheckInterval '120' NewRelicLicenseKey 'ABC123' NewRelicMonitoringEnabled 'true' ServicenowCloseIncidentAction 'activity' ServicenowIntegrationId '1' ServicenowMonitoringEnabled 'true' ServicenowNewIncidentAction 'create' ServicenowSeverityCriticalImpact 'low' ServicenowSeverityInfoImpact 'high' ServicenowSeverityWarningImpact 'high'"

// RenderSettingMonitoringConfig generates a Terraform configuration for the tenant resource.
// It accepts optional overrides for field values. Default values are used if not overridden.
func RenderSettingMonitoringConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"MorpheusAutoCreateChecks":         "true",
		"MorpheusAvailabilityPrecision":    "4",
		"MorpheusAvailabilityTimeFrame":    "30",
		"MorpheusDefaultCheckInterval":     "120",
		"NewRelicLicenseKey":               "ABC123",
		"NewRelicMonitoringEnabled":        "true",
		"ServicenowCloseIncidentAction":    "activity",
		"ServicenowIntegrationId":          "1",
		"ServicenowMonitoringEnabled":      "true",
		"ServicenowNewIncidentAction":      "create",
		"ServicenowSeverityCriticalImpact": "low",
		"ServicenowSeverityInfoImpact":     "high",
		"ServicenowSeverityWarningImpact":  "high",
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
	templatePath := filepath.Join(dir, "setting_monitoring_resource.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath,
		args...,
	)
}
