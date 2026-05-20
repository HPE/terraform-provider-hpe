resource "hpe_morpheus_library_container_type" "example" {
  name                = "App Node"
  short_name          = "app-node"
  container_version   = "1.0"
  provision_type_code = "vmware"
  description         = "Application server node type"
  virtual_image_id    = 1
  server_type         = "vm"
}
