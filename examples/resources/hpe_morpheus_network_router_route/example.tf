resource "hpe_morpheus_network_router_route" "example" {
  router_id   = 1
  source      = "10.0.0.0/24"
  destination = "192.168.1.0/24"
}
