// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package blueprint

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/testhelpers"
)

//go:generate ../../../../../bin/render -out examples/resources/morpheus_instance_type/resource.tf morpheus_instance_type_resource.tf.tmpl Name "tf_example_instance" Code "tf_example_instance" Description "Terraform Example Instance Type" Labels "[\"demo\", \"instance\", \"terraform\"]" Category "web" Visibility "private" ImagePath "tfexample.png" ImageName "tfexample.png" Featured "false" EnableDeployments "true" EnableScaling "true" EnableSettings "true" EnvironmentPrefix "TFEXAMPLE_DEMO" OptionTypeIds "[1910, 1912]" EvarFirstName "first" EvarFirstValue "first" EvarFirstExport "true" EvarSecondName "second" EvarSecondMaskedValue "second" EvarSecondExport "false"

func RenderInstanceTypeConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Name":                  "Example",
		"Code":                  "example",
		"Description":           "Terraform Example Instance Type",
		"Labels":                `["demo", "instance", "terraform"]`,
		"Category":              "web",
		"Visibility":            "private",
		"ImagePath":             "tfexample.png",
		"ImageName":             "tfexample.png",
		"Featured":              "false",
		"EnableDeployments":     "true",
		"EnableScaling":         "true",
		"EnableSettings":        "true",
		"EnvironmentPrefix":     "TFEXAMPLE_DEMO",
		"OptionTypeIds":         "[1910, 1912]",
		"EvarFirstName":         "first",
		"EvarFirstValue":        "first",
		"EvarFirstExport":       "true",
		"EvarSecondName":        "second",
		"EvarSecondMaskedValue": "second",
		"EvarSecondExport":      "false",
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
	templatePath := filepath.Join(dir, "morpheus_instance_type_resource.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath,
		args...,
	)
}
