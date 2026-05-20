resource "hpe_morpheus_job" "example" {
  name          = "Nightly Cleanup"
  workflow_id   = 1
  schedule_mode = "scheduled"
  target_type   = "appliance"
  enabled       = true
}
