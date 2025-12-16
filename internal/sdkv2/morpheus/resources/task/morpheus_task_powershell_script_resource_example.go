// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package task

import (
	"strings"

	"fmt"
	"path/filepath"
	"runtime"
	"testing"
	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/testhelpers"
)

//go:generate go run ../../../../../cmd/render -out examples/resources/morpheus_task_powershell_script/resource.tf morpheus_task_powershell_script_resource.tf.tmpl Name tfexample_powershell_local Code tfexample_powershell_local Labels "\"demo\", \"terraform\"" SourceType local ScriptContent "Write-Output \"testing\"" ElevatedShell true Retryable true RetryCount 1 RetryDelaySeconds 10 AllowCustomConfig true
//go:generate go run ../../../../../cmd/render -out examples/resources/morpheus_task_powershell_script/resource_git.tf morpheus_task_powershell_script_resource_git.tf.tmpl Name tfexample_powershell_git Code tfexample_powershell_git Labels "\"demo\", \"terraform\"" SourceType repository ResultType json ScriptPath example.ps VersionRef master RepositoryId 1 ElevatedShell true Retryable true RetryCount 1 RetryDelaySeconds 10 AllowCustomConfig true
//go:generate go run ../../../../../cmd/render -out examples/resources/morpheus_task_powershell_script/resource_url.tf morpheus_task_powershell_script_resource_url.tf.tmpl Name tfexample_powershell_url Code tfexample_powershell_url Labels "\"demo\", \"terraform\"" SourceType url ResultType json ScriptPath https://example.com/example.ps ElevatedShell true Retryable true RetryCount 1 RetryDelaySeconds 10 AllowCustomConfig true

func RenderMorpheusTaskGroovyScriptConfig(
	t *testing.T,
	name string,
	overrides map[string]string,
) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Name":              name,
		"Code":              name,
		"Labels":            "[\"demo\", \"terraform\"]",
		"SourceType":        "local",
		"ScriptContent":     "println \"hello\"",
		"Retryable":         "true",
		"RetryCount":        "1",
		"RetryDelaySeconds": "10",
		"AllowCustomConfig": "true",
	}

	for key, value := range overrides {
		defaults[key] = value
	}

	args := []string{}
	for key, value := range defaults {
		args = append(args, key, value)
	}

	resourceConfig, err := testhelpers.RenderExample(
		t,
		"morpheus_task_groovy_script_resource.tf.tmpl",
		args...,
	)
	if err != nil {
		return "", err
	}

	return resourceConfig, nil
}

func RenderTaskRubyScriptConfig(
	t *testing.T,
	name string,
	overrides map[string]string,
) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Name":              name,
		"Code":              name,
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

	resourceConfig, err := testhelpers.RenderExample(
		t,
		"task_ruby_script_resource.tf.tmpl",
		"Name", defaults["Name"],
		"Code", defaults["Code"],
		"Labels", defaults["Labels"],
		"SourceType", defaults["SourceType"],
		"ScriptContent", defaults["ScriptContent"],
		"Retryable", defaults["Retryable"],
		"RetryCount", defaults["RetryCount"],
		"RetryDelaySeconds", defaults["RetryDelaySeconds"],
		"AllowCustomConfig", defaults["AllowCustomConfig"],
	)

	return resourceConfig, err
}

func RenderTaskShellScriptConfig(
	t *testing.T,
	name string,
	overrides map[string]string,
) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Name":              name,
		"Code":              name,
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
		"Name", defaults["Name"],
		"Code", defaults["Code"],
		"Labels", defaults["Labels"],
		"SourceType", defaults["SourceType"],
		"ScriptContent", defaults["ScriptContent"],
		"Sudo", defaults["Sudo"],
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

