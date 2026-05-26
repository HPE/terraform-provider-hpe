resource "hpe_morpheus_library_container_type" "example" {
  name                = "Example Container Type"
  short_name          = "example"
  container_version   = "1.0"
  provision_type_code = "docker"
}
