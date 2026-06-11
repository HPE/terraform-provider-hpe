resource "hpe_morpheus_subnet" "example" {
  name       = "Example Subnet"
  type_id    = 1
  network_id = 1
  visibility = "private"
}
