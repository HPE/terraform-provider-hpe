resource "hpe_morpheus_os_type_image" "example" {
  os_type_id        = 1
  virtual_image_id  = 42
  cloud_id          = 10
  provision_type_id = 3
}
