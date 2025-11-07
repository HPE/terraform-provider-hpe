resource "hpe_morpheus_job_workflow" "tf_example_workflow_job_date_and_time" {
  name           = "TF Example Workflow Job Manual"
  enabled        = true
  labels         = ["aws", "demo"]
  workflow_id    = 1
  schedule_mode  = "manual"
  context_type   = "instance-label"
  instance_label = "demo"
  custom_options = {
    "demo" = "testing"
  }
}
