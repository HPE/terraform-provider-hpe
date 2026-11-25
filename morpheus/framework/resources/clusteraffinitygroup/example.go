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

//go:generate ../../../../bin/render -out examples/resources/morpheus_cluster_affinity_group/example_members.tf example_members.tf.tmpl ClusterName "hvm-cluster" ServerNames "\"web-1\", \"web-2\"" ClusterId "data.hpe_morpheus_cluster.hvm.id" Name "web-keep-separate" AffinityType "KEEP_SEPARATE" Servers "[for s in data.hpe_morpheus_compute_server.member : s.id]"

// RenderClusterAffinityGroupMembersConfig renders the example of a group that
// owns its membership: servers is set, so the group is created with those
// members and every apply re-asserts the set.
//
// ClusterId and Servers are HCL expressions, not values, so a caller can point
// them at objects defined elsewhere in its configuration -- an instance's
// compute_servers, say. Servers has no default: left out, the group is
// rendered without it, which is how a test shows what dropping it does. The
// data sources that the published example resolves its members with are
// rendered only when ClusterName is given.
func RenderClusterAffinityGroupMembersConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"ClusterId":    "1",
		"Name":         "Example Affinity Group",
		"AffinityType": "KEEP_TOGETHER",
	}

	merged := testhelpers.MergeOverrides(defaults, overrides)

	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("unable to get current file path")
	}
	dir := filepath.Dir(filename)
	templatePath := filepath.Join(dir, "example_members.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath,
		testhelpers.RenderArgs(merged)...,
	)
}
