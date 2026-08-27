data "hpe_morpheus_load_balancer_virtual_server" "example" {
  load_balancer_id = 1
  vip_name         = "my-web-vs"
}
