
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

resource "hpe_opsramp_site" "vmware_site_spain" {
  name    = "Spain Site"
  country = "Spain"
}

# Create individual resources
resource "hpe_opsramp_resource" "resource1" {
  resource_name = "TestResource1"
  resource_type = "Linux"
}

resource "hpe_opsramp_resource" "resource2" {
  resource_name = "TestResource2"
  resource_type = "Linux"
}

resource "hpe_opsramp_resource" "resource3" {
  resource_name = "TestResource3"
  resource_type = "Linux"
}

resource "hpe_opsramp_site" "vmware_site_valencia" {
  parent_id = hpe_opsramp_site.vmware_site_spain.id
  name      = "Valencia Data Center"
  address   = "Av. del General Avilés, 35-37, Benicalap"
  country   = "Spain"
  zip       = "46035"
  state     = "Comunitat Valenciana"
  city      = "València"
  resources = [hpe_opsramp_resource.resource1.uuid]
}

resource "hpe_opsramp_site" "vmware_site_madrid" {
  parent_id    = hpe_opsramp_site.vmware_site_spain.id
  name         = "Madrid Data Center"
  address      = "Calle Vicente Aleixandre, 1"
  country      = "Spain"
  zip          = "28232"
  state        = "Madrid"
  city         = "Las Rozas de Madrid"
  search_query = format("uuid = \"%s\"", hpe_opsramp_resource.resource2.uuid)
}

resource "hpe_opsramp_site" "vmware_site_barcelona" {
  parent_id    = hpe_opsramp_site.vmware_site_spain.id
  name         = "Barcelona Data Center"
  address      = "Carrer de Tànger, 66"
  country      = "Spain"
  zip          = "08018"
  state        = "Barcelona"
  city         = "Sant Martí"
  search_query = format("uuid = \"%s\"", hpe_opsramp_resource.resource2.uuid)
  resources = [
    hpe_opsramp_resource.resource3.uuid
  ]
}