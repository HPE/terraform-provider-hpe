// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package blueprint

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
)

//go:generate sh -c "../../../../bin/render -out examples/resources/morpheus_node_type/resource.tf node_type_resource.tf.tmpl Category 'tfexample' Labels '[\"demo\", \"nodeType\", \"terraform\"]' Name 'tf_example_node_type' ServicePortName1 'web' ServicePortName2 'secureweb' ServicePortPort1 '8080' ServicePortPort2 '8443' ServicePortProtocol1 'HTTP' ServicePortProtocol2 'HTTPS' ShortName 'tfexamplenodetype' Technology 'vmware' Version '2.0' VirtualImageId '10'"

// RenderNodeTypeConfig generates a Terraform configuration for the tenant resource.
// It accepts optional overrides for field values. Default values are used if not overridden.
func RenderNodeTypeConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Category": "tfexample",
		"Labels":   "[\"demo\", \"nodeType\", \"terraform\"]",
		"Name":     "tf_example_node_type",
		// FileTemplateIds / ScriptTemplateIds are intentionally unset here: the
		// shared example and TestAccMorpheusNodeTypeExampleOk do not define the
		// template fixtures to reference. The MORPH-11006 regression test in
		// node_type_resource_test.go passes them explicitly with self-contained
		// fixtures via the tmpl's optional blocks.
		"ServicePortName1":     "web",
		"ServicePortName2":     "secureweb",
		"ServicePortPort1":     "8080",
		"ServicePortPort2":     "8443",
		"ServicePortProtocol1": "HTTP",
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
