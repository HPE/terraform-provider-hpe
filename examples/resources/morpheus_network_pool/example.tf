resource "hpe_morpheus_network_pool" "example" {
  name           = "App Pool"
  type_code      = "morpheus"
  subnet_address = "10.0.1.0"
  netmask        = "255.255.255.0"
  gateway        = "10.0.1.1"
  dns_domain     = "example.com"

  # Every range of the pool, in the order you want them listed. Each range is
  # edited in place (by its id) when its addresses change.
  ip_address_ranges = [
    {
      starting_address = "10.0.1.10"
      ending_address   = "10.0.1.50"
    },
    {
      starting_address = "10.0.1.100"
      ending_address   = "10.0.1.150"
    },
  ]
}
