resource "hpe_morpheus_workflow_provisioning" "tf_example_provisioning_workflow" {
  name        = "tf_example_provisioning_workflow"
  description = "Terraform provisioning workflow example"
  labels      = ["demo", "terraform"]
  platform    = "all"
  visibility  = "private"
  task {
    task_id    = 18
    task_phase = "configure"
  }
}