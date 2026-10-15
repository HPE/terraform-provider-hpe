data "hpe_morpheus_cloud" "hvm_cloud" {
  name = "hvm"
}

# The cluster to provision into: its resource pool is the instance's, and its
# networks are where the network name is looked up (every HVM cluster has a
# Management network of its own, so a name alone is ambiguous as soon as a
# cloud has two clusters).
data "hpe_morpheus_cluster" "hvm" {
  name = "MyCluster"
}

data "hpe_morpheus_group" "hvm" {
  name = "TestGroup"
}

data "hpe_morpheus_service_plan" "hvm_plan" {
  name                = "1 CPU, 512MB Memory"
  provision_type_code = "kvm"
}

# The layout name is not unique on its own -- appliances carry one per guest
# version -- so the version is needed to resolve a single layout. The layout
# resolves its instance type with it.
data "hpe_morpheus_instance_type_layout" "hvm" {
  name    = "Single KVM VM"
  version = "11"
}

data "hpe_morpheus_network" "hvm" {
  name       = "MyNetwork"
  cluster_id = data.hpe_morpheus_cluster.hvm.id
}

# The image to provision from. Looking it up by name keeps the configuration
# portable between appliances, where the numeric id differs.
data "hpe_morpheus_image" "hvm" {
  name = "morpheus-central"
}

resource "hpe_morpheus_instance" "example" {
  name             = "TestInstance"
  cloud_id         = data.hpe_morpheus_cloud.hvm_cloud.id
  layout_id        = data.hpe_morpheus_instance_type_layout.hvm.id
  instance_type_id = data.hpe_morpheus_instance_type_layout.hvm.instance_type.id

  group_id = data.hpe_morpheus_group.hvm.id
  plan_id  = data.hpe_morpheus_service_plan.hvm_plan.id

  instance_context = "dev"

  network_interfaces = [
    {
      network_id = data.hpe_morpheus_network.hvm.id
    }
  ]

  volumes = [
    {
      root_volume              = true
      name                     = "root"
      size                     = 10
      storage_type_id          = 1
      datastore_auto_selection = "auto"
    }
  ]

  tags = [
    {
      name  = "sweepable"
      value = "true"
    },
    {
      name  = "managed_by"
      value = "terraform"
    }
  ]

  config_hvm = {
    resource_pool_id      = "pool-${data.hpe_morpheus_cluster.hvm.permissions.resource_pool.id}"
    nested_virtualization = "off"
    no_agent              = false
    create_user           = false

    # image_id overrides the image configured on the layout. Omit it to take the
    # layout default. Changing it replaces the instance, because the image is
    # only applied at provision time.
    image_id = data.hpe_morpheus_image.hvm.id
  }

  timeouts = {
    create = "1h"
    delete = "20m"
    update = "20m"
    read   = "10m"
  }
}
