// (C) Copyright 2025-2026 Hewlett Packard Enterprise Development LP

package credential_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/resources/credential"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
)

func TestAccMorpheusCredentialResourceTenantUsernameKeypairExampleOk(t *testing.T) {
	if capabilities.Missing(t, capabilities.All) {
		t.Log("Skipping test due to missing capabilities")

		return
	}
	t.Parallel()

	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()

	name := acctest.RandomWithPrefix(t.Name())

	resourceConfig, err := credential.RenderCredentialTenantUsernameKeypairConfig(t, map[string]string{
		"Name": name,
	})
	if err != nil {
		t.Fatal(err)
	}

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"hpe_morpheus_credential.tf_example_credential_tenant_username_keypair",
			"name",
			name,
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_credential.tf_example_credential_tenant_username_keypair",
			"description",
			"terraform credential example for tenant username keypair",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_credential.tf_example_credential_tenant_username_keypair",
			"enabled",
			"true",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_credential.tf_example_credential_tenant_username_keypair",
			"type",
			"tenant-username-keypair",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_credential.tf_example_credential_tenant_username_keypair",
			"tenant",
			"tenant123",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_credential.tf_example_credential_tenant_username_keypair",
			"username",
			"admin",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_credential.tf_example_credential_tenant_username_keypair",
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
