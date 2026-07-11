resource "hpe_opsramp_device_group" "device_group_query" {
  parent_id    = hpe_opsramp_device_group.device_group_root.id
  name         = "Test Queries"
  search_query = format("resourceType = \"Server\" AND uuid = \"%s\"", hpe_opsramp_resource.resource3.uuid)
}