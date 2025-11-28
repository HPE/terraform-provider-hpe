# Storage Server Storage Quota Policy - Limits storage on specific storage server
# Allowed associated_resource_types: Global
# Tenant specification: allowed (can specify tenants array)
resource "hpe_morpheus_policy" "storage_server_quota" {
  name                     = "Storage Server Storage Quota Policy"
  description              = "Limit storage usage on specific storage server"
  associated_resource_type = "Global"
  associated_resource_id   = 9969
  enabled                  = true

  policy_type = {
    code = "storageServerQuota"
  }

  config = {
    storageServerId = "1"    # ID of the storage server
    maxStorage      = "1000" # Maximum storage in GB
  }
}
