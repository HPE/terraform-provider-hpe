// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package task

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
)

//go:generate /bin/sh -c " ../../../../../bin/render -out examples/resources/morpheus_task_email/resource.tf morpheus_task_email_resource.tf.tmpl Name tfexample_email Code tfexample_email Labels '[\"demo\",\"terraform\"]' EmailAddress '<%=instance.createdByEmail%>' Subject '<%=instance.hostname%> provisioning complete' Source local Content 'Your instance <%=instance.hostname%> was provisioned.' SkipWrappedEmailTemplate false Retryable true RetryCount 1 RetryDelaySeconds 10 AllowCustomConfig true"
//go:generate /bin/sh -c " ../../../../../bin/render -out examples/resources/morpheus_task_email/resource_git.tf morpheus_task_email_resource_git.tf.tmpl Name tfexample_email_git Code tfexample_email_git Labels '[\"demo\",\"terraform\"]' EmailAddress '<%=instance.createdByEmail%>' Subject '<%=instance.hostname%> provisioning complete' Source repository ContentPath example.txt RepositoryId 1 VersionRef main SkipWrappedEmailTemplate false Retryable true RetryCount 1 RetryDelaySeconds 10 AllowCustomConfig true"
//go:generate /bin/sh -c " ../../../../../bin/render -out examples/resources/morpheus_task_email/resource_url.tf morpheus_task_email_resource_url.tf.tmpl Name tfexample_email_url Code tfexample_email_url Labels '[\"demo\",\"terraform\"]' EmailAddress '<%=instance.createdByEmail%>' Subject '<%=instance.hostname%> provisioning complete' Source url ContentUrl https://example.com/example.txt SkipWrappedEmailTemplate false Retryable true RetryCount 1 RetryDelaySeconds 10 AllowCustomConfig true"

func RenderTaskEmailConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Name":                     "Example",
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
	templatePath := filepath.Join(dir, "morpheus_task_email_resource.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath,
		args...,
	)
}

func RenderTaskEmailGitConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Name":                     "Example",
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
	templatePath := filepath.Join(dir, "morpheus_task_email_resource_git.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath,
		args...,
	)
}

func RenderTaskEmailUrlConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Name":                     "Example",
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
	templatePath := filepath.Join(dir, "morpheus_task_email_resource_url.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath,
		args...,
	)
}
