---
page_title: "Migrate Morpheus Provider Instance to HPE Provider Instance"
subcategory: "Migration"
---

# Migrate Morpheus Instance to HPE Instance

This guide provides step-by-step instructions on how to migrate a Morpheus instance to an HPE instance using the `hpe_morpheus_instance` resource in Terraform.
The guide covers the following types of Morpheus instance:
- VMware Virtual Machine
- HVM Virtual Machine

We strongly recommend that users familiarize themselves with the [importing existing resources](https://developer.hashicorp.com/terraform/language/import)
documentation provided by HashiCorp before proceeding with the migration process.

The following instructions describe importing a single Morpheus instance.

The process involves the following steps:
1. Do a trial run of the migration in a separate directory
1. Get the Morpheus ID of the instance to be migrated
1. In the test directory create a Terraform configuration file with an `import` block for the instance, the block will look like this:
   ```HCL
   import {
     id = 125623
     to = hpe_morpheus_instance.example
   }
   ```
1. Create a Terraform configuration file with a `provider` and `terraform` block, for example:
   ```HCL
   terraform {
     required_providers {
       hpe = {
         source  = "HPE/hpe"
         version = "= 1.0.0"
       }
     }
   }
   
   provider "hpe" {
     morpheus {
       url      = var.morpheus_url
       username = var.username
       password = var.password
     }
   }
   ```
   Note that the above example includes some variable definitions, which are in another configuration file:
   ```HCL
    variable "morpheus_url" {
      type = string
    }
    
    variable "username" {
      type = string
    }
    
    variable "password" {
      type      = string
      sensitive = true
    }
   ```
1. With these configuration files in place run `terraform plan` to generate a file with draft HCL to be used for the actual import,
   the variable file `local_test.tfvars` contains the actual values for the variables defined in the configuration file.
   ```bash
   $ terraform plan -var-file local_test.tfvars -generate-config-out=generated_instance_example.tf
   ```
1. Review the generated file `generated_instance_example.tf` and make any necessary adjustments as explained in the following sections.
1. Once satisfied with the generated file, proceed to the actual migration by running `terraform apply` in the test directory.
   ```bash
   $ terraform apply -var-file local_test.tfvars
   ```
1. Terraform will report errors related to the `network_interface` block, which is expected.  These errors will look similar to
   the following:
   ```bash
   Plan: 1 to import, 0 to add, 1 to change, 0 to destroy.
   hpe_morpheus_instance.example: Importing... [id=125623]
   hpe_morpheus_instance.example: Import complete [id=125623]
   hpe_morpheus_instance.example: Modifying... [name=EOTVMWareInstance-2]
   ╷
   │ Error: Provider produced inconsistent result after apply
   │ 
   │ When applying changes to hpe_morpheus_instance.example_vmware, provider "provider[\"registry.terraform.io/hpe/hpe\"]" produced an unexpected new value: .network_interfaces[0].name: was null, but now
   │ cty.StringVal("eth0").
   │ 
   │ This is a bug in the provider, which should be reported in the provider's own issue tracker.
   ╵
   ╷
   │ Error: Provider produced inconsistent result after apply
   │ 
   │ When applying changes to hpe_morpheus_instance.example_vmware, provider "provider[\"registry.terraform.io/hpe/hpe\"]" produced an unexpected new value: .network_interfaces[0].primary_interface: was null, but
   │ now cty.True.
   │ 
   │ This is a bug in the provider, which should be reported in the provider's own issue tracker.
   ╵
   ╷
   │ Error: Provider produced inconsistent result after apply
   │ 
   │ When applying changes to hpe_morpheus_instance.example_vmware, provider "provider[\"registry.terraform.io/hpe/hpe\"]" produced an unexpected new value: .network_interfaces[0].ip_address: was
   │ cty.StringVal(""), but now cty.StringVal("10.32.151.143").
   │ 
   │ This is a bug in the provider, which should be reported in the provider's own issue tracker.
   ```
1. Remove the relevant `import` block from the configuration file
1. Check that a subsequent `terraform plan` will complete successfully and show no changes - ignore the warning about provider development overrides,
   that is because in my testing I am using a locally built version of the provider:
   ```bash
   $ terraform plan -var-file local_test.tfvars
   ╷
   │ Warning: Provider development overrides are in effect
   │ 
   │ The following provider development overrides are set in the CLI configuration:
   │  - hpe/hpe in /Users/eamonn/go/bin
   │ 
   │ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible with published releases.
   ╵
   hpe_morpheus_instance.example: Refreshing state... [name=EOTVMWareInstance-2]
   
   No changes. Your infrastructure matches the configuration.
   
   Terraform has compared your real infrastructure against your configuration and found no differences, so no changes are needed.
   ```
1. Once the test migration is successful, you can use the generated configuration file along with the `import` block in
    your actual working directory to perform the migration, repeating the `terraform apply` step as needed.  Note that you
    will see the same expected errors related to the `network_interface` block during the apply step.

## HVM Specific Adjustments

After step 5 above, the generated file `generated_instance_example.tf` will need some adjustments specific to HVM instances.
This is an example generated file for an HVM instance:
```HCL
# __generated__ by Terraform
# Please review these resources and move them into your main configuration files.

# __generated__ by Terraform from "125624"
resource "hpe_morpheus_instance" "example_hvm" {
  cloud_id = 7714
  config = {
    backup = {
      AdditionalProperties = {
        jobName     = ""
        jobSchedule = 2
        name        = ""
      }
      CreateBackup       = null
      JobAction          = "new"
      JobRetentionCount  = "3"
      ProviderBackupType = 14
    }
    createBackup         = true
    createUser           = false
    customOptions        = {}
    kvmHostId            = 346676
    layoutSize           = 1
    memoryDisplay        = "MB"
    nestedVirtualization = "off"
    noAgent = {
      Bool   = true
      String = null
    }
    poolProviderType = "mvm"
    resourcePoolId = {
      Int64  = null
      String = "pool-62299"
    }
  }
  evars            = null
  group_id         = 1
  instance_context = "dev"
  instance_type_id = 9
  layout_id        = 5385
  layout_size      = 1
  name             = "EOTTestInstance-2"
  network_interfaces = [
    {
      child_virtual_networks = [
        {
          ip_address       = null
          ip_mode          = null
          ip_pool          = 3251
          network_group_id = null
          network_id       = 103481
          network_type_id  = null
        },
      ]
      ip_address       = null
      ip_mode          = null
      ip_pool          = 3251
      network_group_id = null
      network_id       = 103481
      network_type_id  = null
    },
    {
      child_virtual_networks = null
      ip_address             = null
      ip_mode                = null
      ip_pool                = null
      network_group_id       = null
      network_id             = 103482
      network_type_id        = null
    },
  ]
  plan_id = 173
  ports   = null
  tags = [
    {
      name  = "acctest"
      value = "true"
    },
    {
      name  = "hpe_morpheus_instance"
      value = "true"
    },
    {
      name  = "managed_by"
      value = "terraform"
    },
    {
      name  = "terraform"
      value = "true"
    },
  ]
  task_set_id = null
  timeouts    = null
  volumes = [
    {
      controller_mount_point   = "1133003:0:9:0"
      datastore_auto_selection = null
      datastore_id             = 38658
      name                     = "root"
      root_volume              = true
      size                     = 10
      size_id                  = null
      storage_type_id          = 1
    },
    {
      controller_mount_point   = "1133003:0:9:1"
      datastore_auto_selection = null
      datastore_id             = 38658
      name                     = "data"
      root_volume              = false
      size                     = 10
      size_id                  = null
      storage_type_id          = 1
    },
  ]
}
```

### `config` Block

This block will need to be edited to match the desired configuration for a HVM instance, which looks like this:
```HCL
  config = {
    resourcePoolId       = "pool-62299"
    poolProviderType     = "mvm"
    nestedVirtualization = "off"
    noAgent              = false
    createUser           = true
  }
```

This is the generated `config` block
```HCL
  config = {
    backup = {
      AdditionalProperties = {
        jobName     = ""
        jobSchedule = 2
        name        = ""
      }
      CreateBackup       = null
      JobAction          = "new"
      JobRetentionCount  = "3"
      ProviderBackupType = 14
    }
    createBackup         = true
    createUser           = false
    customOptions        = {}
    kvmHostId            = 346676
    layoutSize           = 1
    memoryDisplay        = "MB"
    nestedVirtualization = "off"
    noAgent = {
      Bool   = true
      String = null
    }
    poolProviderType = "mvm"
    resourcePoolId = {
      Int64  = null
      String = "pool-62299"
    }
  }
```

The following changes are needed:
- Remove the entire `backup` block
- Remove the `createBackup` line
- Remove the `customOptions` line
- Remove the `layoutSize` line
- Remove the `memoryDisplay` line
- Change the `noAgent` block to `noAgent = true` which is the setting in the `Bool` field of the block
- Change the `resourcePoolId` block to `resourcePoolId = "pool-62299"` which is the value in the `String` field of the block

### VMware Specific Adjustments
After step 5 above, the generated file `generated_instance_example.tf` will need some adjustments specific to VMware instances.
This is an example generated file for an VMware instance:
```HCL
# __generated__ by Terraform
# Please review these resources and move them into your main configuration files.

# __generated__ by Terraform from "125623"
resource "hpe_morpheus_instance" "example_vmware" {
   cloud_id = 2
   config = {
      backup = {
         AdditionalProperties = {
            jobName = ""
            name    = ""
         }
         CreateBackup       = null
         JobAction          = "new"
         JobRetentionCount  = null
         ProviderBackupType = 50
      }
      createBackup         = true
      createUser           = false
      customOptions        = {}
      layoutSize           = 1
      memoryDisplay        = "MB"
      nestedVirtualization = "off"
      noAgent = {
         Bool   = true
         String = null
      }
      resourcePoolId = {
         Int64  = null
         String = "pool-1"
      }
      vmwareFolderId = "group-v79"
   }
   evars            = null
   group_id         = 28
   instance_context = "dev"
   instance_type_id = 9
   layout_id        = 1108
   layout_size      = 1
   name             = "EOTVMWareInstance-2"
   network_interfaces = [
      {
         child_virtual_networks = null
         ip_address             = null
         ip_mode                = null
         ip_pool                = 3251
         network_group_id       = null
         network_id             = 86657
         network_type_id        = null
      },
   ]
   plan_id = 241
   ports   = null
   tags = [
      {
         name  = "acctest"
         value = "true"
      },
      {
         name  = "hpe_morpheus_instance"
         value = "true"
      },
      {
         name  = "managed_by"
         value = "terraform"
      },
      {
         name  = "terraform"
         value = "true"
      },
   ]
   task_set_id = null
   timeouts    = null
   volumes = [
      {
         controller_mount_point   = "1133000:0:4:0"
         datastore_auto_selection = null
         datastore_id             = 1
         name                     = "root"
         root_volume              = true
         size                     = 10
         size_id                  = null
         storage_type_id          = 1
      },
      {
         controller_mount_point   = "1133000:0:4:1"
         datastore_auto_selection = null
         datastore_id             = 1
         name                     = "data"
         root_volume              = false
         size                     = 10
         size_id                  = null
         storage_type_id          = 1
      },
   ]
}
```

### `config` Block

This block will need to be edited to match the desired configuration for a VMware instance, which looks like this:
```HCL
  config = {
   createUser           = false
   nestedVirtualization = "off"
   noAgent = true
   resourcePoolId = "pool-1"
   vmwareFolderId = "group-v79"
}
```

This is the generated `config` block
```HCL
  config = {
    backup = {
      AdditionalProperties = {
        jobName = ""
        name    = ""
      }
      CreateBackup       = null
      JobAction          = "new"
      JobRetentionCount  = null
      ProviderBackupType = 50
    }
    createBackup         = true
    createUser           = false
    customOptions        = {}
    layoutSize           = 1
    memoryDisplay        = "MB"
    nestedVirtualization = "off"
    noAgent = {
      Bool   = true
      String = null
    }
    resourcePoolId = {
      Int64  = null
      String = "pool-1"
    }
    vmwareFolderId = "group-v79"
  }
```

The following changes are needed:
- Remove the entire `backup` block
- Remove the `createBackup` line
- Remove the `customOptions` line
- Remove the `layoutSize` line
- Remove the `memoryDisplay` line
- Change the `noAgent` block to `noAgent = true` which is the setting in the `Bool` field of the block
- Change the `resourcePoolId` block to `resourcePoolId = "pool-1"` which is the value in the `String` field of the block