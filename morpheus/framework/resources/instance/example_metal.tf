data "hpe_morpheus_cloud" "test" {
  name = "aCloud"
}

data "hpe_morpheus_environment" "test" {
  name = "anEnvironment"
}

data "hpe_morpheus_group" "test" {
  name = "aGroup"
}

data "hpe_morpheus_instance_type_layout" "test" {
  name = "Single ILO Server"
}

data "hpe_morpheus_role" "test" {
  name = "aRole"
}

data "hpe_morpheus_service_plan" "tp" {
    name                = "G3i"
    provision_type_code = "hpe-baremetal-plugin.provision"
}

resource "hpe_morpheus_instance" "example" {
  name             = "TestInstance"
  cloud_id         = data.hpe_morpheus_cloud.test.id # BM Cloud
  layout_id        = data.hpe_morpheus_instance_type_layout.test.id
  instance_type_id = 56 # BM Instance

  group_id = data.hpe_morpheus_group.test.id
  plan_id  = data.hpe_morpheus_service_plan.tp.id

  instance_context = "dev"

  network_interfaces = [
    {
      network_id      = 21
      ipMode          = ""
      network_type_id = 18
    },
    {
      network_id      = 21
      ipMode          = ""
      network_type_id = 18
    }
  ]

  volumes = [
    {
      root_volume     = true
      name            = "root"
      size            = 0
      storage_type_id = 76
      datastore_id    = null
    },
    {
      root_volume     = false
      name            = "data"
      size            = 16 # GB
      storage_type_id = 84
      datastore_id    = 11
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
      name  = "mytag"
      value = "true"
    },
    {
      name  = "jesskey"
      value = "terraform"
    }
  ]

  config_bmaas = {
    image_id                 = 231
    resource_pool_id         = "pool-1"
    no_agent                 = true
    create_user              = false
    enforce_raid_boot_volume = true

    # Optionally pin specific bare metal host id(s) instead of auto-selecting
    # from the resource pool:
    # selected_hosts = [155]
  }

  timeouts = {
    create = "2h"
  }
}