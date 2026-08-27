resource "hpe_morpheus_job_task" "example" {
  name           = "TF Example Job Task Manual"
  enabled        = true
  labels         = ["aws", "demo"]
  task_id        = 1
  schedule_mode  = "manual"
  context_type   = "instance-label"
  instance_label = "demo"
}
