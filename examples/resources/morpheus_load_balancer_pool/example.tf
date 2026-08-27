resource "hpe_morpheus_load_balancer_pool" "nsxt" {
  load_balancer_id = 1
  name             = "NSX-T Pool"
  description      = "An NSX-T load balancer pool"
  vip_balance      = "ROUND_ROBIN"
  min_active       = 1

  config_nsxt = {
    snat_translation_type   = "LBSnatAutoMap"
    tcp_multiplexing        = true
    tcp_multiplexing_number = 6
  }
}
