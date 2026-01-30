// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package group

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/testhelpers"
)

//go:generate ../../../../../../bin/render -out examples/resources/morpheus_group/example.tf example.tf.tmpl Name "TestGroup" Location "here" Code "aCode" Label "aLabel"

func RenderGroupConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Name":     "ExampleUserRole",
		"Location": "here",
		"Code":     "aCode",
		"Label":    "aLabel",
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
	templatePath := filepath.Join(dir, "example.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath,
		args...,
	)
}
