resource "hpe_morpheus_os_type" "example" {
  name               = "Example OS Type"
  code               = "example.os.type"
  platform           = "linux"
  bit_count          = 64
  description        = "An example OS type"
  os_family          = "debian"
  os_version         = "22.04"
  install_agent      = true
  cloud_init_version = "2"
}
