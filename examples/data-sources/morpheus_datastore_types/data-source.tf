# List all datastore types.
data "hpe_morpheus_datastore_types" "all" {}

# Filter datastore types by code.
data "hpe_morpheus_datastore_types" "alletra" {
  filter {
    name   = "code"
    values = ["hpedatastore-alletra-mp-bmaas"]
  }
}

# Filter datastore types that are creatable.
data "hpe_morpheus_datastore_types" "creatable" {
  filter {
    name   = "creatable"
    values = ["true"]
  }
}
