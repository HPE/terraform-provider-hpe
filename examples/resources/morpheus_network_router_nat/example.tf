resource "hpe_morpheus_network_router_nat" "example" {
  router_id      = 1
  name           = "Example NAT Rule"
  source_network = "10.0.0.0/24"
  description    = "Example SNAT rule"
}
