// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package form_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus"
	sdkv2morpheus "github.com/HPE/terraform-provider-hpe/morpheus/sdkv2"
	"github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/resources/form"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
)

func TestAccMorpheusFormVmwFoldersOk(t *testing.T) {
	if capabilities.Missing(t, capabilities.All) {
		t.Log("Skipping test due to missing capabilities")

		return
	}
	t.Parallel()

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	defer testhelpers.RecordResult(t)
	providerConfig := testhelpers.ProviderBlock()
	name := acctest.RandomWithPrefix(t.Name())
	code := toCode(name)
	optTypeCode := code + "-ot"
	optTypeName := name + " option type"

	resourceConfig, err := form.RenderVmwFoldersConfig(t, map[string]string{
		"Name":           name,
		"Code":           code,
		"OptionTypeCode": optTypeCode,
		"OptionTypeName": optTypeName,
	})
	if err != nil {
		t.Fatal(err)
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), sdkv2morpheus.Provider()),
		Steps: []resource.TestStep{
			{
				Config:             providerConfig + resourceConfig,
				ExpectNonEmptyPlan: false,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hpe_morpheus_form.example", "code", code),
					resource.TestCheckResourceAttr("hpe_morpheus_form.example", "description", "demo"),
					resource.TestCheckResourceAttr("hpe_morpheus_form.example", "option_type.0.type", "vmwFolders"),
					resource.TestCheckResourceAttr("hpe_morpheus_form.example", "option_type.0.code", optTypeCode),
					resource.TestCheckResourceAttr("hpe_morpheus_form.example", "option_type.0.cloud_field_type", "value"),
					resource.TestCheckResourceAttr("hpe_morpheus_form.example", "option_type.0.cloud_id", "1"),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_form.example",
						"option_type.0.description",
						"Terraform vmwFolders example",
					),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_form.example",
						"option_type.0.display_value_on_details",
						"true",
					),
					resource.TestCheckResourceAttr("hpe_morpheus_form.example", "option_type.0.exclude_from_search", "true"),
					resource.TestCheckResourceAttr("hpe_morpheus_form.example", "option_type.0.export_meta", "true"),
					resource.TestCheckResourceAttr("hpe_morpheus_form.example", "option_type.0.field_label", "VmwFolders"),
					resource.TestCheckResourceAttr("hpe_morpheus_form.example", "option_type.0.field_name", "vmwFolders"),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_form.example",
						"option_type.0.group_field_type",
						"value",
					),
					resource.TestCheckResourceAttr("hpe_morpheus_form.example", "option_type.0.group_id", "1"),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_form.example",
						"option_type.0.help_block",
						"Select a vmwFolder",
					),
					resource.TestCheckResourceAttr("hpe_morpheus_form.example", "option_type.0.hidden", "false"),
					resource.TestCheckResourceAttr("hpe_morpheus_form.example", "option_type.0.locked", "true"),
					resource.TestCheckResourceAttr("hpe_morpheus_form.example", "option_type.0.name", optTypeName),
					resource.TestCheckResourceAttr("hpe_morpheus_form.example", "option_type.0.plan_field_type", "value"),
					resource.TestCheckResourceAttr("hpe_morpheus_form.example", "option_type.0.plan_id", "1"),
					resource.TestCheckResourceAttr("hpe_morpheus_form.example", "option_type.0.required", "true"),
				),
			},
			{
				Config:             providerConfig + resourceConfig,
				ExpectNonEmptyPlan: false,
				PlanOnly:           true,
			},
		},
	})
}
