# NSX-T Network Server
#
# A generic network controller integration. Its type-specific option values go
# in the free-form config map; the common connection fields are attributes. A
# non-SDN network server is attached to the cloud it serves (cloud_id).
resource "hpe_morpheus_network_server" "nsxt" {
  name                = "NSX-T"
  type_code           = "nsx-t"
  cloud_id            = 1
  service_url         = "https://nsxt.example.com"
  service_username    = "admin"
  service_password_wo = "example"
}

# Using a stored credential instead of a username and password:
#
# resource "hpe_morpheus_network_server" "nsxt" {
#   name          = "NSX-T"
#   type_code     = "nsx-t"
#   cloud_id      = 1
#   service_url   = "https://nsxt.example.com"
#   credential_id = 42
# }
