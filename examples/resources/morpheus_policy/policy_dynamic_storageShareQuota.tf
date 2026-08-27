# File Share Storage Quota Policy - Limits file share storage
# Allowed associated_resource_types: User, Global
# Tenant specification: allowed (can specify tenants array)
resource "hpe_morpheus_policy" "file_share_quota" {
  name                     = "File Share Storage Quota Policy"
  description              = "Limit file share storage usage"
  associated_resource_type = "User"
  associated_resource_id   = 9969
  enabled                  = true

  policy_type = {
    code = "storageShareQuota"
  }

  config = {
    # Required
    maxStorage = "1000" # Maximum storage in GB
  }
}
