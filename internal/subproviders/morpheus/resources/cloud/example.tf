resource "hpe_morpheus_cloud" "example" {
  name      = "TestCloud"
  tenant_id = 1
  group_id  = 1

  code                     = "aCode"
  external_id              = "aCode"
  labels                   = ["aLabel1", "aLabel2"]
  agent_install_mode       = "ssh"
  appliance_url            = "https://somewhere.com"
  auto_recover_power_state = true
  costing_mode             = "costing"
  data_center_name         = "aDatacenter"
  enabled                  = true
  guidance_mode            = "off"
  import_existing_vms      = "off"
  keyboard_layout          = "us"
  location                 = "somewhere"
  security_mode            = "off"
  visibility               = "public"

  config_hvm = {
    certificate_provider          = "internal"
    enable_network_type_selection = false
  }
}
