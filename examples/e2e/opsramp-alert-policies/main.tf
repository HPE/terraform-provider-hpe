
terraform {
  required_providers {
    hpe = {
      source  = "HPE/hpe"
      version = ">= 1.6.0"
    }
  }
}

provider "hpe" {
  opsramp {
    client_id     = "abcdefghijklmnopqrstuvwxyz123456"
    client_secret = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890ab"
    endpoint      = "tenant.api.pov.opsramp.com"
    tenant        = "abcdefgh-1234-5678-90ab-cdefghijklmn"
  }
}

# User Group
resource "hpe_opsramp_user_group" "client_user_group" {
  name        = "User Group"
  description = "User group"
}

resource "hpe_opsramp_alert_correlation_policy" "topology_correlation_policy" {
  name = "Topology-based"

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

resource "hpe_opsramp_alert_correlation_policy" "similarity_correlation_policy" {
  name = "Similarity-based"

  enabled_mode    = "OBSERVED"
  filter_query    = ""
  inference_query = ""
  type            = "CO_OCCURRENCE"
  machine_learning = {
    continuous_learning = true
    topology            = false
    matching_conditions = [
      {
        property   = "service_group"
        match_type = "Identical"
      }
    ]
  }

  inference_subject = ""
}

resource "hpe_opsramp_first_response_policy" "seasonality_first_response_policy" {
  name = "Seasonal-based Suppression"

  enabled_mode = "OBSERVED"
  filter_query = ""

  pattern_actions = {
    seasonality_time_frame = "7D"
    suppress = {
      seasonal_alerts = true
    }
  }
}

resource "hpe_opsramp_alert_prediction_policy" "default_alert_prediction_policy" {
  name = "Default Prediction Policy"

  enabled_mode = "OFF"
  filter_query = ""

  seasonality_time_frame    = "7D"
  generate_prediction_alert = true
}

resource "hpe_opsramp_kb_category" "kb_category_default" {
  name        = "Default articles"
  description = "Category for default articles"
}

resource "hpe_opsramp_kb_article" "kb_article_default" {
  subject     = "Default article"
  content     = "Default article's description"
  category_id = hpe_opsramp_kb_category.kb_category_default.id
}


resource "hpe_opsramp_servicemap" "sm_root" {
  name = "Default Application"
  type = "Service"
}

# Create multiple categories
resource "hpe_opsramp_servicedesk_category" "category1" {
  name        = "Category1"
  description = "Category1 Description"
  ticket_type = "serviceRequests"
}

# Create  business impacts
resource "hpe_opsramp_servicedesk_business_impact" "business_impact1" {
  name        = "Business Impact1"
  description = "Business Impact1 Description"
}

# Create urgencies
resource "hpe_opsramp_servicedesk_urgency" "urgency1" {
  name        = "Urgency1"
  description = "Urgency1 Description"
}

resource "hpe_opsramp_alert_escalation_policy" "default_alert_escalation_policy" {
  name            = "Default escalation policy"
  precedence      = 1
  enabled_mode    = "OBSERVED"
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
        cc                    = "enrique.larriba@hpe.com"
      }
      update_incident = {
        update_incident_mode                = "UpdateWhenAlertStateChange"
        update_incident_subject_mode        = "UpdateIncidentSubject"
        auto_resolve_incident_mode          = "AutoResolveIncident"
        auto_heal_wait_time                 = 0
        update_priority_by_ml_configuration = false
        priority_rules = [
          {
            alert_state        = "CRITICAL"
            business_impact_id = hpe_opsramp_servicedesk_business_impact.business_impact1.id
            urgency_id         = hpe_opsramp_servicedesk_urgency.urgency1.id
            priority           = "Urgent"
          }
        ]
      }
    }
  ]
  search_query          = "subject CONTAINS \"test\""
  resource_search_query = "name CONTAINS \"test\""
}