// SPDX-FileCopyrightText: Copyright Hewlett Packard Enterprise Development LP
// SPDX-License-Identifier: Apache-2.0

package resources_test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/HPE/terraform-provider-opsramp/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccAlertEscalationPolicyResource(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		policyName := acctest.RandomName("esc-policy")
		groupName := acctest.RandomName("esc-group")

		resource.ParallelTest(t, resource.TestCase{
			PreCheck:                 acctest.PreCheck(t),
			ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories(),
			CheckDestroy:             testAccCheckAlertEscalationPolicyDestroy(t),
			Steps: []resource.TestStep{
				{
					Config: testAccAlertEscalationPolicyConfig(policyName, groupName),
					Check: resource.ComposeAggregateTestCheckFunc(
						testAccEnsureAlertEscalationPolicyExists(t, "opsramp_alert_escalation_policy.test_policy"),
						resource.TestCheckResourceAttrSet("opsramp_alert_escalation_policy.test_policy", "id"),
						resource.TestCheckResourceAttr("opsramp_alert_escalation_policy.test_policy", "name", policyName),
						resource.TestCheckResourceAttr("opsramp_alert_escalation_policy.test_policy", "enabled_mode", "OBSERVED"),
					),
				},
			},
		})
	})
}

func testAccAlertEscalationPolicyConfig(policyName string, groupName string) string {
	return fmt.Sprintf(`
%s
resource "opsramp_user_group" "esc_test_group" {
	name        = "%s"
	description = "Group for escalation policy test"
}

resource "opsramp_alert_escalation_policy" "test_policy" {
	name         = "%s"
	precedence   = 1
	enabled_mode = "OBSERVED"

	escalation_type = "AUTOMATIC_UNTIL_ACKNOWLEDGED_CLOSED_SUPPRESSED_TICKETED"
	policy_type     = "ESCALATION_POLICY"

	escalations = [
		{
			wait_mins          = 0
			priority           = "Normal"
			repeat_frequency   = 5
			notify_limit_count = 2
			action             = "NOTIFICATION"
			recipients = [
				{
					id   = opsramp_user_group.esc_test_group.unique_id
					type = "USERGROUP"
				}
			]
			notification_type        = "basic"
			notification_template_id = "ae6d595e-77a1-5262-a674-ea4c5afa6320"
		}
	]
	search_query = "subject CONTAINS \"test\""
}
`, acctest.ProviderConfigHCL(), groupName, policyName)
}

func testAccEnsureAlertEscalationPolicyExists(t *testing.T, resourceName string) resource.TestCheckFunc {
	t.Helper()

	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource not found in state: %s", resourceName)
		}

		id := strings.TrimSpace(rs.Primary.ID)
		if id == "" {
			return fmt.Errorf("resource id is empty in state for %s", resourceName)
		}

		tenantID := os.Getenv("OPSRAMP_TENANT")
		if clientID, ok := rs.Primary.Attributes["client"]; ok && strings.TrimSpace(clientID) != "" {
			tenantID = clientID
		}

		apiClient, err := acctest.APIClient(t)
		if err != nil {
			return fmt.Errorf("failed to initialize opsramp api client: %w", err)
		}

		_, err = apiClient.GetAlertEscalationPolicy(tenantID, id)
		if err != nil {
			return fmt.Errorf("alert escalation policy %s (%s) was not found in opsramp api: %w", resourceName, id, err)
		}

		return nil
	}
}

func testAccCheckAlertEscalationPolicyDestroy(t *testing.T) resource.TestCheckFunc {
	t.Helper()

	return func(s *terraform.State) error {
		apiClient, err := acctest.APIClient(t)
		if err != nil {
			return fmt.Errorf("failed to initialize opsramp api client: %w", err)
		}

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "opsramp_alert_escalation_policy" {
				continue
			}

			tenantID := os.Getenv("OPSRAMP_TENANT")
			if clientID, ok := rs.Primary.Attributes["client"]; ok && strings.TrimSpace(clientID) != "" {
				tenantID = clientID
			}

			_, err := apiClient.GetAlertEscalationPolicy(tenantID, rs.Primary.ID)
			if err == nil {
				return fmt.Errorf("alert escalation policy still exists: %s", rs.Primary.ID)
			}

			errText := strings.ToLower(err.Error())
			if !strings.Contains(errText, "404") && !strings.Contains(errText, "no escalate alert policy details found") {
				return fmt.Errorf("unexpected error checking deleted alert escalation policy %s: %w", rs.Primary.ID, err)
			}
		}

		return nil
	}
}
