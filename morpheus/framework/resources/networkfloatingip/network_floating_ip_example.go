// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkfloatingip

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
)

//go:generate ../../../../bin/render -out examples/resources/morpheus_network_floating_ip/example.tf example.tf.tmpl FloatingIpPoolId "data.hpe_morpheus_option_source_values.floating_ip_pools.values[0].value" ResourcePoolId "1" NetworkServerName "HPE Morpheus SDN"

func renderExample(t *testing.T, tmpl string, defaults, overrides map[string]string) (string, error) {
	t.Helper()

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

	templatePath := filepath.Join(filepath.Dir(filename), tmpl)

	return testhelpers.RenderExample(t, templatePath, args...)
}

// RenderNetworkFloatingIpConfig renders the floating IP allocation example.
func RenderNetworkFloatingIpConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	return renderExample(t, "example.tf.tmpl", map[string]string{
		"FloatingIpPoolId":  "data.hpe_morpheus_option_source_values.floating_ip_pools.values[0].value",
		"ResourcePoolId":    "1",
		"NetworkServerName": "HPE Morpheus SDN",
	}, overrides)
}
