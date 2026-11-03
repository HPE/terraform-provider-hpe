data "hpe_morpheus_cloud" "hvm" {
  name = "HVM Cloud"
}

data "hpe_morpheus_group" "example" {
  name = "Example Group"
}

# An HVM data network attaches a VLAN on a host interface to the cluster. Its
# cluster is the resource pool; its network server is the cluster's own, which
# is hidden from the ordinary listing (GET /api/networks/servers?showHidden=true).
resource "hpe_morpheus_network" "hvm_data" {
  name              = "example-terraform-hvm-data"
  cloud_id          = data.hpe_morpheus_cloud.hvm.id
  group_id          = data.hpe_morpheus_group.example.id
  type_code         = "hvmPortGroupDataNetwork"
  zone_pool_id      = 153047
  network_server_id = 7834
  vlan_id           = 2104
  visibility        = "private"
  tenant_ids        = [1]

  # Write-only. The layout and interface cannot be changed after create; the
  # MTU can (1500, or 9000 for jumbo frames).
  config_hvm_data = {
    cluster_network_layout = "HOST_INTERFACE"
    host_network_interface = "bond0"
    mtu                    = 1500
  }
}
