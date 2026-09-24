data "hpe_morpheus_cloud" "openstack" {
  name = "OpenStack Cloud"
}

data "hpe_morpheus_group" "example" {
  name = "Example Group"
}

data "hpe_morpheus_network_server" "neutron" {
  name = "OpenStack Cloud"
}

resource "hpe_morpheus_network" "openstack" {
  name              = "example-terraform-openstack"
  cloud_id          = data.hpe_morpheus_cloud.openstack.id
  group_id          = data.hpe_morpheus_group.example.id
  type_code         = "openstackPrivate"
  zone_pool_id      = 46930
  network_server_id = data.hpe_morpheus_network_server.neutron.id
  cidr              = "192.168.177.0/24"
  gateway           = "192.168.177.1"
  dhcp_server       = true
  visibility        = "private"
  tenant_ids        = [1]

  # Write-only. The DHCP allocation pool is updated in place; host routes
  # cannot be changed after create.
  config_openstack = {
    dhcp_range = "192.168.177.10,192.168.177.100"
  }
}
