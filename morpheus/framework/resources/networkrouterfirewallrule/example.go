// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkrouterfirewallrule

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
)

//go:generate ../../../../bin/render -out examples/resources/morpheus_network_router_firewall_rule/example.tf example.tf.tmpl RouterId "1" ParentGroupName "Example Firewall Rule Group" ParentId "data.hpe_morpheus_network_router_firewall_rule_group.example.external_id" Name "Example Firewall Rule" Policy "accept" Enabled "true"

func RenderNetworkRouterFirewallRuleConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"RouterId": "1",
		// parent_id is an expression, not a literal. The published example
		// (see the go:generate directive above) resolves it through the
		// hpe_morpheus_network_router_firewall_rule_group data source, which is
		// how a practitioner attaches a rule to a group that already exists.
		// Acceptance tests create their own group, so they reference that
		// resource directly and leave ParentGroupName unset, which omits the
		// data source block entirely.
		"ParentId": "hpe_morpheus_network_router_firewall_rule_group.example.external_id",
		"Name":     "Example Firewall Rule",
		"Policy":   "accept",
		"Enabled":  "true",
	}

	args := testhelpers.RenderArgs(testhelpers.MergeOverrides(defaults, overrides))

	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("unable to get current file path")
	}
	dir := filepath.Dir(filename)
	templatePath := filepath.Join(dir, "example.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath,
		args...,
	)
}
