data "hpe_morpheus_images" "example" {
  sort_ascending = true
  source         = "Synced"
  filter {
    name   = "name"
    values = ["Test*"]
  }
}
