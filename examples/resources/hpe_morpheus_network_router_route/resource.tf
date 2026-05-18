resource "hpe_morpheus_network_router_route" "example" {
  router_id     = 42
  name          = "example-route"
  network       = "10.0.0.0/24"
  next_hop      = "10.0.0.1"
  description   = "Example route"
  mtu           = 1500
  enabled       = true
  default_route = false
}
