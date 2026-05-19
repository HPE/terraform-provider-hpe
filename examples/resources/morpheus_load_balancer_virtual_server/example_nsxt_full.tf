resource "hpe_morpheus_load_balancer_virtual_server" "nsxt" {
  load_balancer_id = 1
  vip_name         = "example-nsxt-vip-ssl-client"
  description      = "Example NSX-T virtual server"
  vip_address      = "10.0.0.5"
  vip_port         = 443
  vip_protocol     = "http"
  pool_id          = 11
  ssl_cert         = 12
  ssl_server_cert  = 0

  config_nsxt = {
    application_profile = 13
    persistence         = "COOKIE"
    persistence_profile = 16
    ssl_client_profile  = 19
    ssl_server_profile  = 0
  }
}
