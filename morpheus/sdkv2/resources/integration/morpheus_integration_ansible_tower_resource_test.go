// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package integration_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	morpheus "github.com/HPE/terraform-provider-hpe/internal/sdk/legacy"

	"github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/resources/integration"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
)

// TestAccMorpheusIntegrationAnsibleTowerExampleOk creates an integration for
// the Tower the environment names (testhelpers/tower_fixture.go). Morpheus
// verifies the Tower on create -- a TCP check, /api/v2/ping/ and /api/v2/me/
// -- so the URL has to answer; the lab's stand-in serves exactly those. The
// password is asserted as configured and the plan after apply is clean: the
// appliance reports only a hash of it, which the resource records when it sets
// the password and compares later reads with (MORPH-20064).
func TestAccMorpheusIntegrationAnsibleTowerExampleOk(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.AnsibleTower)

	towerURL, towerUsername, towerPassword := testhelpers.AnsibleTowerCredentials(t)

	t.Parallel()

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()

	name := acctest.RandomWithPrefix(t.Name())

	// The template quotes its placeholders itself, so the values are escaped
	// for HCL's quoted form, not quoted again.
	resourceConfig, err := integration.RenderIntegrationAnsibleTowerConfig(t, map[string]string{
		"Name":     name,
		"Url":      testhelpers.HCLEscape(towerURL),
		"Username": testhelpers.HCLEscape(towerUsername),
		"Password": testhelpers.HCLEscape(towerPassword),
	})
	if err != nil {
		t.Fatal(err)
	}

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"hpe_morpheus_integration_ansible_tower.tf_example_ansible_tower_integration",
			"name",
			name,
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_integration_ansible_tower.tf_example_ansible_tower_integration",
			"enabled",
			"true",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_integration_ansible_tower.tf_example_ansible_tower_integration",
			"url",
			towerURL,
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_integration_ansible_tower.tf_example_ansible_tower_integration",
			"username",
			towerUsername,
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_integration_ansible_tower.tf_example_ansible_tower_integration",
			"password",
			towerPassword,
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

// TestAccMorpheusIntegrationAnsibleTowerPasswordChangedOutOfBand covers the
// other half of MORPH-20064: a password changed on the appliance behind
// Terraform's back is detected. The appliance reports a secret only as its
// hash; the resource recorded the hash it read back when it set the password,
// so after the change the read sees a different one and puts it in the
// password's place in state, the next plan wants to put the configured
// password back, the apply does, and the plan after that is clean.
// The change is made through the API with the sdkv2 resources' own SDK; the
// mock Tower accepts any password, so Morpheus's re-verification on the update
// passes.
func TestAccMorpheusIntegrationAnsibleTowerPasswordChangedOutOfBand(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.AnsibleTower)

	towerURL, towerUsername, towerPassword := testhelpers.AnsibleTowerCredentials(t)

	t.Parallel()

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()
	name := acctest.RandomWithPrefix(t.Name())

	resourceConfig, err := integration.RenderIntegrationAnsibleTowerConfig(t, map[string]string{
		"Name":     name,
		"Url":      testhelpers.HCLEscape(towerURL),
		"Username": testhelpers.HCLEscape(towerUsername),
		"Password": testhelpers.HCLEscape(towerPassword),
	})
	if err != nil {
		t.Fatal(err)
	}

	const resourceName = "hpe_morpheus_integration_ansible_tower.tf_example_ansible_tower_integration"

	// The integration's id, taken from state after the create, for the
	// out-of-band change.
	var integrationID int64

	rememberID := func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("%s not in state", resourceName)
		}

		id, err := strconv.ParseInt(rs.Primary.ID, 10, 64)
		if err != nil {
			return fmt.Errorf("%s has id %q: %w", resourceName, rs.Primary.ID, err)
		}

		integrationID = id

		return nil
	}

	changePasswordOutOfBand := func() {
		client := testhelpers.NewLegacyClient(t)
		resp, err := client.UpdateIntegration(integrationID, &morpheus.Request{
			Body: map[string]any{
				"integration": map[string]any{
					"servicePassword": towerPassword + "-changed-out-of-band",
				},
			},
		})
		if err := testhelpers.LegacyError(resp, err); err != nil {
			t.Fatalf("changing the password of integration %d out of band: %v", integrationID, err)
		}
	}

	passwordIsConfigured := resource.TestCheckResourceAttr(resourceName, "password", towerPassword)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + resourceConfig,
				Check:  resource.ComposeAggregateTestCheckFunc(passwordIsConfigured, rememberID),
			},
			{
				// Changed behind Terraform's back: the plan must want it back (an
				// in-place update of this one resource), and the apply puts the
				// configured password back.
				PreConfig: changePasswordOutOfBand,
				Config:    providerConfig + resourceConfig,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionUpdate),
					},
				},
				Check: passwordIsConfigured,
			},
			{
				// Nothing is left to do.
				Config:             providerConfig + resourceConfig,
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
		},
	})
}
