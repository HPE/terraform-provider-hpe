resource "hpe_morpheus_cluster" "example_hvm" {
  name        = "TestCluster"
  description = "A test HVM cluster"
  cloud_id    = 1
  group_id    = 1
  layout_id   = 2

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

    # HVM 2.0 network configuration. Only honoured when the selected cluster
    # layout has clusterVersion 2.0 or greater; otherwise Morpheus silently
    # drops these fields.
    virt_switch_name = "vSwitch0"
    dns1             = "8.8.8.8"
    dns2             = "8.8.4.4"
    dns3             = "1.1.1.1"
    search_domain    = "example.com"
    ntp1             = "0.pool.ntp.org"
    ntp2             = "1.pool.ntp.org"
    proxy            = "http://proxy.example.com:3128"
    no_proxy         = "localhost,127.0.0.1,.example.com"
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
