# Expiration Policy - Sets instance expiration and renewal options
# Allowed associated_resource_types: Group, Cloud, User, Global
# Tenant specification: allowed (can specify tenants array)
resource "hpe_morpheus_policy" "expiration" {
  name                     = "Expiration Policy"
  description              = "Set instance expiration and renewal policies"
  associated_resource_type = "User"
  associated_resource_id   = 9969
  enabled                  = true

  policy_type = {
    code = "lifecycle"
  }

  config = {
    lifecycleType                     = "user"                      # Options: "user" (user configurable), "fixed" (fixed expiration)
    lifecycleAge                      = "30"                        # Days until expiration
    lifecycleRenewal                  = "7"                         # Days for renewal window
    lifecycleNotify                   = "1"                         # Days before expiration to notify
    lifecycleMessage                  = "Instance will expire soon" # Notification message
    lifecycleAutoRenew                = "on"                        # Options: "on", "off" - auto renewal lifecycle
    lifecycleAllowExtend              = "off"                       # Options: "on", "off" - allow users to extend
    lifecycleExtensionsBeforeApproval = "0"                         # Number of extensions before requiring approval
    lifecycleHideFixed                = false                       # Hide fixed expiration date from users
    # accountIntegrationId = "1"                                    # ID of your ServiceNow or approval integration
    # workflowType = "workflow"                                     # Options: "workflow" (legacy workflow), "flow" (ServiceNow Flow)
    # lifecycleWorkflowId = "123"                                   # ID of legacy ServiceNow workflow (set if workflowType is 'workflow')
    # flowId = "456"                                                # ID of ServiceNow Flow (set if workflowType is 'flow')
  }
}
