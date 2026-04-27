data "hpe_morpheus_cloud" "vme_cloud" {
  name = "MyCloud"
}

data "hpe_morpheus_service_plan" "vme_512mb" {
  name                = "1 CPU, 1GB Memory"
  provision_type_code = "kvm"
}

resource "hpe_morpheus_instance" "example" {
  name             = "TestInstance"
  cloud_id         = data.hpe_morpheus_cloud.vme_cloud.id 
  layout_id        = 644 # Single KVM VM
  instance_type_id = 9 # (HVM) mvm-cluster

  group_id = 1
  plan_id  = data.hpe_morpheus_service_plan.vme_512mb.id # kvm-vm-512

  instance_context = "dev"
  network_interfaces = [
    {
      network_id = 755
    }
  ]

  volumes = [
    {
      root_volume     = true
      name            = "root"
      size            = 10
      storage_type_id = 1
      datastore_id    = 555
    },
    {
      root_volume     = false
      name            = "data"
      size            = 10
      storage_type_id = 1
      datastore_id    = 555
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
    resource_pool_id      = "pool-700"
    nested_virtualization = "off"
    no_agent              = true
    create_user           = false
  }
}
