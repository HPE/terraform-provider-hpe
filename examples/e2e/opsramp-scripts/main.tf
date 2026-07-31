
terraform {
  required_providers {
    hpe = {
      source  = "HPE/hpe"
      version = ">= 1.6.0"
    }
  }
}

provider "hpe" {
  opsramp {
    client_id     = "abcdefghijklmnopqrstuvwxyz123456"
    client_secret = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890ab"
    endpoint      = "tenant.api.pov.opsramp.com"
    tenant        = "abcdefgh-1234-5678-90ab-cdefghijklmn"
  }
}

resource "hpe_opsramp_script_category" "automation" {
  name = "Test Automation Scripts"
}

resource "hpe_opsramp_script_category" "linux" {
  name      = "Test Linux Scripts"
  parent_id = hpe_opsramp_script_category.automation.uuid
}

resource "hpe_opsramp_script" "restart_service" {
  category_id     = hpe_opsramp_script_category.linux.uuid
  name            = "Restart Service"
  description     = "Restart a service on a Linux machine."
  platforms       = ["LINUX"]
  execution_type  = "SHELL"
  install_timeout = 120

  attachment = {
    name = "restart_service_linux.sh"
    file = file("./restart_service_linux.sh")
  }

  parameters = [
    {
      name          = "service_name"
      description   = ""
      default_value = ""
      type          = "REQUIRED"
      data_type     = "STRING"
    },
  ]
}