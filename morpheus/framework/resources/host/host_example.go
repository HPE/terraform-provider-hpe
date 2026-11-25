// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package host

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
)

//go:generate ../../../../bin/render -out examples/resources/morpheus_host/example.tf example.tf.tmpl Name "example-host" CloudId 15174 GroupId 1 ResourcePoolId "pool-1" ComputeServerType "omega.baremetal.stub-server" ConsoleHost "10.0.0.1" ConsoleUsername "admin" NumNics 1

// RenderHostOmegaConfig renders the Omega bare-metal stub example used by the
// acceptance tests.
func RenderHostOmegaConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Name":              "example-host",
		"CloudId":           "15174",
		"GroupId":           "1",
		"ResourcePoolId":    "pool-1",
		"ComputeServerType": "omega.baremetal.stub-server",
		"ConsoleHost":       "10.0.0.1",
		"ConsoleUsername":   "admin",
		"NumNics":           "1",
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

	return testhelpers.RenderExample(t, templatePath, args...)
}
