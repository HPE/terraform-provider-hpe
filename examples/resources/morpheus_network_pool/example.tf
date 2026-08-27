resource "hpe_morpheus_network_pool" "example" {
  name           = "App Pool"
  type_code      = "morpheus"
  subnet_address = "10.0.1.0"
  netmask        = "255.255.255.0"
  gateway        = "10.0.1.1"
  dns_domain     = "example.com"

  ip_ranges = {
    starting_address = "10.0.1.10"
    ending_address   = "10.0.1.50"
  }
}
