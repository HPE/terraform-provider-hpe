// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package clusterlayout

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
)

//go:generate ../../../../bin/render -out examples/data-sources/morpheus_cluster_layout/example-id.tf example-id.tf.tmpl Id 7
//go:generate ../../../../bin/render -out examples/data-sources/morpheus_cluster_layout/example-name.tf example-name.tf.tmpl Name "Kubernetes Cluster"

// RenderClusterLayoutDataSourceByIDConfig renders the by-id cluster layout data
// source example. It is exported so other resource and data source tests can
// compose a cluster layout lookup block into their own configurations.
func RenderClusterLayoutDataSourceByIDConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Id": "7",
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

// RenderClusterLayoutDataSourceByNameConfig renders the by-name cluster layout
// data source example. It is exported so other resource and data source tests
// can compose a cluster layout lookup block into their own configurations.
func RenderClusterLayoutDataSourceByNameConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Name": "Kubernetes Cluster",
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
