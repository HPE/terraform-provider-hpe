# Power Scheduling Policy - Enforces power schedules
# Allowed associated_resource_types: Group, Cloud, User, Global
# Tenant specification: allowed (can specify tenants array)
resource "hpe_morpheus_policy" "power_schedule" {
  name                     = "Power Scheduling Policy"
  description              = "Enforce power schedules for instances"
  associated_resource_type = "User"
  associated_resource_id   = 9969
  enabled                  = true

  policy_type = {
    code = "powerSchedule"
  }

  config = {
    powerScheduleType      = "user" # Options: "user" (user configurable), "fixed" (strict schedule)
    powerSchedule          = "1"    # ID of the power schedule
    powerScheduleHideFixed = false  # Hide fixed schedule from users
  }
}
