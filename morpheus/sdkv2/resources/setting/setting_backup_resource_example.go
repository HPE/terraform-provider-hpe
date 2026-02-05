// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package setting

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
)

//go:generate sh -c "../../../../bin/render -out examples/resources/morpheus_setting_backup/resource.tf setting_backup_resource.tf.tmpl BackupAppliance 'false' CreateBackups 'true' DefaultBackupScheduleId '3' DefaultBackupStorageBucketId '17' RetentionDays '21' ScheduledBackups 'true'"

// RenderSettingBackupConfig generates a Terraform configuration for the tenant resource.
// It accepts optional overrides for field values. Default values are used if not overridden.
func RenderSettingBackupConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"BackupAppliance":              "false",
		"CreateBackups":                "true",
		"DefaultBackupScheduleId":      "3",
		"DefaultBackupStorageBucketId": "17",
		"RetentionDays":                "21",
		"ScheduledBackups":             "true",
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
	templatePath := filepath.Join(dir, "setting_backup_resource.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath,
		args...,
	)
}
