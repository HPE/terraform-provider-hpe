# EfficientIP (SOLIDserver) Network Pool Server
#
# EfficientIP is provided by the external morpheus-efficientip-plugin and uses the
# "solidserver" type code. The plugin must be installed for this type to be available.
#
# Applicable attributes for EfficientIP:
#   name, type_code, enabled, service_url, service_username, service_password_wo,
#   service_password_wo_version, credential_id, ignore_ssl, inventory_existing,
#   service_throttle_rate
resource "hpe_morpheus_network_pool_server" "efficientip" {
  name                        = "EfficientIP IPAM"
  type_code                   = "solidserver"
  enabled                     = true
  service_url                 = "https://solidserver.example.com"
  service_username            = "admin"
  service_password_wo         = "changeme"
  service_password_wo_version = 1
  ignore_ssl                  = true
  inventory_existing          = false
  service_throttle_rate       = 0
}