// RenderTaskJavascriptConfig generates a terraform configuration string
// for task javascript resource.
// It accepts a name and a map to override default field values.
func RenderTaskJavascriptConfig(
	t *testing.T,
	name string,
	overrides map[string]string,
) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Name":              name,
		"Code":              name,
		"Labels":            `["demo","terraform"]`,
		"ScriptContent":     `console.log("testing")`,
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
	templatePath := filepath.Join(dir, "morpheus_task_javascript_resource.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath,
		"Name", defaults["Name"],
		"Code", defaults["Code"],
		"Labels", defaults["Labels"],
		"ScriptContent", defaults["ScriptContent"],
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

func RenderMorpheusTaskEmailGitConfig(
	t *testing.T,
	name string,
	overrides map[string]string,
) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Name":                     name,
		"Code":                     "tfexample_email_git",
		"Labels":                   `["demo","terraform"]`,
		"EmailAddress":             "<%=instance.createdByEmail%>",
		"Subject":                  "<%=instance.hostname%> provisioning complete",
		"Source":                   "repository",
		"ContentPath":              "example.txt",
		"RepositoryId":             "0",
		"VersionRef":               "main",
		"SkipWrappedEmailTemplate": "false",
		"Retryable":                "true",
		"RetryCount":               "1",
		"RetryDelaySeconds":        "10",
		"AllowCustomConfig":        "true",
	}

	for key, value := range overrides {
		defaults[key] = value
	}

	args := []string{}
	for key, value := range defaults {
		args = append(args, key, value)
	}

	// Get the directory where this source file is located
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("unable to get current file path")
	}
	dir := filepath.Dir(filename)
	templatePath := filepath.Join(dir, "morpheus_task_email_resource_git.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath, args...)
}

// RenderTaskWriteAttributesConfig generates a Terraform configuration for testing
// the task_write_attributes resource. It accepts overrides to customize field values.
func RenderTaskWriteAttributesConfig(
	t *testing.T,
	name string,
	overrides map[string]string,
) (string, error) {
	t.Helper()

	code := strings.ToLower(name)

	defaults := map[string]string{
		"Label1":            "demo",
		"Label2":            "terraform",
		"Attributes":        `{"demo":"test"}`,
		"Retryable":         "true",
		"RetryCount":        "1",
		"RetryDelaySeconds": "10",
		"AllowCustomConfig": "true",
	}

	// Apply overrides
	for key, value := range overrides {
		defaults[key] = value
	}

	// Get the directory where this source file is located
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("unable to get current file path")
	}
	dir := filepath.Dir(filename)
	templatePath := filepath.Join(dir, "../task/morpheus_task_write_attributes_resource.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath,
		"Name", name,
		"Code", code,
		"Label1", defaults["Label1"],
		"Label2", defaults["Label2"],
		"Attributes", defaults["Attributes"],
		"Retryable", defaults["Retryable"],
		"RetryCount", defaults["RetryCount"],
		"RetryDelaySeconds", defaults["RetryDelaySeconds"],
		"AllowCustomConfig", defaults["AllowCustomConfig"],
	)
}

func RenderTaskShellScriptUrlConfig(
	t *testing.T,
	name string,
	overrides map[string]string,
) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Name":              name,
		"Code":              name,
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
		"Name", defaults["Name"],
		"Code", defaults["Code"],
		"Labels", defaults["Labels"],
		"SourceType", defaults["SourceType"],
		"ResultType", defaults["ResultType"],
		"ScriptPath", defaults["ScriptPath"],
		"Sudo", defaults["Sudo"],
		"Retryable", defaults["Retryable"],
		"RetryCount", defaults["RetryCount"],
		"RetryDelaySeconds", defaults["RetryDelaySeconds"],
		"AllowCustomConfig", defaults["AllowCustomConfig"],
	)
}

func RenderMorpheusTaskPythonScriptGitConfig(
	t *testing.T,
	name string,
	overrides map[string]string,
) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Name":               name,
		"Code":               strings.ToLower(name),
		"Labels":             "[\"demo\", \"terraform\"]",
		"SourceType":         "repository",
		"ResultType":         "json",
		"ScriptPath":         "example.py",
		"VersionRef":         "master",
		"RepositoryId":       "0",
		"CommandArguments":   "example",
		"AdditionalPackages": "pyyaml",
		"PythonBinary":       "/usr/bin/python3",
		"Retryable":          "true",
		"RetryCount":         "1",
		"RetryDelaySeconds":  "10",
		"AllowCustomConfig":  "true",
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
	templatePath := filepath.Join(dir, "morpheus_task_python_script_resource_git.tf.tmpl")

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
		"CommandArguments", defaults["CommandArguments"],
		"AdditionalPackages", defaults["AdditionalPackages"],
		"PythonBinary", defaults["PythonBinary"],
		"Retryable", defaults["Retryable"],
		"RetryCount", defaults["RetryCount"],
		"RetryDelaySeconds", defaults["RetryDelaySeconds"],
		"AllowCustomConfig", defaults["AllowCustomConfig"],
	)
}

