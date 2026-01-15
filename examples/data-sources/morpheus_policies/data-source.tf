data "hpe_morpheus_policies" "example" {
  sort_ascending = true

  filter {
    name   = "name"
    values = [".*"]
  }

  filter {
    name   = "type"
    values = ["Max VMs", "Workflow"]
  }
}
