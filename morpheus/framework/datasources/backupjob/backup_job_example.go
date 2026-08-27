// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package backupjob

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
)

//go:generate ../../../../bin/render -out examples/data-sources/morpheus_backup_job/example-id.tf example-id.tf.tmpl Id 12
//go:generate ../../../../bin/render -out examples/data-sources/morpheus_backup_job/example-name.tf example-name.tf.tmpl Name "Nightly VM Backup"

// RenderBackupJobDataSourceByIDConfig renders the by-id backup job data source
// example. It is exported so other resource and data source tests can compose a
// backup job lookup block into their own configurations.
func RenderBackupJobDataSourceByIDConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Id": "12",
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
	templatePath := filepath.Join(dir, "example-id.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath,
		args...,
	)
}

// RenderBackupJobDataSourceByNameConfig renders the by-name backup job data
// source example. It is exported so other resource and data source tests can
// compose a backup job lookup block into their own configurations.
func RenderBackupJobDataSourceByNameConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Name": "Nightly VM Backup",
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
	templatePath := filepath.Join(dir, "example-name.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath,
		args...,
	)
}
