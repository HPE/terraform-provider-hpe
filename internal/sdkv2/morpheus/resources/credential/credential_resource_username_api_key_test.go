// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package credential_test

import (
	"testing"

	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/testhelpers"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func RenderCredentialUsernameApiKeyConfig(
	t *testing.T,
	overrides map[string]string,
) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Description": "terraform credential example for username api key",
		"Enabled":     "true",
		"Type":        "username-api-key",
		"Username":    "admin",
		"ApiKey":      "MFIEIWEIFINEF",
	}

	for key, value := range overrides {
		defaults[key] = value
	}

	return testhelpers.RenderExample(
		t,
		"credential_resource_username_api_key.tf.tmpl",
		"Name", defaults["Name"],
		"Description", defaults["Description"],
		"Enabled", defaults["Enabled"],
		"Type", defaults["Type"],
		"Username", defaults["Username"],
		"ApiKey", defaults["ApiKey"],
	)
}

func TestAccMorpheusCredentialResourceUsernameApiKeyExampleOk(t *testing.T) {
	t.Parallel()

	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()

	name := acctest.RandomWithPrefix(t.Name())

	resourceConfig, err := RenderCredentialUsernameApiKeyConfig(t, map[string]string{
		"Name": name,
	})
	if err != nil {
		t.Fatal(err)
	}

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"hpe_morpheus_credential.tf_example_credential_username_api_key",
			"name",
			name,
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_credential.tf_example_credential_username_api_key",
			"description",
			"terraform credential example for username api key",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_credential.tf_example_credential_username_api_key",
			"enabled",
			"true",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_credential.tf_example_credential_username_api_key",
			"type",
			"username-api-key",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_credential.tf_example_credential_username_api_key",
			"username",
			"admin",
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:             providerConfig + resourceConfig,
				ExpectNonEmptyPlan: true,
				Check:              checkFn,
				PlanOnly:           true,
			},
			{
				Config: providerConfig + resourceConfig,
				Check:  checkFn,
			},
			{
				Config:             providerConfig + resourceConfig,
				ExpectNonEmptyPlan: false,
				Check:              checkFn,
				PlanOnly:           true,
			},
		},
	})
}
