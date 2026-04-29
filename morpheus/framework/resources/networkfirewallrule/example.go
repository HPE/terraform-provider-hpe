// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkfirewallrule

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
)

//go:generate ../../../../bin/render -out examples/resources/hpe_morpheus_network_firewall_rule/example.tf example.tf.tmpl NetworkIntegrationId "1" Name "Example Firewall Rule" Direction "Ingress" Policy "Accept" Enabled "true" RuleGroupId "1"

func RenderNetworkFirewallRuleConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"NetworkIntegrationId": "128",
		"Name":            "Example Firewall Rule",
		"Direction":       "Ingress",
		"Policy":          "Accept",
		"Enabled":         "true",
		"RuleGroupId":     "233264",
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
	templatePath := filepath.Join(dir, "example.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath,
		args...,
	)
}
