// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package environment_test

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

// RenderMorpheusEnvironmentConfig renders the environment resource configuration with default values
// that can be overridden by providing a map of field name to value.
func RenderMorpheusEnvironmentConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Active":      "true",
		"Code":        "tfexample",
		"Description": "Terraform Example",
		"Name":        acctest.RandomWithPrefix(t.Name()),
	}

	for key, value := range overrides {
		defaults[key] = value
	}

	return testhelpers.RenderExample(t, "morpheus_environment_resource.tf.tmpl",
		"Active", defaults["Active"],
		"Code", defaults["Code"],
		"Description", defaults["Description"],
		"Name", defaults["Name"],
	)
}

func TestAccMorpheusEnvironmentExampleOk(t *testing.T) {
	t.Parallel()

	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()

	name := acctest.RandomWithPrefix(t.Name())

	resourceConfig, err := RenderMorpheusEnvironmentConfig(t, map[string]string{
		"Name": name,
	})
	if err != nil {
		t.Fatal(err)
	}

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"hpe_morpheus_environment.tf_example_environment",
			"active",
			"true",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_environment.tf_example_environment",
			"code",
			"tfexample",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_environment.tf_example_environment",
			"description",
			"Terraform Example",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_environment.tf_example_environment",
			"name",
			name,
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
