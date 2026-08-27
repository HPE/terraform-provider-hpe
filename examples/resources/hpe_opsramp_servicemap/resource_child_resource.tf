resource "hpe_opsramp_servicemap" "child_resource" {
  name      = "Resource Test"
  type      = "Resource"
  parent    = hpe_opsramp_servicemap.root.id
  resources = [hpe_opsramp_resource.resource1.uuid]
}