// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package optiontype_test

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

// RenderMorpheusOptionTypeTextConfig generates a Terraform configuration
// for the morpheus_option_type_text resource
func RenderMorpheusOptionTypeTextConfig(
	t *testing.T,
	overrides map[string]string,
) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Name":                  acctest.RandomWithPrefix(t.Name()),
		"Description":           "Terraform text option type example",
		"Labels":                "[\"demo\", \"terraform\"]",
		"FieldName":             "test1",
		"ExportMeta":            "true",
		"DependentField":        "dependent_example",
		"VisibilityField":       "visibility_example",
		"RequireField":          "require_example",
		"ShowOnEdit":            "true",
		"Editable":              "true",
		"DisplayValueOnDetails": "true",
		"FieldLabel":            "numbers",
		"Placeholder":           "fewf",
		"DefaultValue":          "testing",
		"HelpBlock":             "fiwefw",
		"Required":              "true",
		"VerifyPattern":         "a\\\\D{4}",
	}

	for key, value := range overrides {
		defaults[key] = value
	}

	return testhelpers.RenderExample(
		t,
		"morpheus_option_type_text_resource.tf.tmpl",
		"Name", defaults["Name"],
		"Description", defaults["Description"],
		"Labels", defaults["Labels"],
		"FieldName", defaults["FieldName"],
		"ExportMeta", defaults["ExportMeta"],
		"DependentField", defaults["DependentField"],
		"VisibilityField", defaults["VisibilityField"],
		"RequireField", defaults["RequireField"],
		"ShowOnEdit", defaults["ShowOnEdit"],
		"Editable", defaults["Editable"],
		"DisplayValueOnDetails", defaults["DisplayValueOnDetails"],
		"FieldLabel", defaults["FieldLabel"],
		"Placeholder", defaults["Placeholder"],
		"DefaultValue", defaults["DefaultValue"],
		"HelpBlock", defaults["HelpBlock"],
		"Required", defaults["Required"],
		"VerifyPattern", defaults["VerifyPattern"],
	)
}

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

func TestAccMorpheusOptionTypeTextExampleOk(t *testing.T) {
	t.Parallel()

	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()

	name := acctest.RandomWithPrefix(t.Name())

	resourceConfig, err := RenderMorpheusOptionTypeTextConfig(t, map[string]string{
		"Name": name,
	})
	if err != nil {
		t.Fatal(err)
	}

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"hpe_morpheus_option_type_text.tf_example_text_option_type",
			"name",
			name,
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_option_type_text.tf_example_text_option_type",
			"description",
			"Terraform text option type example",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_option_type_text.tf_example_text_option_type",
			"labels.0",
			"demo",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_option_type_text.tf_example_text_option_type",
			"labels.1",
			"terraform",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_option_type_text.tf_example_text_option_type",
			"field_name",
			"test1",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_option_type_text.tf_example_text_option_type",
			"export_meta",
			"true",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_option_type_text.tf_example_text_option_type",
			"dependent_field",
			"dependent_example",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_option_type_text.tf_example_text_option_type",
			"visibility_field",
			"visibility_example",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_option_type_text.tf_example_text_option_type",
			"require_field",
			"require_example",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_option_type_text.tf_example_text_option_type",
			"show_on_edit",
			"true",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_option_type_text.tf_example_text_option_type",
			"editable",
			"true",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_option_type_text.tf_example_text_option_type",
			"display_value_on_details",
			"true",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_option_type_text.tf_example_text_option_type",
			"field_label",
			"numbers",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_option_type_text.tf_example_text_option_type",
			"placeholder",
			"fewf",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_option_type_text.tf_example_text_option_type",
			"default_value",
			"testing",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_option_type_text.tf_example_text_option_type",
			"help_block",
			"fiwefw",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_option_type_text.tf_example_text_option_type",
			"required",
			"true",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_option_type_text.tf_example_text_option_type",
			"verify_pattern",
			"a\\D{4}",
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
				Config: providerConfig + resourceConfig,
				Check:  checkFn,
			},
			// Plan after apply
			{
				Config:             providerConfig + resourceConfig,
				ExpectNonEmptyPlan: false,
				PlanOnly:           true,
			},
		},
	})
}
