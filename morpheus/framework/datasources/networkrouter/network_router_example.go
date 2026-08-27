// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkrouter

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
)

//go:generate ../../../../bin/render -out examples/data-sources/morpheus_network_router/example-id.tf example-id.tf.tmpl Id 99
//go:generate ../../../../bin/render -out examples/data-sources/morpheus_network_router/example-name.tf example-name.tf.tmpl Name "Example name"

func RenderNetworkRouterByNameConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Name": "acceptance test router",
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

func RenderNetworkRouterByIdConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Id": "1",
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
