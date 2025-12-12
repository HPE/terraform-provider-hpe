// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package credential_test

import (
	"testing"

	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/testhelpers"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func RenderCredentialUsernameKeypairConfig(
	t *testing.T,
	overrides map[string]string,
) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Description": "terraform credential example for username key pair",
		"Enabled":     "true",
		"Type":        "username-keypair",
		"Username":    "admin",
		"KeyPairId":   "2",
	}

	for key, value := range overrides {
		defaults[key] = value
	}

	return testhelpers.RenderExample(
		t,
		"credential_resource_username_keypair.tf.tmpl",
		"Name", defaults["Name"],
		"Description", defaults["Description"],
		"Enabled", defaults["Enabled"],
		"Type", defaults["Type"],
		"Username", defaults["Username"],
		"KeyPairId", defaults["KeyPairId"],
	)
}

func TestAccMorpheusCredentialResourceUsernameKeypairExampleOk(t *testing.T) {
	t.Parallel()

	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()

	name := acctest.RandomWithPrefix(t.Name())

	resourceConfig, err := RenderCredentialUsernameKeypairConfig(t, map[string]string{
		"Name": name,
	})
	if err != nil {
		t.Fatal(err)
	}

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"hpe_morpheus_credential.tf_example_credential_username_keypair",
			"name",
			name,
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_credential.tf_example_credential_username_keypair",
			"description",
			"terraform credential example for username key pair",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_credential.tf_example_credential_username_keypair",
			"enabled",
			"true",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_credential.tf_example_credential_username_keypair",
			"type",
			"username-keypair",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_credential.tf_example_credential_username_keypair",
			"username",
			"admin",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_credential.tf_example_credential_username_keypair",
			"key_pair_id",
			"2",
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
