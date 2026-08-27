# Filter by cloud name. Values are Go regular expressions; a storage volume
# matches the block if its cloud_name matches ANY value.
data "hpe_morpheus_storage_volumes" "by_cloud" {
  filter {
    name   = "cloud_name"
    values = ["Production"]
  }
}

# Filter by status.
data "hpe_morpheus_storage_volumes" "online" {
  filter {
    name   = "status"
    values = ["provisioned"]
  }
}

# Filter by the storage server the volume is provisioned on, by name.
data "hpe_morpheus_storage_volumes" "alletra" {
  filter {
    name   = "storage_server_name"
    values = ["^alletra-"]
  }
}

# Multiple filter blocks are ANDed together.
data "hpe_morpheus_storage_volumes" "prod_provisioned" {
  filter {
    name   = "cloud_name"
    values = ["Production"]
  }

  filter {
    name   = "status"
    values = ["provisioned"]
  }
}

# No filter blocks returns all storage volumes (up to 250).
data "hpe_morpheus_storage_volumes" "all" {}

# Consume the full objects returned in `storage_volumes`.
output "storage_volume_ids" {
  value = [for v in data.hpe_morpheus_storage_volumes.by_cloud.storage_volumes : v.id]
}
