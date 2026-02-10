// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package template_test

import (
	"context"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-mux/tf5to6server"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	sdkv2morpheus "github.com/HPE/terraform-provider-hpe/morpheus/sdkv2"
	"github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/resources/template"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/systemoverride"
)

func TestMain(m *testing.M) {
	systemoverride.ParseFlags()
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

func TestAccMorpheusClusterPackageExampleOk(t *testing.T) {
	t.Parallel()

	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	testSystem := systemoverride.GetPreferred(t, "feature")
	providerConfig := testhelpers.ProviderBlock(testSystem)

	name := acctest.RandomWithPrefix(t.Name())
	resourceConfig, err := template.RenderClusterPackageConfig(t, map[string]string{
		"Name": name,
	})
	if err != nil {
		t.Fatal(err)
	}

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"hpe_morpheus_cluster_package.tf_example_cluster_package",
			"name",
			name,
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_cluster_package.tf_example_cluster_package",
			"code",
			"tf-example-cluster-package",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_cluster_package.tf_example_cluster_package",
			"description",
			"Terraform example cluster package",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_cluster_package.tf_example_cluster_package",
			"package_version",
			"1.2.3",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_cluster_package.tf_example_cluster_package",
			"type",
			"apps",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_cluster_package.tf_example_cluster_package",
			"package_type",
			"example",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_cluster_package.tf_example_cluster_package",
			"enabled",
			"true",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_cluster_package.tf_example_cluster_package",
			"repeat_install",
			"true",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_cluster_package.tf_example_cluster_package",
			"spec_template_ids.#",
			"2",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_cluster_package.tf_example_cluster_package",
			"spec_template_ids.0",
			"1",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_cluster_package.tf_example_cluster_package",
			"spec_template_ids.1",
			"2",
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
