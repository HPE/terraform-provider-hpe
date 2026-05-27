// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package container_script

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
)

//go:generate ../../../../bin/render -out examples/resources/morpheus_container_script/example.tf example.tf.tmpl Name "Install Dependencies" ScriptPhase "provision" ScriptType "bash" Script file("/path-to-file") SudoUser "true" FailOnError "true"

func RenderContainerScriptConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Name":        "Install Dependencies",
		"ScriptPhase": "provision",
		"ScriptType":  "bash",
		"Script":      "\"apt-get update && apt-get install -y curl wget\"",
		"SudoUser":    "true",
		"FailOnError": "true",
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
