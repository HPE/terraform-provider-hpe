resource "hpe_morpheus_storage_volume" "example" {
  name              = "Example Storage Volume"
  type_code         = "hpealletraMPLUN"
  storage_server_id = 1
  max_storage       = 10 # GiB
}
