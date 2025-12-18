data "hpe_morpheus_cloud" "example" {
  name = "Google Cloud"
}

data "hpe_morpheus_group" "example" {
  name = "Examle Group"
}

data "hpe_morpheus_tenant" "example" {
  name = "Master Tenant"
}

resource "hpe_morpheus_network" "gcp" {
  name        = "example-terraform-gcp"
  description = "GCP network"
  cloud_id    = data.hpe_morpheus_cloud.example.id
  pool_id     = 1
  group_id    = data.hpe_morpheus_group.example.id
  type_id     = 38
  config = {
    mtu        = "1460"
    autoCreate = true
  }
  active                     = true
  dhcp_server                = false
  appliance_url_proxy_bypass = true
  tenant_ids = [
    data.hpe_morpheus_tenant.example.id,
  ]
  visibility   = "private"
  cidr         = "10.0.0.0/8"
  zone_pool_id = 85990
  labels       = ["terraform", "example"]
}
