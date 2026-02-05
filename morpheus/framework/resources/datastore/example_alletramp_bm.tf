data "hpe_morpheus_cloud" "metal" {
  name = "Metal"
}

resource "hpe_morpheus_datastore" "example" {
  name = "TestAlletraDatastore"
  datastore_type = {
    id   = 12
    code = "hpedatastore-alletra-mp-bmaas"
  }
  associated_resource_type = "Cloud"
  visibility               = "private"
  active                   = true
  associated_resource_id   = data.hpe_morpheus_cloud.metal.id

  config = {
    protocol_type     = "iSCSI"
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