data "hpe_morpheus_cluster_affinity_groups" "example" {
  cluster_id = 1

  filter {
    name   = "affinity_type"
    values = ["KEEP_TOGETHER"]
  }
}
