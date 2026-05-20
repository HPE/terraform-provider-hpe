resource "hpe_morpheus_task" "example_task" {
    name = "Example Generic Task"
    task_type_code = "nestedWorkflow"

    config = {
        operationalWorkflowId = "4090"
        operationalWorkflowName = "Example Workflow"
    }

    execute_target = "local"
    retryable = false
}
