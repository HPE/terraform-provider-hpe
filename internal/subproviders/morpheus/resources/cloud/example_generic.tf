resource "hpe_morpheus_cloud" "example" {
  # Required fields
  name      = "TestCloud"
  tenant_id = 1
  group_id  = 1

  # General configuration
  code             = "aCode"
  labels           = ["aLabel1", "aLabel2"]
  enabled          = true
  location         = "somewhere"
  visibility       = "public"
  cloud_type_code  = "standard"

  # Agent and provisioning settings
  agent_install_mode       = "ssh"
  auto_recover_power_state = true

  # Cost and guidance settings
  costing_mode  = "costing"
  guidance_mode = "off"

  # Security settings
  security_mode = "off"

  # Shove into config
  appliance_url            = "https://somewhere.com"
  keyboard_layout = "us"
  data_center_name = "aDatacenter"
  external_id      = "aCode"
  import_existing_vms      = "off"

  # Generic configuration (HVM cloud)
  config = {
    certificateProvider          = "internal"
    enableNetworkTypeSelection = false
  }
}
