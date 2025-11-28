# Shutdown Policy - Auto-shutdown idle instances
# Allowed associated_resource_types: Group, Cloud, User, Global
# Tenant specification: allowed (can specify tenants array)
resource "hpe_morpheus_policy" "shutdown" {
  name                     = "Shutdown Policy"
  description              = "Auto-shutdown idle instances"
  associated_resource_type = "User"
  associated_resource_id   = 9969
  enabled                  = true

  policy_type = {
    code = "shutdown"
  }

  config = {
    shutdownType                     = "user"                        # Options: "user" (user configurable), "fixed" (strict shutdown)
    shutdownAge                      = "30"                          # Days instance is allowed to run before shutdown
    shutdownRenewal                  = "7"                           # If the instance is renewed, this is the number of day increments the shutdown date is increased by.
    shutdownNotify                   = "1"                           # Days before shutdown to notify via email
    shutdownMessage                  = "Instance will shutdown soon" # Notification message
    shutdownAutoRenew                = "on"                          # Options: "on", "off"
    shutdownAllowExtend              = "off"                         # Options: "on", "off" - allow users to extend
    shutdownExtensionsBeforeApproval = "0"                           # Number of extensions before requiring approval
    shutdownHideFixed                = false                         # Hide shutdown if fixed value
  }
}
