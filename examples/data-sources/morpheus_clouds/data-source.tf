data "hpe_morpheus_clouds" "example" {
  sort_ascending = true
  filter {
    name   = "name"
    values = ["Test*"]
  }
}
