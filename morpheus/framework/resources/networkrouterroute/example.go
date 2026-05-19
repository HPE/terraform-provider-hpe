// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

//go:generate ../../../../bin/render -out examples/resources/hpe_morpheus_network_router_route/resource.tf example.tf.tmpl RouterId "42" Name "example-route" Network "10.0.0.0/24" NextHop "10.0.0.1" Description "Example route" Mtu "1500" Enabled "true" DefaultRoute "false"

package networkrouterroute

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
)

func RenderRouteConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"RouterId":     "1",
		"Name":         "example-route",
		"Network":      "10.0.0.0/24",
		"NextHop":      "10.0.0.1",
		"Description":  "example-route",
		"Mtu":          "1500",
		"Enabled":      "true",
		"DefaultRoute": "false",
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
