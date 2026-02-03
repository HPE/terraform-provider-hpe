// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package template

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
)

//go:generate ../../../../../bin/render -out examples/resources/morpheus_spec_template_kubernetes/resource_git.tf morpheus_spec_template_kubernetes_resource_git.tf.tmpl Name tf-kubernetes-spec-example-git SourceType repository RepositoryId 2 VersionRef main SpecPath ./spec.yaml

//go:generate ../../../../../bin/render -out examples/resources/morpheus_spec_template_kubernetes/resource_local.tf morpheus_spec_template_kubernetes_resource_local.tf.tmpl Name tf-terraform-spec-example-local SourceType local SpecContent "---\napiVersion: apps/v1\nkind: Deployment\nmetadata:\n  name: nginx-deployment\n  labels:\n    app: nginx\nspec:\n  replicas: 3\n  selector:\n    matchLabels:\n      app: nginx\n  template:\n    metadata:\n      labels:\n        app: nginx\n    spec:\n      containers:\n      - name: nginx\n        image: nginx:1.14.2\n        ports:\n        - containerPort: 80\n"

//go:generate ../../../../../bin/render -out examples/resources/morpheus_spec_template_kubernetes/resource_url.tf morpheus_spec_template_kubernetes_resource_url.tf.tmpl Name tf-kubernetes-spec-example-url SourceType url SpecPath http://example.com/spec.yaml

// RenderSpecTemplateKubernetesLocalConfig renders the configuration for the local
// Kubernetes spec template resource. Pass overrides as a map to customize field values.
func RenderSpecTemplateKubernetesLocalConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaultSpecContent := `---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-deployment
  labels:
    app: nginx
spec:
  replicas: 3
  selector:
    matchLabels:
      app: nginx
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
      - name: nginx
        image: nginx:1.14.2
        ports:
        - containerPort: 80`

	defaults := map[string]string{
		"Name":        "Example",
		"SourceType":  "local",
		"SpecContent": defaultSpecContent,
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
	templatePath := filepath.Join(dir, "spec_template_kubernetes_resource_local.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath,
		args...,
	)
}

// RenderSpecTemplateKubernetesGitConfig renders the configuration for the Git-based
// Kubernetes spec template resource. Pass overrides as a map to customize field values.
func RenderSpecTemplateKubernetesGitConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Name":         "Example",
		"SourceType":   "repository",
		"RepositoryId": "2",
		"VersionRef":   "main",
		"SpecPath":     "./spec.yaml",
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
	templatePath := filepath.Join(dir, "spec_template_kubernetes_resource_git.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath,
		args...,
	)
}

// RenderSpecTemplateKubernetesUrlConfig renders the configuration for the URL-based
// Kubernetes spec template resource. Pass overrides as a map to customize field values.
func RenderSpecTemplateKubernetesUrlConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Name":       "Example",
		"SourceType": "url",
		"SpecPath":   "http://example.com/spec.yaml",
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
	templatePath := filepath.Join(dir, "spec_template_kubernetes_resource_url.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath,
		args...,
	)
}
