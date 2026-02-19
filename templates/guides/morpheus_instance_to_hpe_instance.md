---
page_title: "Morpheus Instance to HPE Instance"
subcategory: "Migration"
---

# Migrate Morpheus Provider Instance to HPE Provider Instance

This guide provides step-by-step instructions on how to migrate a Morpheus provider instance to a HPE provider instance using the `hpe_morpheus_instance` resource in Terraform.
The guide covers the following types of Morpheus instance:
- VMware Virtual Machine
- HVM Virtual Machine

We strongly recommend that users familiarise themselves with the [importing existing resources](https://developer.hashicorp.com/terraform/language/import)
documentation provided by HashiCorp before proceeding with the migration process.

The following instructions describe importing a single Morpheus instance.  We recommend doing a trial run of the migration
in a separate directory with a single instance to ensure that you are comfortable with the process before performing the
actual migration in your working directory, which may contain many instances and other resources.

The process involves the following steps:
1. Get the Morpheus ID of the instance to be migrated
2. In the test directory create a Terraform configuration file with an `import` block for the instance, the block will look like this:
    ```HCL
    import {
      id = 125623
      to = hpe_morpheus_instance.example
    }
    ```
3. Create a Terraform configuration file with a `provider` and `terraform` block, for example:
    ```HCL
    terraform {
      required_providers {
        hpe = {
          source  = "HPE/hpe"
          version = ">= 1.1.0"
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
4. With these configuration files in place run `terraform plan` to generate a file with draft HCL to be used for the actual import,
   the variable file `local_test.tfvars` contains the actual values for the variables defined in the configuration file.
    ```bash
    $ terraform plan -var-file local_test.tfvars -generate-config-out=generated_instance_example.tf
    ```
    This is the content of `locat_test.tfvars`:
    ```HCL
    morpheus_url     = < url >
    password         = < password >
    username         = < username >
    ```
5. Review the generated file `generated_instance_example.tf`, see the following sections for example generated files for both HVM and VMware instances.
6. Once satisfied with the generated file, proceed to the actual migration by running `terraform apply` in the test directory.
    ```bash
    $ terraform apply -var-file local_test.tfvars
    ```
7. Terraform *may* report errors related to the `network_interface` block, which can be ignored.  These errors will look similar to
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
    If there aren't any errors related to the network interfaces, then it may be necessary to execute a [--refesh-only](https://developer.hashicorp.com/terraform/tutorials/state/refresh)
    apply to update the state with the correct values from the provider before proceeding to the next step.  First execute the plan with the `--refresh-only` flag to confirm that there are no changes, then execute the apply with the same flag:
    ```bash
    $ terraform plan -var-file local_test.tfvars --refresh-only
    hpe_morpheus_instance.example_vmware: Refreshing state... [name=EOTVMWareInstance-2]
    
    Note: Objects have changed outside of Terraform
    
    Terraform detected the following changes made outside of Terraform since the last "terraform apply" which may have affected this plan:
    # hpe_morpheus_instance.example_vmware has changed
    ~ resource "hpe_morpheus_instance" "example_vmware" {
          id                 = 132546
          name               = "EOTVMWareInstance-2"
        ~ network_interfaces = [
            ~ {
                + ip_address        = "10.32.150.39"
                + name              = "eth0"
                + primary_interface = true
                  # (3 unchanged attributes hidden)
              },
          ]
          tags               = [
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
          # (10 unchanged attributes hidden)
      }
 
 
 
 
     This is a refresh-only plan, so Terraform will not take any actions to undo these. If you were expecting these changes then you can apply this plan to record the updated values in the Terraform state without
     changing any remote objects.
     
     ──────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────
     
     Note: You didn't use the -out option to save this plan, so Terraform can't guarantee to take exactly these actions if you run "terraform apply" now.
     ```
     If the plan looks safe - like above, where only the expected attributes related to the network interfaces are changing - then execute the apply with the `--refresh-only` flag:
     ```bash
     $ terraform apply -var-file local_test.tfvars --refresh-only
     ```
    
8. Remove the relevant `import` block from the configuration file
9. Check that a subsequent `terraform plan` will complete successfully and show no changes:
     ```bash
     $ terraform plan -var-file local_test.tfvars
     
     hpe_morpheus_instance.example: Refreshing state... [name=EOTVMWareInstance-2]
     
     No changes. Your infrastructure matches the configuration.
     
     Terraform has compared your real infrastructure against your configuration and found no differences, so no changes are needed.
     ```
10. Once the test migration is successful, you can use the generated configuration file along with the `import` block in
    your actual working directory to perform the migration, repeating the `terraform apply` step as needed.  Note that you
    will see the same expected errors related to the `network_interface` block during the apply step.

## HVM Example Generated HCL

After step 5 above, the generated file `generated_instance_example.tf` will look similar to the following:
```HCL
# __generated__ by Terraform from "131782"
resource "hpe_morpheus_instance" "example_hvm" {
   cloud_id = 7714
   config_hvm = {
      create_user           = false
      kvm_host_id           = null
      nested_virtualization = "off"
      no_agent              = true
      resource_pool_id      = "pool-62299"
   }
   config_vmware    = null
   group_id         = 1
   instance_context = "dev"
   instance_type_id = 9
   layout_id        = 5385
   layout_size      = 1
   name             = "EOTTestInstance-4"
   network_interfaces = [
      {
         child_virtual_networks = [
            {
               ip_pool    = 3251
               network_id = 103481
            },
         ]
         ip_pool    = 3251
         network_id = 103481
      },
      {
         child_virtual_networks = null
         network_id             = 103482
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
         controller_mount_point   = "1183619:0:9:0"
         datastore_auto_selection = null
         datastore_id             = 38658
         name                     = "root"
         root_volume              = true
         size                     = 10
         size_id                  = null
         storage_type_id          = 1
      },
      {
         controller_mount_point   = "1183619:0:9:1"
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

## VMware Example Generated HCL
After step 5 above, the generated file `generated_instance_example.tf` will look similar to the following:
```HCL
# __generated__ by Terraform from "131783"
resource "hpe_morpheus_instance" "example_vmware" {
   cloud_id   = 2
   config_hvm = null
   config_vmware = {
      create_user           = false
      nested_virtualization = "off"
      no_agent              = true
      resource_pool_id      = "pool-1"
      vmware_folder_id      = "group-v79"
   }
   group_id         = 28
   instance_context = "dev"
   instance_type_id = 9
   layout_id        = 1108
   layout_size      = 1
   name             = "EOTVMWareInstance-2"
   network_interfaces = [
      {
         child_virtual_networks = null
         ip_pool                = 3251
         network_id             = 86657
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
         controller_mount_point   = "1183621:0:4:0"
         datastore_auto_selection = null
         datastore_id             = 1
         name                     = "root"
         root_volume              = true
         size                     = 10
         size_id                  = null
         storage_type_id          = 1
      },
      {
         controller_mount_point   = "1183621:0:4:1"
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