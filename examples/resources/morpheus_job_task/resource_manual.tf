data "hpe_morpheus_task" "example_task" {
  name = "Deploy app"
}

resource "hpe_morpheus_job_task" "tf_example_job_task_manual" {
  name           = "TF Example Task Job Manual"
  enabled        = true
  labels         = ["aws", "demo"]
  task_id        = data.hpe_morpheus_task.example_task.id
  schedule_mode  = "manual"
  context_type   = "instance-label"
  instance_label = "demo"
}
