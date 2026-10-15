# Everything the instance is built from is looked up by name. Ids are allocated
# per appliance, so a configuration that carries one only ever works on the
# appliance it was read from; a name resolves wherever the object exists.

data "hpe_morpheus_cloud" "vme_cloud" {
  name = "MyCloud"
}

# The cluster to provision into. Its resource pool is where the instance goes
# (config_hvm.resource_pool_id), and its networks and datastores are where the
# names below are looked up: every HVM cluster has a Management network and a
# local datastore, so a name alone is ambiguous as soon as a cloud has two.
data "hpe_morpheus_cluster" "vme" {
  name = "MyCluster"
}

# The layout resolves its instance type with it, so the two cannot be mismatched.
# Layout names repeat once per guest version; the version pins one down.
data "hpe_morpheus_instance_type_layout" "vme" {
  name    = "Single HVM"
  version = "1.0"
}

data "hpe_morpheus_network" "vme" {
  name       = "MyNetwork"
  cluster_id = data.hpe_morpheus_cluster.vme.id
}

data "hpe_morpheus_datastore" "vme" {
  name       = "MyDatastore"
  cluster_id = data.hpe_morpheus_cluster.vme.id
}

data "hpe_morpheus_service_plan" "vme_plan" {
  name                = "1 CPU, 1GB Memory"
  provision_type_code = "kvm"
}

data "hpe_morpheus_group" "example" {
  name = "TestGroup"
}

resource "hpe_morpheus_instance" "example" {
  name             = "TestInstance"
  cloud_id         = data.hpe_morpheus_cloud.vme_cloud.id
  layout_id        = data.hpe_morpheus_instance_type_layout.vme.id
  instance_type_id = data.hpe_morpheus_instance_type_layout.vme.instance_type.id

  group_id = data.hpe_morpheus_group.example.id
  plan_id  = data.hpe_morpheus_service_plan.vme_plan.id

  instance_context = "dev"
  network_interfaces = [
    {
      network_id = data.hpe_morpheus_network.vme.id
    }
  ]

  # storage_type_id 1 is the seeded "standard" disk type, the same on every
  # appliance; it is a type, not an object, and has no data source.
  volumes = [
    {
      root_volume     = true
      name            = "root"
      size            = 10
      storage_type_id = 1
      datastore_id    = data.hpe_morpheus_datastore.vme.id
    },
    {
      root_volume     = false
      name            = "data"
      size            = 10
      storage_type_id = 1
      datastore_id    = data.hpe_morpheus_datastore.vme.id
    }
  ]

  tags = [
    {
      name  = "terraform"
      value = "true"
    },
    {
      name  = "acctest"
      value = "true"
    },
    {
      name  = "hpe_morpheus_instance"
      value = "true"
    },
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
    resource_pool_id      = "pool-${data.hpe_morpheus_cluster.vme.permissions.resource_pool.id}"
    nested_virtualization = "off"
    no_agent              = false
    create_user           = false
  }
}
