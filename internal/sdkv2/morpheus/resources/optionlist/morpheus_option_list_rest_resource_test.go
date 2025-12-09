// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package optionlist_test

import (
	"context"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-mux/tf5to6server"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/testhelpers"
	sdkv2morpheus "github.com/HPE/terraform-provider-hpe/internal/sdkv2/morpheus"
)

func TestMain(m *testing.M) {
	code := m.Run()

	testhelpers.WriteMergedResults()

	os.Exit(code)
}

func newProviderWithError() (tfprotov6.ProviderServer, error) {
	return tf5to6server.UpgradeServer(context.Background(), sdkv2morpheus.Provider().GRPCProvider)
}

var testAccProtoV6ProviderFactories = map[string]func() (
	tfprotov6.ProviderServer, error,
){
	"hpe": newProviderWithError,
}

// RenderMorpheusOptionListRestConfig generates a Terraform configuration for the morpheus_option_list_rest resource.
// It accepts an optional map of field overrides to customize the default values.
// Supported override keys: "Name", "Description", "Visibility", "SourceUrl", "RealTime", "IgnoreSslErrors",
// "SourceMethod", "InitialDataset", "TranslationScript", "SourceHeaderName1", "SourceHeaderValue1",
// "SourceHeaderName2", "SourceHeaderValue2"
func RenderMorpheusOptionListRestConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Name":            acctest.RandomWithPrefix(t.Name()),
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
	return testhelpers.RenderExample(t, "morpheus_option_list_rest_resource.tf.tmpl",
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

func TestAccMorpheusOptionListRestExampleOk(t *testing.T) {
	t.Parallel()

	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()

	name := acctest.RandomWithPrefix(t.Name())

	resourceConfig, err := RenderMorpheusOptionListRestConfig(t, map[string]string{"Name": name})
	if err != nil {
		t.Fatal(err)
	}

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"hpe_morpheus_option_list_rest.tf_example_rest_option_list",
			"name",
			name,
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_option_list_rest.tf_example_rest_option_list",
			"description",
			"Terraform REST option list example",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_option_list_rest.tf_example_rest_option_list",
			"visibility",
			"private",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_option_list_rest.tf_example_rest_option_list",
			"source_url",
			"https://api.github.com/repos/hashicorp/consul/releases",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_option_list_rest.tf_example_rest_option_list",
			"real_time",
			"true",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_option_list_rest.tf_example_rest_option_list",
			"ignore_ssl_errors",
			"true",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_option_list_rest.tf_example_rest_option_list",
			"source_method",
			"GET",
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Plan
			{
				Config:             providerConfig + resourceConfig,
				ExpectNonEmptyPlan: true,
				Check:              checkFn,
				PlanOnly:           true,
			},
			// Apply
			{
				Config:             providerConfig + resourceConfig,
				Check:              checkFn,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}
