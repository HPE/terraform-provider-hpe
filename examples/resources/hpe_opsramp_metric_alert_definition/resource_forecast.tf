resource "hpe_opsramp_metric_alert_definition" "example_forecast_alert" {
	name = "Example Forecast Alert"
    alert_type = "METRICS"
	query = "metrics_samples_count"
	
    alert_threshold_type = "FORECAST"
    alert_threshold_data = {
        warning_condition = "5"
        critical_condition = "3"
        limit = 90
    }

    alert_trigger_duration = "0s"

    subject = "$$__name__ alert for $$resource.name$$ - $$component.name$$ - $$metric.value$$ ($$threshold)"
    description = "This is an example metric alert definition created for testing purposes."

	entity_type = ["RESOURCE"]
	component = ["$$__name__"]
	status = true

	labels = []
	attributes = []
}