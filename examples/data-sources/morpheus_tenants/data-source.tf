data "hpe_morpheus_tenants" "example" {
  sort_ascending = true
  filter {
    name   = "name"
    values = [".*"]
  }
}
