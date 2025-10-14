resource "hpe_morpheus_cloud" "example" {
  name      = "TestCloud"
  tenant_id = 1
  group_id  = 1

  code             = "aCode"
  external_id      = "aCode"
  labels           = ["terraform", "acctest", "hpe_morpheus_cloud", "sweepable", "aLabel1", "aLabel2"]
  data_center_name = "aDatacenter"
  enabled          = true
  location         = "somewhere"
  visibility       = "public"

  agent_install_mode       = "ssh"
  appliance_url            = "https://somewhere.com"
  auto_recover_power_state = true
  import_existing_vms      = "off"

  costing_mode  = "costing"
  guidance_mode = "off"

  security_mode = "off"

  keyboard_layout = "us"

  config_hvm = {
    certificate_provider          = "internal"
    enable_network_type_selection = false
  }
}
