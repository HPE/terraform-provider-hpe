// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package taskpowershellscript

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
)

//go:generate ../../../../bin/render -out examples/resources/morpheus_task_powershell_script/resource.tf example.tf.tmpl Name tfexample_powershell_local Code tfexample_powershell_local Labels "[\"demo\", \"terraform\"]" SourceType local ScriptContent "Write-Output \"testing\"" ElevatedShell true Retryable true RetryCount 1 RetryDelaySeconds 10 AllowCustomConfig true
//go:generate ../../../../bin/render -out examples/resources/morpheus_task_powershell_script/resource_git.tf example_git.tf.tmpl Name tfexample_powershell_git Code tfexample_powershell_git Labels "[\"demo\", \"terraform\"]" SourceType repository ResultType json ScriptPath example.ps VersionRef master RepositoryId 1 ElevatedShell true Retryable true RetryCount 1 RetryDelaySeconds 10 AllowCustomConfig true
//go:generate ../../../../bin/render -out examples/resources/morpheus_task_powershell_script/resource_url.tf example_url.tf.tmpl Name tfexample_powershell_url Code tfexample_powershell_url Labels "[\"demo\", \"terraform\"]" SourceType url ResultType json ScriptPath https://example.com/example.ps ElevatedShell true Retryable true RetryCount 1 RetryDelaySeconds 10 AllowCustomConfig true

func renderConfig(t *testing.T, tmpl string, defaults, overrides map[string]string) (string, error) {
	t.Helper()

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
	templatePath := filepath.Join(dir, tmpl)

	return testhelpers.RenderExample(t, templatePath, args...)
}

// RenderTaskPowershellScriptConfig renders the local-script example.
func RenderTaskPowershellScriptConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Name":              "Example",
		"Code":              "example",
		"Labels":            "[\"demo\", \"terraform\"]",
		"SourceType":        "local",
		"ScriptContent":     "Write-Output \"testing\"",
		"ElevatedShell":     "true",
		"Retryable":         "true",
		"RetryCount":        "1",
		"RetryDelaySeconds": "10",
		"AllowCustomConfig": "true",
	}

	return renderConfig(t, "example.tf.tmpl", defaults, overrides)
}

// RenderTaskPowershellScriptGitConfig renders the repository-script example.
func RenderTaskPowershellScriptGitConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Name":              "Example",
		"Code":              "example",
		"Labels":            "[\"demo\", \"terraform\"]",
		"SourceType":        "repository",
		"ResultType":        "json",
		"ScriptPath":        "example.ps",
		"VersionRef":        "master",
		"RepositoryId":      "0",
		"ElevatedShell":     "true",
		"Retryable":         "true",
		"RetryCount":        "1",
		"RetryDelaySeconds": "10",
		"AllowCustomConfig": "true",
	}

	return renderConfig(t, "example_git.tf.tmpl", defaults, overrides)
}

// RenderTaskPowershellScriptUrlConfig renders the url-script example.
func RenderTaskPowershellScriptUrlConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Name":              "Example",
		"Code":              "example",
		"Labels":            "[\"demo\", \"terraform\"]",
		"SourceType":        "url",
		"ResultType":        "json",
		"ScriptPath":        "https://example.com/example.ps",
		"ElevatedShell":     "true",
		"Retryable":         "true",
		"RetryCount":        "1",
		"RetryDelaySeconds": "10",
		"AllowCustomConfig": "true",
	}

	return renderConfig(t, "example_url.tf.tmpl", defaults, overrides)
}
