resource "hpe_morpheus_cluster_mks_vsphere" "example" {
  name                    = "tfvsphere"
  resource_prefix         = "vmpre"
  hostname_prefix         = "ospre"
  description             = "Terraform MKS cluster example"
  cloud_id                = data.hpe_morpheus_cloud.morpheus_vsphere.id
  group_id                = data.hpe_morpheus_group.morpheus_lab.id
  cluster_layout_id       = 1070
  pod_cidr                = "172.20.0.0/16"
  service_cidr            = "172.30.0.0/16"
  cluster_repo_account_id = 1

  master_node_pool {
    plan_id          = data.hpe_morpheus_service_plan.master_nodes.id
    resource_pool_id = data.hpe_morpheus_resource_pool.vsphere_resource_pool.id

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

    tags = {
      "app" = "mksmaster"
    }
  }

  worker_node_pool {
    count            = 3
    plan_id          = data.hpe_morpheus_service_plan.worker_nodes.id
    resource_pool_id = data.hpe_morpheus_resource_pool.vsphere_resource_pool.id

    network_interface {
      network_id = data.hpe_morpheus_network.vm_network.id
    }

    storage_volume {
      root         = true
      size         = 20
      name         = "root"
      storage_type = 1
      datastore_id = 2
    }

    storage_volume {
      root         = false
      size         = 20
      name         = "data"
      storage_type = 1
      datastore_id = 2
    }

    tags = {
      "app" = "mksworker"
    }
  }
}
