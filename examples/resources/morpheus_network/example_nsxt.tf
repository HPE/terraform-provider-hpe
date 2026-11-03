data "hpe_morpheus_cloud" "nsxt" {
  name = "NSX-T Cloud"
}

data "hpe_morpheus_group" "example" {
  name = "Example Group"
}

# The NSX-T integration the segment is created through.
data "hpe_morpheus_network_server" "nsxt" {
  name = "Labs NSX"
}

resource "hpe_morpheus_network" "nsxt_segment" {
  name              = "example-terraform-nsxt-segment"
  display_name      = "Example NSX-T Segment"
  description       = "NSX-T overlay segment managed by Terraform"
  cloud_id          = data.hpe_morpheus_cloud.nsxt.id
  group_id          = data.hpe_morpheus_group.example.id
  type_code         = "nsxtLogicalSwitch"
  network_server_id = data.hpe_morpheus_network_server.nsxt.id
  cidr              = "172.16.10.0/24"
  gateway           = "172.16.10.1"
  active            = true
  dhcp_server       = true

  # Write-only: sent to Morpheus, never stored in state. Read back in config_from_api.
  config_nsxt = {
    connected_gateway          = "/infra/tier-1s/my-tier1-gw"
    vlan_ids                   = ""
    subnet_ip_management_type  = "dhcpLocal"
    subnet_dhcp_server_address = "172.16.10.2/24"
    dhcp_range                 = "172.16.10.100-172.16.10.200"
    subnet_dhcp_lease_time     = 86400
  }
}
