
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
    client_id = "abcdefghijklmnopqrstuvwxyz123456"
    client_secret = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890ab"
    endpoint      = "tenant.api.pov.opsramp.com"
    tenant        = "abcdefgh-1234-5678-90ab-cdefghijklmn"
  }
}

# Create individual resources
resource "hpe_opsramp_resource" "resource1" {
  resource_name = "Test1"
  resource_type = "Linux"
}

resource "hpe_opsramp_resource" "resource2" {
  resource_name = "Test2"
  resource_type = "Linux"
}

resource "hpe_opsramp_resource" "resource3" {
  resource_name = "Test3"
  resource_type = "Linux"
}

resource "hpe_opsramp_device_group" "device_group_root" {
  name      = "Test Groups"
  resources = []
}

resource "hpe_opsramp_device_group" "device_group_resources" {
  parent_id = hpe_opsramp_device_group.device_group_root.id
  name      = "Test Resources"
  resources = [hpe_opsramp_resource.resource1.uuid]
}

resource "hpe_opsramp_device_group" "device_group_query" {
  parent_id    = hpe_opsramp_device_group.device_group_root.id
  name         = "Test Queries"
  search_query = format("resourceType = \"Linux\" AND uuid = \"%s\"", hpe_opsramp_resource.resource2.uuid)
}

resource "hpe_opsramp_device_group" "device_group_mixed" {
  parent_id    = hpe_opsramp_device_group.device_group_root.id
  name         = "Test Queries Mixed"
  search_query = format("resourceType = \"Linux\" AND uuid = \"%s\"", hpe_opsramp_resource.resource2.uuid)
  resources    = [hpe_opsramp_resource.resource3.uuid]
}
