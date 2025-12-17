// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package task

import (
	"strings"

	"fmt"
	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/testhelpers"
	"path/filepath"
	"runtime"
	"testing"
)

//go:generate go run ../../../../../cmd/render -out examples/resources/morpheus_task_powershell_script/resource.tf morpheus_task_powershell_script_resource.tf.tmpl Name tfexample_powershell_local Code tfexample_powershell_local Labels "\"demo\", \"terraform\"" SourceType local ScriptContent "Write-Output \"testing\"" ElevatedShell true Retryable true RetryCount 1 RetryDelaySeconds 10 AllowCustomConfig true
//go:generate go run ../../../../../cmd/render -out examples/resources/morpheus_task_powershell_script/resource_git.tf morpheus_task_powershell_script_resource_git.tf.tmpl Name tfexample_powershell_git Code tfexample_powershell_git Labels "\"demo\", \"terraform\"" SourceType repository ResultType json ScriptPath example.ps VersionRef master RepositoryId 1 ElevatedShell true Retryable true RetryCount 1 RetryDelaySeconds 10 AllowCustomConfig true
//go:generate go run ../../../../../cmd/render -out examples/resources/morpheus_task_powershell_script/resource_url.tf morpheus_task_powershell_script_resource_url.tf.tmpl Name tfexample_powershell_url Code tfexample_powershell_url Labels "\"demo\", \"terraform\"" SourceType url ResultType json ScriptPath https://example.com/example.ps ElevatedShell true Retryable true RetryCount 1 RetryDelaySeconds 10 AllowCustomConfig true

func RenderTaskPowershellScriptConfig(
	t *testing.T,
	name string,
	overrides map[string]string,
) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Name":              name,
		"Code":              name,
		"Labels":            `"demo", "terraform"`,
		"SourceType":        "local",
		"ScriptContent":     `Write-Output \"testing\"`,
		"ElevatedShell":     "true",
		"Retryable":         "true",
		"RetryCount":        "1",
		"RetryDelaySeconds": "10",
		"AllowCustomConfig": "true",
	}

	for key, value := range overrides {
		defaults[key] = value
	}

	// Get the directory where this source file is located
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("unable to get current file path")
	}
	dir := filepath.Dir(filename)
	templatePath := filepath.Join(dir, "morpheus_task_powershell_script_resource.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath,
		"Name", defaults["Name"],
		"Code", defaults["Code"],
		"Labels", defaults["Labels"],
		"SourceType", defaults["SourceType"],
		"ScriptContent", defaults["ScriptContent"],
		"ElevatedShell", defaults["ElevatedShell"],
		"Retryable", defaults["Retryable"],
		"RetryCount", defaults["RetryCount"],
		"RetryDelaySeconds", defaults["RetryDelaySeconds"],
		"AllowCustomConfig", defaults["AllowCustomConfig"],
	)
}

func RenderTaskPowershellScriptGitConfig(
	t *testing.T,
	name string,
	overrides map[string]string,
) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Name":              name,
		"Code":              strings.ToLower(name),
		"Labels":            `"demo", "terraform"`,
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

	for key, value := range overrides {
		defaults[key] = value
	}

	// Get the directory where this source file is located
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("unable to get current file path")
	}
	dir := filepath.Dir(filename)
	templatePath := filepath.Join(dir, "morpheus_task_powershell_script_resource_git.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath,
		"Name", defaults["Name"],
		"Code", defaults["Code"],
		"Labels", defaults["Labels"],
		"SourceType", defaults["SourceType"],
		"ResultType", defaults["ResultType"],
		"ScriptPath", defaults["ScriptPath"],
		"VersionRef", defaults["VersionRef"],
		"RepositoryId", defaults["RepositoryId"],
		"ElevatedShell", defaults["ElevatedShell"],
		"Retryable", defaults["Retryable"],
		"RetryCount", defaults["RetryCount"],
		"RetryDelaySeconds", defaults["RetryDelaySeconds"],
		"AllowCustomConfig", defaults["AllowCustomConfig"],
	)
}

func RenderTaskPowershellScriptUrlConfig(
	t *testing.T,
	name string,
	overrides map[string]string,
) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Name":              name,
		"Code":              name,
		"Labels":            `"demo", "terraform"`,
		"SourceType":        "url",
		"ResultType":        "json",
		"ScriptPath":        "https://example.com/example.ps",
		"ElevatedShell":     "true",
		"Retryable":         "true",
		"RetryCount":        "1",
		"RetryDelaySeconds": "10",
		"AllowCustomConfig": "true",
	}

	for key, value := range overrides {
		defaults[key] = value
	}

	// Get the directory where this source file is located
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("unable to get current file path")
	}
	dir := filepath.Dir(filename)
	templatePath := filepath.Join(dir, "morpheus_task_powershell_script_resource_url.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath,
		"Name", defaults["Name"],
		"Code", defaults["Code"],
		"Labels", defaults["Labels"],
		"SourceType", defaults["SourceType"],
		"ResultType", defaults["ResultType"],
		"ScriptPath", defaults["ScriptPath"],
		"ElevatedShell", defaults["ElevatedShell"],
		"Retryable", defaults["Retryable"],
		"RetryCount", defaults["RetryCount"],
		"RetryDelaySeconds", defaults["RetryDelaySeconds"],
		"AllowCustomConfig", defaults["AllowCustomConfig"],
	)
}
