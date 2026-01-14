resource "hpe_morpheus_task_library_script" "example" {
  name                = "Example Terraform Library Script Task"
  code                = "tf-example-library-script-task"
  labels              = ["demo", "library", "terraform"]
  execute_target      = "resource"
  script_template     = "My script template"
  script_template_id  = 1
  retryable           = true
  retry_count         = 1
  retry_delay_seconds = 10
  allow_custom_config = true
}
