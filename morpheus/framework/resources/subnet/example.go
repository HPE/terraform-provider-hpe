// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package subnet

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
)

//go:generate ../../../../bin/render -out examples/resources/morpheus_subnet/example.tf example.tf.tmpl Name "Example Subnet" TypeId "1" NetworkId "1" Visibility "private"

func RenderSubnetConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Name":       "Example Subnet",
		"TypeId":     "8",
		"NetworkId":  "88",
		"Visibility": "private",
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
