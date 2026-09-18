// (C) Copyright 2025-2026 Hewlett Packard Enterprise Development LP

package credential_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/resources/credential"
	"github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/resources/trust"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
)

func TestAccMorpheusCredentialResourceEmailPrivateKeyExampleOk(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)

	t.Parallel()

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()

	name := acctest.RandomWithPrefix(t.Name())

	// MORPH-16494: the credential fixtures previously pinned the master-tenant
	// key pair id "2", which does not exist as a sub-tenant. Create a key pair
	// per-run and reference its id so the credential is self-contained.
	dependenciesConfig, err := trust.RenderKeyPairConfig(t, map[string]string{
		"Name": name,
	})
	if err != nil {
		t.Fatal(err)
	}

	resourceConfig, err := credential.RenderCredentialEmailPrivateKeyConfig(t, map[string]string{
		"Name":      name,
		"KeyPairId": "hpe_morpheus_key_pair.example.id",
	})
	if err != nil {
		t.Fatal(err)
	}

	config := providerConfig + dependenciesConfig + resourceConfig

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
	}

	appliedChecks := append([]resource.TestCheckFunc{}, checks...)
	appliedChecks = append(appliedChecks, resource.TestCheckResourceAttrPair(
		"hpe_morpheus_credential.tf_example_credential_email_private_key",
		"key_pair_id",
		"hpe_morpheus_key_pair.example",
		"id",
	))

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)
	appliedCheckFn := resource.ComposeAggregateTestCheckFunc(appliedChecks...)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:             config,
				ExpectNonEmptyPlan: true,
				Check:              checkFn,
				PlanOnly:           true,
			},
			{
				Config: config,
				Check:  appliedCheckFn,
			},
			{
				Config:             config,
				ExpectNonEmptyPlan: false,
				Check:              appliedCheckFn,
				PlanOnly:           true,
			},
		},
	})
}
