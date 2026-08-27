resource "hpe_opsramp_client" "example_client" {
  name      = "Example Client"
  address   = "Valencia, Spain"
  country   = "Spain"
  time_zone = "Europe/Paris"

  addons = [
    "Adapter Integrations",
    "Extended Data Retention",
    "Knowledgebase Management",
    "OS Service Start/Stop Actions",
    "Offline Alerts",
    "Process Automation",
    "Remote Access Management",
    "SMS and Voice",
    "Alert Problem Area"
  ]

  packages = [
    "Hybrid Discovery and Monitoring",
    "Event and Incident Management",
    "Remediation and Automation"
  ]
}