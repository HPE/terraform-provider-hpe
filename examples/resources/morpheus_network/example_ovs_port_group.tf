data "hpe_morpheus_cloud" "example" {
  name = "Morpheus Standard Cloud"
}

data "hpe_morpheus_group" "example" {
  name = "ExampleGroup"
}

data "hpe_morpheus_tenant" "example" {
  name = "Master Tenant"
}

resource "hpe_morpheus_network" "ovs_port_group" {
  name                       = "Terraform OVS Port Group"
  description                = "OVS Port Group network"
  cloud_id                   = data.hpe_morpheus_cloud.example.id
  pool_id                    = 3251
  group_id                   = data.hpe_morpheus_group.example.id
  type_id                    = 63
  switch_id                  = Compute
  config                     = {}
  active                     = true
  dhcp_server                = false
  appliance_url_proxy_bypass = true
  tenant_ids                 = [
    data.hpe_morpheus_tenant.example.id,
  ]
  visibility                 = "public"
  cidr                       = "10.32.148.0/22"
  zone_pool_id               = 62299
  vlan_id                    = 43
  labels                     = ["terraform", "example"]

  lifecycle {
    ignore_changes = [name, display_name, description]
  }
}
