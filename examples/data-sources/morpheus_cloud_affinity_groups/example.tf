data "hpe_morpheus_cloud_affinity_groups" "example" {
  cloud_id = 1

  filter {
    name   = "affinity_type"
    values = ["KEEP_TOGETHER"]
  }
}
