resource "hpe_morpheus_network_firewall_rule" "example" {
  network_server_id = 1
  name              = "Example Firewall Rule"
  direction         = "Ingress"
  policy            = "Accept"
  enabled           = true

  rule_group_id = {
    id = 1
  }
}
