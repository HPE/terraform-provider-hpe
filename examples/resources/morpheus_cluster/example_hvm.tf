resource "hpe_morpheus_cluster" "example_hvm" {
  name        = "TestCluster"
  description = "A test HVM cluster"
  cloud_id    = 1
  group_id    = 1
  layout_id   = 2

  # Cluster provisioning builds every worker node, so it can comfortably run
  # past the 45m default on a busy or lower-spec appliance. Raise this if
  # creation is timing out while the cluster is still reported "provisioning".
  timeouts = {
    create = "55m"
  }

  labels = [
    "terraform",
    "example",
  ]

  config_hvm = {
    create_user       = false
    dynamic_placement = false
    cpu_arch          = "x86_64"
    cpu_model         = "host-model"
    power_policy      = "balanced"
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
