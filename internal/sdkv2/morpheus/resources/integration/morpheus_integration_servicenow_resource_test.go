// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package integration_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/testhelpers"
)

func RenderIntegrationServicenowConfig(
	t *testing.T,
	name string,
	overrides map[string]string,
) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Name":                     name,
		"Enabled":                  "true",
		"Url":                      "https://servicenowprod.service-now.com",
		"Username":                 "my-snow-username",
		"Password":                 "my-snow-password",
		"DefaultCmdbBusinessClass": "demo",
	}

	for key, value := range overrides {
		defaults[key] = value
	}

	return testhelpers.RenderExample(
		t,
		"morpheus_integration_servicenow_resource.tf.tmpl",
		"Name", defaults["Name"],
		"Enabled", defaults["Enabled"],
		"Url", defaults["Url"],
		"Username", defaults["Username"],
		"Password", defaults["Password"],
		"DefaultCmdbBusinessClass", defaults["DefaultCmdbBusinessClass"],
	)
}

func TestAccMorpheusIntegrationServicenowExampleOk(t *testing.T) {
	t.Skip("Skipping due to lack of available resources to test against")

	t.Parallel()

	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()

	name := acctest.RandomWithPrefix(t.Name())

	resourceConfig, err := RenderIntegrationServicenowConfig(t, name, nil)
	if err != nil {
		t.Fatal(err)
	}

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"hpe_morpheus_integration_servicenow.tf_example_servicenow_integration",
			"name",
			name,
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_integration_servicenow.tf_example_servicenow_integration",
			"enabled",
			"true",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_integration_servicenow.tf_example_servicenow_integration",
			"url",
			"https://servicenowprod.service-now.com",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_integration_servicenow.tf_example_servicenow_integration",
			"username",
			"my-snow-username",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_integration_servicenow.tf_example_servicenow_integration",
			"password",
			"my-snow-password",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_integration_servicenow.tf_example_servicenow_integration",
			"default_cmdb_business_class",
			"demo",
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
