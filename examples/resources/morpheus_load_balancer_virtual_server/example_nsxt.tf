resource "hpe_morpheus_load_balancer_virtual_server" "nsxt" {
  load_balancer_id = 1
  vip_name         = "example-nsxt-vip"
  description      = "Example NSX-T virtual server"
  vip_address      = "10.0.0.2"
  vip_port         = 443
  vip_protocol     = "https"

  config_nsxt = {
    application_profile = "/infra/lb-app-profiles/default-http-lb-app-profile"
  }
}
