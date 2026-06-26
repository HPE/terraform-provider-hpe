resource "hpe_morpheus_storage_volume" "generic" {
  name              = "Example Storage Volume"
  type_id           = 1
  storage_server_id = 1
  max_storage       = 30 # GiB

  # config is a generic, write-only configuration map for storage volume types
  # that do not have a typed config block. Like config_alletramp_bmaas, its
  # values are sent to the API on create but are never stored in state.
  # Increment config_wo_version to recreate the volume with the new
  # configuration. The keys are the storage plugin's native config keys.
  config = {
    hpe_storage_datastore = 5
  }
  config_wo_version = 1
}
