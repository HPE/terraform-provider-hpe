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