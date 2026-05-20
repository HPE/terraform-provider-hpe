resource "hpe_morpheus_virtual_image" "example" {
  name          = "Ubuntu 22.04 LTS"
  image_type    = "vmware"
  os_type_id    = 1
  is_cloud_init = true
  install_agent = true
  min_ram       = 1073741824
}
