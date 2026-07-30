// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package resources_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/opsramp/acctest"
)

// TestAccE2ESimpleAlertPolicies exercises the simple-alert-policies e2e scenario:
// alert correlation, prediction, first response, escalation policies with dependencies.
func TestAccE2ESimpleAlertPolicies(t *testing.T) {
	t.Run("full alert policy stack", func(t *testing.T) {
		corrPolicy := acctest.RandomName("e2e-corr")
		predPolicy := acctest.RandomName("e2e-pred")
		frpPolicy := acctest.RandomName("e2e-frp")
		escPolicy := acctest.RandomName("e2e-esc")
		userGroup := acctest.RandomName("e2e-esc-grp")
		smName := acctest.RandomName("e2e-sm")
		kbCat := acctest.RandomName("e2e-kb-cat")
		kbArticle := acctest.RandomName("e2e-kb-art")
		sdCat := acctest.RandomName("e2e-sd-cat")
		sdImpact := acctest.RandomName("e2e-sd-imp")
		sdUrgency := acctest.RandomName("e2e-sd-urg")

		resource.ParallelTest(t, resource.TestCase{
			PreCheck:                 acctest.PreCheck(t),
			ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories(),
			Steps: []resource.TestStep{
				{
					Config: testAccE2ESimpleAlertPoliciesConfig(
						corrPolicy, predPolicy, frpPolicy, escPolicy,
						userGroup, smName, kbCat, kbArticle,
						sdCat, sdImpact, sdUrgency,
					),
					Check: resource.ComposeAggregateTestCheckFunc(
						testAccEnsureAlertCorrelationPolicyExists(t, "hpe_opsramp_alert_correlation_policy.topology"),
						testAccEnsureAlertPredictionPolicyExists(t, "hpe_opsramp_alert_prediction_policy.prediction"),
						testAccEnsureFirstResponsePolicyExists(t, "hpe_opsramp_first_response_policy.frp"),
						testAccEnsureAlertEscalationPolicyExists(t, "hpe_opsramp_alert_escalation_policy.escalation"),
						resource.TestCheckResourceAttr("hpe_opsramp_alert_correlation_policy.topology", "name", corrPolicy),
						resource.TestCheckResourceAttr("hpe_opsramp_alert_prediction_policy.prediction", "name", predPolicy),
						resource.TestCheckResourceAttr("hpe_opsramp_first_response_policy.frp", "name", frpPolicy),
						resource.TestCheckResourceAttr("hpe_opsramp_alert_escalation_policy.escalation", "name", escPolicy),
					),
				},
			},
		})
	})
}

func testAccE2ESimpleAlertPoliciesConfig(corrPolicy, predPolicy, frpPolicy, escPolicy, userGroup, smName, kbCat, kbArticle, sdCat, sdImpact, sdUrgency string) string {
	return fmt.Sprintf(`
%s
resource "hpe_opsramp_user_group" "alert_test_group" {
	name        = "%s"
	description = "User group for alert policy test"
}

resource "hpe_opsramp_alert_correlation_policy" "topology" {
	name = "%s"

	enabled_mode    = "OBSERVED"
	filter_query    = ""
	inference_query = ""
	type            = "CO_OCCURRENCE"
	machine_learning = {
		continuous_learning = true
		topology            = true
		topology_depth      = 3
	}

	inference_subject = ""
}

resource "hpe_opsramp_first_response_policy" "frp" {
	name = "%s"

	enabled_mode = "OBSERVED"
	filter_query = ""

	pattern_actions = {
		seasonality_time_frame = "7D"
		suppress = {
			seasonal_alerts = true
		}
	}
}

resource "hpe_opsramp_alert_prediction_policy" "prediction" {
	name = "%s"

	enabled_mode = "OFF"
	filter_query = ""

	seasonality_time_frame    = "7D"
	generate_prediction_alert = true
}

resource "hpe_opsramp_kb_category" "alert_kb_cat" {
	name        = "%s"
	description = "KB category for alert policy test"
}

resource "hpe_opsramp_kb_article" "alert_kb_article" {
	subject     = "%s"
	content     = "Article for alert policy testing"
	category_id = hpe_opsramp_kb_category.alert_kb_cat.id
}

resource "hpe_opsramp_servicemap" "alert_sm" {
	name = "%s"
	type = "Service"
}

resource "hpe_opsramp_servicedesk_category" "alert_sd_cat" {
	name        = "%s"
	description = "Category for escalation test"
	ticket_type = "serviceRequests"
}

resource "hpe_opsramp_servicedesk_business_impact" "alert_sd_impact" {
	name        = "%s"
	description = "Business impact for escalation test"
}

resource "hpe_opsramp_servicedesk_urgency" "alert_sd_urgency" {
	name        = "%s"
	description = "Urgency for escalation test"
}

resource "hpe_opsramp_alert_escalation_policy" "escalation" {
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
					id   = hpe_opsramp_user_group.alert_test_group.unique_id
					type = "USERGROUP"
				}
			]
			notification_type        = "basic"
			notification_template_id = "ae6d595e-77a1-5262-a674-ea4c5afa6320"
		},
		{
			wait_mins = 5
			action    = "INCIDENT"
			incident = {
				priority              = "Normal"
				subject               = "Event $alert.subject have been found"
				description           = "Event description $alert.description"
				assignee_group_id     = hpe_opsramp_user_group.alert_test_group.unique_id
				category_id           = hpe_opsramp_servicedesk_category.alert_sd_cat.id
				sub_category_id       = ""
				business_impact_id    = hpe_opsramp_servicedesk_business_impact.alert_sd_impact.id
				urgency_id            = hpe_opsramp_servicedesk_urgency.alert_sd_urgency.id
				knowledge_article_ids = [hpe_opsramp_kb_article.alert_kb_article.id]
				cc                    = "test@example.com"
			}
			update_incident = {
				update_incident_mode             = "UpdateWhenAlertStateChange"
				update_incident_subject_mode     = "UpdateIncidentSubject"
				auto_resolve_incident_mode       = "AutoResolveIncident"
				auto_heal_wait_time              = 0
				update_priority_by_ml_configuration = false
				priority_rules                   = []
			}
		}
	]
	search_query          = "subject CONTAINS \"test\""
	resource_search_query = "serviceGroups.uniqueId = \"${hpe_opsramp_servicemap.alert_sm.id}\""
}
`, acctest.ProviderConfigHCL(), userGroup, corrPolicy, frpPolicy, predPolicy, kbCat, kbArticle, smName, sdCat, sdImpact, sdUrgency, escPolicy)
}
