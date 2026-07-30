
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
  resource_name = "Test1"
  resource_type = "Linux"
}

resource "hpe_opsramp_resource" "resource2" {
  resource_name = "Test2"
  resource_type = "Linux"
}

# Create a service with multiple resources
resource "hpe_opsramp_servicemap" "servicemap_root" {
  name = "TestRoot"
  type = "Service"
}

resource "hpe_opsramp_servicemap" "servicemap_child1" {
  name   = "Test1"
  type   = "Service"
  parent = hpe_opsramp_servicemap.servicemap_root.id
}

resource "hpe_opsramp_servicemap" "servicemap_child2" {
  name   = "Test2"
  type   = "Service"
  parent = hpe_opsramp_servicemap.servicemap_root.id
}

resource "hpe_opsramp_servicemap" "servicemap_child21" {
  name      = "Test21"
  type      = "Resource"
  parent    = hpe_opsramp_servicemap.servicemap_child2.id
  resources = [hpe_opsramp_resource.resource1.uuid]
}

resource "hpe_opsramp_servicemap" "servicemap_child22" {
  name         = "Test22"
  type         = "Resource"
  parent       = hpe_opsramp_servicemap.servicemap_child2.id
  search_query = "resourceType = \"Server\" AND name CONTAINS \"Test\""
}

# Service map links
resource "hpe_opsramp_servicemap" "servicemap_linked_root" {
  name = "Test Linked Root"
  type = "Service"
}

resource "hpe_opsramp_servicemap_link" "servicemap_link" {
  parent = hpe_opsramp_servicemap.servicemap_root.id
  link   = hpe_opsramp_servicemap.servicemap_linked_root.id
}