// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package optionlist

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/testhelpers"
)

//go:generate go run ../../../../../cmd/render -out examples/resources/morpheus_option_list_api/resource.tf morpheus_option_list_api_resource.tf.tmpl Name tf_example_api_option_list Description "Terraform Morpheus API option list example" Visibility private OptionList instances TranslationScript "var i=0;\nresults = [];\nfor(i; i<data.length; i++) {\n  results.push({name: data[i].name, value: data[i].name});\n}"

func RenderOptionListApiConfig(
	t *testing.T,
	name string,
	overrides map[string]string,
) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Name":        name,
		"Description": "Terraform Morpheus API option list example",
		"Visibility":  "private",
		"OptionList":  "instances",
		"TranslationScript": "var i=0;\nresults = [];\nfor(i; i<data.length; i++) {" +
			"\n  results.push({name: data[i].name, value: data[i].name});\n}",
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
	templatePath := filepath.Join(dir, "morpheus_option_list_api_resource.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath,
		"Name", defaults["Name"],
		"Description", defaults["Description"],
		"Visibility", defaults["Visibility"],
		"OptionList", defaults["OptionList"],
		"TranslationScript", defaults["TranslationScript"],
	)
}
