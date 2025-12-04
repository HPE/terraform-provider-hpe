// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package template_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/testhelpers"
)

func TestAccMorpheusHpeMorpheusSpecTemplateCloudFormationResourceUrlExampleOk(t *testing.T) {
	t.Parallel()

	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()

	name := acctest.RandomWithPrefix(t.Name())

	resourceConfig, err := testhelpers.RenderExample(t, "hpe_morpheus_spec_template_cloud_formation_resource_url.tf.tmpl",
		"Name", name,
		"SourceType", "url",
		"SpecPath", "http://example.com/spec.yaml",
		"CapabilityIam", "true",
		"CapabilityNamedIam", "true",
		"CapabilityAutoExpand", "true",
	)
	if err != nil {
		t.Fatal(err)
	}

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"hpe_morpheus_spec_template_cloud_formation.tfexample_cloud_formation_spec_template_url",
			"name",
			name,
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_spec_template_cloud_formation.tfexample_cloud_formation_spec_template_url",
			"source_type",
			"url",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_spec_template_cloud_formation.tfexample_cloud_formation_spec_template_url",
			"spec_path",
			"http://example.com/spec.yaml",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_spec_template_cloud_formation.tfexample_cloud_formation_spec_template_url",
			"capability_iam",
			"true",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_spec_template_cloud_formation.tfexample_cloud_formation_spec_template_url",
			"capability_named_iam",
			"true",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_spec_template_cloud_formation.tfexample_cloud_formation_spec_template_url",
			"capability_auto_expand",
			"true",
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
