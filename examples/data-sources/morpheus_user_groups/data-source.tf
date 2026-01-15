data "hpe_morpheus_user_groups" "example" {
  sort_ascending = true

  filter {
    name   = "name"
    values = [".*"]
  }
}
