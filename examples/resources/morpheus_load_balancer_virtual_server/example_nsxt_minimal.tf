resource "hpe_morpheus_load_balancer_virtual_server" "nsxt_minimal" {
  load_balancer_id = 1
  vip_name         = "example-nsxt-vip"
  description      = "Example NSX-T virtual server"
  vip_address      = "10.0.0.4"
  vip_port         = 443
  vip_protocol     = "http"
  pool_id          = 11

  config_nsxt = {
    application_profile = 13
    persistence         = ""
  }
}
