data "hpe_morpheus_cloud" "example" {
  name = "Standard Cloud"
}

data "hpe_morpheus_group" "example" {
  name = "Example Group"
}

data "hpe_morpheus_tenant" "example" {
  name = "Master Tenant"
}

resource "hpe_morpheus_network" "host" {
  name                       = "example-terraform-host"
  description                = "A host network"
  cloud_id                   = data.hpe_morpheus_cloud.example.id
  pool_id                    = 1
  group_id                   = data.hpe_morpheus_group.example.id
  type_id                    = 1
  config                     = {}
  active                     = true
  dhcp_server                = false
  appliance_url_proxy_bypass = true
  tenant_ids                 = [
    data.hpe_morpheus_tenant.example.id,
  ]
  visibility                 = "private"
  cidr                       = "10.0.0.0/8"
  labels                     = [terraform, example]
}
