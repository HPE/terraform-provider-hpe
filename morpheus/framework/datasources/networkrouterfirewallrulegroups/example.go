// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkrouterfirewallrulegroups

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
)

//go:generate ../../../../bin/render -out examples/data-sources/morpheus_network_router_firewall_rule_groups/data-source.tf data-source.tf.tmpl RouterId 5

// RenderConfig renders the data source config template with the given overrides.
// Use from acceptance tests to produce valid HCL referencing the data source.
func RenderConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"RouterId": "5",
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
	templatePath := filepath.Join(dir, "data-source.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath,
		args...,
	)
}
