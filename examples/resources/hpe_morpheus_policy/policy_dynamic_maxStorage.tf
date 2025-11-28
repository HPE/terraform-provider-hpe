# Max Storage Policy - Limits storage allocation
resource "hpe_morpheus_policy" "max_storage" {
  name                     = "Max Storage Policy"
  description              = "Limit maximum storage allocation"
  associated_resource_type = "User"
  associated_resource_id   = 9969
  enabled                  = true

  policy_type = {
    code = "maxStorage"
  }

  config = {
    maxStorage        = "1000" # Maximum storage in GB
    excludeContainers = "off"  # Options: "on", "off" - exclude containers from count
  }
}
