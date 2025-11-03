data "hpe_morpheus_group" "morpheus_lab" {
  name = "MORPHEUS"
}

data "hpe_morpheus_cloud" "morpheus_vsphere" {
  name = "MORPHEUSVCENTER"
}

data "hpe_morpheus_network" "vm_network" {
  name = "VM Network"
}

data "hpe_morpheus_network" "internal_network" {
  name = "Internal Network"
}

data "hpe_morpheus_service_plan" "master_nodes" {
  name                = "2 CPU, 8GB Memory"
  provision_type_code = "vmware"
}

data "hpe_morpheus_service_plan" "worker_nodes" {
  name                = "2 CPU, 16GB Memory"
  provision_type_code = "vmware"
}

resource "hpe_morpheus_cluster_mks_vsphere" "tf_example_vsphere_instance" {
  name                    = "tfvsphere"
  resource_prefix         = "vmpre"
  hostname_prefix         = "ospre"
  description             = "Terraform MKS cluster example"
  cloud_id                = data.morpheus_cloud.morpheus_vsphere.id
  group_id                = data.morpheus_group.morpheus_lab.id
  cluster_layout_id       = 244
  pod_cidr                = "172.20.0.0/16"
  service_cidr            = "172.30.0.0/16"
  workflow_id             = 100
  api_proxy_id            = 1
  cluster_repo_account_id = 1

  master_node_pool {
    plan_id          = data.morpheus_plan.master_nodes
    resource_pool_id = 20

    network_interface {
      network_id = data.morpheus_network.vm_network.id
    }

    storage_volume {
      root         = true
      size         = 30
      name         = "root"
      storage_type = 1
      datastore_id = 29
    }

    tags = {
      "app" = "mksmaster"
    }
  }

  worker_node_pool {
    count            = 3
    plan_id          = data.morpheus_plan.worker_nodes
    resource_pool_id = 20

    network_interface {
      network_id = data.morpheus_network.vm_network.id
    }

    network_interface {
      network_id = data.morpheus_network.internal_network.id
    }

    storage_volume {
      root         = true
      size         = 30
      name         = "root"
      storage_type = 1
      datastore_id = 21
    }

    storage_volume {
      root         = false
      size         = 15
      name         = "data"
      storage_type = 1
      datastore_id = 21
    }

    tags = {
      "app" = "mksworker"
    }
  }
}

