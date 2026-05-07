resource "hpe_morpheus_network_router" "example" {
  name                   = "TestRouter"
  type_id                = 1
  group_id               = 1
  network_integration_id = 1

  config_nsx_gateway_tier1 = {
    ip_management_type = "dhcpLocal"
  }
}
