resource "hpe_opsramp_servicemap" "child_search" {
  name         = "Search Test"
  type         = "Resource"
  parent       = hpe_opsramp_servicemap.root.id
  search_query = "resourceType = \"Server\" AND name CONTAINS \"Test\""
}