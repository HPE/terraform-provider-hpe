resource "hpe_morpheus_cloud" "example" {
  # Required fields
  name      = "TestCloud"
  tenant_id = 1
  group_id  = 1


  # General configuration
  code             = "aCode"
  external_id      = "aCode"
  labels           = ["aLabel1", "aLabel2"]
  data_center_name = "aDatacenter"
  enabled          = true
  location         = "somewhere"
  visibility       = "public"

  # Agent and provisioning settings
  agent_install_mode       = "ssh"
  appliance_url            = "https://somewhere.com"
  auto_recover_power_state = true
  import_existing_vms      = "off"

  # Cost and guidance settings
  costing_mode  = "costing"
  guidance_mode = "off"

  # Security settings
  security_mode = "off"

  # Console settings
  keyboard_layout = "us"

  # HVM-specific configuration
  config_hvm = {
    certificate_provider          = "internal"
    enable_network_type_selection = false
  }
}
