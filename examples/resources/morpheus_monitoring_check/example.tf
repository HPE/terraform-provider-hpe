resource "hpe_morpheus_monitoring_check" "example" {
  name           = "Website Health"
  check_type_id  = 1
  description    = "HTTP health check for production website"
  check_interval = 60
  active         = true
  severity       = "critical"
}
