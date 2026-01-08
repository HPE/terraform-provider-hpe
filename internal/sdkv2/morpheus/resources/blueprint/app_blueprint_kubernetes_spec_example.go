// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package blueprint

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/testhelpers"
)

//go:generate go run ../../../../../cmd/render -out examples/resources/app_blueprint_kubernetes/app_blueprint_kubernetes_resource_spec.tf app_blueprint_kubernetes_resource_spec.tf.tmpl Category "k8s" Description "tf example kubernetes app blueprint" Name "tf-kubernetes-app-blueprint-example-spec" SourceType "spec" SpecTemplateIds "[2, 3]"

// RenderAppBlueprintKubernetesSpecConfig generates a Terraform configuration for the tenant resource.
// It accepts optional overrides for field values. Default values are used if not overridden.

func RenderAppBlueprintKubernetesSpecConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Category":        "k8s",
		"Description":     "tf example kubernetes app blueprint",
		"Name":            "tf-kubernetes-app-blueprint-example-spec",
		"SourceType":      "spec",
		"SpecTemplateIds": "[2, 3]",
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
	templatePath := filepath.Join(dir, "app_blueprint_kubernetes_resource_spec.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath,
		args...,
	)
}
