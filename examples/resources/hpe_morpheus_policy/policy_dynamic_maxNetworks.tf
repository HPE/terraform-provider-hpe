# Network Quota Policy - Limits network count
resource "hpe_morpheus_policy" "network_quota" {
  name                     = "Network Quota Policy"
  description              = "Limit maximum network count"
  associated_resource_type = "User"
  associated_resource_id   = 9969
  enabled                  = true

  policy_type = {
    code = "maxNetworks"
  }

  config = {
    maxNetworks = "10" # Maximum number of networks
  }
}
