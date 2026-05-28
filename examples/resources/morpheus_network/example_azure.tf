data "hpe_morpheus_cloud" "example" {
  name = "Azure Cloud"
}

data "hpe_morpheus_group" "example" {
  name = "Example Group"
}

data "hpe_morpheus_tenant" "example" {
  name = "Master Tenant"
}

resource "hpe_morpheus_network" "azure" {
  name                       = "example-terraform-azure"
  description                = "Azure network"
  cloud_id                   = data.hpe_morpheus_cloud.example.id
  pool_id                    = 1
  group_id                   = data.hpe_morpheus_group.example.id
  type_id                    = 35
  cidr                       = "10.100.0.0/16"
  visibility                 = "public"
  active                     = true
  dhcp_server                = true
  appliance_url_proxy_bypass = false
  labels                     = ["terraform", "example"]
  config = {
    "resourceGroupId" = all-attrs-resource-group
    "subnetName"      = "all-attrs-subnet"
    "subnetCidr"      = "10.100.1.0/24"
  }
  tenant_ids = [
    data.hpe_morpheus_tenant.example.id,
  ]
}
