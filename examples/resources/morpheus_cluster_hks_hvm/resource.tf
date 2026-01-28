resource "hpe_morpheus_cluster_hks_hvm" "example" {
  name              = "tfhvm"
  resource_prefix   = "vmpre"
  hostname_prefix   = "ospre"
  description       = "Terraform HKS cluster example"
  cloud_id          = data.hpe_morpheus_cloud.morpheus_hvm.id
  group_id          = data.hpe_morpheus_group.morpheus_lab.id
  cluster_layout_id = 1070
  pod_cidr          = "172.20.0.0/16"
  service_cidr      = "172.30.0.0/16"
  workers           = 3

  server {
    plan_id          = data.hpe_morpheus_service_plan.master_nodes.id
    resource_pool_id = data.hpe_morpheus_resource_pool.hvm.id

    network_interface {
      network_id = data.hpe_morpheus_network.vm_network.id
    }

    storage_volume {
      root         = true
      size         = 20
      name         = "root"
      storage_type = 1
      datastore_id = 9
    }

    storage_volume {
      root         = false
      size         = 20
      name         = "data"
      storage_type = 1
      datastore_id = 2
    }

    tags = {
      "app" = "hksmaster"
    }
  }
}
