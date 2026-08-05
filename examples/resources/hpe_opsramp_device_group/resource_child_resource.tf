resource "hpe_opsramp_device_group" "device_group_resources" {
  parent_id = hpe_opsramp_device_group.device_group_root.id
  name      = "Test Resources"
  resources = [hpe_opsramp_resource.resource2.uuid]
}