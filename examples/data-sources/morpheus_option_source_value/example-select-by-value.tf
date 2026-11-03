# Selecting by value: confirm at plan time that a value the configuration
# already holds is one the source offers, and learn its name.

data "hpe_morpheus_cloud" "hvm" {
  name = "my-hvm-cloud"
}

data "hpe_morpheus_instance_type_layout" "hvm" {
  name    = "Single KVM VM"
  version = "22.04"
}

# If pool-2 is not selectable for this cloud and layout, the plan fails and
# names the pools that are.
data "hpe_morpheus_option_source_value" "pool_2" {
  source    = "zonePools"
  cloud_id  = data.hpe_morpheus_cloud.hvm.id
  layout_id = data.hpe_morpheus_instance_type_layout.hvm.id
  value     = "pool-2"
}

output "pool_2" {
  value = {
    name          = data.hpe_morpheus_option_source_value.pool_2.name
    provider_type = data.hpe_morpheus_option_source_value.pool_2.attributes["providerType"]
    cluster_type  = data.hpe_morpheus_option_source_value.pool_2.attributes["type"]
  }
}
