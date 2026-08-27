resource "hpe_opsramp_servicemap_link" "servicemap_link" {
  parent = hpe_opsramp_servicemap.servicemap_root.id
  link   = hpe_opsramp_servicemap.servicemap_linked_root.id
}