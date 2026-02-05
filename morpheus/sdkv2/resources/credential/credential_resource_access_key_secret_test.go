// (C) Copyright 2025-2026 Hewlett Packard Enterprise Development LP

package credential_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/resources/credential"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
)

func TestAccMorpheusCredentialResourceAccessKeySecretExampleOk(t *testing.T) {
	t.Parallel()

	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()

	name := acctest.RandomWithPrefix(t.Name())

	resourceConfig, err := credential.RenderCredentialAccessKeySecretConfig(t, map[string]string{
		"Name": name,
	})
	if err != nil {
		t.Fatal(err)
	}

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"hpe_morpheus_credential.tf_example_credential_access_key_secret",
			"name",
			name,
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_credential.tf_example_credential_access_key_secret",
			"description",
			"terraform credential example for access key and secret key",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_credential.tf_example_credential_access_key_secret",
			"enabled",
			"true",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_credential.tf_example_credential_access_key_secret",
			"type",
			"access-key-secret",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_credential.tf_example_credential_access_key_secret",
			"access_key",
			"FIEFMIQNQ",
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
