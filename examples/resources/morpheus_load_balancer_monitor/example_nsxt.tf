resource "hpe_morpheus_load_balancer_monitor" "nsxt" {
  load_balancer_id    = 1
  name                = "NSX-T HTTP Monitor"
  description         = "An NSX-T HTTP health check monitor"
  monitor_type        = "http"
  monitor_interval    = 5
  monitor_timeout     = 15
  monitor_destination = "/"
  fall_count          = 3
  rise_count          = 3
  alias_port          = 8080
  send_data           = "GET / HTTP/1.1"
  send_type           = "GET"
  send_version        = "HTTP_VERSION_1_1"
  receive_data        = ""
  receive_code        = "200"
  data_length         = 0
  max_retry           = 3
}
