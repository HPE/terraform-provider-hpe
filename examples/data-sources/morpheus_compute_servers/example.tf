data "hpe_morpheus_compute_servers" "example" {
  instance_id = 1

  filter {
    name   = "status"
    values = ["provisioned"]
  }
}
