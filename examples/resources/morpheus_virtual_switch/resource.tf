data "hpe_morpheus_cluster_available_nics" "nics" {
  cluster_id = 1
}

# The data-nfs address of each host, by host name.
locals {
  nfs_ips = {
    "hvm-host-1" = "10.20.0.11"
    "hvm-host-2" = "10.20.0.12"
  }
}

resource "hpe_morpheus_virtual_switch" "example" {
  cluster_id = 1
  name       = "vsw-hvm"
  bond_mode  = "active-backup"
  mtu        = 9000

  # One uplink per host, fanned out over the available-nics data source. As
  # hosts join the cluster they appear here and the next apply adds their uplink.
  uplinks = [
    for host in data.hpe_morpheus_cluster_available_nics.nics.available_nics : {
      server_id = host.server_id
      nic1      = "eth2"
      nic2      = "eth3"
    }
  ]

  # management lives on the cluster's first switch (vs0); a second switch carries
  # a vm segment and a tagged data-nfs segment with a distinct per-host IP.
  segments = {
    vm = {
      traffic_type = "vm"
    }
    nfs = {
      traffic_type = "data-nfs"
      vlan_id      = 100
      netmask      = "255.255.255.0"
      # Address each host by identity, not by its position in available_nics:
      # the API does not promise an order, and a host that is briefly offline
      # drops out of the list, which would shift every later address.
      host_ips = [
        for host in data.hpe_morpheus_cluster_available_nics.nics.available_nics : {
          server_id = host.server_id
          ip        = local.nfs_ips[host.server_name]
        }
      ]
    }
  }
}
