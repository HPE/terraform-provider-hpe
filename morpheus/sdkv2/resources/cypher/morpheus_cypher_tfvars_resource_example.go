// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package cypher

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
)

//go:generate sh -c "../../../../bin/render -out examples/resources/morpheus_cypher_tfvars/resource.tf hpe_morpheus_cypher_tfvars_resource.tf.tmpl Key securetfvars Value 'account=12345\npassword=supersecure' Ttl 86400"

func RenderCypherTfvarsConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Key":   "Example",
		"Ttl":   "86400",
		"Value": "account=12345\npassword=supersecure",
	}

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
	templatePath := filepath.Join(dir, "hpe_morpheus_cypher_tfvars_resource.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath,
		args...,
	)
}
