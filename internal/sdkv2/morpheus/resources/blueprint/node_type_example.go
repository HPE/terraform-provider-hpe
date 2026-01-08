// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package blueprint

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/testhelpers"
)

//go:generate go run ../../../../../cmd/render -out examples/resources/node_type/node_type_resource.tf node_type_resource.tf.tmpl Category "tfexample" FileTemplateIds "[data.hpe_morpheus_file_template.tfexample.id, 113]" Labels "[\"demo\", \"nodeType\", \"terraform\"]" Name "tf_example_node_type" ScriptTemplateIds "[data.hpe_morpheus_script_template.tfscript1.id, data.hpe_morpheus_script_template.tfscript2.id]" ServicePortName "web" ServicePortName2 "secureweb" ServicePortPort "8080" ServicePortPort2 "8443" ServicePortProtocol "HTTP" ServicePortProtocol2 "HTTPS" ShortName "tfexamplenodetype" Technology "vmware" Version "2.0" VirtualImageId "10"

// RenderNodeTypeConfig generates a Terraform configuration for the tenant resource.
// It accepts optional overrides for field values. Default values are used if not overridden.
func RenderNodeTypeConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Category":             "tfexample",
		"FileTemplateIds":      "[113]",
		"Labels":               `["demo", "nodeType", "terraform"]`,
		"Name":                 "tf_example_node_type",
		"ScriptTemplateIds":    "[456, 789]",
		"ServicePortName":      "web",
		"ServicePortName2":     "secureweb",
		"ServicePortPort":      "8080",
		"ServicePortPort2":     "8443",
		"ServicePortProtocol":  "HTTP",
		"ServicePortProtocol2": "HTTPS",
		"ShortName":            "tfexamplenodetype",
		"Technology":           "vmware",
		"Version":              "2.0",
		"VirtualImageId":       "10",
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
	templatePath := filepath.Join(dir, "node_type_resource.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath,
		args...,
	)
}
