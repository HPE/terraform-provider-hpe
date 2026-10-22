resource "hpe_opsramp_client" "example_client" {
  name      = "Example Client"
  address   = "Valencia, Spain"
  country   = "Spain"
  time_zone = "Europe/Paris"

  addons = [
    "Adapter Integrations",
    "Alert Problem Area",
    "Application Management",
    "Batch Exports",
    "Knowledgebase Management",
    "Offline Alerts",
    "OS Service Start/Stop Actions",
    "Process Automation",
    "Remote Access Management",
    "Scheduled Task Management",
    "SKU Management",
    "SLA Management",
    "Stream Exports",
    "Log Management",
    "Trace Management"
  ]

  packages = [
    "Hybrid Discovery and Monitoring",
    "Event and Incident Management",
    "Remediation and Automation"
  ]
}