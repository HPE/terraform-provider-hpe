data "hpe_morpheus_task" "example_task" {
  name = "Deploy app"
}


resource "hpe_morpheus_job_task" "tf_example_job_task_schedule" {
  name                  = "TF Example Task Job Schedule"
  enabled               = true
  labels                = ["aws", "demo"]
  task_id               = data.morpheus_task.example_task.id
  schedule_mode         = "scheduled"
  execution_schedule_id = 1
  context_type          = "instance"
  instance_ids          = [91]
  custom_config         = "{\"test\":\"new\"}"
}
