# Network Pool Server using a stored credential
#
# Instead of providing service_username and service_password_wo inline,
# reference a stored credential by ID. credential_id conflicts with
# service_username and service_password_wo.
resource "hpe_morpheus_network_pool_server" "with_credential" {
  name          = "Infoblox with Credential"
  type_id       = 1
  enabled       = true
  service_url   = "https://infoblox.example.com/wapi/v2.12"
  credential_id = 42
  ignore_ssl    = true
}
