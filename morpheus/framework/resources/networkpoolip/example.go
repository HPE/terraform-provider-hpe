// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkpoolip

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
)

//go:generate ../../../../bin/render -out examples/resources/morpheus_network_pool_ip/example.tf example.tf.tmpl NetworkPoolId "1" Hostname "host01" IpAddress "10.0.0.50" RefType "ComputeServer" RefId "1"

func RenderNetworkPoolIpConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"NetworkPoolId": "1",
		"Hostname":      "host01",
		"IpAddress":     "10.0.0.50",
		"RefType":       "ComputeServer",
		"RefId":         "1",
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
