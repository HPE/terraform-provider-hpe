---
page_title: "hpe_morpheus_instance Resource - terraform-provider-hpe"
subcategory: "morpheus"
description: |-
  
---
# hpe_morpheus_instance (Resource)



Instance is a virtual machine, bare metal machine or container deployed and managed by HPE Morpheus.
Morpheus oversees its entire lifecycle, from initial provisioning to scaling, 
monitoring, and eventual decommissioning.

-> Currently HVM, VMware, AWS, Azure and BMaaS instances are supported.  We have static `config` schema for the following:<br>
&nbsp;&nbsp;&nbsp;&nbsp;- Azure: `config_azure`<br>
&nbsp;&nbsp;&nbsp;&nbsp;- AWS: `config_aws`<br>
&nbsp;&nbsp;&nbsp;&nbsp;- HVM: `config_hvm`<br>
&nbsp;&nbsp;&nbsp;&nbsp;- VMware: `config_vmware`<br><br>
Some general issues<br>
&nbsp;&nbsp;&nbsp;&nbsp;- With Morpheus versions prior to 8.0.11, make sure the root volume is the first defined.<br>
&nbsp;&nbsp;&nbsp;&nbsp;- Updates fail when removing optional fields.<br>
&nbsp;&nbsp;&nbsp;&nbsp;- Updates fail when removing `evars`.<br>
These will be addressed in a future release.

-> Some general notes:<br><br>
When an instance is created, it is marked as "ready" before DHCP has assigned IP addresses to all
`network_interfaces` and any `child_virtual_networks`.  A `terraform plan` will report that no changes
will be made.  Eventually, when all IP addresses have been assigned (this can be seen in the UI) a [refresh-only](https://developer.hashicorp.com/terraform/tutorials/state/refresh)
`terraform plan --refresh-only` will report that one or more `network_interfaces` have changed. A subsequent
`terraform apply --refresh-only` will update the State with the missing IP addresses.<br><br>
`layout_size` is optional and at the moment the only supported value is `1` which is also the default.
In other words we only support the creation of one VM per instance.  We may relax this restriction in a future release.<br><br>
We support `timeouts` using the Hashicorp Framework [timeouts package](https://developer.hashicorp.com/terraform/plugin/framework/resources/timeouts).
If the `timeouts` settings are changed in HCL an
`Update` will be triggered.  If the only change detected is for `timeouts` then the State will be updated with
the new settings but no `Morpheus` `Update` API calls will be made.  The default timeout for `create`, `delete`
`read` and `update` is 45 minutes.<br><br>
**Update** is supported but note the following:<br>
&nbsp;&nbsp;&nbsp;&nbsp;- **Do not** change the order of volumes or network_interfaces in HCL<br>
&nbsp;&nbsp;&nbsp;&nbsp;- The addition and removal of volumes is supported<br>
&nbsp;&nbsp;&nbsp;&nbsp;- Existing volumes can have their size increased but not decreased<br>
&nbsp;&nbsp;&nbsp;&nbsp;- Network interface update is supported for Morpheus versions >= 8.1.2, for earlier<br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  versions network_interfaces changes will force a new instance to be created.<br>
&nbsp;&nbsp;&nbsp;&nbsp;- Updates to `service_plan_options` for service plans that support the setting of options<br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp; are supported for Morpheus versions >= 8.1.2, for earlier versions updates<br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp; will force a new instance to be created<br>
&nbsp;&nbsp;&nbsp;&nbsp;- Updates fail when removing optional fields or `evars`<br><br>
We've added a `connection_info` section (read-only) which contains the IP address(es) by which the instance
can be accessed<br><br>
The `datastore_auto_selection` attribute is not supported for BMaaS instances.<br><br>
When creating an instance with network bonding and/or LAGs we cannot reconcile the created list of `network_interfaces`
with the HCL supplied.  In these cases the `connection_info` section will contain IP address(es).  To access the full
network configuration use the `hpe_morpheus_instance` `data-source` to read back the created instance.

-> Some of the examples below have the following settings in their `config` blocks:<br>
&nbsp;&nbsp;&nbsp;&nbsp;- `no_agent` in static config (equivalent to `noAgent` in dynamic config) is set to `true`<br>
&nbsp;&nbsp;&nbsp;&nbsp;- `create_user` in static config (equivalent to `createUser` in dynamic config) is set to<br>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  `false`<br><br>
These settings can be changed as required.

## HVM Instance

```terraform
data "hpe_morpheus_cloud" "vme_cloud" {
  name = "MyCloud"
}

data "hpe_morpheus_service_plan" "vme_512mb" {
  name                = "1 CPU, 1GB Memory"
  provision_type_code = "kvm"
}

resource "hpe_morpheus_instance" "example" {
  name             = "TestInstance"
  cloud_id         = data.hpe_morpheus_cloud.vme_cloud.id 
  layout_id        = 644 # Single KVM VM
  instance_type_id = 9 # (HVM) mvm-cluster

  group_id = 1
  plan_id  = data.hpe_morpheus_service_plan.vme_512mb.id # kvm-vm-512

  instance_context = "dev"
  network_interfaces = [
    {
      network_id = 755
    }
  ]

  volumes = [
    {
      root_volume     = true
      name            = "root"
      size            = 10
      storage_type_id = 1
      datastore_id    = 555
    },
    {
      root_volume     = false
      name            = "data"
      size            = 10
      storage_type_id = 1
      datastore_id    = 555
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

  config_hvm = {
    resource_pool_id      = "pool-700"
    nested_virtualization = "off"
    no_agent              = true
    create_user           = false
  }
}
```

## HVM Instance with two networks

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

  config_hvm = {
    resource_pool_id      = "pool-62299"
    nested_virtualization = "off"
    no_agent              = false
    create_user           = true
  }
}
```

## HVM Instance with timeouts

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

  config_hvm = {
    resource_pool_id      = "pool-62299"
    nested_virtualization = "off"
    no_agent              = false
    create_user           = true
  }

  timeouts = {
    create = "1h"
    delete = "20m"
    update = "20m"
    read   = "10m"
  }
}
```

## VMware VM Instance

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

  config_vmware = {
    resource_pool_id      = "pool-1"
    nested_virtualization = "off"
    no_agent              = true
    create_user           = false
    vmware_folder_id      = "group-v79"
  }

  timeouts = {
    create = "1h"
    delete = "20m"
    update = "20m"
    read   = "10m"
  }
}
```

## VMware VM Instance with Service Plan Options and Network Domain

Some Service Plans have options that must be set at the time of creation of the instance.

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

  network_domain_id = 27

  service_plan_options = {
    max_cores = 2
    max_memory = 2048
  }

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

  config_vmware = {
    resource_pool_id      = "pool-1"
    nested_virtualization = "off"
    no_agent              = true
    create_user           = false
    vmware_folder_id      = "group-v79"
  }

  timeouts = {
    create = "1h"
    delete = "20m"
    update = "20m"
    read   = "10m"
  }
}
```

