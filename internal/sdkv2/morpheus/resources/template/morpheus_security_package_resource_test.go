// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package template_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/testhelpers"
)

// RenderMorpheusSecurityPackageConfig generates a test configuration for security package resource.
// It accepts a name and a map of field overrides to customize the default values.
func RenderMorpheusSecurityPackageConfig(t *testing.T, name string, overrides map[string]string) string {
	t.Helper()

	defaults := map[string]string{
		"Name":        name,
		"Description": "Terraform security package example",
		"Labels":      "[\"demo\", \"terraform\"]",
		"Enabled":     "true",
		"Url": "https://github.com/ComplianceAsCode/content/releases/download/v0.1.59/" +
			"scap-security-guide-0.1.59.zip",
	}

	for key, value := range overrides {
		defaults[key] = value
	}

	resourceConfig, err := testhelpers.RenderExample(
		t,
		"morpheus_security_package_resource.tf.tmpl",
		"Name", defaults["Name"],
		"Description", defaults["Description"],
		"Labels", defaults["Labels"],
		"Enabled", defaults["Enabled"],
		"Url", defaults["Url"],
	)
	if err != nil {
		t.Fatal(err)
	}

	return resourceConfig
}

func TestAccMorpheusSecurityPackageExampleOk(t *testing.T) {
	t.Parallel()

	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()

	name := acctest.RandomWithPrefix(t.Name())

	resourceConfig := RenderMorpheusSecurityPackageConfig(t, name, map[string]string{})

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"hpe_morpheus_security_package.tf_example_security_package",
			"name",
			name,
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_security_package.tf_example_security_package",
			"description",
			"Terraform security package example",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_security_package.tf_example_security_package",
			"labels.#",
			"2",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_security_package.tf_example_security_package",
			"enabled",
			"true",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_security_package.tf_example_security_package",
			"url",
			"https://github.com/ComplianceAsCode/content/releases/download/v0.1.59/"+
				"scap-security-guide-0.1.59.zip",
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
