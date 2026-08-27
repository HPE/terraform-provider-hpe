resource "hpe_opsramp_servicemap" "child_service" {
  name   = "Service Test"
  type   = "Service"
  parent = hpe_opsramp_servicemap.root.id
}