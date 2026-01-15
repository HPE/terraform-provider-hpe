data "hpe_morpheus_networks" "example" {
  cloud_id       = 3
  sort_ascending = true
  filter {
    name   = "name"
    values = ["Test*"]
  }
}
