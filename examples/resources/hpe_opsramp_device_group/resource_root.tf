resource "hpe_opsramp_device_group" "device_group_root" {
  name      = "Test Group"
  resources = [hpe_opsramp_resource.resource1.uuid]
}