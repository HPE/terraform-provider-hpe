resource "hpe_morpheus_host" "example" {
  name                = "example-host"
  cloud_id            = 15174
  group_id            = 1
  compute_server_type = "omega.baremetal.stub-server"

  config_omega_bmaas = {
    console_host     = "10.0.0.1"
    console_username = "admin"
    num_nics         = 1
    resource_pool_id = "pool-1"
  }
}
