# Max Memory Policy - Limits memory allocation
resource "hpe_morpheus_policy" "max_memory" {
  name                     = "Max Memory Policy"
  description              = "Limit maximum memory allocation"
  associated_resource_type = "User"
  associated_resource_id   = 9969
  enabled                  = true

  policy_type = {
    code = "maxMemory"
  }

  config = {
    maxMemory         = "8192" # Maximum memory in MB
    excludeContainers = "off"  # Options: "on", "off" - exclude containers from count
  }
}
