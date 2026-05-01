data "hpe_morpheus_load_balancer_virtual_server" "example" {
  vip_name         = "Example virtual server"
  load_balancer_id = 1
}
