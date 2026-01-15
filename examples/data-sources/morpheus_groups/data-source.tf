data "hpe_morpheus_groups" "example" {
  sort_ascending = false
  filter {
    name   = "location"
    values = ["denver"]
  }
}
