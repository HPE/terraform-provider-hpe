resource "hpe_opsramp_metric_alert_definition" "example_dynamic_change_alert" {
	name = "Example Dynamic Change Alert"
    alert_type = "METRICS"
	query = "metrics_samples_count"
	
    alert_threshold_type = "DYNAMIC_CHANGE_DETECTION"
    
    alert_threshold_data = {
        direction = "increaseordecrease"
        learning_period = "4h"
        standard_deviation = 2
    }
    alert_trigger_duration = "5m"

    subject = "$$__name__ alert for $$resource.name$$ - $$component.name$$ - $$metric.value$$ ($$threshold)"
    description = "This is an example metric alert definition created for testing purposes."

	entity_type = ["RESOURCE"]
	component = ["$$__name__"]
	status = true

	labels = [
        {
            name = "environment"
            value = "production"
        },
        {
            name = "team"
            value = "ops"
        }
    ]
	attributes = [
        {
            name = "host"
            value = "$$__name__"
        }
    ]
}