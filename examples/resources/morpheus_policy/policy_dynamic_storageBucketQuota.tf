# Object Storage Quota Policy - Limits object storage
# Allowed associated_resource_types: User, Global
# Tenant specification: allowed (can specify tenants array)
resource "hpe_morpheus_policy" "object_storage_quota" {
  name                     = "Object Storage Quota Policy"
  description              = "Limit object storage usage"
  associated_resource_type = "User"
  associated_resource_id   = 9969
  enabled                  = true

  policy_type = {
    code = "storageBucketQuota"
  }

  config = {
    # Required
    maxStorage = "1000" # Maximum storage in GB
  }
}
