# HPE Morpheus SDN (Juniper CN2) HTTP health-check monitor
#
# The SDN load balancer's HTTP monitor carries its check in the typed config
# block: the HTTP method, URI, port and expected return code. The controller
# requires a timeout and a retry count greater than zero, so set
# monitor_timeout and max_retry explicitly.
resource "hpe_morpheus_load_balancer_monitor" "sdn" {
  load_balancer_id = 1
  name             = "SDN HTTP Monitor"
  description      = "An HPE Morpheus SDN HTTP health check monitor"
  monitor_type     = "HTTP"
  monitor_interval = 5
  monitor_timeout  = 15
  max_retry        = 3

  config = {
    http_method      = "GET"
    http_uri         = "/health"
    http_port        = 8080
    http_return_code = "200"
  }
}
