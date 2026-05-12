resource "hpe_morpheus_network_router" "example" {
  name                   = "TestRouter"
  group_id               = 1
  network_integration_id = 1

  config_nsxt_gateway_tier0 = {
    ha_mode      = "ACTIVE_ACTIVE"
    restart_mode = "HELPER_ONLY"
  }
}
