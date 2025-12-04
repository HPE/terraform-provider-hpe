// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package blueprint_test

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

func TestAccMorpheusHpeMorpheusInstanceTypeExampleOk(t *testing.T) {
	t.Parallel()

	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()

	name := acctest.RandomWithPrefix(t.Name())

	resourceConfig, err := testhelpers.RenderExample(t, "hpe_morpheus_instance_type_resource.tf.tmpl",
		"Name", name,
		"Code", "tf_example_instance",
		"Description", "Terraform Example Instance Type",
		"Labels", `["demo", "instance", "terraform"]`,
		"Category", "web",
		"Visibility", "private",
		"ImagePath", "tfexample.png",
		"ImageName", "tfexample.png",
		"Featured", "false",
		"EnableDeployments", "true",
		"EnableScaling", "true",
		"EnableSettings", "true",
		"EnvironmentPrefix", "TFEXAMPLE_DEMO",
		"OptionTypeIds", "[1910, 1912]",
		"Evar1Name", "first",
		"Evar1Value", "first",
		"Evar1Export", "true",
		"Evar2Name", "second",
		"Evar2MaskedValue", "second",
		"Evar2Export", "false",
	)
	if err != nil {
		t.Fatal(err)
	}

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"hpe_morpheus_instance_type.tf_example_instance_type",
			"name",
			name,
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_instance_type.tf_example_instance_type",
			"code",
			"tf_example_instance",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_instance_type.tf_example_instance_type",
			"description",
			"Terraform Example Instance Type",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_instance_type.tf_example_instance_type",
			"category",
			"web",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_instance_type.tf_example_instance_type",
			"visibility",
			"private",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_instance_type.tf_example_instance_type",
			"image_path",
			"tfexample.png",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_instance_type.tf_example_instance_type",
			"image_name",
			"tfexample.png",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_instance_type.tf_example_instance_type",
			"featured",
			"false",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_instance_type.tf_example_instance_type",
			"enable_deployments",
			"true",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_instance_type.tf_example_instance_type",
			"enable_scaling",
			"true",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_instance_type.tf_example_instance_type",
			"enable_settings",
			"true",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_instance_type.tf_example_instance_type",
			"environment_prefix",
			"TFEXAMPLE_DEMO",
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
