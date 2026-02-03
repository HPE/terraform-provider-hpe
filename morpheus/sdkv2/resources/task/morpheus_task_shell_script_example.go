// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package task

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
)

//go:generate sh -c "../../../../bin/render -out examples/resources/morpheus_task_shell_script/resource.tf morpheus_task_shell_script_resource.tf.tmpl Name tfexample_shell_local Code tfexample_shell_local Labels '[\"demo\", \"terraform\"]' SourceType local ScriptContent '  echo \"testing\"' Sudo true Retryable true RetryCount 1 RetryDelaySeconds 10 AllowCustomConfig true"
//go:generate sh -c "../../../../bin/render -out examples/resources/morpheus_task_shell_script/resource_git.tf morpheus_task_shell_script_resource_git.tf.tmpl Name tfexample_shell_git Code tfexample_shell_git Labels '[\"demo\", \"terraform\"]' SourceType repository ResultType json ScriptPath example.sh VersionRef master RepositoryId 1 Sudo true Retryable true RetryCount 1 RetryDelaySeconds 10 AllowCustomConfig true"
//go:generate sh -c "../../../../bin/render -out examples/resources/morpheus_task_shell_script/resource_url.tf morpheus_task_shell_script_resource_url.tf.tmpl Name tfexample_shell_url Code tfexample_shell_url Labels '[\"demo\", \"terraform\"]' SourceType url ResultType json ScriptPath https://example.com/example.sh Sudo true Retryable true RetryCount 1 RetryDelaySeconds 10 AllowCustomConfig true"

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
	templatePath := filepath.Join(dir, "morpheus_task_shell_script_resource.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath,
		args...,
	)
}

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
	templatePath := filepath.Join(dir, "morpheus_task_shell_script_resource_git.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath,
		args...,
	)
}

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
	templatePath := filepath.Join(dir, "morpheus_task_shell_script_resource_url.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath,
		args...,
	)
}
