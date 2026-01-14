resource "hpe_morpheus_workflow_provisioning" "example" {
  name        = "tf_example_provisioning_workflow"
  description = "Terraform provisioning workflow example"
  labels      = ["demo", "terraform"]
  platform    = "all"
  visibility  = "private"
}
