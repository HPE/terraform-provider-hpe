// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package optionlist

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/testhelpers"
)

//go:generate go run ../../../../../cmd/render -out examples/resources/morpheus_option_list_rest/resource.tf morpheus_option_list_rest_resource.tf.tmpl Name tf_example_rest_option_list Description "Terraform REST option list example" Visibility private SourceUrl https://api.github.com/repos/hashicorp/consul/releases RealTime true IgnoreSslErrors true SourceMethod GET InitialDataset "  [{\"name\": \"Level 1\",\"value\":\"level1\"},\n  {\"name\": \"Level 2\",\"value\":\"level2\"},\n  {\"name\": \"Level 3\",\"value\":\"level3\"}\n  ]" TranslationScript "      for(var x=0;x < 5; x++) {\n          results.push({name: data[x].name,value:data[x].name});\n        }" SourceHeaderName1 Accept SourceHeaderValue1 application/json SourceHeaderName2 Authorization SourceHeaderValue2 "Basic YWRtaW46YWRtaW4="

// RenderOptionListRestConfig generates a Terraform configuration for the morpheus_option_list_rest resource.
// It accepts an optional map of field overrides to customize the default values.
// Supported override keys: "Name", "Description", "Visibility", "SourceUrl", "RealTime", "IgnoreSslErrors",
// "SourceMethod", "InitialDataset", "TranslationScript", "SourceHeaderName1", "SourceHeaderValue1",
// "SourceHeaderName2", "SourceHeaderValue2"
func RenderOptionListRestConfig(
	t *testing.T,
	name string,
	overrides map[string]string,
) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Name":            name,
		"Description":     "Terraform REST option list example",
		"Visibility":      "private",
		"SourceUrl":       "https://api.github.com/repos/hashicorp/consul/releases",
		"RealTime":        "true",
		"IgnoreSslErrors": "true",
		"SourceMethod":    "GET",
		"InitialDataset": "  [{\"name\": \"Level 1\",\"value\":\"level1\"},\n" +
			"  {\"name\": \"Level 2\",\"value\":\"level2\"},\n" +
			"  {\"name\": \"Level 3\",\"value\":\"level3\"}\n  ]",
		"TranslationScript": "      for(var x=0;x < 5; x++) {\n" +
			"          results.push({name: data[x].name,value:data[x].name});\n" +
			"        }",
		"SourceHeaderName1":  "Accept",
		"SourceHeaderValue1": "application/json",
		"SourceHeaderName2":  "Authorization",
		"SourceHeaderValue2": "Basic YWRtaW46YWRtaW4=",
	}

	// Apply overrides
	for key, value := range overrides {
		defaults[key] = value
	}

	//nolint: lll
	// Get the directory where this source file is located
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("unable to get current file path")
	}
	dir := filepath.Dir(filename)
	templatePath := filepath.Join(dir, "morpheus_option_list_rest_resource.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath,
		"Name", defaults["Name"],
		"Description", defaults["Description"],
		"Visibility", defaults["Visibility"],
		"SourceUrl", defaults["SourceUrl"],
		"RealTime", defaults["RealTime"],
		"IgnoreSslErrors", defaults["IgnoreSslErrors"],
		"SourceMethod", defaults["SourceMethod"],
		"InitialDataset", defaults["InitialDataset"],
		"TranslationScript", defaults["TranslationScript"],
		"SourceHeaderName1", defaults["SourceHeaderName1"],
		"SourceHeaderValue1", defaults["SourceHeaderValue1"],
		"SourceHeaderName2", defaults["SourceHeaderName2"],
		"SourceHeaderValue2", defaults["SourceHeaderValue2"],
	)
}
