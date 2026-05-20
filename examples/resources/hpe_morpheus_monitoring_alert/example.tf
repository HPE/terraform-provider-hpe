resource "hpe_morpheus_monitoring_alert" "example" {
  name         = "High Severity Alert"
  min_severity = "critical"
  min_duration = 5
  active       = true
  all_checks   = true
}
