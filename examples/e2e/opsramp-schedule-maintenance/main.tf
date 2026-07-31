
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

# Create individual resources
resource "hpe_opsramp_resource" "resource1" {
  alias_name    = "MyFirstResource"
  resource_name = "TestResource1"
  resource_type = "Other"
}

resource "hpe_opsramp_resource" "resource2" {
  alias_name    = "MySecondResource"
  hostname      = "testresource2.local"
  resource_type = "Other"
}

resource "hpe_opsramp_device_group" "example_group" {
  name = "Example Device Group"
  resources = [
    hpe_opsramp_resource.resource1.uuid,
    hpe_opsramp_resource.resource2.uuid
  ]
}

resource "hpe_opsramp_scheduled_maintenance" "example" {
  name        = "Example Maintenance Window"
  description = "This is an example maintenance window created for testing purposes."

  schedule = {
    type       = "recurring"
    start_time = "2026-06-10T15:30:00+0100"
    end_time   = "2026-06-10T20:00:00+0100"
    end_by     = "Never"
    pattern = {
      type             = "daily"
      day_frequency    = "everyday"
      repeat_frecuency = 2
    }
    timezone = "Europe/London"

  }

  run_rba             = true
  install_patch       = true
  correlate_alerts    = true
  run_escalate_action = true

  alert_conditions = {
    matching_type = "ALL"
    rules = [
      {
        key      = "subject"
        operator = "endswith"
        value    = "test"
      },
      {
        key      = "description"
        operator = "startswith"
        value    = "test"
      },
      {
        key      = "serviceName"
        operator = "contains"
        value    = "test"
      },
      {
        key      = "resourceName"
        operator = "regex"
        value    = "test"
      }
    ]
  }

  device_ids       = []
  device_group_ids = [hpe_opsramp_device_group.example_group.id]
  site_ids         = []

  notify_before_end_time   = "0"
  notify_before_start_time = "0"
  user_ids                 = []
  user_group_ids           = []
}