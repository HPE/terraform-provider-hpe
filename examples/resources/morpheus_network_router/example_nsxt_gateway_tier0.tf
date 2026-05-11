resource "hpe_morpheus_network_router" "example" {
  name                   = "TestRouter"
  type_id                = 1
  group_id               = 1
  network_integration_id = 1

  config_nsxt_gateway_tier0 = {
    ha_mode      = "ACTIVE_ACTIVE"
    ip_server_id = 1
    restart_mode = "HELPER_ONLY"
  }
}
