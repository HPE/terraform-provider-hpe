resource "hpe_morpheus_network_pool_server" "example" {
  name             = "InfoBlox IPAM"
  type_id          = 1
  service_url      = "https://ipam.example.com/api"
  service_username = "admin"
  service_password = "changeme"
  enabled          = true
}
