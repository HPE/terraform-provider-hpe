resource "hpe_morpheus_load_balancer" "lb" {
  name       = "example-terraform-nsxt-lb"
  type_code  = "nsx-t"
  visibility = "public"
  network_server_id = 5

  config_nsxt = {
    admin_state   = true
    log_level     = "INFO"
    size          = "SMALL"
    tier1_gateway = "tier1-gateway"
  }
}
