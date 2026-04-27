data "hpe_morpheus_cloud" "vmware_cloud" {
  name = "QA VMware"
}

data "hpe_morpheus_service_plan" "vmware_512mb" {
  name                = "1 CPU, 1GB Memory"
  provision_type_code = "vmware"
}

data "hpe_morpheus_instance_type_layout" "vmware" {
  name    = "VMware VM"
  version = "22.04"
}

resource "hpe_morpheus_instance" "example" {
  name             = "TestInstance"
  cloud_id         = data.hpe_morpheus_cloud.vmware_cloud.id
  layout_id        = data.hpe_morpheus_instance_type_layout.vmware.id
  instance_type_id = 9

  group_id = 28
  plan_id  = data.hpe_morpheus_service_plan.vmware_512mb.id

  instance_context = "dev"
  network_interfaces = [
    {
      network_id = 86657
    }
  ]

  network_domain_id = 27

  service_plan_options = {
    max_cores = 2
    max_memory = 2048
  }

  volumes = [
    {
      root_volume              = true
      name                     = "root"
      size                     = 10
      storage_type_id          = 1
      datastore_auto_selection = "auto"
    },
    {
      root_volume              = false
      name                     = "data"
      size                     = 10
      storage_type_id          = 1
      datastore_auto_selection = "auto"
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

  config_vmware = {
    resource_pool_id      = "pool-1"
    nested_virtualization = "off"
    no_agent              = true
    create_user           = false
    vmware_folder_id      = "group-v79"
  }

  timeouts = {
    create = "1h"
    delete = "20m"
    update = "20m"
    read   = "10m"
  }
}
