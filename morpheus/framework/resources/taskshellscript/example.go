// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package taskshellscript

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
)

//go:generate ../../../../bin/render -out examples/resources/morpheus_task_shell_script/resource.tf example.tf.tmpl Name tfexample_shell_local Code tfexample_shell_local Labels "[\"demo\", \"terraform\"]" SourceType local ScriptContent "  echo \"testing\"" Sudo true Retryable true RetryCount 1 RetryDelaySeconds 10 AllowCustomConfig true
//go:generate ../../../../bin/render -out examples/resources/morpheus_task_shell_script/resource_git.tf example_git.tf.tmpl Name tfexample_shell_git Code tfexample_shell_git Labels "[\"demo\", \"terraform\"]" SourceType repository ResultType json ScriptPath example.sh VersionRef master RepositoryId 1 Sudo true Retryable true RetryCount 1 RetryDelaySeconds 10 AllowCustomConfig true
//go:generate ../../../../bin/render -out examples/resources/morpheus_task_shell_script/resource_url.tf example_url.tf.tmpl Name tfexample_shell_url Code tfexample_shell_url Labels "[\"demo\", \"terraform\"]" SourceType url ResultType json ScriptPath https://example.com/example.sh Sudo true Retryable true RetryCount 1 RetryDelaySeconds 10 AllowCustomConfig true

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

// RenderTaskShellScriptConfig renders the local-script example (resource
// address hpe_morpheus_task_shell_script.tfexample_shell_local). Kept with the
// same signature and semantics as the deleted SDKv2 task.RenderTaskShellScriptConfig
// so existing callers (sdkv2 job task tests) keep working.
func RenderTaskShellScriptConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Name":              "Example",
		"Code":              "example",
		"Labels":            "[\"demo\", \"terraform\"]",
		"SourceType":        "local",
		"ScriptContent":     "  echo \"testing\"",
		"Sudo":              "true",
		"Retryable":         "true",
		"RetryCount":        "1",
		"RetryDelaySeconds": "10",
		"AllowCustomConfig": "true",
	}

	return renderConfig(t, "example.tf.tmpl", defaults, overrides)
}

// RenderTaskShellScriptGitConfig renders the repository-script example.
func RenderTaskShellScriptGitConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Name":              "Example",
		"Code":              "example",
		"Labels":            "[\"demo\", \"terraform\"]",
		"SourceType":        "repository",
		"ResultType":        "json",
		"ScriptPath":        "example.sh",
		"VersionRef":        "master",
		"RepositoryId":      "0",
		"Sudo":              "true",
		"Retryable":         "true",
		"RetryCount":        "1",
		"RetryDelaySeconds": "10",
		"AllowCustomConfig": "true",
	}

	return renderConfig(t, "example_git.tf.tmpl", defaults, overrides)
}

// RenderTaskShellScriptUrlConfig renders the url-script example.
func RenderTaskShellScriptUrlConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Name":              "Example",
		"Code":              "example",
		"Labels":            "[\"demo\", \"terraform\"]",
		"SourceType":        "url",
		"ResultType":        "json",
		"ScriptPath":        "https://example.com/example.sh",
		"Sudo":              "true",
		"Retryable":         "true",
		"RetryCount":        "1",
		"RetryDelaySeconds": "10",
		"AllowCustomConfig": "true",
	}

	return renderConfig(t, "example_url.tf.tmpl", defaults, overrides)
}
