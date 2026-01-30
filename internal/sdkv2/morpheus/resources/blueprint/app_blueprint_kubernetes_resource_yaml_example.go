// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package blueprint

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/testhelpers"
)

//go:generate sh -c " ../../../../../bin/render -out examples/resources/morpheus_app_blueprint_kubernetes/resource_yaml.tf app_blueprint_kubernetes_resource_yaml.tf.tmpl BlueprintContent '---\napiVersion: apps/v1\nkind: Deployment\nmetadata:\n name: nginx-deployment\n labels:\n app: nginx\nspec:\n replicas: 3\n selector:\n matchLabels:\n app: nginx\n template:\n metadata:\n labels:\n app: nginx\n spec:\n containers:\n - name: nginx\n image: nginx:1.14.2\n ports:\n - containerPort: 80' Category 'k8s' Description 'tf example kubernetes app blueprint' Name 'tf-kubernetes-app-blueprint-example-yaml' SourceType 'yaml'"

// RenderAppBlueprintKubernetesYamlConfig generates a Terraform configuration for the tenant resource.
// It accepts optional overrides for field values. Default values are used if not overridden.
func RenderAppBlueprintKubernetesYamlConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"BlueprintContent": "---\napiVersion: apps/v1\nkind: Deployment\nmetadata:\n name: nginx-deployment\n" +
			" labels:\n app: nginx\nspec:\n replicas: 3\n selector:\n matchLabels:\n app: nginx\n" +
			" template:\n metadata:\n labels:\n app: nginx\n spec:\n containers:\n - name: nginx\n" +
			" image: nginx:1.14.2\n ports:\n - containerPort: 80",
		"Category":    "k8s",
		"Description": "tf example kubernetes app blueprint",
		"Name":        "tf-kubernetes-app-blueprint-example-yaml",
		"SourceType":  "yaml",
	}

	// Apply overrides to defaults
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
	templatePath := filepath.Join(dir, "app_blueprint_kubernetes_resource_yaml.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath,
		args...,
	)
}
