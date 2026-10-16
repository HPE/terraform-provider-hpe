# An HVM VM with two network interfaces whose second interface must not take
# over the default route or DNS, fixed from cloud-init user data before the
# Morpheus agent is installed.
#
# user_data is cloud-config YAML (or a plain script) merged over the user data
# Morpheus generates from the cloud, OS type, image, layout and type set:
#   - your bootcmd and runcmd entries run BEFORE the ones Morpheus adds (which
#     install the agent), so bootcmd is where guest network fix-ups belong;
#   - map keys merge; list values (packages, write_files, users, ...) REPLACE
#     Morpheus's list of the same name rather than appending to it;
#   - all runcmd entries run as one shell script: a failing command of yours
#     stops Morpheus's, so guard optional steps with `|| true`;
#   - a non-YAML value is treated as shell lines appended to runcmd.
#   - a `packages` entry the guest cannot fetch (no route to a mirror) blocks
#     cloud-init, and with it Morpheus's own runcmd and IP detection, until the
#     provision times out -- only list packages the guest network can reach.
# It is applied once, when the VM first boots. Changing it replaces the
# instance; add `lifecycle { ignore_changes = [config_hvm.user_data] }` to
# update only the stored value (which a later clone inherits).
#
# The same attribute exists on hpe_morpheus_instance_clone. Only HVM honours
# it as a typed attribute; for other clouds put `userData` in the generic
# `config` block, where support depends on the appliance's cloud plugin.

data "hpe_morpheus_cloud" "hvm_cloud" {
  name = "hvm"
}

data "hpe_morpheus_cluster" "hvm" {
  name = "MyCluster"
}

data "hpe_morpheus_group" "hvm" {
  name = "TestGroup"
}

data "hpe_morpheus_service_plan" "hvm_plan" {
  name                = "1 CPU, 512MB Memory"
  provision_type_code = "kvm"
}

data "hpe_morpheus_instance_type_layout" "hvm" {
  name    = "Single KVM VM"
  version = "11"
}

data "hpe_morpheus_network" "management" {
  name       = "Management"
  cluster_id = data.hpe_morpheus_cluster.hvm.id
}

data "hpe_morpheus_network" "vmnet" {
  name       = "vmnet"
  cluster_id = data.hpe_morpheus_cluster.hvm.id
}

resource "hpe_morpheus_instance" "example" {
  name             = "TestInstance"
  cloud_id         = data.hpe_morpheus_cloud.hvm_cloud.id
  layout_id        = data.hpe_morpheus_instance_type_layout.hvm.id
  instance_type_id = data.hpe_morpheus_instance_type_layout.hvm.instance_type.id
  group_id         = data.hpe_morpheus_group.hvm.id
  plan_id          = data.hpe_morpheus_service_plan.hvm_plan.id
  instance_context = "dev"

  # eth0: the VM network the workload lives on; eth1: management, which
  # Morpheus would otherwise also give a default route and DNS servers.
  network_interfaces = [
    { network_id = data.hpe_morpheus_network.vmnet.id },
    { network_id = data.hpe_morpheus_network.management.id },
  ]

  volumes = [
    {
      root_volume              = true
      name                     = "root"
      size                     = 10
      storage_type_id          = 1
      datastore_auto_selection = "auto"
    }
  ]

  config_hvm = {
    resource_pool_id      = "pool-${data.hpe_morpheus_cluster.hvm.permissions.resource_pool.id}"
    nested_virtualization = "off"
    no_agent              = false
    create_user           = false

    # Morpheus trims the value, so trimspace() keeps the plan empty. Use
    # templatefile() instead of a heredoc to parameterise it.
    user_data = trimspace(<<-EOT
      #cloud-config
      # Keep eth1 (management) off the default route and out of DNS, and make
      # eth0 the only default -- written and applied in bootcmd so it is in
      # place before Morpheus's runcmd reaches the appliance to fetch the agent.
      bootcmd:
        - |
          cat > /etc/netplan/60-tf-routes.yaml <<'NP'
          network:
            version: 2
            ethernets:
              eth1:
                dhcp4-overrides: { use-routes: false, use-dns: false }
              eth0:
                routes:
                  - to: default
                    via: 192.168.100.1
          NP
        - chmod 600 /etc/netplan/60-tf-routes.yaml
        - netplan apply || true
      # `packages` REPLACES Morpheus's list (it only asks for curl), so keep curl.
      packages:
        - curl
        - jq
      # Do NOT set `users` here unless you mean to drop the accounts Morpheus
      # creates (the cloud-init user and the provisioning user's Linux
      # account): like every list, it REPLACES Morpheus's. Add an operator
      # account from runcmd instead, where it is additive.
      runcmd:
        - useradd -m -s /bin/bash -G sudo operator || true
        - install -d -m 700 -o operator -g operator /home/operator/.ssh
        - echo 'ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIExampleKeyReplaceMe operator@example' > /home/operator/.ssh/authorized_keys
        - chown operator:operator /home/operator/.ssh/authorized_keys && chmod 600 /home/operator/.ssh/authorized_keys
      write_files:
        - path: /etc/profile.d/proxy.sh
          permissions: "0644"
          content: |
            export http_proxy=http://proxy.example.com:8080
            export https_proxy=http://proxy.example.com:8080
    EOT
    )
  }
}
