data "hpe_morpheus_cloud" "aws_cloud" {
  name = "QA Amazon"
}

data "hpe_morpheus_instance_type_layout" "aws" {
  name    = "Amazon VM"
  version = "22.04"
}

resource "hpe_morpheus_instance" "example" {
  name             = "TestInstance"
  cloud_id         = data.hpe_morpheus_cloud.aws_cloud.id
  layout_id        = data.hpe_morpheus_instance_type_layout.aws.id
  instance_type_id = 9

  group_id = 28
  plan_id  = 622

  instance_context = "dev"
  network_interfaces = [
    {
      network_id = 28
    }
  ]

  volumes = [
    {
      root_volume              = true
      name                     = "root"
      size                     = 100
      storage_type_id          = 23
      datastore_auto_selection = "auto"
    },
    {
      root_volume              = false
      name                     = "data"
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

  config_aws = {
    resource_pool_id      = "pool-12284"
    no_agent              = true
    create_user           = false
    security_groups = [
      { id = "sg-4eaf812b" },
    ]
  }

  timeouts = {
    create = "1h"
    delete = "20m"
    update = "20m"
    read   = "10m"
  }
}
