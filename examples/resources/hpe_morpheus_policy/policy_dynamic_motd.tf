# Message of the Day (MOTD) Policy - Displays login messages
# Allowed associated_resource_types: Global
# Tenant specification: allowed (can specify tenants array)
resource "hpe_morpheus_policy" "motd" {
  name                     = "MOTD Policy"
  description              = "Display message of the day on login"
  associated_resource_type = "Global"
  enabled                  = true

  policy_type = {
    code = "motd"
  }

  config = {
    "motd.title"     = "Welcome"                          # Message title
    "motd.message"   = "Welcome to the Morpheus platform" # Message content
    "motd.type"      = "info"                             # Options: "info", "warning", "danger"
    "motd._fullPage" = "off"                              # Options: "on", "off" - display full page
  }
}
