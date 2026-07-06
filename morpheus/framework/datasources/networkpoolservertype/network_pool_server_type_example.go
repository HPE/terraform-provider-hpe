// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkpoolservertype

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
)

//go:generate ../../../../bin/render -out examples/data-sources/morpheus_network_pool_server_type/example-id.tf example-id.tf.tmpl Id 3
//go:generate ../../../../bin/render -out examples/data-sources/morpheus_network_pool_server_type/example-name.tf example-name.tf.tmpl Name "Infoblox"

// RenderNetworkPoolServerTypeDataSourceByIDConfig renders the by-id network pool
// server type data source example. It is exported so other resource and data
// source tests can compose a pool server type lookup block into their configs.
func RenderNetworkPoolServerTypeDataSourceByIDConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Id": "3",
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

// RenderNetworkPoolServerTypeDataSourceByNameConfig renders the by-name network
// pool server type data source example. It is exported so other resource and
// data source tests can compose a pool server type lookup block into configs.
func RenderNetworkPoolServerTypeDataSourceByNameConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Name": "Infoblox",
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
