// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

//go:generate ../../../../bin/render -out examples/resources/morpheus_network_router_bgp_neighbor/resource.tf example.tf.tmpl RouterId "42" IpAddress "10.0.0.1" Description "Example BGP neighbor" RemoteAs "65001" Weight "100" KeepAlive "60" HoldDown "180"

package networkrouterbgpneighbor

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
)

func RenderBgpNeighborConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"RouterId":    "1",
		"IpAddress":   "192.168.1.1",
		"Description": "example-bgp-neighbor",
		"RemoteAs":    "65001",
		"Weight":      "60",
		"KeepAlive":   "60",
		"HoldDown":    "180",
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
