data "hpe_morpheus_cloud" "nsxt" {
  name = "NSX-T Cloud"
}

data "hpe_morpheus_group" "example" {
  name = "Example Group"
}

resource "hpe_morpheus_network" "nsxt_segment" {
  name         = "example-terraform-nsxt-segment"
  display_name = "Example NSX-T Segment"
  description  = "NSX-T overlay segment managed by Terraform"
  cloud_id     = data.hpe_morpheus_cloud.nsxt.id
  group_id     = data.hpe_morpheus_group.example.id
  type_id      = 7
  cidr         = "172.16.10.0/24"
  gateway      = "172.16.10.1"
  active       = true
  dhcp_server  = true
  config = {
    "connectedGateway"        = "/infra/tier-1s/my-tier1-gw"
    "vlanIDs"                 = ""
    "subnetIpManagementType"  = "dhcpLocal"
    "subnetDhcpServerAddress" = "172.16.10.2/24"
    "dhcpRange"               = "172.16.10.100-172.16.10.200"
    "subnetDhcpLeaseTime"     = "86400"
  }
}
