// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package resources_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/HPE/terraform-provider-hpe/opsramp/acctest"
)

func TestAccAlertEscalationPolicyResource(t *testing.T) {
	clientOverride := acctest.OptionalClientOverride(t)

	t.Run("create", func(t *testing.T) {
		policyName := acctest.RandomName("esc-policy")
		groupName := acctest.RandomName("esc-group")

		resource.ParallelTest(t, resource.TestCase{
			PreCheck:                 acctest.PreCheck(t),
			ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories(),
			CheckDestroy:             testAccCheckAlertEscalationPolicyDestroy(t),
			Steps: []resource.TestStep{
				{
					Config: testAccAlertEscalationPolicyConfig(policyName, groupName, clientOverride),
					Check: resource.ComposeAggregateTestCheckFunc(
						testAccEnsureAlertEscalationPolicyExists(t, "hpe_opsramp_alert_escalation_policy.test_policy"),
						resource.TestCheckResourceAttrSet("hpe_opsramp_alert_escalation_policy.test_policy", "id"),
						resource.TestCheckResourceAttr("hpe_opsramp_alert_escalation_policy.test_policy", "name", policyName),
						resource.TestCheckResourceAttr(
							"hpe_opsramp_alert_escalation_policy.test_policy",
							"enabled_mode",
							"OBSERVED",
						),
					),
				},
			},
		})
	})
}

func testAccAlertEscalationPolicyConfig(policyName string, groupName string, clientOverride string) string {
	clientAttr := acctest.ClientAttrHCL(clientOverride)

	return fmt.Sprintf(`
%s
resource "hpe_opsramp_user_group" "esc_test_group" {
	name        = "%s"
	description = "Group for escalation policy test"
	%s
}

resource "hpe_opsramp_alert_escalation_policy" "test_policy" {
	name         = "%s"
	precedence   = 1
	enabled_mode = "OBSERVED"
	%s

	escalation_type = "AUTOMATIC_UNTIL_ACKNOWLEDGED_CLOSED_SUPPRESSED_TICKETED"
	policy_type     = "ESCALATION_POLICY"

	escalations = [
		{
			wait_mins = 5
			action    = "INCIDENT"
			incident = {
				priority              = "Normal"
				subject               = "Event $alert.subject have been found"
				description           = "Event description $alert.description"
				assignee_group_id     = ""
				category_id           = ""
				sub_category_id       = ""
				business_impact_id    = ""
				urgency_id            = ""
				knowledge_article_ids = []
			}
			update_incident = {
				update_incident_mode         = "UpdateWhenAlertStateChange"
				update_incident_subject_mode = "UpdateIncidentSubject"
				auto_resolve_incident_mode   = "AutoResolveIncident"
				auto_heal_wait_time          = 0

				update_priority_by_ml_configuration = false
				priority_rules                      = []
			}
		}
	]
	search_query = "subject CONTAINS \"test\""
	resource_search_query = "name CONTAINS \"test\""
}
`, acctest.ProviderConfigHCL(), groupName, clientAttr, policyName, clientAttr)
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

		tenantID, _ := acctest.LookupProviderEnv("tenant")

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
			if rs.Type != "hpe_opsramp_alert_escalation_policy" {
				continue
			}

			tenantID, _ := acctest.LookupProviderEnv("tenant")

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
