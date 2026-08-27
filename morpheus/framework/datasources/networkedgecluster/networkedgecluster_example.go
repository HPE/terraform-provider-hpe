// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkedgecluster

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
)

//go:generate ../../../../bin/render -out examples/data-sources/morpheus_network_edge_cluster/example-id.tf example-id.tf.tmpl NetworkServerId 1 Id 99
//go:generate ../../../../bin/render -out examples/data-sources/morpheus_network_edge_cluster/example-name.tf example-name.tf.tmpl NetworkServerId 1 Name "edge-cluster-01"

func RenderNetworkEdgeClusterDataSourceByIDConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"NetworkServerId": "1",
		"Id":              "99",
	}

	args := testhelpers.RenderArgs(testhelpers.MergeOverrides(defaults, overrides))

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

func RenderNetworkEdgeClusterDataSourceByNameConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"NetworkServerId": "1",
		"Name":            "edge-cluster-01",
	}

	args := testhelpers.RenderArgs(testhelpers.MergeOverrides(defaults, overrides))

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
