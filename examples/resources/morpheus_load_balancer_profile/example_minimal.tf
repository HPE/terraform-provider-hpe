resource "hpe_morpheus_load_balancer_profile" "fast_tcp" {
  load_balancer_id = 1
  name             = "Fast TCP Profile"
  service_type     = "LBFastTcpProfile"

  config_fast_tcp = {
    fast_tcp_idle_timeout    = 1800
    connection_close_timeout = 8
  }
}
