// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package settingbackup

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
)

//go:generate sh -c "../../../../bin/render -out examples/resources/morpheus_setting_backup/resource.tf example.tf.tmpl BackupAppliance 'false' CreateBackups 'true' DefaultBackupScheduleId '3' DefaultBackupStorageBucketId '17' RetentionCount '21' ScheduledBackups 'true'"

// ExampleValues returns the values RenderSettingBackupConfig renders unless
// overridden, so a test can assert what it rendered.
func ExampleValues() map[string]string {
	return map[string]string{
		"BackupAppliance":              "false",
		"CreateBackups":                "true",
		"DefaultBackupScheduleId":      "3",
		"DefaultBackupStorageBucketId": "17",
		"RetentionCount":               "21",
		"ScheduledBackups":             "true",
	}
}

// RenderSettingBackupConfig renders the example configuration, with
// ExampleValues overridden by overrides.
func RenderSettingBackupConfig(t *testing.T, overrides map[string]string) (string, error) {
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
