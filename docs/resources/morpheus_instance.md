---
page_title: "hpe_morpheus_instance Resource - terraform-provider-hpe"
subcategory: "morpheus"
description: |-
  
---
# hpe_morpheus_instance (Resource)



Instance is a virtual machine or container deployed and managed by HPE Morpheus.
Morpheus oversees its entire lifecycle, from initial provisioning to scaling, 
monitoring, and eventual decommissioning.

-> Currently HVM and VMware instances are supported.<br>
With Morpheus versions prior to 8.0.11, make sure the root volume is the first defined.<br>
The addition and removal of volumes is not supported during updates.<br>
Updates fail when removing optional fields.<br>
Updates fail when removing `evars`.<br>
These will be addressed in a future release.

-> When an instance is created, it is marked as "ready" before DHCP has assigned IP addresses to all
`network_interfaces` and any `child_virtual_networks`.  A `terraform plan` will report that no changes
will be made.  Eventually, when all IP addresses have been assigned (this can be seen in the UI) a
`terraform apply` will report that no changes have been made but will update the state-file to include
the missing IP addresses.

-> We have removed `layout_size` from the Schema.  For technical reasons we have decided to only allow
the creation of one VM per instance.  When executing terraform an error will be raised stating that
`layout_size` is unsupported.  It is safe to remove the attribute from HCL, a `plan` will show no changes
to infrastructure after removal and on the next `apply` the attribute will be removed from the state-file.

