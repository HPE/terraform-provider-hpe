# List all datastores.
data "hpe_morpheus_datastores" "all" {}

# List only NFS datastores.
data "hpe_morpheus_datastores" "nfs_only" {
  filter {
    name   = "type"
    values = ["nfs"]
  }
}

# Combine filters (AND logic): NFS datastores that are provisioned.
data "hpe_morpheus_datastores" "nfs_provisioned" {
  filter {
    name   = "type"
    values = ["nfs"]
  }
  filter {
    name   = "status"
    values = ["provisioned"]
  }
}
