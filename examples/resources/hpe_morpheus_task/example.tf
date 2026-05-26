resource "hpe_morpheus_task" "example" {
  name           = "Example Task"
  task_type_code = "shellTask"
  execute_target = "local"
}
