resource "hpe_morpheus_library_layout" "example" {
  instance_type_id    = 1
  name                = "Example Layout"
  instance_version    = "1.0"
  provision_type_code = "docker"
}
