data "hpe_opsramp_resource_lookup" "on_prem" {
  query = "tag['application'] CONTAINS '${var.app_id}'AND nativeType IS NULL"
}