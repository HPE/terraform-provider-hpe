resource "hpe_morpheus_cluster" "example_generic_hvm" {
  name              = "TestCluster"
  description       = "A HVM cluster created with a dynamic config"
  cloud_id          = 1
  group_id          = 1
  layout_id         = 2
  cluster_type_code = "mvm-cluster"

  labels = [
    "terraform",
    "example",
  ]

  config = {
    createUser           = false
    cpuArch              = "x86_64"
    cpuModel             = "host-model"
    dynamicPlacementMode = "off"
    powerPolicy          = "balanced"
  }

  server = {
    service_plan_id = 1

    ssh_port                 = 22
    ssh_username             = "user"
    ssh_key_pair_id          = 1
    management_net_interface = "eth0"

    ssh_hosts = [
      {
        name = "host1"
        ip   = "10.0.0.1"
      },
      {
        name = "host2"
        ip   = "10.0.0.2"
      },
      {
        name = "host3"
        ip   = "10.0.0.3"
      }
    ]

    visibility = "private"

    tags = [
      {
        name  = "source"
        value = "terraform"
      },
      {
        name  = "environment"
        value = "example"
      },
    ]
  }
}
