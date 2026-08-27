resource "hpe_morpheus_network_router_route" "example" {
  router_id     = 42
  name          = "example-route"
  source        = "10.0.0.0/24"
  destination   = "10.0.0.1"
  description   = "Example route"
  network_mtu   = 1500
  enabled       = true
  default_route = false
}
