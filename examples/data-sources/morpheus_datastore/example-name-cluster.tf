# Every HVM cluster has datastores of its own (local, morpheus-cloud-init,
# morpheus-images), so on an appliance with more than one cluster a lookup by
# name alone finds one per cluster. cluster_id looks the name up among one
# cluster's datastores only.
data "hpe_morpheus_cluster" "hvm" {
  name = "Example cluster"
}

data "hpe_morpheus_datastore" "example" {
  name       = "local"
  cluster_id = data.hpe_morpheus_cluster.hvm.id
}
