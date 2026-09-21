data "hpe_morpheus_cloud" "vme_cloud" {
  name = "HPE Alletra VME"
}

data "hpe_morpheus_service_plan" "vme_512mb" {
  name                = "1 CPU, 1GB Memory"
  provision_type_code = "kvm"
}

resource "hpe_morpheus_instance" "example" {
  name             = "TestInstance"
  cloud_id         = data.hpe_morpheus_cloud.vme_cloud.id # HPE Alletra VME
  layout_id        = 5385                                 # Single KVM VM
  instance_type_id = 9                  # (HVM) mvm-cluster

  group_id = 1
  plan_id  = data.hpe_morpheus_service_plan.vme_512mb.id # kvm-vm-512

  instance_context = "dev"

  # Environment variables can be added, changed and removed in place. Each
  # variable takes a name and value; export and masked are optional and default
  # to true and false respectively when omitted.
  evars = [
    {
      name  = "APP_ENV"
      value = "production"
    },
    {
      # Hidden in the Morpheus UI and API.
      name   = "DB_PASSWORD"
      value  = "s3cr3t"
      masked = true
    },
    {
      # Kept on the instance record but not exported to the environment.
      name   = "BUILD_NOTE"
      value  = "internal"
      export = false
    }
  ]

  # Ports can only be set when the instance is created. Changing or removing
  # them forces a new instance to be created.
  ports = [
    {
      name = "http"
      port = 8080
    }
  ]

  network_interfaces = [
    {
      network_id = 103481
    }
  ]

  volumes = [
    {
      root_volume     = true
      name            = "root"
      size            = 10
      storage_type_id = 1
      datastore_id    = 38658
    }
  ]

  config_hvm = {
    resource_pool_id      = "pool-62299"
    nested_virtualization = "off"
    no_agent              = true
    create_user           = false
  }
}
