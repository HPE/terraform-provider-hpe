resource "hpe_morpheus_network_pool" "example" {
  name           = "App Pool"
  type_id        = 1
  subnet_address = "10.0.1.0"
  netmask        = "255.255.255.0"
  gateway        = "10.0.1.1"
  dns_domain     = "example.com"
}
