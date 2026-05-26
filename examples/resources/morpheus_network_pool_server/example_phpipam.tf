# phpIPAM Network Pool Server
#
# Applicable attributes for phpIPAM:
#   name, type_id, enabled, service_url, service_username, service_password_wo,
#   service_password_wo_version, credential_id, ignore_ssl, network_filter,
#   service_throttle_rate
resource "hpe_morpheus_network_pool_server" "phpipam" {
  name                        = "phpIPAM"
  type_id                     = 3
  enabled                     = true
  service_url                 = "https://phpipam.example.com/api/app"
  service_username            = "admin"
  service_password_wo         = "changeme"
  service_password_wo_version = 1
  ignore_ssl                  = false
  network_filter              = "172.16.0.0/12"
  service_throttle_rate       = 0
}
