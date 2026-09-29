// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package blueprint_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	sdkv2morpheus "github.com/HPE/terraform-provider-hpe/morpheus/sdkv2"
	"github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/resources/blueprint"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
	"github.com/HPE/terraform-provider-hpe/provider/adapter"
)

func TestAccMorpheusInstanceTypeLayoutExampleOk(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.VMware)

	t.Parallel()

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()
	dependenciesConfig := testhelpers.WhoamiBlock()

	name := acctest.RandomWithPrefix(t.Name())

	// MORPH-16395 / MORPH-10326 / MORPH-11007 regression: exercise every ID-list
	// attribute of the layout with a known element count via self-contained
	// fixtures, then assert exact `.#` counts AND a clean post-apply plan (the
	// PlanOnly step below). Before the fixes these lists either never persisted
	// (vmware spec templates -- MORPH-16395) or read back zero-padded/doubled
	// (MORPH-10326 / MORPH-11007), producing a perpetual, never-converging diff.
	token := acctest.RandString(8)

	resourceConfig, err := blueprint.RenderInstanceTypeLayoutConfig(t, map[string]string{
		"Name":           name,
		"InstanceTypeId": "data.hpe_morpheus_instance_type.example.id",
		// exactly one element each, pointing at the in-config fixtures below.
		"SpecTemplateIds": "[hpe_morpheus_spec_template_terraform.regression.id]",
		"OptionTypeIds":   "[hpe_morpheus_option_type_text.regression.id]",
		"NodeTypeIds":     "[hpe_morpheus_node_type.regression.id]",
		"PriceSetIds":     "[hpe_morpheus_price_set.regression.id]",
	})
	if err != nil {
		t.Fatal(err)
	}

	// Self-contained fixtures + the system KVM instance type. Names are salted
	// with a per-run token so parallel runs do not collide on unique-name
	// constraints (see MORPH-16400).
	resourceConfig += fmt.Sprintf(`
	data "hpe_morpheus_instance_type" "example" {
	  name = "KVM"
	}

	resource "hpe_morpheus_spec_template_terraform" "regression" {
	  name         = "tf-layout-spec-%[1]s"
	  source_type  = "local"
	  spec_content = <<EOF
resource "local_file" "foo" {
  content  = "regression"
  filename = "/tmp/regression"
}
EOF
	}

	resource "hpe_morpheus_option_type_text" "regression" {
	  name        = "tf-layout-opt-%[1]s"
	  field_name  = "tf_layout_opt_%[1]s"
	  field_label = "tf layout opt %[1]s"
	}

	resource "hpe_morpheus_node_type" "regression" {
	  name             = "tf-layout-node-%[1]s"
	  short_name       = "tflayoutnode%[1]s"
	  technology       = "vmware"
	  version          = "1.0"
	  virtual_image_id = 10
	}

	resource "hpe_morpheus_price" "regression" {
	  name          = "tf-layout-price-%[1]s"
	  code          = "tf-layout-price-%[1]s"
	  tenant_id     = data.hpe_morpheus_whoami.current.tenant_id
	  price_type    = "fixed"
	  price_unit    = "minute"
	  incur_charges = "always"
	  currency      = "USD"
	  cost          = 1.00
	}

	resource "hpe_morpheus_price_set" "regression" {
	  name        = "tf-layout-priceset-%[1]s"
	  code        = "tf-layout-priceset-%[1]s"
	  region_code = "us-west-2"
	  price_unit  = "minute"
	  type        = "fixed"
	  price_ids   = [hpe_morpheus_price.regression.id]
	}
	`, token)

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttrSet(
			"hpe_morpheus_instance_type_layout.example",
			"instance_type_id",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_instance_type_layout.example",
			"labels.0",
			"demo",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_instance_type_layout.example",
			"labels.1",
			"layout",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_instance_type_layout.example",
			"labels.2",
			"terraform",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_instance_type_layout.example",
			"name",
			name,
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_instance_type_layout.example",
			"technology",
			"vmware",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_instance_type_layout.example",
			"version",
			"1.0",
		),
		// MORPH-16395 regression: a vmware layout must round-trip spec_template_ids.
		resource.TestCheckResourceAttr(
			"hpe_morpheus_instance_type_layout.example",
			"spec_template_ids.#",
			"1",
		),
		// MORPH-10326 regression: option_type_ids must read back exactly as
		// declared (no zero-pad, no doubling).
		resource.TestCheckResourceAttr(
			"hpe_morpheus_instance_type_layout.example",
			"option_type_ids.#",
			"1",
		),
		// MORPH-10326 / MORPH-11007 regression: node_type_ids must read back
		// exactly as declared and the post-apply plan must be empty (convergence).
		resource.TestCheckResourceAttr(
			"hpe_morpheus_instance_type_layout.example",
			"node_type_ids.#",
			"1",
		),
		// MORPH-10326 regression: price_set_ids must read back exactly as declared.
		resource.TestCheckResourceAttr(
			"hpe_morpheus_instance_type_layout.example",
			"price_set_ids.#",
			"1",
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), sdkv2morpheus.Provider()),
		Steps: []resource.TestStep{
			// Apply
			{
				Config:             providerConfig + dependenciesConfig + resourceConfig,
				ExpectNonEmptyPlan: false,
				Check:              checkFn,
			},
			// Plan after apply -- MORPH-11007 convergence: no perpetual diff.
			{
				Config:             providerConfig + dependenciesConfig + resourceConfig,
				ExpectNonEmptyPlan: false,
				PlanOnly:           true,
			},
		},
	})
}
