// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package subnettype

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
)

//go:generate ../../../../bin/render -out examples/data-sources/morpheus_subnet_type/example-id.tf example-id.tf.tmpl Id 6
//go:generate ../../../../bin/render -out examples/data-sources/morpheus_subnet_type/example-name.tf example-name.tf.tmpl Name "VLAN"

// RenderSubnetTypeDataSourceByIDConfig renders the by-id subnet type data source
// example. It is exported so other resource and data source tests can compose a
// subnet type lookup block into their own configurations.
func RenderSubnetTypeDataSourceByIDConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Id": "6",
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

// RenderSubnetTypeDataSourceByNameConfig renders the by-name subnet type data
// source example. It is exported so other resource and data source tests can
// compose a subnet type lookup block into their own configurations.
func RenderSubnetTypeDataSourceByNameConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Name": "VLAN",
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