func RenderMorpheusTaskGroovyScriptUrlConfig(
	t *testing.T,
	name string,
	overrides map[string]string,
) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Name":              name,
		"Code":              name,
		"Labels":            "[\"demo\", \"terraform\"]",
		"SourceType":        "url",
		"ResultType":        "json",
		"ScriptPath":        "https://example.com/example.groovy",
		"Retryable":         "true",
		"RetryCount":        "1",
		"RetryDelaySeconds": "10",
		"AllowCustomConfig": "true",
	}

	for key, value := range overrides {
		defaults[key] = value
	}

	args := []string{}
	for key, value := range defaults {
		args = append(args, key, value)
	}

	resourceConfig, err := testhelpers.RenderExample(
		t,
		"morpheus_task_groovy_script_resource_url.tf.tmpl",
		args...,
	)
	if err != nil {
		return "", err
	}

	return resourceConfig, nil
}

func RenderMorpheusTaskPythonScriptConfig(
	t *testing.T, name string,
	overrides map[string]string,
) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Name":               name,
		"Code":               strings.ToLower(name),
		"Labels":             "[\"demo\", \"terraform\"]",
		"SourceType":         "local",
		"ScriptContent":      "print('morpheus')\\nprint('python')",
		"CommandArguments":   "example",
		"AdditionalPackages": "pyyaml",
		"PythonBinary":       "/usr/bin/python3",
		"Retryable":          "true",
		"RetryCount":         "1",
		"RetryDelaySeconds":  "10",
		"AllowCustomConfig":  "true",
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
	templatePath := filepath.Join(dir, "morpheus_task_python_script_resource.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath,
		"Name", defaults["Name"],
		"Code", defaults["Code"],
		"Labels", defaults["Labels"],
		"SourceType", defaults["SourceType"],
		"ScriptContent", defaults["ScriptContent"],
		"CommandArguments", defaults["CommandArguments"],
		"AdditionalPackages", defaults["AdditionalPackages"],
		"PythonBinary", defaults["PythonBinary"],
		"Retryable", defaults["Retryable"],
		"RetryCount", defaults["RetryCount"],
		"RetryDelaySeconds", defaults["RetryDelaySeconds"],
		"AllowCustomConfig", defaults["AllowCustomConfig"],
	)
}

func RenderTaskRubyScriptGitConfig(
	t *testing.T,
	name string,
	overrides map[string]string,
) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Name":              name,
		"Code":              name,
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

	resourceConfig, err := testhelpers.RenderExample(
		t,
		"task_ruby_script_resource_git.tf.tmpl",
		"Name", defaults["Name"],
		"Code", defaults["Code"],
		"Labels", defaults["Labels"],
		"SourceType", defaults["SourceType"],
		"ResultType", defaults["ResultType"],
		"ScriptPath", defaults["ScriptPath"],
		"VersionRef", defaults["VersionRef"],
		"RepositoryId", defaults["RepositoryId"],
		"Retryable", defaults["Retryable"],
		"RetryCount", defaults["RetryCount"],
		"RetryDelaySeconds", defaults["RetryDelaySeconds"],
		"AllowCustomConfig", defaults["AllowCustomConfig"],
	)

	return resourceConfig, err
}

func RenderTaskShellScriptGitConfig(
	t *testing.T,
	name string,
	overrides map[string]string,
) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Name":              name,
		"Code":              name,
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
		"Name", defaults["Name"],
		"Code", defaults["Code"],
		"Labels", defaults["Labels"],
		"SourceType", defaults["SourceType"],
		"ResultType", defaults["ResultType"],
		"ScriptPath", defaults["ScriptPath"],
		"VersionRef", defaults["VersionRef"],
		"RepositoryId", defaults["RepositoryId"],
		"Sudo", defaults["Sudo"],
		"Retryable", defaults["Retryable"],
		"RetryCount", defaults["RetryCount"],
		"RetryDelaySeconds", defaults["RetryDelaySeconds"],
		"AllowCustomConfig", defaults["AllowCustomConfig"],
	)
}

