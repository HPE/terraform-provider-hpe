# SolarWinds Network Pool Server
#
# Applicable attributes for SolarWinds:
#   name, type_id, enabled, service_url, service_username, service_password_wo,
#   service_password_wo_version, credential_id, ignore_ssl, service_throttle_rate
resource "hpe_morpheus_network_pool_server" "solarwinds" {
  name                        = "SolarWinds IPAM"
  type_id                     = 4
  enabled                     = true
  service_url                 = "https://solarwinds.example.com:17778/SolarWinds/InformationService/v3/Json"
  service_username            = "admin"
  service_password_wo         = "changeme"
  service_password_wo_version = 1
  ignore_ssl                  = true
  service_throttle_rate       = 100
}
