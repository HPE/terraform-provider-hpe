resource "hpe_opsramp_alert_escalation_policy" "default_alert_escalation_policy" {
  name       = "Default escalation policy"
  precedence = 1

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
          id   = hpe_opsramp_user_group.client_user_group.unique_id
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
        assignee_group_id     = hpe_opsramp_user_group.client_user_group.unique_id
        category_id           = hpe_opsramp_servicedesk_category.category1.id
        sub_category_id       = ""
        business_impact_id    = hpe_opsramp_servicedesk_business_impact.business_impact1.id
        urgency_id            = hpe_opsramp_servicedesk_urgency.urgency1.id
        knowledge_article_ids = [hpe_opsramp_kb_article.kb_article_default.id]
        cc                    = "mail@example.com"
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
  search_query          = "subject CONTAINS \"test\""
  resource_search_query = "serviceGroups.uniqueId = \"${hpe_opsramp_servicemap.sm_root.id}\""
}