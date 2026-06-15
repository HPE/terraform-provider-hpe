// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package form_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus"
	sdkv2morpheus "github.com/HPE/terraform-provider-hpe/morpheus/sdkv2"
	"github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/resources/form"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
)

func TestAccMorpheusFormCloudOk(t *testing.T) {
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

	resourceConfig, err := form.RenderCloudConfig(t, map[string]string{
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
					resource.TestCheckResourceAttr("hpe_morpheus_form.example", "option_type.0.type", "cloud"),
					resource.TestCheckResourceAttr("hpe_morpheus_form.example", "option_type.0.code", optTypeCode),
					resource.TestCheckResourceAttr("hpe_morpheus_form.example", "option_type.0.cloud_type", "4"),
					resource.TestCheckResourceAttr("hpe_morpheus_form.example", "option_type.0.default_value", "test123"),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_form.example",
						"option_type.0.description",
						"Terraform cloud example",
					),
					resource.TestCheckResourceAttr("hpe_morpheus_form.example", "option_type.0.display_value_on_details", "true"),
					resource.TestCheckResourceAttr("hpe_morpheus_form.example", "option_type.0.exclude_from_search", "true"),
					resource.TestCheckResourceAttr("hpe_morpheus_form.example", "option_type.0.export_meta", "true"),
					resource.TestCheckResourceAttr("hpe_morpheus_form.example", "option_type.0.field_label", "cloud input"),
					resource.TestCheckResourceAttr("hpe_morpheus_form.example", "option_type.0.field_name", "cloudInput"),
					resource.TestCheckResourceAttr("hpe_morpheus_form.example", "option_type.0.filter_from_resource", "true"),
					resource.TestCheckResourceAttr("hpe_morpheus_form.example", "option_type.0.group_field_type", "value"),
					resource.TestCheckResourceAttr("hpe_morpheus_form.example", "option_type.0.group_id", "1"),
					resource.TestCheckResourceAttr("hpe_morpheus_form.example", "option_type.0.help_block", "Select a cloud"),
					resource.TestCheckResourceAttr("hpe_morpheus_form.example", "option_type.0.hidden", "false"),
					resource.TestCheckResourceAttr("hpe_morpheus_form.example", "option_type.0.instance_type_code", "apache"),
					resource.TestCheckResourceAttr("hpe_morpheus_form.example", "option_type.0.instance_type_field_type", "value"),
					resource.TestCheckResourceAttr("hpe_morpheus_form.example", "option_type.0.locked", "true"),
					resource.TestCheckResourceAttr("hpe_morpheus_form.example", "option_type.0.name", optTypeName),
					resource.TestCheckResourceAttr("hpe_morpheus_form.example", "option_type.0.placeholder", "Select cloud"),
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

// TestAccMorpheusFormCloudGroupCascadeOk verifies the group→cloud cascade:
// group_field on a cloud option type maps to config.group (not groupField),
// and round-trips cleanly without drift.
func TestAccMorpheusFormCloudGroupCascadeOk(t *testing.T) {
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

	// A form with a group field + cloud field using field-mode cascade.
	cascadeConfig := fmt.Sprintf(`
resource "hpe_morpheus_form" "cascade" {
  name        = %q
  code        = %q
  description = "cloud cascade test"

  option_type {
    name        = "Group Selector"
    code        = "%s-grp"
    type        = "group"
    field_label = "Group"
    field_name  = "fGroups"
  }

  option_type {
    name             = "Cloud Selector"
    code             = "%s-cld"
    type             = "cloud"
    field_label      = "Cloud"
    field_name       = "fClouds"
    group_field_type = "field"
    group_field      = "fGroups"
  }
}
`, name, code, code, code)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), sdkv2morpheus.Provider()),
		Steps: []resource.TestStep{
			// Apply
			{
				Config:             providerConfig + cascadeConfig,
				ExpectNonEmptyPlan: false,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hpe_morpheus_form.cascade", "option_type.0.type", "group"),
					resource.TestCheckResourceAttr("hpe_morpheus_form.cascade", "option_type.1.type", "cloud"),
					resource.TestCheckResourceAttr("hpe_morpheus_form.cascade", "option_type.1.group_field_type", "field"),
					resource.TestCheckResourceAttr("hpe_morpheus_form.cascade", "option_type.1.group_field", "fGroups"),
				),
			},
			// Plan after apply — no drift
			{
				Config:             providerConfig + cascadeConfig,
				ExpectNonEmptyPlan: false,
				PlanOnly:           true,
			},
		},
	})
}
