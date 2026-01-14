resource "hpe_morpheus_job_task" "example" {
  name                    = "TF Example Job Task Date and Time"
  enabled                 = true
  labels                  = ["aws", "demo"]
  task_id                 = 1
  schedule_mode           = "date_and_time"
  scheduled_date_and_time = "2022-12-30T06:00:00Z"
  context_type            = "instance"
  instance_ids            = [1, 2]
}
