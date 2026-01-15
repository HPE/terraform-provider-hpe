resource "hpe_morpheus_datastore" "example" {
  name = "TestAlletraDatastore"
  datastore_type = {
    id   = 8
    code = "hpedatastore-alletra-mp"
  }
  associated_resource_type = "Cluster"
  visibility               = "private"
  active                   = true
  associated_resource_id   = 1

  config_alletramp_hvm = {
    protocol_type     = "iSCSI"
    enable_ransomware = false
  }

  storage_server = {
    id = 1
  }

  resource_permissions = {
    groups = [
      {
        id = 1
      }
    ]
  }
  tenants = [
    {
      id = 1
    }
  ]
}