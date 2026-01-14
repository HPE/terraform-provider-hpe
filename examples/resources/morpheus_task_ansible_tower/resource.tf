resource "hpe_morpheus_task_ansible_tower" "example" {
  name                         = "tfexample_task_ansible_tower"
  code                         = "tfexample-ansible-tower-task"
  labels                       = ["demo", "terraform"]
  ansible_tower_integration_id = 1
  ansible_tower_inventory_id   = 5
  group                        = "demo"
  job_template_id              = 3
  scm_override                 = "main"
  execute_mode                 = "executeAll"
  execute_target               = "local"
  retryable                    = true
  retry_count                  = 5
  retry_delay_seconds          = 10
  allow_custom_config          = true
  visibility                   = "public"
}
