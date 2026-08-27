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