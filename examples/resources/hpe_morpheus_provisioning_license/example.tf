resource "hpe_morpheus_provisioning_license" "example" {
  name         = "Windows Server 2022"
  license_type = "win"
  license_key  = "XXXXX-XXXXX-XXXXX-XXXXX-XXXXX"
  description  = "Windows Server 2022 Standard license"
}
