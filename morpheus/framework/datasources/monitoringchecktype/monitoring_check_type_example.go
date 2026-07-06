// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package monitoringchecktype

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
)

//go:generate ../../../../bin/render -out examples/data-sources/morpheus_monitoring_check_type/example-id.tf example-id.tf.tmpl Id 5
//go:generate ../../../../bin/render -out examples/data-sources/morpheus_monitoring_check_type/example-name.tf example-name.tf.tmpl Name "Service Monitor"

// RenderMonitoringCheckTypeDataSourceByIDConfig renders the by-id monitoring
// check type data source example. It is exported so other resource and data
// source tests can compose a check type lookup block into their configurations.
func RenderMonitoringCheckTypeDataSourceByIDConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Id": "5",
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

// RenderMonitoringCheckTypeDataSourceByNameConfig renders the by-name monitoring
// check type data source example. It is exported so other resource and data
// source tests can compose a check type lookup block into their configurations.
func RenderMonitoringCheckTypeDataSourceByNameConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Name": "Service Monitor",
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
