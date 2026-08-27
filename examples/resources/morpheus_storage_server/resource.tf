resource "hpe_morpheus_storage_server" "example" {
  name                        = "nfs-storage-01"
  type                        = "nfs"
  description                 = "Primary NFS storage server"
  enabled                     = true
  visibility                  = "private"
  service_host                = "10.0.1.50"
  service_url                 = "nfs://10.0.1.50/exports"
  service_username            = "admin"
  service_password_wo         = var.storage_password
  service_password_wo_version = 1

  tenants = [1, 2]
}

# Using a stored credential instead of inline credentials
resource "hpe_morpheus_storage_server" "with_credential" {
  name          = "nfs-storage-02"
  type          = "nfs"
  description   = "NFS storage using stored credential"
  enabled       = true
  visibility    = "private"
  service_host  = "10.0.1.51"
  service_url   = "nfs://10.0.1.51/exports"
  credential_id = 5
}
