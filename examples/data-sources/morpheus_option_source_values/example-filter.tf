# Narrowing a list with filter blocks: blocks are ANDed, values within a block
# are ORed, and values are Go regular expressions.

data "hpe_morpheus_cloud" "hvm" {
  name = "my-hvm-cloud"
}

data "hpe_morpheus_instance_type_layout" "hvm" {
  name    = "Single KVM VM"
  version = "22.04"
}

# The HVM cluster pools whose name begins "prod-". attributes.<key> filters on
# a source-specific attribute; a value that lacks it does not match.
data "hpe_morpheus_option_source_values" "prod_pools" {
  source    = "zonePools"
  cloud_id  = data.hpe_morpheus_cloud.hvm.id
  layout_id = data.hpe_morpheus_instance_type_layout.hvm.id

  filter {
    name   = "attributes.providerType"
    values = ["^mvm$"]
  }

  filter {
    name   = "name"
    values = ["^prod-"]
  }
}

# The networks in this cloud whose name mentions management, either spelling.
data "hpe_morpheus_option_source_values" "management_networks" {
  source   = "networks"
  cloud_id = data.hpe_morpheus_cloud.hvm.id

  filter {
    name   = "name"
    values = ["(?i)management", "(?i)mgmt"]
  }
}

output "prod_pools" {
  value = { for v in data.hpe_morpheus_option_source_values.prod_pools.values : v.name => v.value }
}

output "management_network_ids" {
  value = data.hpe_morpheus_option_source_values.management_networks.values[*].value
}
