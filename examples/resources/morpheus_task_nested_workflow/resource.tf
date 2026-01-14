resource "hpe_morpheus_task_nested_workflow" "example" {
  name                      = "tfexample_nested_workflow"
  code                      = "tfexample_nested_workflow"
  labels                    = ["demo", "terraform"]
  operational_workflow_id   = 1
  operational_workflow_name = "Example workflow"
}
