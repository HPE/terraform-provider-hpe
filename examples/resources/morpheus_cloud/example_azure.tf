resource "hpe_morpheus_cloud" "example" {
  name      = "TestCloud"
  tenant_id = 1
  group_id  = 1

  code             = "aCode"
  external_id      = "aCode"
  labels           = ["aLabel1", "aLabel2"]
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

  # Inventory discovery defaults
  default_datastore_sync_active      = true
  default_folder_sync_active         = true
  default_network_sync_active        = true
  default_plan_sync_active           = true
  default_pool_sync_active           = true
  default_security_group_sync_active = true

  config_azure = {
    azure_region    = "eastus"
    subscriber_id   = "sub-12345"
    tenant_id       = "tenant-67890"
    client_id       = "client-abc"
    client_secret   = "secret-xyz"
    resource_group  = "my-rg"
    cloud_type      = "global"
    import_existing = "on"
    rpc_mode        = "guestExec"
  }
}
