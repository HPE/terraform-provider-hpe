resource "hpe_morpheus_job_workflow" "example" {
  name                  = "TF Example Workflow Job Schedule"
  enabled               = true
  labels                = ["aws", "demo"]
  workflow_id           = 1
  schedule_mode         = "scheduled"
  execution_schedule_id = 1
  context_type          = "instance"
  instance_ids          = [91]
  custom_options        = { "demo" = "testing" }
}
