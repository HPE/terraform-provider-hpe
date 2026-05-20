resource "hpe_morpheus_network_dhcp_server" "example" {
  network_integration_id = 16
  name                   = "Example DHCP Server"
  server_ip_address      = "192.168.1.1/24"
  lease_time             = 86400

  config_nsxt = {
    edge_cluster = "qa-edge-cluster-01"
  }
}
