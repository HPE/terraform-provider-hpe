resource "hpe_morpheus_task_library_template" "example" {
  name                = "Example Terraform Library Template Task"
  code                = "tf-example-library-template-task"
  labels              = ["demo", "library", "terraform"]
  execute_target      = "resource"
  file_template       = "My file template"
  file_template_id    = 1
  retryable           = true
  retry_count         = 1
  retry_delay_seconds = 10
  allow_custom_config = true
}
