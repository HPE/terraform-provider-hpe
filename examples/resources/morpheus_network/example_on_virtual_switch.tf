data "hpe_morpheus_cloud" "hvm" {
  name = "HVM Cloud"
}

data "hpe_morpheus_cluster" "hvm" {
  name = "hvm-cluster"
}

data "hpe_morpheus_group" "example" {
  name = "Example Group"
}

data "hpe_morpheus_tenant" "master" {
  name = "Master Tenant"
}

data "hpe_morpheus_cluster_available_nics" "nics" {
  cluster_id = data.hpe_morpheus_cluster.hvm.id
}

# The cluster's own network server, hidden from the ordinary listing; every
# network on the cluster is created through it.
data "hpe_morpheus_network_server" "hvm" {
  name = "${data.hpe_morpheus_cluster.hvm.name} HVM Network Server"
}

# The data-nfs address of each host, by host name: available_nics does not
# promise an order, so addressing by position would reassign IPs on a reorder.
locals {
  nfs_ips = {
    "hvm-host-1" = "10.20.0.11"
    "hvm-host-2" = "10.20.0.12"
  }
}

# A second virtual switch on a Morpheus 9.1 HVM cluster, with a VM traffic
# segment for standard networks and a tagged data-nfs segment for a data network.
resource "hpe_morpheus_virtual_switch" "vs1" {
  cluster_id = data.hpe_morpheus_cluster.hvm.id
  name       = "vs1"
  mtu        = 9000

  uplinks = [
    for host in data.hpe_morpheus_cluster_available_nics.nics.available_nics : {
      server_id = host.server_id
      nic1      = "eth2"
    }
  ]

  segments = {
    vm = {
      traffic_type = "vm"
    }
    nfs = {
      traffic_type = "data-nfs"
      vlan_id      = 100
      netmask      = "255.255.255.0"
      host_ips = [
        for host in data.hpe_morpheus_cluster_available_nics.nics.available_nics : {
          server_id = host.server_id
          ip        = local.nfs_ips[host.server_name]
        }
      ]
    }
  }
}

# A standard network names its switch by the switch's libvirt network name
# (hvmnetN), which only the hvm-vswitch-or-router option source exposes; look it
# up by the switch's name. Going through the resource's name also makes
# Terraform destroy the network before the switch, which Morpheus requires.
data "hpe_morpheus_option_source_value" "vs1" {
  source       = "hvm-vswitch-or-router"
  service_type = "com.hpe.morpheus.hvm"
  cloud_id     = data.hpe_morpheus_cloud.hvm.id
  name         = hpe_morpheus_virtual_switch.vs1.name
  parameters = {
    "network.zonePool.id" = tostring(data.hpe_morpheus_cluster.hvm.permissions.resource_pool.id)
  }
}

# The network's VLAN, and every trunked VLAN, must not collide with a segment
# VLAN on that switch (the nfs segment above uses 100).
resource "hpe_morpheus_network" "vm" {
  name              = "example-terraform-vm-on-vs1"
  cloud_id          = data.hpe_morpheus_cloud.hvm.id
  group_id          = data.hpe_morpheus_group.example.id
  type_code         = "ovsPortGroup"
  switch_id         = data.hpe_morpheus_option_source_value.vs1.value
  zone_pool_id      = data.hpe_morpheus_cluster.hvm.permissions.resource_pool.id
  network_server_id = data.hpe_morpheus_network_server.hvm.id
  vlan_id           = 2103
  cidr              = "10.209.103.0/24"
  visibility        = "private"
  tenant_ids        = [data.hpe_morpheus_tenant.master.id]

  config_ovs_port_group = {
    vlan_trunks = "200,201"
    trunk_only  = false
  }
}

# An HVM data network on the switch's data-nfs segment: the VSWITCH layout
# takes the segment's id, and the segment decides the VLAN and MTU.
resource "hpe_morpheus_network" "nfs" {
  name              = "example-terraform-nfs-on-vs1"
  cloud_id          = data.hpe_morpheus_cloud.hvm.id
  group_id          = data.hpe_morpheus_group.example.id
  type_code         = "hvmPortGroupDataNetwork"
  zone_pool_id      = data.hpe_morpheus_cluster.hvm.permissions.resource_pool.id
  network_server_id = data.hpe_morpheus_network_server.hvm.id
  visibility        = "private"
  tenant_ids        = [data.hpe_morpheus_tenant.master.id]

  config_hvm_data = {
    cluster_network_layout    = "VSWITCH"
    virtual_switch_segment_id = hpe_morpheus_virtual_switch.vs1.segments["nfs"].id
  }
}
