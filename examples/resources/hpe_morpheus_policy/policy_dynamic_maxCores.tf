# Max Cores Policy - Limits CPU cores
resource "hpe_morpheus_policy" "max_cores" {
  name                     = "Max Cores Policy"
  description              = "Limit maximum CPU cores"
  associated_resource_type = "User"
  associated_resource_id   = 9969
  enabled                  = true

  policy_type = {
    code = "maxCores"
  }

  config = {
    maxCores          = "32"  # Maximum number of CPU cores
    excludeContainers = "off" # Options: "on", "off" - exclude containers from count
  }
}
