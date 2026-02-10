// (C) Copyright 2025-2026 Hewlett Packard Enterprise Development LP

package credential_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/resources/credential"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/systemoverride"
)

func TestAccMorpheusCredentialResourceEmailPrivateKeyExampleOk(t *testing.T) {
	t.Parallel()

	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	testSystem := systemoverride.GetPreferred(t, "feature")
	providerConfig := testhelpers.ProviderBlock(testSystem)

	name := acctest.RandomWithPrefix(t.Name())

	resourceConfig, err := credential.RenderCredentialEmailPrivateKeyConfig(t, map[string]string{
		"Name": name,
	})
	if err != nil {
		t.Fatal(err)
	}

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"hpe_morpheus_credential.tf_example_credential_email_private_key",
			"name",
			name,
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_credential.tf_example_credential_email_private_key",
			"description",
			"terraform credential example for email private key",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_credential.tf_example_credential_email_private_key",
			"enabled",
			"true",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_credential.tf_example_credential_email_private_key",
			"type",
			"email-private-key",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_credential.tf_example_credential_email_private_key",
			"email",
			"test@example.local",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_credential.tf_example_credential_email_private_key",
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
