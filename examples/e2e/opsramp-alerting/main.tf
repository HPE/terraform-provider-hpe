
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
    client_id = "abcdefghijklmnopqrstuvwxyz123456"
    client_secret = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890ab"
    endpoint      = "tenant.api.pov.opsramp.com"
    tenant        = "abcdefgh-1234-5678-90ab-cdefghijklmn"
  }
}

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

resource "hpe_opsramp_metric_alert_definition" "example_dynamic_threshold_alert" {
	name = "Example Forecast Alert"
    alert_type = "METRICS"
	query = "metrics_samples_count"
	
    alert_threshold_type = "DYNAMIC_THRESHOLD"
    alert_threshold_data = {
        limit = 2
    }
    alert_trigger_duration = "1m"

    no_data_condition = "WARNING_ALERT"

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


resource "hpe_opsramp_log_alert_definition" "example_log_alert" {
    name = "Example Log Alert"
    query = "{__resource.CODE_FILE = \"src/core/unit.c\"}"
    
    alert_no_data = "noalert"

    status = "disabled"

    heal_query = "noAutoHeal"

    entity_type = "RESOURCE"
    component = "$$name"

    subject = "Log alert for $$resource.name$$ - $$component.name$$"
    description = "This is an example log alert definition created for testing purposes."

    conditions = [
        {
            operator = ">="
            severity = "critical"
            value = 70
        },
        {
            operator = ">="
            severity = "warning"
            value = 50
        }
    ]

    labels = {
        "environment" = "production",
        "team" = "ops"
    }
    
    resource_attributes = {
        "host" = "$$__name__"
    }

    schedule = {
        timezone = "Europe/Paris"

        pattern = {
            type = "second"
            repeat_frequency = 60
            //week_days = "Monday,Tuesday,Wednesday,Thursday,Friday,Saturday,Sunday"
            //day_of_month = "1,15"
            //week_index = "First,Third,Last"
            //day_of_week = "Monday,Tuesday,Wednesday,Thursday,Friday,Saturday,Sunday"
        }

        
    }

}