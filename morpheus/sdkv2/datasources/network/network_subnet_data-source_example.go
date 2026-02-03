// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package network

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
)

//go:generate sh -c " ../../../../../bin/render -out examples/data-sources/morpheus_network_subnet/data-source.tf network_subnet_data-source.tf.tmpl Name '\"TF Example Network Subnet\"' NetworkId '166333'"

// RenderNetworkSubnetConfig generates a Terraform configuration for the network_subnet resource.
// It accepts optional overrides for field values. Default values are used if not overridden.
func RenderNetworkSubnetConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Name":      "\"TF Example Network Subnet\"",
		"NetworkId": "166333",
	}

	// Apply overrides to defaults
	for key, value := range overrides {
		defaults[key] = value
	}

	var args []string
	for key, value := range defaults {
		args = append(args, key, value)
	}

	// Get the directory where this source file is located
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("unable to get current file path")
	}
	dir := filepath.Dir(filename)
	templatePath := filepath.Join(dir, "network_subnet_data-source.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath,
		args...,
	)
}
