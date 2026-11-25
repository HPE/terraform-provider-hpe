data "hpe_morpheus_network_floating_ips" "example" {
  ip_status = "assigned"

  filter {
    name   = "cloud_name"
    values = ["prod-.*"]
  }
}
