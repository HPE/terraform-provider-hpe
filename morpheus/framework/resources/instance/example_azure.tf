data "hpe_morpheus_cloud" "azure_cloud" {
  name = "QA Azure"
}

data "hpe_morpheus_instance_type_layout" "azure" {
  name    = "Azure VM"
  version = "22.04"
}

resource "hpe_morpheus_instance" "example" {
  name             = "TestInstance"
  cloud_id         = data.hpe_morpheus_cloud.azure_cloud.id
  layout_id        = data.hpe_morpheus_instance_type_layout.azure.id
  instance_type_id = 9

  group_id = 28
  plan_id  = 622

  instance_context = "dev"
  network_interfaces = [
    {
      network_id = 28
      # Some clouds (such as Azure) provision onto a specific subnet. Use subnet_id
      # instead of network_id (the two, along with network_group_id, are mutually
      # exclusive). The subnet id is available from the Morpheus
      # /api/options/zoneNetworkOptions endpoint, e.g.:
      #   subnet_id = 5
    }
  ]

  volumes = [
    {
      root_volume              = true
      name                     = "root"
      size                     = 100
      storage_type_id          = 23
      datastore_auto_selection = "auto"
    }
  ]

  tags = [
    {
      name  = "terraform"
      value = "true"
    },
    {
      name  = "managed_by"
      value = "terraform"
    }
  ]

  config_azure = {
    resource_pool_id = "pool-12284"
    create_user      = false
    azure_region     = "eastus"
  }

  timeouts = {
    create = "1h"
    delete = "20m"
    update = "20m"
    read   = "10m"
  }
}
