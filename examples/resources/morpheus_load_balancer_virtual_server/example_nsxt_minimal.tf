resource "hpe_morpheus_load_balancer_virtual_server" "nsxt_minimal" {
  load_balancer_id = 1
  vip_name         = "example-nsxt-minimal"
  description      = "Minimal NSX-T virtual server"
  vip_address      = "10.0.0.3"
  vip_port         = 80
  vip_protocol     = "http"
  pool_id          = 42

  config_nsxt = {
    application_profile = 85
    persistence         = ""
  }
}
