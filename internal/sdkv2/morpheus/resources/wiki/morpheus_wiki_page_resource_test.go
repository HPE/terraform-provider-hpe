// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package wiki_test

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

// RenderMorpheusWikiPageConfig generates a Terraform configuration for the wiki page resource.
// It accepts a map of field overrides to customize default values.
func RenderMorpheusWikiPageConfig(t *testing.T, overrides map[string]string) string {
	t.Helper()

	// Default field values
	defaults := map[string]string{
		"Name":     acctest.RandomWithPrefix(t.Name()),
		"Category": "morpheus-terraform",
	}

	// Apply overrides
	for key, value := range overrides {
		defaults[key] = value
	}

	resourceConfig, err := testhelpers.RenderExample(t, "morpheus_wiki_page_resource.tf.tmpl",
		"Name", defaults["Name"],
		"Category", defaults["Category"],
	)
	if err != nil {
		t.Fatal(err)
	}

	return resourceConfig
}

func TestAccMorpheusWikiPageExampleOk(t *testing.T) {
	t.Parallel()

	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()

	name := acctest.RandomWithPrefix(t.Name())

	resourceConfig := RenderMorpheusWikiPageConfig(t, map[string]string{
		"Name": name,
	})

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"hpe_morpheus_wiki_page.tfexample_wiki_page",
			"name",
			name,
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_wiki_page.tfexample_wiki_page",
			"category",
			"morpheus-terraform",
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
