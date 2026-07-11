resource "hpe_opsramp_script" "restart_service" {
  category_id     = hpe_opsramp_script_category.os_management.uuid
  name            = "Restart Service"
  description     = "Restarts a named OS service."
  platforms       = ["LINUX"]
  execution_type  = "SHELL"
  install_timeout = 120

  attachment = {
    name        = "restart_service_linux.sh"
    file = file("restart_service_linux.sh")
  }

  parameters = [
    {
      name          = "service_name"
      description   = "Name of the service to restart"
      default_value = "opsramp-agent"
      type          = "REQUIRED"
      data_type     = "STRING"
    },
  ]
}