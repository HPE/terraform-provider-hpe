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
)

func TestAccMorpheusFormCloudOk(t *testing.T) {
	t.Parallel()

	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()
	name := acctest.RandomWithPrefix(t.Name())

	resourceConfig, err := form.RenderCloudConfig(t, map[string]string{"Name": name})
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
					resource.TestCheckResourceAttr("hpe_morpheus_form.example", "code", "demo"),
					resource.TestCheckResourceAttr("hpe_morpheus_form.example", "description", "demo"),
					resource.TestCheckResourceAttr("hpe_morpheus_form.example", "option_type1_type", "cloud"),
					resource.TestCheckResourceAttr("hpe_morpheus_form.example", "option_type1_code", "cloud-input"),
					resource.TestCheckResourceAttr("hpe_morpheus_form.example", "option_type1_cloud_type", "4"),
					resource.TestCheckResourceAttr("hpe_morpheus_form.example", "option_type1_default_value", "test123"),
					resource.TestCheckResourceAttr("hpe_morpheus_form.example", "option_type1_description", "Terraform cloud example"),
					resource.TestCheckResourceAttr("hpe_morpheus_form.example", "option_type1_display_value_on_details", "true"),
					resource.TestCheckResourceAttr("hpe_morpheus_form.example", "option_type1_exclude_from_search", "true"),
					resource.TestCheckResourceAttr("hpe_morpheus_form.example", "option_type1_export_meta", "true"),
					resource.TestCheckResourceAttr("hpe_morpheus_form.example", "option_type1_field_label", "cloud input"),
					resource.TestCheckResourceAttr("hpe_morpheus_form.example", "option_type1_field_name", "cloudInput"),
					resource.TestCheckResourceAttr("hpe_morpheus_form.example", "option_type1_filter_from_resource", "true"),
					resource.TestCheckResourceAttr("hpe_morpheus_form.example", "option_type1_group_field_type", "value"),
					resource.TestCheckResourceAttr("hpe_morpheus_form.example", "option_type1_group_id", "1"),
					resource.TestCheckResourceAttr("hpe_morpheus_form.example", "option_type1_help_block", "Select a cloud"),
					resource.TestCheckResourceAttr("hpe_morpheus_form.example", "option_type1_hidden", "false"),
					resource.TestCheckResourceAttr("hpe_morpheus_form.example", "option_type1_instance_type_code", "apache"),
					resource.TestCheckResourceAttr("hpe_morpheus_form.example", "option_type1_instance_type_field_type", "value"),
					resource.TestCheckResourceAttr("hpe_morpheus_form.example", "option_type1_locked", "true"),
					resource.TestCheckResourceAttr("hpe_morpheus_form.example", "option_type1_name", "tf cloud example"),
					resource.TestCheckResourceAttr("hpe_morpheus_form.example", "option_type1_placeholder", "Select cloud"),
					resource.TestCheckResourceAttr("hpe_morpheus_form.example", "option_type1_required", "true"),
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
