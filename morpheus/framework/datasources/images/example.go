// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package images

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
)

//go:generate ../../../../bin/render -out examples/data-sources/morpheus_images/example.tf example.tf.tmpl ImageTypes "[\"qcow2\", \"raw\"]" NamePattern "^ubuntu" Status "active"

// RenderConfig renders the example with the given overrides, so acceptance
// tests exercise the same template the published example is generated from and
// the two cannot drift apart.
func RenderConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		// HVM provisions from either of these, and asking for both in one
		// request is what the repeatable image_type filter is for.
		"ImageTypes":  `["qcow2", "raw"]`,
		"NamePattern": "^ubuntu",
		"Status":      "active",
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
