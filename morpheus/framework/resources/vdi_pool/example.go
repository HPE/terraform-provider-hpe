// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package vdi_pool

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
)

//go:generate ../../../../bin/render -out examples/resources/morpheus_vdi_pool/example.tf example.tf.tmpl Name "Developer Desktops" Description "VDI pool for development team" MaxPoolSize "10" MinIdle "2" InitialPoolSize "3" Enabled "true" PersistentUser "true" IdleTimeout "30"

func RenderVdiPoolConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Name":            "Developer Desktops",
		"Description":     "VDI pool for development team",
		"MaxPoolSize":     "10",
		"MinIdle":         "2",
		"InitialPoolSize": "3",
		"Enabled":         "true",
		"PersistentUser":  "true",
		"IdleTimeout":     "30",
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
