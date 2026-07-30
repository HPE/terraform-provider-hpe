resource "hpe_opsramp_metric_alert_definition" "example_dynamic_threshold_alert" {
  name       = "Example Forecast Alert"
  alert_type = "METRICS"
  query      = "metrics_samples_count"

  alert_threshold_type = "DYNAMIC_THRESHOLD"
  alert_threshold_data = {
    limit = 2
  }
  alert_trigger_duration = "1m"

  no_data_condition = "WARNING_ALERT"

  subject     = "$$__name__ alert for $$resource.name$$ - $$component.name$$ - $$metric.value$$ ($$threshold)"
  description = "This is an example metric alert definition created for testing purposes."

  entity_type = ["RESOURCE"]
  component   = ["$$__name__"]
  status      = true

  labels     = []
  attributes = []
}