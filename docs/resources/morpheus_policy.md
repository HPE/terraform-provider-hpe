---
page_title: "hpe_morpheus_policy Resource - terraform-provider-hpe"
subcategory: "morpheus"
description: |-
  
---
# hpe_morpheus_policy (Resource)



Policies are different rules that can be applied to various Morpheus resources.

-> An `associated_resource_type` must be set (can be one of 'Global', 'Group', 'Cloud', 'User', 'Role', 'Network', or 'Plan').
-> If the `associated_resource_type` is not 'Global', then an `associated_resource_id` must also be set.

-> A `policy_type_code` must be set corresponding to the respective policy type.

-> A `config` block must be set corresponding to the respective policy type. Example configs for respective policy types are included below:

- [Approve Delete (deleteApproval)](#approve-delete-deleteapproval)
- [Approve Provision (provisionApproval)](#approve-provision-provisionapproval)
- [Approve Reconfigure (reconfigureApproval)](#approve-reconfigure-reconfigureapproval)
- [Approve Workflow Execute (workflowApproval)](#approve-workflow-execute-workflowapproval)
- [Backup Creation (createBackup)](#backup-creation-createbackup)
- [Backup Targets (backupStorage)](#backup-targets-backupstorage) Not currently supported
- [Budget (maxPrice)](#budget-maxprice)
- [Cluster Resource Name (serverNaming)](#cluster-resource-name-servernaming)
- [Cypher Access (cypher)](#cypher-access-cypher)
- [Delayed Delete (delayedRemoval)](#delayed-delete-delayedremoval)
- [Expiration (lifecycle)](#expiration-lifecycle)
- [File Share Storage Quota (storageShareQuota)](#file-share-storage-quota-storagesharequota)
- [Hostname (hostNaming)](#hostname-hostnaming)
- [Instance Name (naming)](#instance-name-naming)
- [Instance Networks (requiredNetwork)](#instance-networks-requirednetwork)
- [Max Containers (maxContainers)](#max-containers-maxcontainers)
- [Max Cores (maxCores)](#max-cores-maxcores)
- [Max Hosts (maxHosts)](#max-hosts-maxhosts)
- [Max Load Balancer Pools (maxPools)](#max-load-balancer-pools-maxpools)
- [Max Memory (maxMemory)](#max-memory-maxmemory)
- [Max Pool Members (maxPoolMembers)](#max-pool-members-maxpoolmembers)
- [Max Snapshots (maxSnapshots)](#max-snapshots-maxsnapshots)
- [Max Storage (maxStorage)](#max-storage-maxstorage)
- [Max Virtual Servers (maxVirtualServers)](#max-virtual-servers-maxvirtualservers)
- [Max VMs (maxVms)](#max-vms-maxvms)
- [Message of the Day (motd)](#message-of-the-day-motd)
- [Network Quota (maxNetworks)](#network-quota-maxnetworks)
- [Object Storage Quota (storageBucketQuota)](#object-storage-quota-storagebucketquota)
- [Power Scheduling (powerSchedule)](#power-scheduling-powerschedule)
- [Router Quota (maxRouters)](#router-quota-maxrouters)
- [Shutdown (shutdown)](#shutdown-shutdown)
- [Storage Server Storage Quota (storageServerQuota)](#storage-server-storage-quota-storageserverquota)
- [Tags (tags)](#tags-tags)
- [User Creation (createUser)](#user-creation-createuser)
- [User Group Creation (createUserGroup)](#user-group-creation-createusergroup)
- [Workflow (workflow)](#workflow-workflow)

## Config Examples

### Approve Delete (deleteApproval)

```terraform
# Approve Delete Policy
# Allowed associated_resource_types: Group, Cloud, User, Global, Label
# Tenant specification: allowed (can specify tenants array)
resource "hpe_morpheus_policy" "approve_delete" {
  name                     = "Approve Delete Policy"
  description              = "Require approval before deleting instances"
  associated_resource_type = "User"
  associated_resource_id   = 9969
  enabled                  = true

  policy_type = {
    code = "deleteApproval"
  }

  config = {
    accountIntegrationId = "1" # ID of your ServiceNow or approval integration
  }
}
```

### Approve Provision (provisionApproval)

```terraform
# Approve Provision Policy
# Allowed associated_resource_types: Group, Cloud, User, Global, Label
# Tenant specification: allowed (can specify tenants array)
resource "hpe_morpheus_policy" "approve_provision" {
  name                     = "Approve Provision Policy"
  description              = "Require approval before provisioning instances"
  associated_resource_type = "User"
  associated_resource_id   = 9969
  enabled                  = true

  policy_type = {
    code = "provisionApproval"
  }

  config = {
    accountIntegrationId = "1" # ID of your ServiceNow or approval integration
  }
}
```

### Approve Reconfigure (reconfigureApproval)

```terraform
# Approve Reconfigure Policy
# Allowed associated_resource_types: Group, Cloud, User, Global, Label
# Tenant specification: allowed (can specify tenants array)
resource "hpe_morpheus_policy" "approve_reconfigure" {
  name                     = "Approve Reconfigure Policy"
  description              = "Require approval before reconfiguring instances"
  associated_resource_type = "User"
  associated_resource_id   = 9969
  enabled                  = true

  policy_type = {
    code = "reconfigureApproval"
  }

  config = {
    accountIntegrationId = "1" # ID of your ServiceNow or approval integration
  }
}
```

### Approve Workflow Execute (workflowApproval)

```terraform
# Approve Workflow Execute Policy
# Allowed associated_resource_types: Group, Cloud, User, Global, Label
# Tenant specification: allowed (can specify tenants array)
resource "hpe_morpheus_policy" "approve_workflow" {
  name                     = "Approve Workflow Execute Policy"
  description              = "Require approval before executing workflows"
  associated_resource_type = "User"
  associated_resource_id   = 9969
  enabled                  = true

  policy_type = {
    code = "workflowApproval"
  }

  config = {
    accountIntegrationId = "1" # ID of your ServiceNow or approval integration
  }
}
```

### Backup Creation (createBackup)

```terraform
# Backup Creation Policy
# Allowed associated_resource_types: Group, Cloud, User, Global
# Tenant specification: allowed (can specify tenants array)
resource "hpe_morpheus_policy" "backup_creation" {
  name                     = "Backup Creation Policy"
  description              = "Enforce backup creation for instances"
  associated_resource_type = "User"
  associated_resource_id   = 9969
  enabled                  = true

  policy_type = {
    code = "createBackup"
  }

  config = {
    createBackupType = "user" # Options: "user" (user configurable), "fixed" (strict pattern)
    createBackup     = true   # Enforce backup creation
  }
}
```

### Backup Targets (backupStorage)
Not currently supported

### Budget (maxPrice)

```terraform
# Budget Policy - Limits instance costs
# Allowed associated_resource_types: Group, Cloud, User, Global, Plan
# Tenant specification: allowed (can specify tenants array)
resource "hpe_morpheus_policy" "budget" {
  name                     = "Budget Policy"
  description              = "Limit maximum instance costs"
  associated_resource_type = "User"
  associated_resource_id   = 9969
  enabled                  = true

  policy_type = {
    code = "maxPrice"
  }

  config = {
    maxPrice         = "1000"  # Maximum price limit
    maxPriceCurrency = "USD"   # Currency code
    maxPriceUnit     = "month" # Options: "hour", "month"
  }
}
```

### Cluster Resource Name (serverNaming)

```terraform
# Cluster Resource Name Policy - Enforces naming conventions for cluster resources
# Allowed associated_resource_types: Group, Cloud, User, Global
# Tenant specification: allowed (can specify tenants array)
resource "hpe_morpheus_policy" "cluster_naming" {
  name                     = "Cluster Resource Naming Policy"
  description              = "Enforce naming for cluster resources"
  associated_resource_type = "User"
  associated_resource_id   = 9969
  enabled                  = true

  policy_type = {
    code = "serverNaming"
  }

  config = {
    serverNamingType     = "user"                                        # Options: "user" (user configurable), "fixed" (strict pattern)
    serverNamingPattern  = "cluster-$${groupCode}-$${type}-$${sequence}" # Naming pattern with variables
    serverNamingConflict = true                                          # Auto-resolve conflicts
  }
}
```

### Cypher Access (cypher)

```terraform
# Cypher Access Policy - Controls access to Cypher secrets
# Allowed associated_resource_types: User, Global, Role
# Tenant specification: allowed (can specify tenants array)
resource "hpe_morpheus_policy" "cypher_access" {
  name                     = "Cypher Access Policy"
  description              = "Control Cypher key access permissions"
  associated_resource_type = "User"
  associated_resource_id   = 9969
  enabled                  = true

  policy_type = {
    code = "cypher"
  }

  config = {
    keyPattern = "secret/*" # Pattern to match Cypher keys (e.g., "secret/*", "password/*")
    read       = true       # Allow read access
    write      = true       # Allow write access
    update     = true       # Allow update access
    delete     = false      # Deny delete access
    list       = true       # Allow list access
  }
}
```

### Delayed Delete (delayedRemoval)

```terraform
# Delayed Delete Policy - Delays instance deletion
# Allowed associated_resource_types: Group, Cloud, User, Global
# Tenant specification: allowed (can specify tenants array)
resource "hpe_morpheus_policy" "delayed_delete" {
  name                     = "Delayed Delete Policy"
  description              = "Delay instance deletion by specified days"
  associated_resource_type = "User"
  associated_resource_id   = 9969
  enabled                  = true

  policy_type = {
    code = "delayedRemoval"
  }

  config = {
    removalAge = "30" # Number of days to delay deletion
  }
}
```

### Expiration (lifecycle)

```terraform
# Expiration Policy - Sets instance expiration and renewal options
# Allowed associated_resource_types: Group, Cloud, User, Global
# Tenant specification: allowed (can specify tenants array)
resource "hpe_morpheus_policy" "expiration" {
  name                     = "Expiration Policy"
  description              = "Set instance expiration and renewal policies"
  associated_resource_type = "User"
  associated_resource_id   = 9969
  enabled                  = true

  policy_type = {
    code = "lifecycle"
  }

  config = {
    lifecycleType                     = "user"                      # Options: "user" (user configurable), "fixed" (fixed expiration)
    lifecycleAge                      = "30"                        # Days until expiration
    lifecycleRenewal                  = "7"                         # Days for renewal window
    lifecycleNotify                   = "1"                         # Days before expiration to notify
    lifecycleMessage                  = "Instance will expire soon" # Notification message
    lifecycleAutoRenew                = "on"                        # Options: "on", "off"
    lifecycleAllowExtend              = "off"                       # Options: "on", "off" - allow users to extend
    lifecycleExtensionsBeforeApproval = "0"                         # Number of extensions before requiring approval
    lifecycleHideFixed                = false                       # Hide fixed expiration date from users
  }
}
```

### File Share Storage Quota (storageShareQuota)

```terraform
# File Share Storage Quota Policy - Limits file share storage
# Allowed associated_resource_types: User, Global
# Tenant specification: allowed (can specify tenants array)
resource "hpe_morpheus_policy" "file_share_quota" {
  name                     = "File Share Storage Quota Policy"
  description              = "Limit file share storage usage"
  associated_resource_type = "User"
  associated_resource_id   = 9969
  enabled                  = true

  policy_type = {
    code = "storageShareQuota"
  }

  config = {
    maxStorage = "1000" # Maximum storage in GB
  }
}
```

### Hostname (hostNaming)

```terraform
# Hostname Policy - Enforces hostname naming conventions
# Allowed associated_resource_types: Group, Cloud, User, Global
# Tenant specification: allowed (can specify tenants array)
resource "hpe_morpheus_policy" "hostname" {
  name                     = "Hostname Policy"
  description              = "Enforce hostname naming conventions"
  associated_resource_type = "User"
  associated_resource_id   = 9969
  enabled                  = true

  policy_type = {
    code = "hostNaming"
  }

  config = {
    hostNamingType    = "user"                                     # Options: "user" (user configurable), "fixed" (strict pattern)
    hostNamingPattern = "host-$${groupCode}-$${type}-$${sequence}" # Naming pattern with variables
  }
}
```

### Instance Name (naming)

```terraform
# Instance Name Policy - Enforces instance naming conventions
# Allowed associated_resource_types: Group, Cloud, User, Global
# Tenant specification: allowed (can specify tenants array)
resource "hpe_morpheus_policy" "instance_naming" {
  name                     = "Instance Name Policy"
  description              = "Enforce instance naming conventions"
  associated_resource_type = "User"
  associated_resource_id   = 9969
  enabled                  = true

  policy_type = {
    code = "naming"
  }

  config = {
    namingType     = "user"                                   # Options: "user" (user configurable), "fixed" (strict pattern)
    namingPattern  = "vm-$${groupCode}-$${type}-$${sequence}" # Naming pattern with variables
    namingConflict = true                                     # Auto-resolve conflicts
  }
}
```

### Instance Networks (requiredNetwork)

```terraform
# Instance Networks Policy - Requires specific networks for instances
# Allowed associated_resource_types: Group, Cloud
# Tenant specification: NOT allowed (allowOnTenant = false)
resource "hpe_morpheus_policy" "required_networks" {
  name                     = "Instance Networks Policy"
  description              = "Require specific networks for instances"
  associated_resource_type = "Cloud"
  associated_resource_id   = 9969
  enabled                  = true

  policy_type = {
    code = "requiredNetwork"
  }

  config = {
    requiredNetworks = [100, 200] # Array of required network IDs
  }
}
```

### Max Containers (maxContainers)

```terraform
# Max Containers Policy - Limits container count
# Allowed associated_resource_types: Group, Cloud, User, Global, Plan
# Tenant specification: allowed (can specify tenants array)
resource "hpe_morpheus_policy" "max_containers" {
  name                     = "Max Containers Policy"
  description              = "Limit maximum container count"
  associated_resource_type = "User"
  associated_resource_id   = 9969
  enabled                  = true

  policy_type = {
    code = "maxContainers"
  }

  config = {
    maxContainers = "50" # Maximum number of containers
  }
}
```

### Max Cores (maxCores)

```terraform
# Max Cores Policy - Limits CPU cores
# Allowed associated_resource_types: Group, Cloud, User, Global, Plan
# Tenant specification: allowed (can specify tenants array)
resource "hpe_morpheus_policy" "max_cores" {
  name                     = "Max Cores Policy"
  description              = "Limit maximum CPU cores"
  associated_resource_type = "User"
  associated_resource_id   = 9969
  enabled                  = true

  policy_type = {
    code = "maxCores"
  }

  config = {
    maxCores          = "32"  # Maximum number of CPU cores
    excludeContainers = "off" # Options: "on", "off" - exclude containers from count
  }
}
```

### Max Hosts (maxHosts)

```terraform
# Max Hosts Policy - Limits host count
# Allowed associated_resource_types: Group, Cloud, User, Global, Plan
# Tenant specification: allowed (can specify tenants array)
resource "hpe_morpheus_policy" "max_hosts" {
  name                     = "Max Hosts Policy"
  description              = "Limit maximum host count"
  associated_resource_type = "User"
  associated_resource_id   = 9969
  enabled                  = true

  policy_type = {
    code = "maxHosts"
  }

  config = {
    maxHosts = "10" # Maximum number of hosts
  }
}
```

### Max Load Balancer Pools (maxPools)

```terraform
# Max Load Balancer Pools Policy - Limits load balancer pools
# Allowed associated_resource_types: Cloud, User, Global
# Tenant specification: allowed (can specify tenants array)
resource "hpe_morpheus_policy" "max_pools" {
  name                     = "Max Load Balancer Pools Policy"
  description              = "Limit maximum load balancer pools"
  associated_resource_type = "User"
  associated_resource_id   = 9969
  enabled                  = true

  policy_type = {
    code = "maxPools"
  }

  config = {
    maxPools = "5" # Maximum number of load balancer pools
  }
}
```

### Max Memory (maxMemory)

```terraform
# Max Memory Policy - Limits memory allocation
# Allowed associated_resource_types: Group, Cloud, User, Global, Plan
# Tenant specification: allowed (can specify tenants array)
resource "hpe_morpheus_policy" "max_memory" {
  name                     = "Max Memory Policy"
  description              = "Limit maximum memory allocation"
  associated_resource_type = "User"
  associated_resource_id   = 9969
  enabled                  = true

  policy_type = {
    code = "maxMemory"
  }

  config = {
    maxMemory         = "8"   # Maximum memory in GB
    excludeContainers = "off" # Options: "on", "off" - exclude containers from count
  }
}
```

### Max Pool Members (maxPoolMembers)

```terraform
# Max Pool Members Policy - Limits pool members
# Allowed associated_resource_types: Cloud, Global
# Tenant specification: allowed (can specify tenants array)
resource "hpe_morpheus_policy" "max_pool_members" {
  name                     = "Max Pool Members Policy"
  description              = "Limit maximum pool members"
  associated_resource_type = "Cloud"
  associated_resource_id   = 9969
  enabled                  = true

  policy_type = {
    code = "maxPoolMembers"
  }

  config = {
    maxPoolMembers = "10" # Maximum number of pool members
  }
}
```

### Max Snapshots (maxSnapshots)

```terraform
# Max Snapshots Policy - Limits snapshots per VM
# Allowed associated_resource_types: Group, Cloud, User, Global
# Tenant specification: allowed (can specify tenants array)
resource "hpe_morpheus_policy" "max_snapshots" {
  name                     = "Max Snapshots Policy"
  description              = "Limit maximum snapshots per VM"
  associated_resource_type = "User"
  associated_resource_id   = 9969
  enabled                  = true

  policy_type = {
    code = "maxSnapshots"
  }

  config = {
    maxSnapshots = "5" # Maximum number of snapshots per VM
  }
}
```

### Max Storage (maxStorage)

```terraform
# Max Storage Policy - Limits storage allocation
# Allowed associated_resource_types: Group, Cloud, User, Global, Plan
# Tenant specification: allowed (can specify tenants array)
resource "hpe_morpheus_policy" "max_storage" {
  name                     = "Max Storage Policy"
  description              = "Limit maximum storage allocation"
  associated_resource_type = "User"
  associated_resource_id   = 9969
  enabled                  = true

  policy_type = {
    code = "maxStorage"
  }

  config = {
    maxStorage        = "1000" # Maximum storage in GB
    excludeContainers = "off"  # Options: "on", "off" - exclude containers from count
  }
}
```

### Max Virtual Servers (maxVirtualServers)

```terraform
# Max Virtual Servers Policy - Limits virtual server count
# Allowed associated_resource_types: Cloud, Global
# Tenant specification: allowed (can specify tenants array)
resource "hpe_morpheus_policy" "max_virtual_servers" {
  name                     = "Max Virtual Servers Policy"
  description              = "Limit maximum virtual server count"
  associated_resource_type = "Cloud"
  associated_resource_id   = 9969
  enabled                  = true

  policy_type = {
    code = "maxVirtualServers"
  }

  config = {
    maxVirtualServers = "10" # Maximum number of virtual servers
  }
}
```

### Max VMs (maxVms)

```terraform
# Max VMs Policy - Limits VM count
# Allowed associated_resource_types: Group, Cloud, User, Global, Network, Plan
# Tenant specification: allowed (can specify tenants array)
resource "hpe_morpheus_policy" "max_vms" {
  name                     = "Max VMs Policy"
  description              = "Limit maximum VM count"
  associated_resource_type = "User"
  associated_resource_id   = 9969
  enabled                  = true

  policy_type = {
    code = "maxVms"
  }

  config = {
    maxVms = "20" # Maximum number of VMs
  }
}
```

### Message of the Day (motd)

```terraform
# Message of the Day (MOTD) Policy - Displays login messages
# Allowed associated_resource_types: Global
# Tenant specification: allowed (can specify tenants array)
resource "hpe_morpheus_policy" "motd" {
  name                     = "MOTD Policy"
  description              = "Display message of the day on login"
  associated_resource_type = "Global"
  enabled                  = true

  policy_type = {
    code = "motd"
  }

  config = {
    "motd.title"     = "Welcome"                          # Message title
    "motd.message"   = "Welcome to the Morpheus platform" # Message content
    "motd.type"      = "info"                             # Options: "info", "warning", "danger"
    "motd._fullPage" = "off"                              # Options: "on", "off" - display full page
  }
}
```

### Network Quota (maxNetworks)

```terraform
# Network Quota Policy - Limits network count
# Allowed associated_resource_types: Group, Cloud, User, Global
# Tenant specification: allowed (can specify tenants array)
resource "hpe_morpheus_policy" "network_quota" {
  name                     = "Network Quota Policy"
  description              = "Limit maximum network count"
  associated_resource_type = "User"
  associated_resource_id   = 9969
  enabled                  = true

  policy_type = {
    code = "maxNetworks"
  }

  config = {
    maxNetworks = "10" # Maximum number of networks
  }
}
```

### Object Storage Quota (storageBucketQuota)

```terraform
# Object Storage Quota Policy - Limits object storage
# Allowed associated_resource_types: User, Global
# Tenant specification: allowed (can specify tenants array)
resource "hpe_morpheus_policy" "object_storage_quota" {
  name                     = "Object Storage Quota Policy"
  description              = "Limit object storage usage"
  associated_resource_type = "User"
  associated_resource_id   = 9969
  enabled                  = true

  policy_type = {
    code = "storageBucketQuota"
  }

  config = {
    maxStorage = "1000" # Maximum storage in GB
  }
}
```

### Power Scheduling (powerSchedule)

```terraform
# Power Scheduling Policy - Enforces power schedules
# Allowed associated_resource_types: Group, Cloud, User, Global
# Tenant specification: allowed (can specify tenants array)
resource "hpe_morpheus_policy" "power_schedule" {
  name                     = "Power Scheduling Policy"
  description              = "Enforce power schedules for instances"
  associated_resource_type = "User"
  associated_resource_id   = 9969
  enabled                  = true

  policy_type = {
    code = "powerSchedule"
  }

  config = {
    powerScheduleType      = "user" # Options: "user" (user configurable), "fixed" (strict schedule)
    powerSchedule          = "1"    # ID of the power schedule
    powerScheduleHideFixed = false  # Hide fixed schedule from users
  }
}
```

### Router Quota (maxRouters)

```terraform
# Router Quota Policy - Limits router count
# Allowed associated_resource_types: Group, Cloud, User, Global
# Tenant specification: allowed (can specify tenants array)
resource "hpe_morpheus_policy" "router_quota" {
  name                     = "Router Quota Policy"
  description              = "Limit maximum router count"
  associated_resource_type = "User"
  associated_resource_id   = 9969
  enabled                  = true

  policy_type = {
    code = "maxRouters"
  }

  config = {
    maxRouters = "5" # Maximum number of routers
  }
}
```

### Shutdown (shutdown)

```terraform
# Shutdown Policy - Auto-shutdown idle instances
# Allowed associated_resource_types: Group, Cloud, User, Global
# Tenant specification: allowed (can specify tenants array)
resource "hpe_morpheus_policy" "shutdown" {
  name                     = "Shutdown Policy"
  description              = "Auto-shutdown idle instances"
  associated_resource_type = "User"
  associated_resource_id   = 9969
  enabled                  = true

  policy_type = {
    code = "shutdown"
  }

  config = {
    shutdownType                     = "user"                        # Options: "user" (user configurable), "fixed" (strict shutdown)
    shutdownAge                      = "30"                          # Days instance is allowed to run before shutdown
    shutdownRenewal                  = "7"                           # If the instance is renewed, this is the number of day increments the shutdown date is increased by.
    shutdownNotify                   = "1"                           # Days before shutdown to notify via email
    shutdownMessage                  = "Instance will shutdown soon" # Notification message
    shutdownAutoRenew                = "on"                          # Options: "on", "off"
    shutdownAllowExtend              = "off"                         # Options: "on", "off" - allow users to extend
    shutdownExtensionsBeforeApproval = "0"                           # Number of extensions before requiring approval
    shutdownHideFixed                = false                         # Hide shutdown if fixed value
  }
}
```

### Storage Server Storage Quota (storageServerQuota)

```terraform
# Storage Server Storage Quota Policy - Limits storage on specific storage server
# Allowed associated_resource_types: Global
# Tenant specification: allowed (can specify tenants array)
resource "hpe_morpheus_policy" "storage_server_quota" {
  name                     = "Storage Server Storage Quota Policy"
  description              = "Limit storage usage on specific storage server"
  associated_resource_type = "Global"
  associated_resource_id   = 9969
  enabled                  = true

  policy_type = {
    code = "storageServerQuota"
  }

  config = {
    storageServerId = "1"    # ID of the storage server
    maxStorage      = "1000" # Maximum storage in GB
  }
}
```

### Tags (tags)

```terraform
# Tags Policy - Enforces instance tagging
# Allowed associated_resource_types: Group, Cloud, User, Global
# Tenant specification: allowed (can specify tenants array)
resource "hpe_morpheus_policy" "tags" {
  name                     = "Tags Policy"
  description              = "Enforce instance tagging requirements"
  associated_resource_type = "User"
  associated_resource_id   = 9969
  enabled                  = true

  policy_type = {
    code = "tags"
  }

  config = {
    strict      = true          # Strict enforcement
    key         = "environment" # Tag key to enforce
    value       = "production"  # Tag value (optional, can be left empty for any value)
    valueListId = ""            # ID of value from value list (optional)
  }
}
```

### User Creation (createUser)

```terraform
# User Creation Policy - Controls user creation on instances
# Allowed associated_resource_types: Group, Cloud, User, Global
# Tenant specification: allowed (can specify tenants array)
resource "hpe_morpheus_policy" "user_creation" {
  name                     = "User Creation Policy"
  description              = "Control user creation on provisioned instances"
  associated_resource_type = "User"
  associated_resource_id   = 9969
  enabled                  = true

  policy_type = {
    code = "createUser"
  }

  config = {
    createUserType = "user" # Options: "user" (user configurable), "fixed"
    createUser     = true   # Enforce user creation
  }
}
```

### User Group Creation (createUserGroup)

```terraform
# User Group Creation Policy - Assigns default user group
# Allowed associated_resource_types: Group, Cloud, User, Global
# Tenant specification: allowed (can specify tenants array)
resource "hpe_morpheus_policy" "user_group_creation" {
  name                     = "User Group Creation Policy"
  description              = "Assign default user group for created users"
  associated_resource_type = "User"
  associated_resource_id   = 9969
  enabled                  = true

  policy_type = {
    code = "createUserGroup"
  }

  config = {
    userGroup = "1" # ID of the user group to assign
  }
}
```

### Workflow (workflow)

```terraform
# Workflow Policy - Executes workflow on provision
# Allowed associated_resource_types: Group, Cloud, User, Global
# Tenant specification: allowed (can specify tenants array)
# Note: This example uses the morpheus external provider to create a workflow resource
# because the hpe provider does not yet have a workflow resource implemented.
# You will need to configure the morpheus provider in your terraform configuration.

terraform {
  required_providers {
    hpe = {
      source  = "HPE/hpe"
      version = "= 0.3.0"
    }
    morpheus = {
      source  = "gomorpheus/morpheus"
      version = "~> 0.13.2"
    }
  }
}

resource "morpheus_operational_workflow" "example" {
  name        = "Example Policy Workflow"
  description = "Example workflow for policy testing"
}

resource "hpe_morpheus_policy" "workflow" {
  name                     = "Workflow Policy"
  description              = "Execute workflow on instance provision"
  associated_resource_type = "User"
  associated_resource_id   = 9969
  enabled                  = true

  policy_type = {
    code = "workflow"
  }

  config = {
    workflowId = morpheus_operational_workflow.example.id # ID of the workflow to execute
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `associated_resource_type` (String) Type of the resource this policy is associated with, can be 'Global', 'Group', 'Cloud', 'User', 'Role', 'Network', 'Plan' or 'Label'
- `name` (String) A name for the policy
- `policy_type` (Attributes) (see [below for nested schema](#nestedatt--policy_type))

### Optional

- `associated_resource_id` (Number) The ID of the resource this policy is associated with, e.g. Group, Cloud, User, Role, Network, Plan, Label
- `config` (Dynamic) Generic Policy Configuration
- `config_approval` (Attributes) Configuration for the following policy types:
	- Approve Delete (deleteApproval)
	- Approve Provision (provisionApproval)
	- Approve Reconfigure (reconfigureApproval)
	- Approve Workflow Execute (workflowApproval) (see [below for nested schema](#nestedatt--config_approval))
- `config_backup_storage` (Attributes) Configuration for the following policy type:
	- Backup Targets (backupStorage) (see [below for nested schema](#nestedatt--config_backup_storage))
- `config_create_backup` (Attributes) Configuration for the following policy type:
	- Backup Creation (createBackup) (see [below for nested schema](#nestedatt--config_create_backup))
- `config_create_user` (Attributes) Configuration for the following policy type:
	- User Creation (createUser) (see [below for nested schema](#nestedatt--config_create_user))
- `config_create_user_group` (Attributes) Configuration for the following policy type:
	- User Group Creation (createUserGroup) (see [below for nested schema](#nestedatt--config_create_user_group))
- `config_cypher` (Attributes) Configuration for the following policy type:
	- Cypher Access (cypher) (see [below for nested schema](#nestedatt--config_cypher))
- `config_delayed_removal` (Attributes) Configuration for the following policy type:
	- Delayed Delete (delayedRemoval) (see [below for nested schema](#nestedatt--config_delayed_removal))
- `config_host_naming` (Attributes) Configuration for the following policy type:
	- Hostname (hostNaming) (see [below for nested schema](#nestedatt--config_host_naming))
- `config_lifecycle` (Attributes) Configuration for the following policy type:
	- Expiration (lifecycle) (see [below for nested schema](#nestedatt--config_lifecycle))
- `config_max_containers` (Attributes) Configuration for the following policy type:
	- Max Containers (maxContainers) (see [below for nested schema](#nestedatt--config_max_containers))
- `config_max_cores` (Attributes) Configuration for the following policy type:
	- Max Cores (maxCores) (see [below for nested schema](#nestedatt--config_max_cores))
- `config_max_hosts` (Attributes) Configuration for the following policy type:
	- Max Hosts (maxHosts) (see [below for nested schema](#nestedatt--config_max_hosts))
- `config_max_memory` (Attributes) Configuration for the following policy type:
	- Max Memory (maxMemory) (see [below for nested schema](#nestedatt--config_max_memory))
- `config_max_networks` (Attributes) Configuration for the following policy type:
	- Network Quota (maxNetworks) (see [below for nested schema](#nestedatt--config_max_networks))
- `config_max_pool_members` (Attributes) Configuration for the following policy type:
	- Max Pool Members (maxPoolMembers) (see [below for nested schema](#nestedatt--config_max_pool_members))
- `config_max_pools` (Attributes) Configuration for the following policy type:
	- Max Load Balancer Pools (maxPools) (see [below for nested schema](#nestedatt--config_max_pools))
- `config_max_price` (Attributes) Configuration for the following policy type:
	- Budget (maxPrice) (see [below for nested schema](#nestedatt--config_max_price))
- `config_max_routers` (Attributes) Configuration for the following policy type:
	- Router Quota (maxRouters) (see [below for nested schema](#nestedatt--config_max_routers))
- `config_max_snapshots` (Attributes) Configuration for the following policy type:
	- Max Snapshots (maxSnapshots) (see [below for nested schema](#nestedatt--config_max_snapshots))
- `config_max_storage` (Attributes) Configuration for the following policy types:
	- Max Storage (maxStorage)
	- Object Storage Quota (storageBucketQuota)
	- File Share Storage Quota (storageShareQuota) (see [below for nested schema](#nestedatt--config_max_storage))
- `config_max_virtual_servers` (Attributes) Configuration for the following policy type:
	- Max Virtual Servers (maxVirtualServers) (see [below for nested schema](#nestedatt--config_max_virtual_servers))
- `config_max_vms` (Attributes) Configuration for the following policy type:
	- Max VMs (maxVms) (see [below for nested schema](#nestedatt--config_max_vms))
- `config_motd` (Attributes) Configuration for the following policy type:
	- Message of the Day (motd) (see [below for nested schema](#nestedatt--config_motd))
- `config_naming` (Attributes) Configuration for the following policy type:
	- Instance Name (naming) (see [below for nested schema](#nestedatt--config_naming))
- `config_power_schedule` (Attributes) Configuration for the following policy type:
	- Power Scheduling (powerSchedule) (see [below for nested schema](#nestedatt--config_power_schedule))
- `config_required_network` (Attributes) Configuration for the following policy type:
	- Instance Networks (requiredNetwork) (see [below for nested schema](#nestedatt--config_required_network))
- `config_server_naming` (Attributes) Configuration for the following policy type:
	- Cluster Resource Name (serverNaming) (see [below for nested schema](#nestedatt--config_server_naming))
- `config_shutdown` (Attributes) Configuration for the following policy type:
	- Shutdown (shutdown) (see [below for nested schema](#nestedatt--config_shutdown))
- `config_storage_server_quota` (Attributes) Configuration for the following policy type:
	- Storage Server Storage Quota (storageServerQuota) (see [below for nested schema](#nestedatt--config_storage_server_quota))
- `config_tags` (Attributes) Configuration for the following policy type:
	- Tags (tags) (see [below for nested schema](#nestedatt--config_tags))
- `config_workflow` (Attributes) Configuration for the following policy type:
	- Workflow (workflow) (see [below for nested schema](#nestedatt--config_workflow))
- `description` (String) A description for the policy
- `each_user` (Boolean) Apply individually to each user in role.  Only when `refType` equals `Role`
- `enabled` (Boolean) Set to false to disable
- `tenants` (Set of Number) Array of tenants to scope the policy to

### Read-Only

- `cloud` (Attributes) (see [below for nested schema](#nestedatt--cloud))
- `group` (Attributes) (see [below for nested schema](#nestedatt--group))
- `id` (Number) The ID of this resource.
- `owner` (Attributes) (see [below for nested schema](#nestedatt--owner))
- `role` (Attributes) (see [below for nested schema](#nestedatt--role))
- `user` (Attributes) (see [below for nested schema](#nestedatt--user))

<a id="nestedatt--policy_type"></a>
### Nested Schema for `policy_type`

Required:

- `code` (String) The policy type code. See `Retrieves all Policy Types` endpoint for listing.

Read-Only:

- `id` (Number)
- `name` (String)


<a id="nestedatt--config_approval"></a>
### Nested Schema for `config_approval`

Optional:

- `account_integration_id` (String) ID of your ServiceNow or approval integration
- `flow_id` (String) ID of ServiceNow Flow (set if workflowType is 'flow')
- `workflow_id` (String) ID of legacy ServiceNow workflow (set if workflowType is 'workflow')
- `workflow_type` (String) Options: "workflow" (legacy workflow), "flow" (ServiceNow Flow)


<a id="nestedatt--config_backup_storage"></a>
### Nested Schema for `config_backup_storage`

Optional:

- `backup_storage_ids` (Set of Number) Array of backup storage IDs to restrict available backup targets


<a id="nestedatt--config_create_backup"></a>
### Nested Schema for `config_create_backup`

Optional:

- `create_backup` (Boolean) Enforce backup creation
- `create_backup_type` (String) Options: "user" (user configurable), "fixed" (strict pattern)


<a id="nestedatt--config_create_user"></a>
### Nested Schema for `config_create_user`

Optional:

- `create_user` (Boolean) Enforce user creation
- `create_user_type` (String) Options: "user" (user configurable), "fixed"


<a id="nestedatt--config_create_user_group"></a>
### Nested Schema for `config_create_user_group`

Optional:

- `user_group` (String) ID of the user group to assign


<a id="nestedatt--config_cypher"></a>
### Nested Schema for `config_cypher`

Optional:

- `delete` (Boolean) Allow delete access
- `key_pattern` (String) Pattern to match Cypher keys (e.g., "secret/*", "password/*")
- `list` (Boolean) Allow list access
- `read` (Boolean) Allow read access
- `update` (Boolean) Allow update access
- `write` (Boolean) Allow write access


<a id="nestedatt--config_delayed_removal"></a>
### Nested Schema for `config_delayed_removal`

Optional:

- `removal_age` (String) Number of days to delay deletion


<a id="nestedatt--config_host_naming"></a>
### Nested Schema for `config_host_naming`

Optional:

- `host_naming_pattern` (String) Name pattern uses ${variable} string interpolation.  Available variables are:<br>groupName, groupCode, cloudName, cloudCode, type, accountId, account, accountType, platform, username, userId, userInitials, provisionType
- `host_naming_type` (String) Options: "user" (user configurable), "fixed" (strict pattern)


<a id="nestedatt--config_lifecycle"></a>
### Nested Schema for `config_lifecycle`

Optional:

- `account_integration_id` (String) ID of your ServiceNow or approval integration
- `flow_id` (String) ID of ServiceNow Flow (set if workflowType is 'flow')
- `lifecycle_age` (String) Days until expiration
- `lifecycle_allow_extend` (Boolean)
- `lifecycle_auto_renew` (Boolean)
- `lifecycle_extensions_before_approval` (String) Number of extensions before requiring approval
- `lifecycle_hide_fixed` (Boolean) Hide fixed expiration from users
- `lifecycle_message` (String) Notification message
- `lifecycle_notify` (String) Days before expiration to notify
- `lifecycle_renewal` (String) Days for renewal window
- `lifecycle_type` (String) Options: "user" (user configurable), "fixed" (fixed expiration)
- `lifecycle_workflow_id` (String) ID of legacy ServiceNow workflow (set if workflowType is 'workflow')
- `workflow_type` (String) Options: "workflow" (legacy workflow), "flow" (ServiceNow Flow)


<a id="nestedatt--config_max_containers"></a>
### Nested Schema for `config_max_containers`

Optional:

- `max_containers` (String) Max Containers


<a id="nestedatt--config_max_cores"></a>
### Nested Schema for `config_max_cores`

Optional:

- `exclude_containers` (Boolean)
- `max_cores` (String) Max Cores


<a id="nestedatt--config_max_hosts"></a>
### Nested Schema for `config_max_hosts`

Optional:

- `max_hosts` (String) Max Hosts


<a id="nestedatt--config_max_memory"></a>
### Nested Schema for `config_max_memory`

Optional:

- `exclude_containers` (Boolean)
- `max_memory` (String) Max Memory (GB)


<a id="nestedatt--config_max_networks"></a>
### Nested Schema for `config_max_networks`

Optional:

- `max_networks` (String) Max Networks


<a id="nestedatt--config_max_pool_members"></a>
### Nested Schema for `config_max_pool_members`

Optional:

- `max_pool_members` (String) Max Pool Members


<a id="nestedatt--config_max_pools"></a>
### Nested Schema for `config_max_pools`

Optional:

- `max_pools` (String) Max Pools


<a id="nestedatt--config_max_price"></a>
### Nested Schema for `config_max_price`

Optional:

- `max_price` (Number) Maximum price limit
- `max_price_currency` (String) Currency code (e.g., USD)
- `max_price_unit` (String) Options: "hour", "month"


<a id="nestedatt--config_max_routers"></a>
### Nested Schema for `config_max_routers`

Optional:

- `max_routers` (String) Max Routers


<a id="nestedatt--config_max_snapshots"></a>
### Nested Schema for `config_max_snapshots`

Optional:

- `max_snapshots` (String) Max Snapshots


<a id="nestedatt--config_max_storage"></a>
### Nested Schema for `config_max_storage`

Optional:

- `exclude_containers` (Boolean)
- `max_storage` (String) Max Storage (GB)


<a id="nestedatt--config_max_virtual_servers"></a>
### Nested Schema for `config_max_virtual_servers`

Optional:

- `max_virtual_servers` (String) Max Virtual Servers


<a id="nestedatt--config_max_vms"></a>
### Nested Schema for `config_max_vms`

Optional:

- `max_vms` (String) Max VMs


<a id="nestedatt--config_motd"></a>
### Nested Schema for `config_motd`

Optional:

- `motddate` (String) Display date for message
- `motdmessage` (String) Message content
- `motdtitle` (String) Message title
- `motdtype` (String) Options: "info", "warning", "critical"


<a id="nestedatt--config_naming"></a>
### Nested Schema for `config_naming`

Optional:

- `naming_conflict` (Boolean) Auto-resolve conflicts
- `naming_pattern` (String) Name pattern uses ${variable} string interpolation.  Available variables are:<br>groupName, groupCode, cloudName, cloudCode, type, accountId, account, accountType, platform, username, userId, userInitials, provisionType
- `naming_type` (String) Options: "user" (user configurable), "fixed" (strict pattern)


<a id="nestedatt--config_power_schedule"></a>
### Nested Schema for `config_power_schedule`

Optional:

- `power_schedule` (String) ID of the power schedule
- `power_schedule_hide_fixed` (Boolean) Hide fixed schedule from users
- `power_schedule_type` (String) Options: "user" (user configurable), "fixed" (strict schedule)


<a id="nestedatt--config_required_network"></a>
### Nested Schema for `config_required_network`

Optional:

- `required_networks` (Set of Number) Array of required network IDs


<a id="nestedatt--config_server_naming"></a>
### Nested Schema for `config_server_naming`

Optional:

- `server_naming_conflict` (Boolean) Auto-resolve conflicts
- `server_naming_pattern` (String) Name pattern uses ${variable} string interpolation.  Available variables are:<br>groupName, groupCode, cloudName, cloudCode, type, accountId, account, accountType, platform, username, userId, userInitials, provisionType
- `server_naming_type` (String) Options: "user" (user configurable), "fixed" (strict pattern)


<a id="nestedatt--config_shutdown"></a>
### Nested Schema for `config_shutdown`

Optional:

- `account_integration_id` (String) ID of your ServiceNow or approval integration
- `flow_id` (String) ID of ServiceNow Flow (set if workflowType is 'flow')
- `shutdown_age` (String) Days instance is allowed to run before shutdown
- `shutdown_allow_extend` (Boolean)
- `shutdown_auto_renew` (Boolean)
- `shutdown_extensions_before_approval` (String) Number of extensions before requiring approval
- `shutdown_hide_fixed` (Boolean) Hide fixed shutdown from users
- `shutdown_message` (String) Notification message
- `shutdown_notify` (String) Days before shutdown to notify via email
- `shutdown_renewal` (String) If the instance is renewed, this is the number of day increments the shutdown date is increased by
- `shutdown_type` (String) Options: "user" (user configurable), "fixed" (strict shutdown)
- `shutdown_workflow_id` (String) ID of legacy ServiceNow workflow (set if workflowType is 'workflow')
- `workflow_type` (String) Options: "workflow" (legacy workflow), "flow" (ServiceNow Flow)


<a id="nestedatt--config_storage_server_quota"></a>
### Nested Schema for `config_storage_server_quota`

Optional:

- `max_storage` (String) Max Storage (GB)
- `storage_server_id` (String) ID of the storage server


<a id="nestedatt--config_tags"></a>
### Nested Schema for `config_tags`

Optional:

- `key` (String) Tag key to enforce
- `strict` (Boolean) Strict enforcement
- `value` (String) Tag value (optional, can be left empty for any value)
- `value_list_id` (String) ID of value from value list (optional)


<a id="nestedatt--config_workflow"></a>
### Nested Schema for `config_workflow`

Optional:

- `workflow_id` (String) ID of the workflow to execute


<a id="nestedatt--cloud"></a>
### Nested Schema for `cloud`

Read-Only:

- `id` (Number)
- `name` (String)


<a id="nestedatt--group"></a>
### Nested Schema for `group`

Read-Only:

- `id` (Number)
- `name` (String)


<a id="nestedatt--owner"></a>
### Nested Schema for `owner`

Read-Only:

- `id` (Number)
- `name` (String)


<a id="nestedatt--role"></a>
### Nested Schema for `role`

Read-Only:

- `authority` (String)
- `id` (Number)


<a id="nestedatt--user"></a>
### Nested Schema for `user`

Read-Only:

- `id` (Number)
- `username` (String)

## Import

Policies can be imported using the `id`, e.g.

```bash
terraform import hpe_morpheus_policy.example 123
```

## Notes
### Associated Resource Types

Policies can be associated with different resource types:

- `Global` - Applies to all resources in the account
- `Group` - Applies to a specific group
- `Cloud` - Applies to a specific cloud
- `User` - Applies to a specific user
- `Role` - Applies to a specific role
- `Network` - Applies to a specific network
- `Plan` - Applies to a specific service plan

### Resource Dependencies

Before creating a policy resource, ensure that:

1. The specified resource `associated_resource_id` refers to exists if `associated_resource_type` is not `Global`
2. Any referenced integrations (e.g., ServiceNow for approval policies) are properly configured
3. Any referenced resources (e.g., workflows, power schedules, networks) exist in Morpheus
