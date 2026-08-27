data "hpe_morpheus_load_balancer_monitor" "example" {
  load_balancer_id = 1
  name             = "HTTP Monitor"
}