-> We support `timeouts` using the Hashicorp Framework [timeouts package](https://developer.hashicorp.com/terraform/plugin/framework/resources/timeouts).
If the `timeouts` settings are changed in HCL an
`Update` will be triggered.  If the only change detected is for `timeouts` then the State will be updated with
the new settings but no `Morpheus` `Update` API calls will be made.  The default timeout for `create`, `delete`
`read` and `update` is 45 minutes

-> We've added a `connection_info` section (read-only) which contains the IP address(es) by which the instance
can be accessed

-> When creating an instance with network bonding and/or LAGs we cannot reconcile the created list of `network_interfaces`
with the HCL supplied.  In these cases the `connection_info` section will contain IP address(es).  To access the full
network configuration use the `hpe_morpheus_instance` `data-source` to read back the created instance.

## Example Usage

### HVM instance

```terraform
data "hpe_morpheus_cloud" "vme_cloud" {
  name = "HPE Alletra VME"
}

data "hpe_morpheus_service_plan" "vme_512mb" {
  name                = "1 CPU, 1GB Memory"
  provision_type_code = "kvm"
}

resource "hpe_morpheus_instance" "example" {
  name             = "TestInstance"
  cloud_id         = data.hpe_morpheus_cloud.vme_cloud.id # HPE Alletra VME
  layout_id        = 5385                                 # Single KVM VM
  instance_type_id = 9                  # (HVM) mvm-cluster

  group_id = 1
  plan_id  = data.hpe_morpheus_service_plan.vme_512mb.id # kvm-vm-512

  instance_context = "dev"
  network_interfaces = [
    {
      network_id = 103481
    }
  ]

  volumes = [
    {
      root_volume     = true
      name            = "root"
      size            = 10
      storage_type_id = 1
      datastore_id    = 38658
    },
    {
      root_volume     = false
      name            = "data"
      size            = 10
      storage_type_id = 1
      datastore_id    = 38658
    }
  ]

  tags = [
    {
      name  = "terraform"
      value = "true"
    },
    {
      name  = "acctest"
      value = "true"
    },
    {
      name  = "hpe_morpheus_instance"
      value = "true"
    },
    {
      name  = "sweepable"
      value = "true"
    },
    {
      name  = "managed_by"
      value = "terraform"
    }
  ]

  config = {
    resourcePoolId       = "pool-62299"
    poolProviderType     = "mvm"
    nestedVirtualization = "off"
    noAgent              = true
    createUser           = false
  }
}
```

### HVM Instance with two networks

One of these networks has a `child_virtual_networks` entry.  Note that only some network types support `child_virtual_networks`.
The Options API "/api/options/zoneNetworkOptions?zoneId=5&provisionTypeId=10" can be used to see which options are available.

```terraform
data "hpe_morpheus_cloud" "vme_cloud" {
  name = "HPE Alletra VME"
}

data "hpe_morpheus_service_plan" "vme_512mb" {
  name                = "1 CPU, 1GB Memory"
  provision_type_code = "kvm"
}

resource "hpe_morpheus_instance" "example" {
  name             = "TestInstance"
  cloud_id         = data.hpe_morpheus_cloud.vme_cloud.id # HPE Alletra VME
  layout_id        = 5385                                 # Single KVM VM
  instance_type_id = 9                  # (HVM) mvm-cluster

  group_id = 1
  plan_id  = data.hpe_morpheus_service_plan.vme_512mb.id # kvm-vm-512

  instance_context = "dev"
  network_interfaces = [
    {
      network_id = 103481
    }
  ]

  network_interfaces = [
    {
      network_id = 103481
      child_virtual_networks = [
        {
          network_id = 103481
        }
      ]
    },
    {
      network_id = 103481
    },
  ]

  volumes = [
    {
      root_volume     = true
      name            = "root"
      size            = 10
      storage_type_id = 1
      datastore_id    = 38658
    },
    {
      root_volume     = false
      name            = "data"
      size            = 10
      storage_type_id = 1
      datastore_id    = 38658
    }
  ]

  tags = [
    {
      name  = "terraform"
      value = "true"
    },
    {
      name  = "acctest"
      value = "true"
    },
    {
      name  = "hpe_morpheus_instance"
      value = "true"
    },
    {
      name  = "sweepable"
      value = "true"
    },
    {
      name  = "managed_by"
      value = "terraform"
    }
  ]

  config = {
    resourcePoolId       = "pool-62299"
    poolProviderType     = "mvm"
    nestedVirtualization = "off"
    noAgent              = true
    createUser           = false
  }
}
```

### HVM Instance with timeouts

We support `timeouts` using the Hashicorp Framework [timeouts package](https://developer.hashicorp.com/terraform/plugin/framework/resources/timeouts).
The following example specifies `timeouts` for `create`, `delete`, `update` and `read`.
If the `timeouts` settings are changed in HCL an `Update` will be triggered.  If the only change
detected is for `timeouts` then the State will be updated with the new settings but no `Morpheus` `Update` API calls
will be made.  The default timeout for `create`, `delete`, `read` and `update` is 45 minutes

```terraform
data "hpe_morpheus_cloud" "vme_cloud" {
  name = "HPE Alletra VME"
}

data "hpe_morpheus_service_plan" "vme_512mb" {
  name                = "1 CPU, 1GB Memory"
  provision_type_code = "kvm"
}

resource "hpe_morpheus_instance" "example" {
  name             = "TestInstance"
  cloud_id         = data.hpe_morpheus_cloud.vme_cloud.id # HPE Alletra VME
  layout_id        = 5385                                 # Single KVM VM
  instance_type_id = 9                  # (HVM) mvm-cluster

  group_id = 1
  plan_id  = data.hpe_morpheus_service_plan.vme_512mb.id # kvm-vm-512

  instance_context = "dev"
  network_interfaces = [
    {
      network_id = 103481
    }
  ]

  volumes = [
    {
      root_volume     = true
      name            = "root"
      size            = 10
      storage_type_id = 1
      datastore_id    = 38658
    },
    {
      root_volume     = false
      name            = "data"
      size            = 10
      storage_type_id = 1
      datastore_id    = 38658
    }
  ]

  tags = [
    {
      name  = "terraform"
      value = "true"
    },
    {
      name  = "acctest"
      value = "true"
    },
    {
      name  = "hpe_morpheus_instance"
      value = "true"
    },
    {
      name  = "sweepable"
      value = "true"
    },
    {
      name  = "managed_by"
      value = "terraform"
    }
  ]

  config = {
    resourcePoolId       = "pool-62299"
    poolProviderType     = "mvm"
    nestedVirtualization = "off"
    noAgent              = true
    createUser           = false
  }

  timeouts = {
    create = "1h"
    delete = "20m"
    update = "20m"
    read   = "10m"
  }
}
```

### VMware VM Instance

```terraform
data "hpe_morpheus_cloud" "vmware_cloud" {
  name = "QA VMware"
}

data "hpe_morpheus_service_plan" "vmware_512mb" {
  name                = "1 CPU, 1GB Memory"
  provision_type_code = "vmware"
}

data "hpe_morpheus_instance_type_layout" "vmware" {
  name    = "VMware VM"
  version = "22.04"
}

resource "hpe_morpheus_instance" "example" {
  name             = "TestInstance"
  cloud_id         = data.hpe_morpheus_cloud.vmware_cloud.id
  layout_id        = data.hpe_morpheus_instance_type_layout.vmware.id
  instance_type_id = 9

  group_id = 28
  plan_id  = data.hpe_morpheus_service_plan.vmware_512mb.id

  instance_context = "dev"
  network_interfaces = [
    {
      network_id = 86657
    }
  ]

  volumes = [
    {
      root_volume              = true
      name                     = "root"
      size                     = 10
      storage_type_id          = 1
      datastore_auto_selection = "auto"
    },
    {
      root_volume              = false
      name                     = "data"
      size                     = 10
      storage_type_id          = 1
      datastore_auto_selection = "auto"
    }
  ]

  tags = [
    {
      name  = "terraform"
      value = "true"
    },
    {
      name  = "acctest"
      value = "true"
    },
    {
      name  = "hpe_morpheus_instance"
      value = "true"
    },
    {
      name  = "sweepable"
      value = "true"
    },
    {
      name  = "managed_by"
      value = "terraform"
    }
  ]

  config = {
    resourcePoolId       = "pool-1"
    nestedVirtualization = "off"
    noAgent              = true
    createUser           = false
    vmwareFolderID       = "group-v79"
  }

  timeouts = {
    create = "1h"
    delete = "20m"
    update = "20m"
    read   = "10m"
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `group_id` (Number) The Group ID to provision the instance into.
- `instance_type_id` (Number) The type of instance by id we want to fetch.
- `layout_id` (Number) The layout id for the instance type that you want to provision. i.e. single process or cluster
- `name` (String) Name of the instance to be created.
- `network_interfaces` (Attributes List) The networkInterfaces parameter is for network configuration.

The Options API "/api/options/zoneNetworkOptions?zoneId=5&provisionTypeId=10" can be used to see which options are available. (see [below for nested schema](#nestedatt--network_interfaces))
- `plan_id` (Number) The id for the memory and storage option pre-configured within Morpheus.

### Optional

- `cloud_id` (Number) The Cloud ID to provision the instance onto.
- `config` (Dynamic) Configuration object. Settings vary by type.
- `evars` (Attributes Set) Environment Variables, an array of objects that have name and value. (see [below for nested schema](#nestedatt--evars))
- `instance_context` (String) Environment
- `ports` (Attributes Set) The ports parameter is for port configuration.

The layout may have default ports, which are defined in node types, that are always configured. This parameter will be for additional custom ports to be opened. (see [below for nested schema](#nestedatt--ports))
- `tags` (Attributes Set) Metadata tags, Array of objects having a name and value. (see [below for nested schema](#nestedatt--tags))
- `task_set_id` (Number) The Workflow ID to execute.
- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))
- `volumes` (Attributes List) Logical Volume configuration to create additional LVs at provision time (see [below for nested schema](#nestedatt--volumes))

### Read-Only

- `connection_info` (List of String) List of IP addresses to use when connecting to instance
- `id` (Number) The ID of this resource.

<a id="nestedatt--network_interfaces"></a>
### Nested Schema for `network_interfaces`

Optional:

- `child_virtual_networks` (Attributes List) The child_virtual_networks parameter is for network configuration of child virtual networks.  Note that this list
cannot be empty, it can either not be specified in HCL or if specified must contain at least one element.

The Options API "/api/options/zoneNetworkOptions?zoneId=5&provisionTypeId=10" can be used to see which options are available. (see [below for nested schema](#nestedatt--network_interfaces--child_virtual_networks))
- `ip_address` (String) The ip address. Not applicable when using DHCP or IP Pools.
- `ip_mode` (String) The mode for determining ip address. Can be 'static', 'dhcp' or ''.  The default is ''.
- `ip_pool` (Number) id of the ip pool to be used with this network
- `network_group_id` (Number) id of the network group to be used. Cannot be used with 'network_id', will be used instead of 'network_id'
- `network_id` (Number) id of the network to be used.  This cannot be used with 'network_group_id'
- `network_type_id` (Number) The id of the type of network interface

Read-Only:

- `name` (String) The name of the interface, e.g. 'eth0', 'eth1'
- `primary_interface` (Boolean) Is this interface the 'primary interface'?

<a id="nestedatt--network_interfaces--child_virtual_networks"></a>
### Nested Schema for `network_interfaces.child_virtual_networks`

Optional:

- `ip_address` (String) The ip address. Not applicable when using DHCP or IP Pools.
- `ip_mode` (String) The mode for determining ip address. Can be 'static', 'dhcp' or ''.  The default is ''.
- `ip_pool` (Number) id of the ip pool to be used with this network
- `network_group_id` (Number) id of the network group to be used. Cannot be used with 'network_id', will be used instead of 'network_id'
- `network_id` (Number) id of the network to be used.  This cannot be used with 'network_group_id'
- `network_type_id` (Number) The id of the type of network interface

Read-Only:

- `name` (String) The name of the interface, e.g. 'eth0', 'eth1'
- `primary_interface` (Boolean) Is this interface the 'primary interface'?



<a id="nestedatt--evars"></a>
### Nested Schema for `evars`

Required:

- `name` (String)
- `value` (String)


<a id="nestedatt--ports"></a>
### Nested Schema for `ports`

Required:

- `port` (Number) Port number.

Optional:

- `load_balancer_protocol` (String) Enable a load balancer and set load balancer protocol. HTTP, HTTPS, or TCP.
- `name` (String) A name for the port.


<a id="nestedatt--tags"></a>
### Nested Schema for `tags`

Required:

- `name` (String)
- `value` (String)


<a id="nestedatt--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).
- `delete` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Setting a timeout for a Delete operation is only applicable if changes are saved into state before the destroy operation occurs.
- `read` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Read operations occur during any refresh or planning operation when refresh is enabled.
- `update` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).


<a id="nestedatt--volumes"></a>
### Nested Schema for `volumes`

Optional:

- `controller_mount_point` (String) The controller mount point specification for this volume in the format:
  "id:busNumber:typeId:unitNumber"
For new storage controllers the id is passed as -1, so an example value would be:
  "-1:1:6:0"
which translates to id: -1 (new), busNumber: 1, storage controller type id: 6 (SCSI VMware Paravirtual), unit number: 0.
The current list of storage controllers is returned for instances and servers for determining existing id values.
Use /api/provision-types?code=vmware to see the available controllerTypes for vmware."
- `datastore_auto_selection` (String) Auto selection can be specified as auto or autoCluster (for clusters).
- `datastore_id` (Number) The ID of the specific datastore.
- `name` (String) Name/type of the LV being created.
- `root_volume` (Boolean) If set to false then a non-root LV will be created.
- `size` (Number) Size of the LV to be created in GBs.  Uses default from service plan.
- `size_id` (Number) Can be used to select pre-existing LV choices from Morpheus.
- `storage_type_id` (Number) Identifier for LV type

Read-Only:

- `id` (Number) The id for the LV configuration being created.
