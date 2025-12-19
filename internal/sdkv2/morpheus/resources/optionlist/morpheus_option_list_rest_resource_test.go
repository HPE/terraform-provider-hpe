// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package optionlist_test

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
	"github.com/HPE/terraform-provider-hpe/internal/sdkv2/morpheus/resources/optionlist"
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

func TestAccMorpheusOptionListRestExampleOk(t *testing.T) {
	t.Parallel()

	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()

	name := acctest.RandomWithPrefix(t.Name())

	resourceConfig, err := optionlist.RenderOptionListRestConfig(t, map[string]string{
		"Name": name,
	})
	if err != nil {
		t.Fatal(err)
	}

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"hpe_morpheus_option_list_rest.tf_example_rest_option_list",
			"name",
			name,
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_option_list_rest.tf_example_rest_option_list",
			"description",
			"Terraform REST option list example",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_option_list_rest.tf_example_rest_option_list",
			"visibility",
			"private",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_option_list_rest.tf_example_rest_option_list",
			"source_url",
			"https://api.github.com/repos/hashicorp/consul/releases",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_option_list_rest.tf_example_rest_option_list",
			"real_time",
			"true",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_option_list_rest.tf_example_rest_option_list",
			"ignore_ssl_errors",
			"true",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_option_list_rest.tf_example_rest_option_list",
			"source_method",
			"GET",
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
				Config:             providerConfig + resourceConfig,
				Check:              checkFn,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}
