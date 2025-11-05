# Message of the Day (MOTD) Policy - Displays login messages
resource "hpe_morpheus_policy" "motd" {
  name                     = "MOTD Policy"
  description              = "Display message of the day on login"
  associated_resource_type = "User"
  associated_resource_id   = 9969
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
