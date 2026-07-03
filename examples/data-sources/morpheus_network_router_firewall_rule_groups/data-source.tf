# List all firewall rule groups for a given network router.
data "hpe_morpheus_network_router_firewall_rule_groups" "example" {
  router_id = 5
}

# Access individual rule group attributes.
output "rule_group_names" {
  value = [for rg in data.hpe_morpheus_network_router_firewall_rule_groups.example.rule_groups : rg.name]
}
