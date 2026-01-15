resource "hpe_morpheus_job_task" "example" {
  name                  = "TF Example Job Task Schedule"
  enabled               = true
  labels                = ["aws", "demo"]
  task_id               = 1
  schedule_mode         = "scheduled"
  execution_schedule_id = 1
  context_type          = "instance"
  instance_ids          = [91]
  custom_config         = "{\"test\":\"new\"}"
}
