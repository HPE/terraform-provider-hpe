// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package cypher_test

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

// RenderHpeMorpheusCypherSecretResourceConfig generates a Terraform configuration
// for the hpe_morpheus_cypher_secret resource from the template file.
func RenderHpeMorpheusCypherSecretResourceConfig(
	t *testing.T,
	overrides map[string]string,
) (string, error) {
	t.Helper()

	// Default field values
	defaults := map[string]string{
		"Key":   "apipassword",
		"Value": "password123",
		"Ttl":   "86400",
	}

	// Apply overrides to defaults
	for key, value := range overrides {
		defaults[key] = value
	}

	// Build arguments for RenderExample
	args := []string{"hpe_morpheus_cypher_secret_resource.tf.tmpl"}
	for key, value := range defaults {
		args = append(args, key, value)
	}

	return testhelpers.RenderExample(t, args[0], args[1:]...)
}

func TestAccMorpheusCypherSecretExampleOk(t *testing.T) {
	t.Parallel()

	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()

	name := acctest.RandomWithPrefix(t.Name())

	resourceConfig, err := RenderHpeMorpheusCypherSecretResourceConfig(t, map[string]string{
		"Key": name,
	})
	if err != nil {
		t.Fatal(err)
	}

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"hpe_morpheus_cypher_secret.tf_example_cypher_secret",
			"key",
			name,
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_cypher_secret.tf_example_cypher_secret",
			"value",
			"password123",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_cypher_secret.tf_example_cypher_secret",
			"ttl",
			"86400",
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