func RenderMorpheusTaskPythonScriptUrlConfig(
	t *testing.T,
	name string,
	overrides map[string]string,
) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Name":               name,
		"Code":               name,
		"Labels":             "[\"demo\", \"terraform\"]",
		"SourceType":         "url",
		"ResultType":         "json",
		"ScriptPath":         "https://example.com/example.py",
		"CommandArguments":   "example",
		"AdditionalPackages": "pyyaml",
		"PythonBinary":       "/usr/bin/python3",
		"Retryable":          "true",
		"RetryCount":         "1",
		"RetryDelaySeconds":  "10",
		"AllowCustomConfig":  "true",
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
	templatePath := filepath.Join(dir, "morpheus_task_python_script_resource_url.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath,
		"Name", defaults["Name"],
		"Code", defaults["Code"],
		"Labels", defaults["Labels"],
		"SourceType", defaults["SourceType"],
		"ResultType", defaults["ResultType"],
		"ScriptPath", defaults["ScriptPath"],
		"CommandArguments", defaults["CommandArguments"],
		"AdditionalPackages", defaults["AdditionalPackages"],
		"PythonBinary", defaults["PythonBinary"],
		"Retryable", defaults["Retryable"],
		"RetryCount", defaults["RetryCount"],
		"RetryDelaySeconds", defaults["RetryDelaySeconds"],
		"AllowCustomConfig", defaults["AllowCustomConfig"],
	)
}

func RenderMorpheusTaskEmailConfig(
	t *testing.T,
	name string,
	overrides map[string]string,
) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Name":                     name,
		"Code":                     "tfexample_email",
		"Labels":                   `["demo","terraform"]`,
		"EmailAddress":             "<%=instance.createdByEmail%>",
		"Subject":                  "<%=instance.hostname%> provisioning complete",
		"Source":                   "local",
		"Content":                  "Your instance <%=instance.hostname%> was provisioned.",
		"SkipWrappedEmailTemplate": "false",
		"Retryable":                "true",
		"RetryCount":               "1",
		"RetryDelaySeconds":        "10",
		"AllowCustomConfig":        "true",
	}

	for key, value := range overrides {
		defaults[key] = value
	}

	args := []string{}
	for key, value := range defaults {
		args = append(args, key, value)
	}

	// Get the directory where this source file is located
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("unable to get current file path")
	}
	dir := filepath.Dir(filename)
	templatePath := filepath.Join(dir, "morpheus_task_email_resource.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath, args...)
}

func RenderMorpheusTaskEmailUrlConfig(
	t *testing.T,
	name string,
	overrides map[string]string,
) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Name":                     name,
		"Code":                     "tfexample_email_url",
		"Labels":                   `["demo","terraform"]`,
		"EmailAddress":             "<%=instance.createdByEmail%>",
		"Subject":                  "<%=instance.hostname%> provisioning complete",
		"Source":                   "url",
		"ContentUrl":               "https://example.com/example.txt",
		"SkipWrappedEmailTemplate": "false",
		"Retryable":                "true",
		"RetryCount":               "1",
		"RetryDelaySeconds":        "10",
		"AllowCustomConfig":        "true",
	}

	for key, value := range overrides {
		defaults[key] = value
	}

	args := []string{}
	for key, value := range defaults {
		args = append(args, key, value)
	}

	// Get the directory where this source file is located
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("unable to get current file path")
	}
	dir := filepath.Dir(filename)
	templatePath := filepath.Join(dir, "morpheus_task_email_resource_url.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath, args...)
}

