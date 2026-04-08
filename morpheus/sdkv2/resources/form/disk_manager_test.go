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

func TestAccMorpheusFormDiskManagerOk(t *testing.T) {
	t.Parallel()

	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()
	name := acctest.RandomWithPrefix(t.Name())

	resourceConfig, err := form.RenderDiskManagerConfig(t, map[string]string{"Name": name})
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
					resource.TestCheckResourceAttr("hpe_morpheus_form.example", "option_type1_type", "diskManager"),
					resource.TestCheckResourceAttr("hpe_morpheus_form.example", "option_type1_cloud_field_type", "value"),
					resource.TestCheckResourceAttr("hpe_morpheus_form.example", "option_type1_cloud_id", "1"),
					resource.TestCheckResourceAttr("hpe_morpheus_form.example", "option_type1_code", "disk-manager-input"),
					resource.TestCheckResourceAttr("hpe_morpheus_form.example", "option_type1_description", "Terraform disk manager example"),
					resource.TestCheckResourceAttr("hpe_morpheus_form.example", "option_type1_display_value_on_details", "true"),
					resource.TestCheckResourceAttr("hpe_morpheus_form.example", "option_type1_enable_datastore_selection", "true"),
					resource.TestCheckResourceAttr("hpe_morpheus_form.example", "option_type1_enable_disk_type_selection", "true"),
					resource.TestCheckResourceAttr("hpe_morpheus_form.example", "option_type1_enable_storage_type_selection", "true"),
					resource.TestCheckResourceAttr("hpe_morpheus_form.example", "option_type1_exclude_from_search", "true"),
					resource.TestCheckResourceAttr("hpe_morpheus_form.example", "option_type1_export_meta", "true"),
					resource.TestCheckResourceAttr("hpe_morpheus_form.example", "option_type1_field_label", "disk manager input"),
					resource.TestCheckResourceAttr("hpe_morpheus_form.example", "option_type1_field_name", "diskManagerInput"),
					resource.TestCheckResourceAttr("hpe_morpheus_form.example", "option_type1_group_field_type", "value"),
					resource.TestCheckResourceAttr("hpe_morpheus_form.example", "option_type1_group_id", "1"),
					resource.TestCheckResourceAttr("hpe_morpheus_form.example", "option_type1_help_block", "Configure disks"),
					resource.TestCheckResourceAttr("hpe_morpheus_form.example", "option_type1_hidden", "false"),
					resource.TestCheckResourceAttr("hpe_morpheus_form.example", "option_type1_image_id", "1"),
					resource.TestCheckResourceAttr("hpe_morpheus_form.example", "option_type1_layout_field_type", "value"),
					resource.TestCheckResourceAttr("hpe_morpheus_form.example", "option_type1_layout_id", "1"),
					resource.TestCheckResourceAttr("hpe_morpheus_form.example", "option_type1_locked", "true"),
					resource.TestCheckResourceAttr("hpe_morpheus_form.example", "option_type1_name", "tf disk manager example"),
					resource.TestCheckResourceAttr("hpe_morpheus_form.example", "option_type1_plan_field_type", "value"),
					resource.TestCheckResourceAttr("hpe_morpheus_form.example", "option_type1_plan_id", "1"),
					resource.TestCheckResourceAttr("hpe_morpheus_form.example", "option_type1_pool_field_type", "value"),
					resource.TestCheckResourceAttr("hpe_morpheus_form.example", "option_type1_pool_id", "1"),
					resource.TestCheckResourceAttr("hpe_morpheus_form.example", "option_type1_required", "true"),
					resource.TestCheckResourceAttr("hpe_morpheus_form.example", "option_type1_virtual_image_field_type", "value"),
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
