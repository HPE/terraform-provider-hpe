// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package integration_test

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

// RenderMorpheusIntegrationVroConfig generates a Terraform configuration for the VRO integration
// resource. It accepts an optional map of field overrides to customize the default values.
func RenderMorpheusIntegrationVroConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	// Default field values
	defaults := map[string]string{
		"Name":     acctest.RandomWithPrefix(t.Name()),
		"Enabled":  "true",
		"Url":      "https://myvro/vco/api",
		"Username": "my-vro-username",
		"Password": "my-vro-password",
		"AuthType": "basic",
		"Tenant":   "vsphere.local",
	}

	// Apply overrides to defaults
	for key, value := range overrides {
		defaults[key] = value
	}

	return testhelpers.RenderExample(
		t,
		"morpheus_integration_vro_resource.tf.tmpl",
		"Name", defaults["Name"],
		"Enabled", defaults["Enabled"],
		"Url", defaults["Url"],
		"Username", defaults["Username"],
		"Password", defaults["Password"],
		"AuthType", defaults["AuthType"],
		"Tenant", defaults["Tenant"],
	)
}

func TestAccMorpheusIntegrationVroExampleOk(t *testing.T) {
	t.Skip("Skipping due to lack of available resources to test against")

	t.Parallel()

	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()

	name := acctest.RandomWithPrefix(t.Name())

	resourceConfig, err := RenderMorpheusIntegrationVroConfig(t, map[string]string{
		"Name": name,
	})
	if err != nil {
		t.Fatal(err)
	}

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"hpe_morpheus_integration_vro.tf_example_vro_integration",
			"name",
			name,
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_integration_vro.tf_example_vro_integration",
			"enabled",
			"true",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_integration_vro.tf_example_vro_integration",
			"url",
			"https://myvro/vco/api",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_integration_vro.tf_example_vro_integration",
			"username",
			"my-vro-username",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_integration_vro.tf_example_vro_integration",
			"password",
			"my-vro-password",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_integration_vro.tf_example_vro_integration",
			"auth_type",
			"basic",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_integration_vro.tf_example_vro_integration",
			"tenant",
			"vsphere.local",
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
