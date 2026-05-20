data "hpe_morpheus_cloud" "example" {
  name = "AWS Cloud"
}

data "hpe_morpheus_group" "example" {
  name = "Example Group"
}

data "hpe_morpheus_tenant" "example" {
  name = "Master Tenant"
}

resource "hpe_morpheus_network" "aws" {
  name        = "example-terraform-aws"
  description = "AWS subnet"
  cloud_id    = data.hpe_morpheus_cloud.example.id
  pool_id     = 1
  group_id    = data.hpe_morpheus_group.example.id
  type_id     = 36
  config = {
    assignPublicIp   = true
    availabilityZone = "us-west-1a"
  }
  active                     = true
  dhcp_server                = true
  appliance_url_proxy_bypass = true
  tenant_ids = [
    data.hpe_morpheus_tenant.example.id,
  ]
  visibility   = "private"
  cidr         = "10.200.99.0/24"
  zone_pool_id = 12329
  labels       = ["terraform", "example"]

  lifecycle {
    ignore_changes = [name, display_name, description]
  }
}
