resource "hpe_morpheus_task" "example_task" {
  name           = "Example Generic Task"
  task_type_code = "nestedWorkflow"

  config = {
    operationalWorkflowId   = "90"
    operationalWorkflowName = "Test 1"
  }

  execute_target = "local"
  retryable      = false
}
