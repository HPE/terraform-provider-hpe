// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package contact_test

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

// RenderContactConfig renders the contact resource configuration with
// optional field overrides
func RenderContactConfig(
	t *testing.T,
	name string,
	overrides map[string]string,
) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Name":         name,
		"EmailAddress": "tfcontact@demo.com",
		"MobileNumber": "123-456-7890",
	}

	for key, value := range overrides {
		defaults[key] = value
	}

	return testhelpers.RenderExample(
		t,
		"contact_resource.tf.tmpl",
		"Name", defaults["Name"],
		"EmailAddress", defaults["EmailAddress"],
		"MobileNumber", defaults["MobileNumber"],
	)
}

func TestAccMorpheusContactExampleOk(t *testing.T) {
	t.Parallel()

	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()

	name := acctest.RandomWithPrefix(t.Name())

	resourceConfig, err := RenderContactConfig(t, name, map[string]string{})
	if err != nil {
		t.Fatal(err)
	}

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"hpe_morpheus_contact.tf_example_contact",
			"name",
			name,
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_contact.tf_example_contact",
			"email_address",
			"tfcontact@demo.com",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_contact.tf_example_contact",
			"mobile_number",
			"123-456-7890",
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
