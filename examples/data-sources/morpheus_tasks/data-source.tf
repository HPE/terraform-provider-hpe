data "hpe_morpheus_tasks" "example" {
  sort_ascending = true

  filter {
    name   = "name"
    values = [".*"]
  }

  filter {
    name   = "type"
    values = ["Shell Script", "Python Script"]
  }
}
