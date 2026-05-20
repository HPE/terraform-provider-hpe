resource "hpe_morpheus_task" "example_task" {
    name = "Example Conditional Workflow Task"
    task_type_code = "conditionalWorkflow"
    config_conditional_workflow = {
        conditional_script = <<-EOT
        if (1 == true) {
            return true;
        }

        return false;
        EOT
        if_operational_workflow_id   = "4090"
        if_operational_workflow_name = "Example If Workflow"

        else_operational_workflow_id   = null
        else_operational_workflow_name = null
    }

    execute_target = "local"
    retryable = false
    allow_custom_config = true
}