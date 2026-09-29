resource "hpe_opsramp_integration" "vmware_integration" {
  display_name = "VMWare Integration 7"
  application  = "VMWARE"

  ip_address     = var.vmware_url_o
  credential_set = hpe_opsramp_credential_set.vmware_credential_set.id

  discovery_profiles = [
    {
      mgmt_profile_uuid = data.hpe_opsramp_management_profile.example.uuid
      scan_now          = true
      policy = {
        entity_type = "ALL"
        match_type  = "ANY"
        rules = [
          {
            filter_type   = "ANY_CLOUD_RESOURCE"
            resource_type = []
          }
        ]
        actions = [
          {
            action = "MANAGE DEVICE"
          }
        ]
      }
    }
  ]
}
