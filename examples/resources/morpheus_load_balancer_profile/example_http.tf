resource "hpe_morpheus_load_balancer_profile" "http" {
  load_balancer_id = 1
  name             = "HTTP Profile"
  description      = "Example HTTP profile"
  service_type     = "LBHttpProfile"

  config_http = {
    http_idle_timeout    = 15
    request_header_size  = 1024
    response_header_size = 4096
    response_timeout     = 60
    https_redirect       = true
    x_forwarded_for      = "INSERT"
  }
}
