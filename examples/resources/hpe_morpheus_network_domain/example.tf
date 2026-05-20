resource "hpe_morpheus_network_domain" "example" {
  name        = "example.com"
  description = "Primary DNS domain"
  public_zone = true
  active      = true
}