// RenderChefBootstrapConfig renders the task chef bootstrap
// resource configuration with the provided overrides applied to default values.
func RenderChefBootstrapConfig(
	t *testing.T,
	name string,
	overrides map[string]string,
) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Name":              name,
		"Code":              name,
		"Labels":            `"demo", "terraform"`,
		"ChefServerId":      "9",
		"Environment":       "dev",
		"RunList":           "role[web]",
		"DataBagKey":        "test123",
		"DataBagKeyPath":    "/etc/chef/databag_secret",
		"NodeName":          "demonode",
		"NodeAttributes":    `"test":"demo"`,
		"Retryable":         "true",
		"RetryCount":        "1",
		"RetryDelaySeconds": "10",
		"AllowCustomConfig": "true",
		"Visibility":        "public",
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
	templatePath := filepath.Join(dir, "morpheus_task_chef_bootstrap_resource.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath,
		"Name", defaults["Name"],
		"Code", defaults["Code"],
		"Labels", defaults["Labels"],
		"ChefServerId", defaults["ChefServerId"],
		"Environment", defaults["Environment"],
		"RunList", defaults["RunList"],
		"DataBagKey", defaults["DataBagKey"],
		"DataBagKeyPath", defaults["DataBagKeyPath"],
		"NodeName", defaults["NodeName"],
		"NodeAttributes", defaults["NodeAttributes"],
		"Retryable", defaults["Retryable"],
		"RetryCount", defaults["RetryCount"],
		"RetryDelaySeconds", defaults["RetryDelaySeconds"],
		"AllowCustomConfig", defaults["AllowCustomConfig"],
		"Visibility", defaults["Visibility"],
	)
}

func RenderMorpheusTaskGroovyScriptGitConfig(
	t *testing.T,
	name string,
	overrides map[string]string,
) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Name":              name,
		"Code":              name,
		"Labels":            "[\"demo\", \"terraform\"]",
		"SourceType":        "repository",
		"ResultType":        "json",
		"ScriptPath":        "example.groovy",
		"VersionRef":        "master",
		"RepositoryId":      "1",
		"Retryable":         "true",
		"RetryCount":        "1",
		"RetryDelaySeconds": "10",
		"AllowCustomConfig": "true",
	}

	for key, value := range overrides {
		defaults[key] = value
	}

	args := []string{}
	for key, value := range defaults {
		args = append(args, key, value)
	}

	resourceConfig, err := testhelpers.RenderExample(
		t,
		"morpheus_task_groovy_script_resource_git.tf.tmpl",
		args...,
	)
	if err != nil {
		return "", err
	}

	return resourceConfig, nil
}

// RenderTaskRestartConfig renders the task restart resource configuration
// with the provided name and field overrides.
func RenderTaskRestartConfig(
	t *testing.T,
	name string,
	overrides map[string]string,
) (string, error) {
	t.Helper()

	// Default values
	defaults := map[string]string{
		"Name":              name,
		"Code":              "tfexample_restart",
		"Labels":            `["demo", "terraform"]`,
		"Retryable":         "true",
		"RetryCount":        "1",
		"RetryDelaySeconds": "10",
		"AllowCustomConfig": "true",
	}

	// Apply overrides
	for key, value := range overrides {
		defaults[key] = value
	}

	resourceConfig, err := testhelpers.RenderExample(
		t,
		"morpheus_task_restart_resource.tf.tmpl",
		"Name", defaults["Name"],
		"Code", defaults["Code"],
		"Labels", defaults["Labels"],
		"Retryable", defaults["Retryable"],
		"RetryCount", defaults["RetryCount"],
		"RetryDelaySeconds", defaults["RetryDelaySeconds"],
		"AllowCustomConfig", defaults["AllowCustomConfig"],
	)
	if err != nil {
		return "", err
	}

	return resourceConfig, nil
}

func RenderTaskRubyScriptUrlConfig(
	t *testing.T,
	name string,
	overrides map[string]string,
) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Name":              name,
		"Code":              name,
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

	resourceConfig, err := testhelpers.RenderExample(
		t,
		"task_ruby_script_resource_url.tf.tmpl",
		"Name", defaults["Name"],
		"Code", defaults["Code"],
		"Labels", defaults["Labels"],
		"SourceType", defaults["SourceType"],
		"ResultType", defaults["ResultType"],
		"ScriptPath", defaults["ScriptPath"],
		"Retryable", defaults["Retryable"],
		"RetryCount", defaults["RetryCount"],
		"RetryDelaySeconds", defaults["RetryDelaySeconds"],
		"AllowCustomConfig", defaults["AllowCustomConfig"],
	)

	return resourceConfig, err
}

