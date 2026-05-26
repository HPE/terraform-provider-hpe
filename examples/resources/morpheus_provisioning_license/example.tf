resource "hpe_morpheus_provisioning_license" "example" {
  name         = "Example License"
  license_type = "win"
  license_key  = "XXXXX-XXXXX-XXXXX-XXXXX-XXXXX"
}
