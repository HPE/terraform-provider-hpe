// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkserver

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
)

//go:generate ../../../../bin/render -out examples/resources/morpheus_network_server/example_juniper_sdn.tf example_juniper_sdn.tf.tmpl Name "HPE Morpheus SDN" TypeCode "morpheus-junipersdn-network" CloudId "1" ControllerAddress "10.0.0.10" ShowIntrospect "false" EnableSubTenants "false"
//go:generate ../../../../bin/render -out examples/resources/morpheus_network_server/example_nsxt.tf example_nsxt.tf.tmpl Name "NSX-T" TypeCode "nsx-t" CloudId "1" ServiceUrl "https://nsxt.example.com" ServiceUsername "admin" ServicePassword "example"

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

func RenderNetworkServerJuniperSdnConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	return renderExample(t, "example_juniper_sdn.tf.tmpl", map[string]string{
		"Name":              "HPE Morpheus SDN",
		"TypeCode":          "morpheus-junipersdn-network",
		"CloudId":           "1",
		"ControllerAddress": "10.0.0.10",
		"ShowIntrospect":    "false",
		"EnableSubTenants":  "false",
	}, overrides)
}

func RenderNetworkServerNsxtConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	return renderExample(t, "example_nsxt.tf.tmpl", map[string]string{
		"Name":            "NSX-T",
		"TypeCode":        "nsx-t",
		"CloudId":         "1",
		"ServiceUrl":      "https://nsxt.example.com",
		"ServiceUsername": "admin",
		"ServicePassword": "example",
	}, overrides)
}
