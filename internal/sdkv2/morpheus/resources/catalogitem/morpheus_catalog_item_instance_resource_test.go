// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package catalogitem_test

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

// RenderCatalogItemInstanceConfig generates a Terraform configuration for catalog item instance resource.
// It accepts a map of field overrides to customize the default values.
func RenderCatalogItemInstanceConfig(
	t *testing.T,
	overrides map[string]string,
) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Name":        "tfexample_instance_catalog",
		"Config":      "{\"name\":\"test\"}",
		"Content":     "{\"name\":\"test\"}",
		"Description": "terraform example instance catalog item",
		"Enabled":     "true",
		"Featured":    "true",
		"ImageName":   "tfexample.png",
		"ImagePath":   "tfexample.png",
		"Visibility":  "private",
	}

	for key, value := range overrides {
		defaults[key] = value
	}

	return testhelpers.RenderExample(
		t,
		"morpheus_catalog_item_instance_resource.tf.tmpl",
		"Name", defaults["Name"],
		"Config", defaults["Config"],
		"Content", defaults["Content"],
		"Description", defaults["Description"],
		"Enabled", defaults["Enabled"],
		"Featured", defaults["Featured"],
		"ImageName", defaults["ImageName"],
		"ImagePath", defaults["ImagePath"],
		"Visibility", defaults["Visibility"],
	)
}

func TestAccMorpheusCatalogItemInstanceExampleOk(t *testing.T) {
	t.Parallel()

	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()

	name := acctest.RandomWithPrefix(t.Name())

	resourceConfig, err := RenderCatalogItemInstanceConfig(t, map[string]string{
		"Name": name,
	})
	if err != nil {
		t.Fatal(err)
	}

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"hpe_morpheus_catalog_item_instance.tf_example_instance_catalog_item",
			"name",
			name,
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_catalog_item_instance.tf_example_instance_catalog_item",
			"description",
			"terraform example instance catalog item",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_catalog_item_instance.tf_example_instance_catalog_item",
			"image_path",
			"tfexample.png",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_catalog_item_instance.tf_example_instance_catalog_item",
			"image_name",
			"tfexample.png",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_catalog_item_instance.tf_example_instance_catalog_item",
			"enabled",
			"true",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_catalog_item_instance.tf_example_instance_catalog_item",
			"featured",
			"true",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_catalog_item_instance.tf_example_instance_catalog_item",
			"content",
			"{\"name\":\"test\"}",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_catalog_item_instance.tf_example_instance_catalog_item",
			"config",
			"{\"name\":\"test\"}",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_catalog_item_instance.tf_example_instance_catalog_item",
			"visibility",
			"private",
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
				Check:              checkFn,
				PlanOnly:           true,
			},
		},
	})
}
