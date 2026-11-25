// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package virtualswitch

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
)

//go:generate ../../../../bin/render -out examples/data-sources/morpheus_virtual_switch/example-id.tf example-id.tf.tmpl Id 99
//go:generate ../../../../bin/render -out examples/data-sources/morpheus_virtual_switch/example-name.tf example-name.tf.tmpl ClusterId 1 Name "vswitch01"

func renderConfig(t *testing.T, tmpl string, defaults, overrides map[string]string) (string, error) {
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

	dir := filepath.Dir(filename)
	templatePath := filepath.Join(dir, tmpl)

	return testhelpers.RenderExample(t, templatePath, args...)
}

// RenderVirtualSwitchByIdConfig renders the id-based lookup example.
func RenderVirtualSwitchByIdConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	return renderConfig(t, "example-id.tf.tmpl", map[string]string{
		"Id": "99",
	}, overrides)
}

// RenderVirtualSwitchByNameConfig renders the cluster_id + name lookup example.
func RenderVirtualSwitchByNameConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	return renderConfig(t, "example-name.tf.tmpl", map[string]string{
		"ClusterId": "1",
		"Name":      "vswitch01",
	}, overrides)
}
