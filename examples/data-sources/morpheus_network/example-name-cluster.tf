# Each HVM cluster has networks of its own — the Management network Morpheus
# creates with it, the port groups on its virtual switches — so two clusters in
# one cloud can both have a network of the same name. cluster_id looks the name
# up among one cluster's networks only.
data "hpe_morpheus_cluster" "hvm" {
  name = "Example cluster"
}

data "hpe_morpheus_network" "example" {
  name       = "Management"
  cluster_id = data.hpe_morpheus_cluster.hvm.id
}
