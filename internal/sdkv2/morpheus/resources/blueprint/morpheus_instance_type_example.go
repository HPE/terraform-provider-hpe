// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package blueprint

import (
	"strings"

	"fmt"
	"path/filepath"
	"runtime"
	"testing"
	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/testhelpers"
)

//go:generate go run ../../../../../cmd/render -out examples/resources/morpheus_instance_type/resource.tf morpheus_instance_type_resource.tf.tmpl Name "tf_example_instance" Code "tf_example_instance" Description "Terraform Example Instance Type" Labels "[\"demo\", \"instance\", \"terraform\"]" Category "web" Visibility "private" ImagePath "tfexample.png" ImageName "tfexample.png" Featured "false" EnableDeployments "true" EnableScaling "true" EnableSettings "true" EnvironmentPrefix "TFEXAMPLE_DEMO" OptionTypeIds "[1910, 1912]" EvarFirstName "first" EvarFirstValue "first" EvarFirstExport "true" EvarSecondName "second" EvarSecondMaskedValue "second" EvarSecondExport "false"

func RenderInstanceTypeConfig(
	t *testing.T,
	name string,
	overrides map[string]string,
) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Name":                  name,
		"Code":                  strings.ToLower(name),
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

	for k, v := range overrides {
		defaults[k] = v
	}

	// Get the directory where this source file is located
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("unable to get current file path")
	}
	dir := filepath.Dir(filename)
	templatePath := filepath.Join(dir, "morpheus_instance_type_resource.tf.tmpl")

	args := make([]string, 0, len(defaults)*2)
	for k, v := range defaults {
		args = append(args, k, v)
	}

	resourceConfig, err := testhelpers.RenderExample(
		t,
		templatePath,
		args...,
	)
	if err != nil {
		return "", err
	}

	return resourceConfig, nil
}

