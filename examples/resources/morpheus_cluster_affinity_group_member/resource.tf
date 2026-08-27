# Place a server in a cluster affinity group without taking ownership of the
# group's other members.
#
# Membership added by other means -- an instance provisioned with
# config_hvm.affinity_group_id, or a node added by hpe_morpheus_instance_node --
# is left untouched.

resource "hpe_morpheus_cluster_affinity_group" "example" {
  cluster_id    = 15056
  name          = "web-tier"
  affinity_type = "KEEP_TOGETHER"
  active        = true
}

resource "hpe_morpheus_cluster_affinity_group_member" "example" {
  cluster_id        = 15056
  affinity_group_id = hpe_morpheus_cluster_affinity_group.example.id
  server_id         = 572308
}
