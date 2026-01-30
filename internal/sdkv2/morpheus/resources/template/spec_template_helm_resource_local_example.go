// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package template

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/testhelpers"
)

//go:generate sh -c " ../../../../../bin/render -out examples/resources/morpheus_spec_template_helm/resource_local.tf spec_template_helm_resource_local.tf.tmpl Name 'tf-helm-spec-example-local' SourceType 'local' SpecContent '<<TFEOF\napiVersion: v1\nkind: Service\nmetadata:\nname: {{ template \"fullname\" . }}\nlabels:\n chart: \"{{ .Chart.Name }}-{{ .Chart.Version | replace \"+\" \"_\" }}\"\nspec:\ntype: {{ .Values.service.type }}\nports:\n- port: {{ .Values.service.externalPort }}\n targetPort: {{ .Values.service.internalPort }}\n protocol: TCP\n name: {{ .Values.service.name }}\nselector:\n app: {{ template \"fullname\" . }}\nTFEOF'"

// RenderSpecTemplateHelmConfigLocal generates a Terraform configuration for the tenant resource.
// It accepts optional overrides for field values. Default values are used if not overridden.
func RenderSpecTemplateHelmConfigLocal(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Name":       "tf-helm-spec-example-local",
		"SourceType": "local",
		"SpecContent": "<<TFEOF\napiVersion: v1\nkind: Service\nmetadata:\nname: {{ template \"fullname\" . }}\nlabels:\n " +
			"chart: \"{{ .Chart.Name }}-{{ .Chart.Version | replace \"+\" \"_\" }}\"\nspec:\n" +
			"type: {{ .Values.service.type }}\nports:\n- port: {{ .Values.service.externalPort }}\n " +
			"targetPort: {{ .Values.service.internalPort }}\n protocol: TCP\n " +
			"name: {{ .Values.service.name }}\nselector:\n app: {{ template \"fullname\" . }}\nTFEOF",
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
	templatePath := filepath.Join(dir, "spec_template_helm_resource_local.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath,
		args...,
	)
}
