data "hpe_morpheus_cluster" "hvm" {
  name = "hvm-cluster"
}

# The members, looked up by name. Each is a compute server -- the machine
# record backing an instance -- of that cluster.
data "hpe_morpheus_compute_server" "member" {
  for_each = toset(["web-1", "web-2"])

  name = each.value
}

resource "hpe_morpheus_cluster_affinity_group" "members" {
  cluster_id    = data.hpe_morpheus_cluster.hvm.id
  name          = "web-keep-separate"
  affinity_type = "KEEP_SEPARATE"

  # Terraform owns the membership: the group is created with these servers and
  # every apply re-asserts exactly this set. Do not also manage members of this
  # group with hpe_morpheus_cluster_affinity_group_member.
  servers = [for s in data.hpe_morpheus_compute_server.member : s.id]
}
