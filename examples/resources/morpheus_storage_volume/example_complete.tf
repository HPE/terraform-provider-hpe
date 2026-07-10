resource "hpe_morpheus_storage_volume" "example" {
  name              = "Example Storage Volume"
  type_code         = "hpealletraMPLUN"
  storage_server_id = 1
  max_storage       = 10 # GiB

  # An Alletra MP LUN requires a data store, supplied via the write-only
  # config_alletramp_bmaas block. Increment config_alletramp_bmaas_wo_version
  # whenever this block changes to recreate the volume.
  config_alletramp_bmaas = {
    datastore_id      = 5
    shared            = false
    compute_server_id = 10
  }
  config_alletramp_bmaas_wo_version = 1
}
