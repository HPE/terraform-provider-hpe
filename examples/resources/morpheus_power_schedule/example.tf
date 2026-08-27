resource "hpe_morpheus_power_schedule" "example" {
  name              = "Business Hours"
  description       = "Power on during business hours"
  schedule_type     = "power"
  schedule_timezone = "America/New_York"
  enabled           = true
  monday_on_time    = "08:00"
  monday_off_time   = "18:00"
  tuesday_on_time   = "08:00"
  tuesday_off_time  = "18:00"
}
