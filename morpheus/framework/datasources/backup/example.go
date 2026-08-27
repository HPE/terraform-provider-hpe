// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package backup

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
)

//go:generate ../../../../bin/render -out examples/data-sources/morpheus_backup/example-id.tf example-id.tf.tmpl Id 99
//go:generate ../../../../bin/render -out examples/data-sources/morpheus_backup/example-name.tf example-name.tf.tmpl Name "Example Backup"

func RenderBackupDataSourceByIDConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	return renderBackupDataSourceConfig(t, "example-id.tf.tmpl", map[string]string{
		"Id": "99",
	}, overrides)
}

func RenderBackupDataSourceByNameConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	return renderBackupDataSourceConfig(t, "example-name.tf.tmpl", map[string]string{
		"Name": "Example Backup",
	}, overrides)
}

func renderBackupDataSourceConfig(
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