## AWS Instance

```terraform
data "hpe_morpheus_cloud" "aws_cloud" {
  name = "QA Amazon"
}

data "hpe_morpheus_instance_type_layout" "aws" {
  name    = "Amazon VM"
  version = "22.04"
}

resource "hpe_morpheus_instance" "example" {
  name             = "TestInstance"
  cloud_id         = data.hpe_morpheus_cloud.aws_cloud.id
  layout_id        = data.hpe_morpheus_instance_type_layout.aws.id
  instance_type_id = 9

  group_id = 28
  plan_id  = 622

  instance_context = "dev"
  network_interfaces = [
    {
      network_id = 28
    }
  ]

  volumes = [
    {
      root_volume              = true
      name                     = "root"
      size                     = 100
      storage_type_id          = 23
      datastore_auto_selection = "auto"
    },
    {
      root_volume              = false
      name                     = "data"
      size                     = 100
      storage_type_id          = 23
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

  config_aws = {
    resource_pool_id      = "pool-12284"
    no_agent              = true
    create_user           = false
    security_groups = [
      { id = "sg-4eaf812b" },
    ]
  }

  timeouts = {
    create = "1h"
    delete = "20m"
    update = "20m"
    read   = "10m"
  }
}
```

## Azure Instance

-> Some clouds (such as Azure) provision onto a specific subnet rather than a network. Use
`subnet_id` in a `network_interfaces` block to target a subnet. `subnet_id`, `network_id` and
`network_group_id` are mutually exclusive. The subnet id is available from the Morpheus
`/api/options/zoneNetworkOptions` endpoint. Note that `subnet_id` is a write-only input: the API
reports the network the subnet resolves to, so on `import` the interface's `network_id` is populated
rather than `subnet_id`.

```terraform
data "hpe_morpheus_cloud" "azure_cloud" {
  name = "QA Azure"
}

data "hpe_morpheus_instance_type_layout" "azure" {
  name    = "Azure VM"
  version = "22.04"
}

resource "hpe_morpheus_instance" "example" {
  name             = "TestInstance"
  cloud_id         = data.hpe_morpheus_cloud.azure_cloud.id
  layout_id        = data.hpe_morpheus_instance_type_layout.azure.id
  instance_type_id = 9

  group_id = 28
  plan_id  = 622

  instance_context = "dev"
  network_interfaces = [
    {
      network_id = 28
      # Some clouds (such as Azure) provision onto a specific subnet. Use subnet_id
      # instead of network_id (the two, along with network_group_id, are mutually
      # exclusive). The subnet id is available from the Morpheus
      # /api/options/zoneNetworkOptions endpoint, e.g.:
      #   subnet_id = 5
    }
  ]

  volumes = [
    {
      root_volume              = true
      name                     = "root"
      size                     = 100
      storage_type_id          = 23
      datastore_auto_selection = "auto"
    }
  ]

  tags = [
    {
      name  = "terraform"
      value = "true"
    },
    {
      name  = "managed_by"
      value = "terraform"
    }
  ]

  config_azure = {
    resource_pool_id = "pool-12284"
    create_user      = false
    azure_region     = "eastus"
  }

  timeouts = {
    create = "1h"
    delete = "20m"
    update = "20m"
    read   = "10m"
  }
}
```

