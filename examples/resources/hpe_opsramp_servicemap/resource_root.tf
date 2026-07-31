resource "hpe_opsramp_servicemap" "root" {
  name = "Root"
  type = "Service"
}

resource "hpe_opsramp_servicemap" "child_resource" {
  name      = "Resource Test"
  type      = "Resource"
  parent    = hpe_opsramp_servicemap.root.id
  resources = [hpe_opsramp_resource.resource1.uuid]
}

resource "hpe_opsramp_servicemap" "child_search" {
  name         = "Search Test"
  type         = "Resource"
  parent       = hpe_opsramp_servicemap.root.id
  search_query = "resourceType = \"Server\" AND name CONTAINS \"Test\""
}