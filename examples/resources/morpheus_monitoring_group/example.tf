resource "hpe_morpheus_monitoring_group" "example" {
  name        = "Production Services"
  description = "Monitoring group for production services"
  min_happy   = 1
  severity    = "critical"
  active      = true
}
