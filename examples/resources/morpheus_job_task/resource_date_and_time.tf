data "hpe_morpheus_task" "example_task" {
  name = "Deploy app"
}

resource "hpe_morpheus_job_task" "tf_example_job_task_date_and_time" {
  name                    = "TF Example Task Job Date and Time"
  enabled                 = true
  labels                  = ["aws", "demo"]
  task_id                 = data.morpheus_task.example_task.id
  schedule_mode           = "date_and_time"
  scheduled_date_and_time = "2022-12-30T06:00:00Z"
  context_type            = "instance"
  instance_ids            = [1, 2]
}
