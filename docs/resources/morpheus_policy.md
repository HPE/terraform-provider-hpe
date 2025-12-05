---
page_title: "hpe_morpheus_policy Resource - terraform-provider-hpe"
subcategory: "morpheus"
description: |-
  
---
# hpe_morpheus_policy (Resource)



Policies are different rules that can be applied to various Morpheus resources.

-> An `associated_resource_type` must be set (can be one of 'Global', 'Group', 'Cloud', 'User', 'Role', 'Network', 'Plan' or 'Label'). The choice of `policy_type` impacts what `associated_resource_type` values are valid and if you can set `tenants`.

-> If the `associated_resource_type` is not 'Global', then an `associated_resource_id` must also be set.

-> A `policy_type.code` must be set corresponding to the respective policy type.

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
    # Required
    accountIntegrationId = "1" # ID of your ServiceNow or approval integration

    # Optional
    # workflowType = "workflow" # Options: "workflow" (legacy workflow), "flow" (ServiceNow Flow)
    # workflowId = "123"      # ID of legacy ServiceNow workflow (set if workflowType is 'workflow')
    # flowId = "456"          # ID of ServiceNow Flow (set if workflowType is 'flow')
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
    # Required
    accountIntegrationId = "1" # ID of your ServiceNow or approval integration

    # Optional
    workflowType = "workflow" # Options: "workflow" (legacy workflow), "flow" (ServiceNow Flow)
    # workflowId = "123"      # ID of legacy ServiceNow workflow (set if workflowType is 'workflow')
    # flowId = "456"          # ID of ServiceNow Flow (set if workflowType is 'flow')
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
    # Required
    accountIntegrationId = "1" # ID of your ServiceNow or approval integration

    # Optional
    workflowType = "workflow" # Options: "workflow" (legacy workflow), "flow" (ServiceNow Flow)
    # workflowId = "123"      # ID of legacy ServiceNow workflow (set if workflowType is 'workflow')
    # flowId = "456"          # ID of ServiceNow Flow (set if workflowType is 'flow')
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
    # Required
    accountIntegrationId = "1" # ID of your ServiceNow or approval integration

    # Optional
    # workflowType = "workflow" # Options: "workflow" (legacy workflow), "flow" (ServiceNow Flow)
    # workflowId = "123"      # ID of legacy ServiceNow workflow (set if workflowType is 'workflow')
    # flowId = "456"          # ID of ServiceNow Flow (set if workflowType is 'flow')
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
    # Required
    createBackupType = "user" # Options: "user" (user configurable), "fixed" (strict pattern)

    # Optional
    createBackup = true # Enforce backup creation
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
    # Required
    maxPrice = "1000" # Maximum price limit

    # Optional
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
    # Required
    serverNamingType = "user" # Options: "user" (user configurable), "fixed" (strict pattern)

    # Optional
    serverNamingPattern  = "cluster-$${groupCode}-$${type}-$${sequence}" # Name pattern uses ${variable} string interpolation. Available variables: groupName, groupCode, cloudName, cloudCode, type, accountId, account, accountType, platform, username, userId, userInitials, provisionType
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
    # Required
    keyPattern = "secret/*" # Pattern to match Cypher keys (e.g., "secret/*", "password/*")

    # Optional
    read   = true  # Allow read access
    write  = true  # Allow write access
    update = true  # Allow update access
    delete = false # Deny delete access
    list   = true  # Allow list access
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
    # Required
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
    # Required
    lifecycleType = "user" # Options: "user" (user configurable), "fixed" (fixed expiration)

    # Optional
    lifecycleAge                      = "30"                        # Days until expiration
    lifecycleRenewal                  = "7"                         # Days for renewal window
    lifecycleNotify                   = "1"                         # Days before expiration to notify
    lifecycleMessage                  = "Instance will expire soon" # Notification message
    lifecycleAutoRenew                = "on"                        # Options: "on", "off" - auto renewal lifecycle
    lifecycleAllowExtend              = "off"                       # Options: "on", "off" - allow users to extend
    lifecycleExtensionsBeforeApproval = "0"                         # Number of extensions before requiring approval
    lifecycleHideFixed                = false                       # Hide fixed expiration date from users
    # accountIntegrationId = "1"                                    # ID of your ServiceNow or approval integration
    # workflowType = "workflow"                                     # Options: "workflow" (legacy workflow), "flow" (ServiceNow Flow)
    # lifecycleWorkflowId = "123"                                   # ID of legacy ServiceNow workflow (set if workflowType is 'workflow')
    # flowId = "456"                                                # ID of ServiceNow Flow (set if workflowType is 'flow')
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
    # Required
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
    # Required
    hostNamingType = "user" # Options: "user" (user configurable), "fixed" (strict pattern)

    # Optional
    hostNamingPattern = "host-$${groupCode}-$${type}-$${sequence}" # Name pattern uses ${variable} string interpolation. Available variables: groupName, groupCode, cloudName, cloudCode, type, accountId, account, accountType, platform, username, userId, userInitials, provisionType
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
    # Required
    namingType = "user" # Options: "user" (user configurable), "fixed" (strict pattern)

    # Optional
    namingPattern  = "vm-$${groupCode}-$${type}-$${sequence}" # Name pattern uses ${variable} string interpolation. Available variables: groupName, groupCode, cloudName, cloudCode, type, accountId, account, accountType, platform, username, userId, userInitials, provisionType
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
    # Required
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
    # Required
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
    # Required
    maxCores = "32" # Maximum number of CPU cores

    # Optional
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
    # Required
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
    # Required
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
    # Required
    maxMemory = "8" # Maximum memory in GB

    # Optional
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
    # Required
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
    # Required
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
    # Required
    maxStorage = "1000" # Maximum storage in GB

    # Optional
    excludeContainers = "off" # Options: "on", "off" - exclude containers from count
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
    # Required
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
    # Required
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
    # Required
    "motd.title"   = "Welcome"                          # Message title
    "motd.message" = "Welcome to the Morpheus platform" # Message content
    "motd.type"    = "info"                             # Options: "info", "warning", "critical"

    # Optional
    "motd.fullPage" = "off" # Options: "on", "off" - display full page
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
    # Required
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
    # Required
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
    # Required
    powerScheduleType = "user" # Options: "user" (user configurable), "fixed" (strict schedule)

    # Optional
    powerSchedule          = "1"   # ID of the power schedule
    powerScheduleHideFixed = false # Hide fixed schedule from users
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
    # Required
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
    # Required
    shutdownType = "user" # Options: "user" (user configurable), "fixed" (strict shutdown)

    # Optional
    shutdownAge                      = "30"                          # Days instance is allowed to run before shutdown
    shutdownRenewal                  = "7"                           # If the instance is renewed, this is the number of day increments the shutdown date is increased by
    shutdownNotify                   = "1"                           # Days before shutdown to notify via email
    shutdownMessage                  = "Instance will shutdown soon" # Notification message
    shutdownAutoRenew                = "on"                          # Options: "on", "off"
    shutdownAllowExtend              = "off"                         # Options: "on", "off" - allow users to extend
    shutdownExtensionsBeforeApproval = "0"                           # Number of extensions before requiring approval
    shutdownHideFixed                = false                         # Hide fixed shutdown from users
    # accountIntegrationId = "1"                                     # ID of your ServiceNow or approval integration
    # workflowType = "workflow"                                      # Options: "workflow" (legacy workflow), "flow" (ServiceNow Flow)
    # shutdownWorkflowId = "123"                                     # ID of legacy ServiceNow workflow (set if workflowType is 'workflow')
    # flowId = "456"                                                 # ID of ServiceNow Flow (set if workflowType is 'flow')
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
    # Required
    storageServerId = "1" # ID of the storage server

    # Optional
    maxStorage = "1000" # Maximum storage in GB
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
    # Required
    strict = true # Strict enforcement

    # Optional
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
    # Required
    createUserType = "user" # Options: "user" (user configurable), "fixed"

    # Optional
    createUser = true # Enforce user creation
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
    # Required
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
    # Required
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
