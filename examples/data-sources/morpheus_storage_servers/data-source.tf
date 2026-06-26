# Filter by name. Values are Go regular expressions; a storage server matches
# the block if its name matches ANY value.
data "hpe_morpheus_storage_servers" "by_name" {
  filter {
    name   = "name"
    values = ["^prod-", "-alletra$"]
  }
}

# Filter by storage server type code (for example, 3par or netapp).
data "hpe_morpheus_storage_servers" "three_par" {
  filter {
    name   = "type_code"
    values = ["3par"]
  }
}

# Filter by the cloud (zone) the storage server is scoped to, by name or id.
data "hpe_morpheus_storage_servers" "by_cloud" {
  filter {
    name   = "cloud_name"
    values = ["Production"]
  }
}

# Multiple filter blocks are ANDed together.
data "hpe_morpheus_storage_servers" "enabled_three_par" {
  filter {
    name   = "type_code"
    values = ["3par"]
  }

  filter {
    name   = "enabled"
    values = ["true"]
  }
}

# No filter blocks returns all storage servers (up to 250).
data "hpe_morpheus_storage_servers" "all" {}

# Consume the full objects returned in `storage_servers`.
output "storage_server_ids" {
  value = [for s in data.hpe_morpheus_storage_servers.by_name.storage_servers : s.id]
}
