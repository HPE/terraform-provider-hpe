resource "hpe_morpheus_load_balancer_monitor" "nsxv" {
  load_balancer_id            = 1
  name                        = "NSX-V HTTP Monitor"
  description                 = "An NSX-V HTTP health check monitor"
  monitor_type                = "http"
  monitor_interval            = 10
  monitor_timeout             = 15
  max_retry                   = 3
  send_data                   = "GET / HTTP/1.0"
  send_type                   = "GET"
  receive_data                = ""
  receive_code                = "200"
  monitor_destination         = "/health"
  monitor_username            = ""
  monitor_password_wo         = ""
  monitor_password_wo_version = 1
}
