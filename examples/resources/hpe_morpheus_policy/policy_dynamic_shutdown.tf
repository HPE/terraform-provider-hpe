# Shutdown Policy - Auto-shutdown idle instances
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
    shutdownType                     = "user"                        # Options: "user" (user can extend), "fixed" (strict shutdown)
    shutdownAge                      = "30"                          # Days of inactivity before shutdown
    shutdownRenewal                  = "7"                           # Days for renewal window
    shutdownNotify                   = "1"                           # Days before shutdown to notify
    shutdownMessage                  = "Instance will shutdown soon" # Notification message
    shutdownAutoRenew                = "on"                          # Options: "on", "off"
    shutdownAllowExtend              = "off"                         # Options: "on", "off" - allow users to extend
    shutdownExtensionsBeforeApproval = "0"                           # Number of extensions before requiring approval
    shutdownHideFixed                = false                         # Hide fixed shutdown date from users
  }
}
