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