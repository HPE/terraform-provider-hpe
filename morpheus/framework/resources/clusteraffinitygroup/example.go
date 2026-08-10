// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package clusteraffinitygroup

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
)

//go:generate ../../../../bin/render -out examples/resources/morpheus_cluster_affinity_group/example.tf example.tf.tmpl ClusterId "1" Name "Example Affinity Group"

func RenderClusterAffinityGroupConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"ClusterId": "1",
		"Name":      "Example Affinity Group",
	}

	merged := testhelpers.MergeOverrides(defaults, overrides)

	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("unable to get current file path")
	}
	dir := filepath.Dir(filename)
	templatePath := filepath.Join(dir, "example.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath,
		testhelpers.RenderArgs(merged)...,
	)
}
