resource "hpe_morpheus_cloud_affinity_group" "example" {
  cloud_id      = 1
  name          = "Example Affinity Group"
  affinity_type = "KEEP_TOGETHER"

  # Morpheus requires a pool on a cloud affinity group, and it must be a
  # resource pool of type Cluster.
  pool_id = 1
}
