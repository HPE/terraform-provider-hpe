# Bluecat Network Pool Server
#
# Applicable attributes for Bluecat:
#   name, type_id, enabled, service_url, service_username, service_password_wo,
#   service_password_wo_version, credential_id, ignore_ssl, network_filter,
#   service_throttle_rate
resource "hpe_morpheus_network_pool_server" "bluecat" {
  name                        = "Bluecat IPAM"
  type_id                     = 2
  enabled                     = true
  service_url                 = "https://bluecat.example.com/api"
  service_username            = "admin"
  service_password_wo         = "changeme"
  service_password_wo_version = 1
  ignore_ssl                  = false
  network_filter              = "192.168.0.0/16"
  service_throttle_rate       = 50
}
