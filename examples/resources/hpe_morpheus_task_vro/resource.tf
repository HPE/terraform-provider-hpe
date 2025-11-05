resource "hpe_morpheus_task_vro" "tf_example_task_vro" {
  name               = "tfexample vro-task"
  code               = "tfexample-vro-task"
  labels             = ["demo", "terraform"]
  vro_integration_id = 1
  vro_workflow_value = 1
  body               = <<EOF
{
    "parameters": [
        {
            "name": "vmName",
            "type": "string",
            "value": {
                "string": {
                    "value": "<%=instance.hostname%>"
                }
            }
        }
    ]
}
EOF
  execute_target     = "local"
  retryable          = false
}
