# HPE Morpheus SDN BGP router (juniperBGPSdnRouter)
#
# A BGP peering between the SDN controller and up to four external peers. The
# router type is selected by the presence of the config_juniper_bgp block; the
# SDN network server is given through network_integration_id. The peering's
# MD5 key is write-only (md5_key_wo): it is never stored in state or returned
# by the API, and a change is detected from a digest kept in private state.
resource "hpe_morpheus_network_router" "example" {
  name                   = "TestRouter"
  group_id               = 1
  network_integration_id = 1

  config_juniper_bgp = {
    vendor           = "juniper"
    asn              = 65000
    peer_address_1   = "10.0.0.1"
    address_families = ["inet"]
    hold_time        = 90
  }
}
