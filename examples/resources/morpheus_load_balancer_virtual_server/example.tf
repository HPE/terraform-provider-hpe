resource "hpe_morpheus_load_balancer_virtual_server" "example" {
  load_balancer_id = 1
  vip_name         = "example-vip"
  description      = "Example virtual server"
  vip_address      = "10.0.0.1"
  vip_port         = 80
  vip_protocol     = "http"
}
