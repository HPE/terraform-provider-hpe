# Checking a network before provisioning: zoneNetworkOptions is the check
# Morpheus applies at create time, read here into a precondition.

data "hpe_morpheus_cloud" "hvm" {
  name = "my-hvm-cloud"
}

data "hpe_morpheus_instance_type_layout" "hvm" {
  name    = "Single KVM VM"
  version = "22.04"
}

data "hpe_morpheus_instance_type" "ubuntu" {
  name = "Ubuntu"
}

data "hpe_morpheus_service_plan" "small" {
  name                = "1 CPU, 1GB Memory"
  provision_type_code = "kvm"
}

# The cluster's pool (see hpe_morpheus_option_source_value).
data "hpe_morpheus_option_source_value" "cluster_pool" {
  source    = "zonePools"
  cloud_id  = data.hpe_morpheus_cloud.hvm.id
  layout_id = data.hpe_morpheus_instance_type_layout.hvm.id
  name      = "hvm-cluster-01"
}

# The networks selectable for that pool. zoneNetworkOptions returns an object of
# lists, so collection names the one to read. Each value is the network id in
# the form the API takes on input: "network-<id>".
data "hpe_morpheus_option_source_values" "eligible_networks" {
  source     = "zoneNetworkOptions"
  cloud_id   = data.hpe_morpheus_cloud.hvm.id
  layout_id  = data.hpe_morpheus_instance_type_layout.hvm.id
  pool_id    = data.hpe_morpheus_option_source_value.cluster_pool.value
  collection = "networks"
}

resource "hpe_morpheus_instance" "vm" {
  name             = "web-01"
  cloud_id         = data.hpe_morpheus_cloud.hvm.id
  group_id         = 1
  instance_type_id = data.hpe_morpheus_instance_type.ubuntu.id
  layout_id        = data.hpe_morpheus_instance_type_layout.hvm.id
  plan_id          = data.hpe_morpheus_service_plan.small.id

  network_interfaces = [
    { network_id = var.network_id }
  ]

  config_hvm = {
    resource_pool_id = data.hpe_morpheus_option_source_value.cluster_pool.value
  }

  lifecycle {
    precondition {
      condition = contains(
        data.hpe_morpheus_option_source_values.eligible_networks.values[*].value,
        "network-${var.network_id}",
      )
      error_message = join(" ", [
        "Network ${var.network_id} is not selectable for ${data.hpe_morpheus_option_source_value.cluster_pool.value}",
        "(${data.hpe_morpheus_option_source_value.cluster_pool.name}) in this cloud. Selectable networks:",
        join(", ", data.hpe_morpheus_option_source_values.eligible_networks.values[*].name),
      ])
    }
  }
}
