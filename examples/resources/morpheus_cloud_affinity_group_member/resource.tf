# Place a server in a cluster affinity group without taking ownership of the
# group's other members.
#
# Membership added by other means -- an instance provisioned with
# config_vmware.affinity_group_id, or a node added by hpe_morpheus_instance_node --
# is left untouched.

resource "hpe_morpheus_cloud_affinity_group" "example" {
  cloud_id      = 2
  pool_id       = 1
  name          = "web-tier"
  affinity_type = "KEEP_TOGETHER"
  active        = true
}

resource "hpe_morpheus_cloud_affinity_group_member" "example" {
  cloud_id          = 2
  affinity_group_id = hpe_morpheus_cloud_affinity_group.example.id
  server_id         = 572308
}