## BMaaS Instance

-> Note that Morpheus version `8.0.13` or later is required for BMaaS instance support.

For the `root` volume, to get the appropriate `storage_type_id` use the [Morpheus cli](https://support.hpe.com/hpesc/public/docDisplay?docId=sd00006978en_us&page=GUID-F3726B48-FFF6-4AAE-ABA4-366F626A544F.html).<br>
```bash
$ morpheus login -u <username> -p <password> --remote-url <morpheus-url>
$ morpheus storage-volume-types list hpeilo-raid1
```
Put that `id` in the `storage_type_id` field for the `root` volume below.  Leave the other fields unchanged.

```terraform
data "hpe_morpheus_cloud" "test" {
  name = "aCloud"
}

data "hpe_morpheus_environment" "test" {
  name = "anEnvironment"
}

data "hpe_morpheus_group" "test" {
  name = "aGroup"
}

data "hpe_morpheus_instance_type_layout" "test" {
  name = "Single ILO Server"
}

data "hpe_morpheus_role" "test" {
  name = "aRole"
}

data "hpe_morpheus_service_plan" "tp" {
    name                = "G3i"
    provision_type_code = "hpe-baremetal-plugin.provision"
}

resource "hpe_morpheus_instance" "example" {
  name             = "TestInstance"
  cloud_id         = data.hpe_morpheus_cloud.test.id # BM Cloud
  layout_id        = data.hpe_morpheus_instance_type_layout.test.id
  instance_type_id = 56 # BM Instance

  group_id = data.hpe_morpheus_group.test.id
  plan_id  = data.hpe_morpheus_service_plan.tp.id

  instance_context = "dev"

  network_interfaces = [
    {
      network_id      = 21
      ipMode          = ""
      network_type_id = 18
    },
    {
      network_id      = 21
      ipMode          = ""
      network_type_id = 18
    }
  ]

  volumes = [
    {
      root_volume     = true
      name            = "root"
      size            = 0
      storage_type_id = 76
      datastore_id    = null
    },
    {
      root_volume     = false
      name            = "data"
      size            = 16 # GB
      storage_type_id = 84
      datastore_id    = 11
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
      name  = "mytag"
      value = "true"
    },
    {
      name  = "jesskey"
      value = "terraform"
    }
  ]

  config = {
    imageId         = 231
    resourcePoolId  = "pool-1"
    isVpcSelectable = true
    serverId        = 155
    noAgent         = false
    createUser      = true
  }

  timeouts = {
    create = "2h"
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
- `config_aws` (Attributes) Configuration options for AWS instances. (see [below for nested schema](#nestedatt--config_aws))
- `config_azure` (Attributes) Configuration options for Azure instances. (see [below for nested schema](#nestedatt--config_azure))
- `config_hvm` (Attributes) Configuration options for HVM instances. (see [below for nested schema](#nestedatt--config_hvm))
- `config_vmware` (Attributes) Configuration options for VMware instances. (see [below for nested schema](#nestedatt--config_vmware))
- `description` (String) A description of the instance.
- `evars` (Attributes Set) Environment Variables, an array of objects that have name and value. (see [below for nested schema](#nestedatt--evars))
- `instance_context` (String) Environment
- `layout_size` (Number) Apply a multiply factor of containers/vms within the instance.
- `network_domain_id` (Number) The Network Domain ID to provision the instance into.
- `ports` (Attributes Set) The ports parameter is for port configuration.

The layout may have default ports, which are defined in node types, that are always configured. This parameter will be for additional custom ports to be opened. (see [below for nested schema](#nestedatt--ports))
- `service_plan_options` (Attributes) Custom options for selected service plan - the supported options depend on the service plan selected (see [below for nested schema](#nestedatt--service_plan_options))
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
- `subnet_id` (Number) id of the subnet to be used (required by some clouds, e.g. Azure). Cannot be used with 'network_id' or 'network_group_id'.

Read-Only:

- `id` (Number) The id of the network interface
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

- `id` (Number) The id of the network interface
- `name` (String) The name of the interface, e.g. 'eth0', 'eth1'
- `primary_interface` (Boolean) Is this interface the 'primary interface'?



<a id="nestedatt--config_aws"></a>
### Nested Schema for `config_aws`

Required:

- `resource_pool_id` (String) The id of the resource group to be used, can be prefixed with 'pool-'.  A resource pool group can be specified instead by prefixing its ID wih 'poolGroup-'.
- `security_groups` (Attributes List) a list of objects containing the ids of the AWS security groups to assign the instance to. (see [below for nested schema](#nestedatt--config_aws--security_groups))

Optional:

- `availability_zone_id` (String) The id of the AWS zone to provision the instance in.
- `create_user` (Boolean) Whether to create a user when provisioning the instance.  The default is 'false'
- `instance_profile` (String) The AWS IAM Profile to use for provisioning.
- `is_ec2` (Boolean) Whether this instance is an EC2 instance.  The default is 'false'.
- `kms_key_id` (String) The AWS KMS Key ID to use for provisioning.
- `no_agent` (Boolean) Whether to skip installing the Morpheus agent on the instance.  The default is 'true'
- `public_ip_type` (String) The type of public IP to associate with the instance.

<a id="nestedatt--config_aws--security_groups"></a>
### Nested Schema for `config_aws.security_groups`

Required:

- `id` (String) id of the AWS security group to assign the instance to.



<a id="nestedatt--config_azure"></a>
### Nested Schema for `config_azure`

Required:

- `resource_pool_id` (String) The id of the Azure resource group to provision the instance in, can be prefixed with 'pool-'. A resource pool group can be specified instead by prefixing its ID with 'poolGroup-'.

Optional:

- `availability_options` (String) The availability option for the instance (zone, set).
- `availability_set` (String) The availability set to use when availability_options is 'set'.
- `availability_zone` (String) The availability zone to use when availability_options is 'zone'.
- `azure_region` (String) The Azure region to provision the instance in.
- `azurefloating_ip` (String) Whether to assign a public IP to the instance (on, off).
- `azuresecurity_group_id` (String) The id of the Azure security group to assign the instance to.
- `boot_diagnostics` (String) Boot diagnostics setting (enable, enable_custom_storage).
- `create_user` (Boolean) Whether to create a user when provisioning the instance.
- `diagnostics_storage_account` (String) The diagnostics storage account to use when boot_diagnostics is 'enable_custom_storage'.
- `os_guest_diagnostics` (String) OS guest diagnostics setting (on, off).


<a id="nestedatt--config_hvm"></a>
### Nested Schema for `config_hvm`

Required:

- `resource_pool_id` (String) The id of the resource group to be used, can be prefixed with 'pool-'.  A resource pool group can be specified instead by prefixing its ID wih 'poolGroup-'.

Optional:

- `create_user` (Boolean) Whether to create a user when provisioning the instance.  The default is 'false'
- `kvm_host_id` (Number) The id of the KVM host to use for provisioning.
- `nested_virtualization` (String) Enable nested virtualization on the instance. Can be a number of valid string values:
   "on", "off", "0", "1", "true", "false", "yes", "no", "".  The default is "off".
- `no_agent` (Boolean) Whether to skip installing the Morpheus agent on the instance.  The default is 'true'


<a id="nestedatt--config_vmware"></a>
### Nested Schema for `config_vmware`

Required:

- `resource_pool_id` (String) The id of the resource group to be used, can be prefixed with 'pool-'.  A resource pool group can be specified instead by prefixing its ID wih 'poolGroup-'.

Optional:

- `create_user` (Boolean) Whether to create a user when provisioning the instance.  The default is 'false'
- `nested_virtualization` (String) Enable nested virtualization on the instance. Can be a number of valid string values:
   "on", "off", "0", "1", "true", "false", "yes", "no", "".  The default is "off".
- `no_agent` (Boolean) Whether to skip installing the Morpheus agent on the instance.  The default is 'true'
- `vmware_folder_id` (String) VMware folder external ID.


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


<a id="nestedatt--service_plan_options"></a>
### Nested Schema for `service_plan_options`

Optional:

- `cores_per_socket` (Number) The number of cores per socket to use for the instance.  This is only applicable for certain service plans.
- `max_cores` (Number) The maximum number of cores to use for the instance.  This is only applicable for certain service plans.
- `max_memory` (Number) The maximum amount of memory (in MB) to use for the instance.  This is only applicable for certain service plans.


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
- `storage_profile` (String) Storage Profile Code for the volume storage profile assignment. eg. `"kvm-cache-none"` or `"kvm-cache-directsync"`.
Use `/api/provision-types?code=kvm` to see the available `storageProfiles` for HVM and KVM.
- `storage_type_id` (Number) Identifier for LV type

Read-Only:

- `id` (Number) The id for the LV configuration being created.
