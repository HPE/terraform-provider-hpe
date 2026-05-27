resource "hpe_morpheus_cloud" "example" {
  name      = "TestCloud"
  tenant_id = 1
  group_id  = 1

  code            = "aCode"
  labels          = ["aLabel1", "aLabel2"]
  enabled         = true
  location        = "somewhere"
  visibility      = "public"
  cloud_type_code = "standard"

  agent_install_mode       = "ssh"
  auto_recover_power_state = true

  costing_mode  = "costing"
  guidance_mode = "off"

  security_mode = "off"

  appliance_url       = "https://somewhere.com"
  keyboard_layout     = "us"
  data_center_name    = "aDatacenter"
  external_id         = "aCode"
  import_existing_vms = "off"

  # Inventory discovery (only security_groups applies to standard clouds)
  default_security_group_sync_active = true

  config = {
    certificateProvider        = "internal"
    enableNetworkTypeSelection = false
  }
}
