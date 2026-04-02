// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package ostypeimage

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
)

//go:generate ../../../../bin/render -out examples/resources/morpheus_os_type_image/example.tf example.tf.tmpl OsTypeId 1 VirtualImageId 42 CloudId 10 ProvisionTypeId 3

func RenderOsTypeImageConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"OsTypeId":        "75",
		"VirtualImageId":  "257",
		"CloudId":         "1",
		"ProvisionTypeId": "22",
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
