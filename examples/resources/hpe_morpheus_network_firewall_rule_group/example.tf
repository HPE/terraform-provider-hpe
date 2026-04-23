resource "hpe_morpheus_network_firewall_rule_group" "example" {
  network_server_id = 128
  name              = "Example Firewall Rule Group"
  external_type     = "SecurityPolicy"
  description       = "An example firewall rule group"
  priority          = 100
  group_layer       = "Application"
}
