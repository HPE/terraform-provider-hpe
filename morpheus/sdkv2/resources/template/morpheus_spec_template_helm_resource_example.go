// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package template

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
)

//go:generate ../../../../../bin/render -out examples/resources/morpheus_spec_template_helm/resource_url.tf morpheus_spec_template_helm_resource_url.tf.tmpl Name tf-helm-spec-example-url SourceType url SpecPath http://example.com/chart.yaml
//go:generate ../../../../../bin/render -out examples/resources/morpheus_spec_template_helm/resource_local.tf morpheus_spec_template_helm_resource_local.tf.tmpl Name tf-helm-spec-example-local SourceType local SpecContent "apiVersion: v1\nkind: Service\nmetadata:\nname: {{ template \"fullname\" . }}\nlabels:\n    chart: \"{{ .Chart.Name }}-{{ .Chart.Version | replace \"+\" \"_\" }}\"\nspec:\ntype: {{ .Values.service.type }}\nports:\n- port: {{ .Values.service.externalPort }}\n    targetPort: {{ .Values.service.internalPort }}\n    protocol: TCP\n    name: {{ .Values.service.name }}\nselector:\n    app: {{ template \"fullname\" . }}"
//go:generate ../../../../../bin/render -out examples/resources/morpheus_spec_template_helm/resource_git.tf morpheus_spec_template_helm_resource_git.tf.tmpl Name tf-helm-spec-example-git SourceType repository RepositoryId 2 VersionRef main SpecPath ./spec.yaml

func RenderSpecTemplateHelmLocalConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaultSpecContent := "apiVersion: v1\nkind: Service\nmetadata:\nname: {{ template \"fullname\" . }}\n" +
		"labels:\n    chart: \"{{ .Chart.Name }}-{{ .Chart.Version | replace \"+\" \"_\" }}\"\nspec:\n" +
		"type: {{ .Values.service.type }}\nports:\n- port: {{ .Values.service.externalPort }}\n" +
		"    targetPort: {{ .Values.service.internalPort }}\n    protocol: TCP\n" +
		"    name: {{ .Values.service.name }}\nselector:\n    app: {{ template \"fullname\" . }}"

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
	templatePath := filepath.Join(dir, "morpheus_spec_template_helm_resource_local.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath,
		args...,
	)
}

func RenderSpecTemplateHelmUrlConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Name":       "Example",
		"SourceType": "url",
		"SpecPath":   "http://example.com/chart.yaml",
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
	templatePath := filepath.Join(dir, "morpheus_spec_template_helm_resource_url.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath,
		args...,
	)
}

func RenderSpecTemplateHelmGitConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Name":         "Example",
		"RepositoryId": "2",
		"SourceType":   "repository",
		"SpecPath":     "./spec.yaml",
		"VersionRef":   "main",
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
	templatePath := filepath.Join(dir, "morpheus_spec_template_helm_resource_git.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath,
		args...,
	)
}
