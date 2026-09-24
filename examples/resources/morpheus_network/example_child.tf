data "hpe_morpheus_cloud" "vsphere" {
  name = "vSphere Cloud"
}

data "hpe_morpheus_group" "example" {
  name = "Example Group"
}

# The network the child is carved from, looked up in the cloud it belongs to.
data "hpe_morpheus_network" "parent" {
  name     = "VM Network"
  cloud_id = data.hpe_morpheus_cloud.vsphere.id
}

# A custom child network has no configuration block: its parent, addressing
# and DNS are attributes of the network. Its name and CIDR are updated in place.
resource "hpe_morpheus_network" "child" {
  name              = "example-terraform-child"
  cloud_id          = data.hpe_morpheus_cloud.vsphere.id
  group_id          = data.hpe_morpheus_group.example.id
  type_code         = "childNetwork"
  parent_network_id = data.hpe_morpheus_network.parent.id
  cidr              = "10.77.0.0/24"
  gateway           = "10.77.0.1"
  visibility        = "private"
  tenant_ids        = [1]
}
