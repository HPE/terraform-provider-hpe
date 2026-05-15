resource "hpe_morpheus_load_balancer_virtual_server" "nsxt" {
  load_balancer_id = 1
  vip_name         = "example-nsxt-vip"
  description      = "Example NSX-T virtual server with full SSL and persistence"
  vip_address      = "10.0.0.2"
  vip_port         = 443
  vip_protocol     = "http"
  pool_id          = 42
  ssl_cert         = 12
  ssl_server_cert  = 0

  config_nsxt = {
    application_profile = 85
    persistence         = "SOURCE_IP"
    persistence_profile = 78
    ssl_client_profile  = 33
    ssl_server_profile  = 0
  }
}
