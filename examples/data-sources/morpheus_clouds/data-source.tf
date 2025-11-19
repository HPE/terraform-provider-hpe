data "hpe_morpheus_clouds" "tf_example_clouds" {
  sort_ascending = true
  filter {
    name   = "name"
    values = ["Test*"]
  }
}