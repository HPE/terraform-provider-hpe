// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package cloudaffinitygroup

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
)

//go:generate ../../../../bin/render -out examples/resources/morpheus_cloud_affinity_group/example.tf example.tf.tmpl CloudId "1" Name "Example Affinity Group" PoolId "1"

func RenderCloudAffinityGroupConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"CloudId": "1",
		"Name":    "Example Affinity Group",
		// Morpheus rejects a cloud affinity group created without a pool, so
		// the template always renders one and the default keeps it valid HCL.
		"PoolId": "1",
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
