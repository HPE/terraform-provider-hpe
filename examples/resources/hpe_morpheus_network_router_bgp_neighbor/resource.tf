resource "hpe_morpheus_network_router_bgp_neighbor" "example" {
  router_id            = 42
  ip_address           = "10.0.0.1"
  description          = "Example BGP neighbor"
  remote_as            = "65001"
  weight               = 100
  keep_alive           = 60
  hold_down            = 180
}
