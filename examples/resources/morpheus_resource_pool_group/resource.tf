resource "hpe_morpheus_resource_pool_group" "example" {
  name              = "TFExample Resource Pool Group"
  description       = "TFExample Resource Pool Group"
  mode              = "roundrobin"
  resource_pool_ids = [1, 2, 3]
  all_group_access  = true
  group_access {
    group_id = 2
    default  = true
  }
  visibility = "public"
  tenant_ids = [1, 2]
}
