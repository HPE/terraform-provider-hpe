// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package networkdhcpserver

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
)

//go:generate ../../../../bin/render -out examples/resources/morpheus_network_dhcp_server/example.tf example.tf.tmpl NetworkIntegrationId "16" Name "Example DHCP Server" ServerIpAddress "192.168.1.1/24" LeaseTime "86400" EdgeCluster "qa-edge-cluster-01" ActiveEdgeNode "" StandbyEdgeNode ""
//go:generate ../../../../bin/render -out examples/resources/morpheus_network_dhcp_server/example_dynamic.tf example_dynamic.tf.tmpl NetworkIntegrationId "16" Name "Example DHCP Server" ServerIpAddress "192.168.1.1/24" LeaseTime "86400" EdgeCluster "qa-edge-cluster-01" ActiveEdgeNode "" StandbyEdgeNode ""

func RenderNetworkDhcpServerConfig(
	t *testing.T,
	overrides map[string]string,
) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"NetworkIntegrationId": "16",
		"Name":                 "Example DHCP Server",
		"ServerIpAddress":      "192.168.1.1/24", // API requires CIDR notation (e.g. "10.0.0.1/24")
		"LeaseTime":            "86400",
		"EdgeCluster":          "qa-edge-cluster-01",
		"ActiveEdgeNode":       "",
		"StandbyEdgeNode":      "",
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

func RenderNetworkDhcpServerDynamicConfig(
	t *testing.T,
	overrides map[string]string,
) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"NetworkIntegrationId": "16",
		"Name":                 "Example DHCP Server",
		"ServerIpAddress":      "192.168.1.1/24",
		"LeaseTime":            "86400",
		"EdgeCluster":          "qa-edge-cluster-01",
		"ActiveEdgeNode":       "",
		"StandbyEdgeNode":      "",
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
	templatePath := filepath.Join(dir, "example_dynamic.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath,
		args...,
	)
}
