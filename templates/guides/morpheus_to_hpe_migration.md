---
page_title: "Morpheus to HPE"
subcategory: "Migration"
---

# Migrate Morpheus Provider Resources to HPE Provider Resources

To migrate resources from the Morpheus Terraform provider to the HPE Terraform provider the basic principle is to
import the existing Morpheus resources created with the Morpheus provider into the HPE provider. HashiCorp documentation
on import can be found [here](https://developer.hashicorp.com/terraform/language/import).

-> [tfmigrator](./tfmigrator_migration.md) has been developed in order to assist with migration. By following the documented process, this tooling should assist in porting `morpheus` resources to their corresponding `hpe` resources - including handling for modules and variables.
Manual guides will continue to be available on our blog, as well as other materials to demonstrate the migration
process using real-world use cases.<br><br>
The section on [Bulk Import](https://developer.hashicorp.com/terraform/language/v1.14.x/import/bulk?page=import&page=bulk) requires provider support for the `List` RPC.  The HPE provider does not currently
support `List`, this is something that we are considering for a future release.<br><br>
Note that some resources in the Morpheus provider have been converted to a more generalised resource in the HPE provider.
See below for details.

## Import Blocks
To import resources we recommend the use of [import blocks](https://developer.hashicorp.com/terraform/language/block/import)
Simple import blocks that specify just an ID and a resource address can be used, and multiple import blocks can be specified,
for example:
```HCL
import {
  id = 125623
  to = hpe_morpheus_instance.example_vmware
}

import {
  id = 125624
  to = hpe_morpheus_instance.example_hvm
}
```

Once resource definitions are supplied a `terraform apply` will import the resources into the state file.  After the
import ensure that a `terraform plan` shows no changes.  If there are changes then the resource definitions need to be
adjusted to match the existing resources.

## Resource Definitions
In addition to the import blocks resource definitions must be supplied.  In general there are two categories of resource
definitions to consider:

### Identical Schema
The following table lists the resources that have identical schema between the Morpheus and HPE providers.

| Morpheus Provider Resource Name           | HPE Provider Resource Name                    |
|-------------------------------------------|-----------------------------------------------|
| morpheus_active_directory_identity_source | hpe_morpheus_identity_source_active_directory |
| morpheus_ansible_integration              | hpe_morpheus_integration_ansible              |
| morpheus_ansible_playbook_task            | hpe_morpheus_task_ansible_playbook            |
| morpheus_ansible_tower_integration        | hpe_morpheus_integration_ansible_tower        |
| morpheus_ansible_tower_task               | hpe_morpheus_task_ansible_tower               |
| morpheus_api_option_list                  | hpe_morpheus_option_list_api                  |
| morpheus_app_blueprint_catalog_item       | hpe_morpheus_catalog_item_app_blueprint       |
| morpheus_appliance_setting                | hpe_morpheus_setting_appliance                |
| morpheus_arm_app_blueprint                | hpe_morpheus_app_blueprint_arm                |
| morpheus_arm_spec_template                | hpe_morpheus_spec_template_arm                |
| morpheus_backup_setting                   | hpe_morpheus_setting_backup                   |
| morpheus_boot_script                      | hpe_morpheus_boot_script                      |
| morpheus_checkbox_option_type             | hpe_morpheus_option_type_checkbox             |
| morpheus_chef_bootstrap_task              | hpe_morpheus_task_chef_bootstrap              |
| morpheus_chef_integration                 | hpe_morpheus_integration_chef                 |
| morpheus_cloud_formation_app_blueprint    | hpe_morpheus_app_blueprint_cloud_formation    |
| morpheus_cloud_formation_spec_template    | hpe_morpheus_spec_template_cloud_formation    |
| morpheus_cluster_layout                   | hpe_morpheus_cluster_layout                   |
| morpheus_cluster_package                  | hpe_morpheus_cluster_package                  |
| morpheus_contact                          | hpe_morpheus_contact                          |
| morpheus_credential                       | hpe_morpheus_credential                       |
| morpheus_cypher_secret                    | hpe_morpheus_cypher_secret                    |
| morpheus_cypher_tfvars                    | hpe_morpheus_cypher_tfvars                    |
| morpheus_docker_registry_integration      | hpe_morpheus_integration_docker_registry      |
| morpheus_email_task                       | hpe_morpheus_task_email                       |
| morpheus_environment                      | hpe_morpheus_environment                      |
| morpheus_execute_schedule                 | hpe_morpheus_execute_schedule                 |
| morpheus_file_template                    | hpe_morpheus_file_template                    |
| morpheus_form                             | hpe_morpheus_form                             |
| morpheus_git_integration                  | hpe_morpheus_integration_git                  |
| morpheus_groovy_script_task               | hpe_morpheus_task_groovy_script               |
| morpheus_guidance_setting                 | hpe_morpheus_setting_guidance                 |
| morpheus_helm_app_blueprint               | hpe_morpheus_app_blueprint_helm               |
| morpheus_helm_spec_template               | hpe_morpheus_spec_template_helm               |
| morpheus_hidden_option_type               | hpe_morpheus_option_type_hidden               |
| morpheus_instance_catalog_item            | hpe_morpheus_catalog_item_instance            |
| morpheus_instance_layout                  | hpe_morpheus_instance_type_layout             |
| morpheus_instance_type                    | hpe_morpheus_instance_type                    |
| morpheus_ipv4_ip_pool                     | hpe_morpheus_ip_pool_ipv4                     |
| morpheus_javascript_task                  | hpe_morpheus_task_javascript                  |
| morpheus_key_pair                         | hpe_morpheus_key_pair                         |
| morpheus_kubernetes_app_blueprint         | hpe_morpheus_app_blueprint_kubernetes         |
| morpheus_kubernetes_spec_template         | hpe_morpheus_spec_template_kubernetes         |
| morpheus_library_script_task              | hpe_morpheus_task_library_script              |
| morpheus_library_template_task            | hpe_morpheus_task_library_template            |
| morpheus_license                          | hpe_morpheus_license                          |
| morpheus_manual_option_list               | hpe_morpheus_option_list_manual               |
| morpheus_monitoring_setting               | hpe_morpheus_setting_monitoring               |
| morpheus_nested_workflow_task             | hpe_morpheus_task_nested_workflow             |
| morpheus_network_domain                   | hpe_morpheus_network_domain                   |
| morpheus_node_type                        | hpe_morpheus_node_type                        |
| morpheus_number_option_type               | hpe_morpheus_option_type_number               |
| morpheus_operational_workflow             | hpe_morpheus_workflow_operational             |
| morpheus_password_option_type             | hpe_morpheus_option_type_password             |
| morpheus_powershell_script_task           | hpe_morpheus_task_powershell_script           |
| morpheus_preseed_script                   | hpe_morpheus_preseed_script                   |
| morpheus_price                            | hpe_morpheus_price                            |
| morpheus_price_set                        | hpe_morpheus_price_set                        |
| morpheus_provisioning_setting             | hpe_morpheus_setting_provisioning             |
| morpheus_provisioning_workflow            | hpe_morpheus_workflow_provisioning            |
| morpheus_puppet_integration               | hpe_morpheus_integration_puppet               |
| morpheus_python_script_task               | hpe_morpheus_task_python_script               |
| morpheus_radio_list_option_type           | hpe_morpheus_option_type_radio_list           |
| morpheus_resource_pool_group              | hpe_morpheus_resource_pool_group              |
| morpheus_rest_option_list                 | hpe_morpheus_option_list_rest                 |
| morpheus_restart_task                     | hpe_morpheus_task_restart                     |
| morpheus_ruby_script_task                 | hpe_morpheus_task_ruby_script                 |
| morpheus_saml_identity_source             | hpe_morpheus_identity_source_saml             |
| morpheus_scale_threshold                  | hpe_morpheus_scale_threshold                  |
| morpheus_script_template                  | hpe_morpheus_script_template                  |
| morpheus_security_package                 | hpe_morpheus_security_package                 |
| morpheus_select_list_option_type          | hpe_morpheus_option_type_select_list          |
| morpheus_servicenow_integration           | hpe_morpheus_integration_servicenow           |
| morpheus_shell_script_task                | hpe_morpheus_task_shell_script                |
| morpheus_task_job                         | hpe_morpheus_job_task                         |
| morpheus_tenant                           | hpe_morpheus_tenant                           |
| morpheus_terraform_app_blueprint          | hpe_morpheus_app_blueprint_terraform          |
| morpheus_terraform_spec_template          | hpe_morpheus_spec_template_terraform          |
| morpheus_text_option_type                 | hpe_morpheus_option_type_text                 |
| morpheus_textarea_option_type             | hpe_morpheus_option_type_textarea             |
| morpheus_typeahead_option_type            | hpe_morpheus_option_type_typeahead            |
| morpheus_user_group                       | hpe_morpheus_user_group                       |
| morpheus_vro_integration                  | hpe_morpheus_integration_vro                  |
| morpheus_vro_task                         | hpe_morpheus_task_vro                         |
| morpheus_vsphere_mks_cluster              | hpe_morpheus_cluster_hks_vsphere              |
| morpheus_wiki_page                        | hpe_morpheus_wiki_page                        |
| morpheus_workflow_catalog_item            | hpe_morpheus_catalog_item_workflow            |
| morpheus_workflow_job                     | hpe_morpheus_job_workflow                     |
| morpheus_write_attributes_task            | hpe_morpheus_task_write_attributes            |

For these resources we can simply copy the resource definition from the Morpheus provider config and change the resource
type to the corresponding HPE provider resource type.  Then an import block must be added for each HPE resource to
import the existing Morpheus resource.  [tfmigrator](./tfmigrator_migration.md) handles this.

### Different Schema
Some resources have differences in schema between the Morpheus and HPE providers.  These resources are generalised
in the HPE provider and several Morpheus provider resources map to a single HPE provider resource:

| Morpheus Provider Resource Name | HPE Provider Resource Name |
|---------------------------------|----------------------------|
| morpheus_aws_cloud | hpe_morpheus_cloud |
| morpheus_azure_cloud | hpe_morpheus_cloud |
| morpheus_standard_cloud | hpe_morpheus_cloud |
| morpheus_vsphere_cloud | hpe_morpheus_cloud |
| morpheus_vsphere_cloud_datastore_configuration | hpe_morpheus_datastore |
| morpheus_aws_instance | hpe_morpheus_instance |
| morpheus_mvm_instance | hpe_morpheus_instance |
| morpheus_vsphere_instance | hpe_morpheus_instance |
| morpheus_backup_creation_policy | hpe_morpheus_policy |
| morpheus_budget_policy | hpe_morpheus_policy |
| morpheus_cluster_resource_name_policy | hpe_morpheus_policy |
| morpheus_cypher_access_policy | hpe_morpheus_policy |
| morpheus_delayed_delete_policy | hpe_morpheus_policy |
| morpheus_delete_approval_policy | hpe_morpheus_policy |
| morpheus_hostname_policy | hpe_morpheus_policy |
| morpheus_instance_name_policy | hpe_morpheus_policy |
| morpheus_max_containers_policy | hpe_morpheus_policy |
| morpheus_max_cores_policy | hpe_morpheus_policy |
| morpheus_max_hosts_policy | hpe_morpheus_policy |
| morpheus_max_memory_policy | hpe_morpheus_policy |
| morpheus_max_storage_policy | hpe_morpheus_policy |
| morpheus_max_vms_policy | hpe_morpheus_policy |
| morpheus_motd_policy | hpe_morpheus_policy |
| morpheus_network_quota_policy | hpe_morpheus_policy |
| morpheus_power_schedule_policy | hpe_morpheus_policy |
| morpheus_provision_approval_policy | hpe_morpheus_policy |
| morpheus_router_quota_policy | hpe_morpheus_policy |
| morpheus_tag_policy | hpe_morpheus_policy |
| morpheus_user_creation_policy | hpe_morpheus_policy |
| morpheus_user_group_creation_policy | hpe_morpheus_policy |
| morpheus_workflow_policy | hpe_morpheus_policy |
| morpheus_tenant_role | hpe_morpheus_role |
| morpheus_user_role | hpe_morpheus_role |
| morpheus_group | hpe_morpheus_group |
| morpheus_service_plan | hpe_morpheus_service_plan |
| morpheus_user | hpe_morpheus_user |

For these resources terraform can be used to generate the HPE resource definition.  This is how the process documented
for [tfmigrator](./tfmigrator_migration.md) handles these resources.