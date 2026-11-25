# An HPE iLO bare-metal host (hpe-baremetal-plugin IloBareMetalUnmanaged type).
resource "hpe_morpheus_host" "ilo_metal" {
  name                = "ilo-metal-1"
  cloud_id            = 20
  group_id            = 1
  compute_server_type = "IloBareMetalUnmanaged"

  config_bmaas = {
    ilo_ip_address   = "10.1.2.3"
    ilo_username     = "Administrator"
    ilo_password_wo  = var.ilo_password
    ilo_mac_address  = "aa:bb:cc:dd:ee:ff"
    resource_pool_id = "1"
  }
}

variable "ilo_password" {
  type      = string
  sensitive = true
}
