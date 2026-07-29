# Firewall rules attach to a rule group on the router. Look the group up and
# use its external_id as the rule's parent_id.
data "hpe_morpheus_network_router_firewall_rule_group" "example" {
  router_id = 1
  name      = "Example Firewall Rule Group"
}

resource "hpe_morpheus_network_router_firewall_rule" "example" {
  router_id = 1
  parent_id = data.hpe_morpheus_network_router_firewall_rule_group.example.external_id
  name      = "Example Firewall Rule"
  policy    = "accept"
  enabled   = true
}
