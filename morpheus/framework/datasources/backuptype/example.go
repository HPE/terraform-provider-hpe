// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package backuptype

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
)

//go:generate ../../../../bin/render -out examples/data-sources/morpheus_backup_type/example-id.tf example-id.tf.tmpl Id 1
//go:generate ../../../../bin/render -out examples/data-sources/morpheus_backup_type/example-name.tf example-name.tf.tmpl Name "File Backup"

func RenderBackupTypeDataSourceByIDConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	return renderBackupTypeDataSourceConfig(t, "example-id.tf.tmpl", map[string]string{
		"Id": "1",
	}, overrides)
}

func RenderBackupTypeDataSourceByNameConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	return renderBackupTypeDataSourceConfig(t, "example-name.tf.tmpl", map[string]string{
		"Name": "File Backup",
	}, overrides)
}

func renderBackupTypeDataSourceConfig(
	t *testing.T,
	templateName string,
	defaults map[string]string,
	overrides map[string]string,
) (string, error) {
	t.Helper()

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
	templatePath := filepath.Join(dir, templateName)

	return testhelpers.RenderExample(
		t,
		templatePath,
		args...,
	)
}
