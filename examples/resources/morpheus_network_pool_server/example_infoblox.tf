# Infoblox Network Pool Server
#
# Applicable attributes for Infoblox:
#   name, type_id, enabled, service_url, service_username, service_password_wo,
#   service_password_wo_version, credential_id, ignore_ssl, network_filter,
#   zone_filter, tenant_match, service_mode, service_throttle_rate
resource "hpe_morpheus_network_pool_server" "infoblox" {
  name                        = "Infoblox IPAM"
  type_id                     = 1
  enabled                     = true
  service_url                 = "https://infoblox.example.com/wapi/v2.12"
  service_username            = "admin"
  service_password_wo         = "changeme"
  service_password_wo_version = 1
  ignore_ssl                  = true
  network_filter              = "10.0.0.0/8"
  zone_filter                 = "example.com"
  tenant_match                = ".*"
  service_mode                = "static"
  service_throttle_rate       = 0
}
