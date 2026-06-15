// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networktransportzone

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
)

//go:generate ../../../../bin/render -out examples/data-sources/morpheus_network_transport_zone/example-id.tf example-id.tf.tmpl NetworkServerId 1 Id 99
//go:generate ../../../../bin/render -out examples/data-sources/morpheus_network_transport_zone/example-name.tf example-name.tf.tmpl NetworkServerId 1 Name "overlay-tz-01"

func RenderNetworkTransportZoneDataSourceByIDConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"NetworkServerId": "1",
		"Id":              "99",
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
	templatePath := filepath.Join(dir, "example-id.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath,
		args...,
	)
}

func RenderNetworkTransportZoneDataSourceByNameConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"NetworkServerId": "1",
		"Name":            "overlay-tz-01",
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
	templatePath := filepath.Join(dir, "example-name.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath,
		args...,
	)
}
