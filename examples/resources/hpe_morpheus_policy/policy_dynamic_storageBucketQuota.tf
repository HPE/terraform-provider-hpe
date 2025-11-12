# Object Storage Quota Policy - Limits object storage
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
    maxStorage = "1000" # Maximum storage in GB
  }
}
