// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package ostype

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
)

//go:generate ../../../../bin/render -out examples/resources/morpheus_os_type/example.tf example.tf.tmpl Name "Example OS Type" Code "example.os.type" Platform "linux" BitCount "64" Description "An example OS type" OsFamily "debian" OsVersion "22.04" InstallAgent "true" CloudInitVersion "2"

func RenderOsTypeConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Name":             "Example OS Type",
		"Code":             "example.os.type",
		"Platform":         "linux",
		"BitCount":         "64",
		"Description":      "An example OS type",
		"OsFamily":         "debian",
		"OsVersion":        "22.04",
		"InstallAgent":     "true",
		"CloudInitVersion": "2",
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
