resource "hpe_opsramp_metric_alert_definition" "example_static_threshold_alert" {
  name       = "Example Static Threshold Alert"
  alert_type = "METRICS"
  query      = "metrics_samples_count"

  alert_threshold_type = "STATIC_THRESHOLD"

  alert_threshold_data = {
    warning_condition  = "50000"
    critical_condition = "30000"
  }

  alert_trigger_duration = "1m"

  subject     = "$$__name__ alert for $$resource.name$$ - $$component.name$$ - $$metric.value$$ ($$threshold)"
  description = "This is an example metric alert definition created for testing purposes."

  entity_type = ["RESOURCE"]
  component   = ["$$__name__"]
  status      = true

  labels = [
    {
      name  = "environment"
      value = "production"
    },
    {
      name  = "team"
      value = "ops"
    }
  ]
  attributes = [
    {
      name  = "host"
      value = "$$__name__"
    }
  ]
}