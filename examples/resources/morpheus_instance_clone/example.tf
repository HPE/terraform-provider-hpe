
# The instance to clone, by name.
data "hpe_morpheus_instance" "source" {
  name = "my-instance"
}

# The network the clone attaches to, looked up within its cluster: every HVM
# cluster has networks of the same name (Management, for one), so the cluster
# is what makes the name unique.
data "hpe_morpheus_cluster" "clone" {
  name = "MyCluster"
}

data "hpe_morpheus_network" "clone" {
  name       = "MyNetwork"
  cluster_id = data.hpe_morpheus_cluster.clone.id
}

resource "hpe_morpheus_instance_clone" "example" {
  source_instance_id = data.hpe_morpheus_instance.source.id
  name               = "my-clone"

  volumes = [
    {
      name        = "root"
      size        = 20
      root_volume = true
    }
  ]

  network_interfaces = [
    {
      network_id = data.hpe_morpheus_network.clone.id
      ip_mode    = ""
    }
  ]
}
