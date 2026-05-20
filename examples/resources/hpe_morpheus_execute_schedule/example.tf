resource "hpe_morpheus_execute_schedule" "example" {
  name              = "Daily Maintenance"
  description       = "Runs daily at midnight"
  schedule_type     = "execute"
  schedule_timezone = "America/New_York"
  cron              = "0 0 * * *"
  enabled           = true
}
