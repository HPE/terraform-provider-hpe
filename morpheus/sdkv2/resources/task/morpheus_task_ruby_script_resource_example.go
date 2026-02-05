// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package task

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
)

//go:generate../../../../bin/render -out examples/resources/morpheus_task_ruby_script/resource.tf task_ruby_script_resource.tf.tmpl Name tfexample_ruby_local Code tfexample_ruby_local Labels "\"demo\", \"terraform\"" SourceType local ScriptContent "puts \"testing\"" Retryable true RetryCount 1 RetryDelaySeconds 10 AllowCustomConfig true
//go:generate../../../../bin/render -out examples/resources/morpheus_task_ruby_script/resource_git.tf task_ruby_script_resource_git.tf.tmpl Name tfexample_ruby_git Code tfexample_ruby_git Labels "\"demo\", \"terraform\"" SourceType repository ResultType json ScriptPath example.rb VersionRef master RepositoryId 1 Retryable true RetryCount 1 RetryDelaySeconds 10 AllowCustomConfig true
//go:generate../../../../bin/render -out examples/resources/morpheus_task_ruby_script/resource_url.tf task_ruby_script_resource_url.tf.tmpl Name tfexample_ruby_url Code tfexample_ruby_url Labels "\"demo\", \"terraform\"" SourceType url ResultType json ScriptPath https://example.com/example.rb Retryable true RetryCount 1 RetryDelaySeconds 10 AllowCustomConfig true

func RenderTaskRubyScriptConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Name":              "Example",
		"Code":              "example",
		"Labels":            "\"demo\", \"terraform\"",
		"SourceType":        "local",
		"ScriptContent":     "puts \"testing\"",
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
	templatePath := filepath.Join(dir, "task_ruby_script_resource.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath,
		args...,
	)
}

func RenderTaskRubyScriptGitConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Name":              "Example",
		"Code":              "example",
		"Labels":            "\"demo\", \"terraform\"",
		"SourceType":        "repository",
		"ResultType":        "json",
		"ScriptPath":        "example.rb",
		"VersionRef":        "master",
		"RepositoryId":      "0",
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
	templatePath := filepath.Join(dir, "task_ruby_script_resource_git.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath,
		args...,
	)
}

func RenderTaskRubyScriptUrlConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Name":              "Example",
		"Code":              "example",
		"Labels":            "\"demo\", \"terraform\"",
		"SourceType":        "url",
		"ResultType":        "json",
		"ScriptPath":        "https://example.com/example.rb",
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
	templatePath := filepath.Join(dir, "task_ruby_script_resource_url.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath,
		args...,
	)
}
