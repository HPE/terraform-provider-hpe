resource "hpe_morpheus_network_router_nat" "example" {
  router_id          = 1
  name               = "Example NAT Rule"
  action             = "SNAT"
  source_network     = "10.0.0.0/24"
  translated_network = "192.168.1.1"
  description        = "Example SNAT rule"
}
