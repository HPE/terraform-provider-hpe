// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package template

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
)

//go:generate../../../../bin/render -out examples/resources/morpheus_script_template/resource.tf morpheus_script_template_resource.tf.tmpl 'Name' 'tf-terraform-script-template' 'Labels' "[\"demo\", \"template\", \"terraform\"]" 'ScriptType' 'bash' 'ScriptPhase' 'provision' 'ScriptContent' "echo \"testing\"" 'RunAsUser' 'root' 'Sudo' 'true'

// RenderScriptTemplateConfig renders the template with provided overrides
func RenderScriptTemplateConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Name":          "Example",
		"Labels":        "[\"demo\", \"template\", \"terraform\"]",
		"ScriptType":    "bash",
		"ScriptPhase":   "provision",
		"ScriptContent": "echo \"testing\"",
		"RunAsUser":     "root",
		"Sudo":          "true",
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
	templatePath := filepath.Join(dir, "morpheus_script_template_resource.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath,
		args...,
	)
}
