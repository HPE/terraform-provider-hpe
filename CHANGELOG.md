# v1.4.0 Release Notes

In this release (v1.4.0) we have added the following resources:

### Networking

- hpe_morpheus_load_balancer_monitor
- hpe_morpheus_load_balancer_virtual_server
- hpe_morpheus_network_dhcp_server
- hpe_morpheus_network_firewall_rule
- hpe_morpheus_network_firewall_rule_group
- hpe_morpheus_network_group
- hpe_morpheus_network_pool
- hpe_morpheus_network_pool_server
- hpe_morpheus_network_router
- hpe_morpheus_network_router_bgp_neighbor
- hpe_morpheus_network_router_firewall_rule
- hpe_morpheus_network_router_nat
- hpe_morpheus_network_router_route
- hpe_morpheus_security_group
- hpe_morpheus_security_group_rule
- hpe_morpheus_subnet

### Monitoring & Operations

- hpe_morpheus_backup
- hpe_morpheus_backup_job
- hpe_morpheus_budget
- hpe_morpheus_monitoring_alert
- hpe_morpheus_monitoring_check
- hpe_morpheus_monitoring_group

### Storage

- hpe_morpheus_storage_bucket
- hpe_morpheus_storage_server
- hpe_morpheus_storage_volume

### Library & Provisioning

- hpe_morpheus_container_script
- hpe_morpheus_option_list
- hpe_morpheus_option_type
- hpe_morpheus_provisioning_license

### VDI

- hpe_morpheus_vdi_app
- hpe_morpheus_vdi_gateway
- hpe_morpheus_vdi_pool

### Compute & Cluster

- hpe_morpheus_cluster_affinity_group
- hpe_morpheus_cluster_namespace

### Identity & Governance

- hpe_morpheus_whitelabel_settings

### Other

- hpe_morpheus_certificate
- hpe_morpheus_deployment
- hpe_morpheus_power_schedule

In this release (v1.4.0) we have added the following data sources:

- hpe_morpheus_load_balancer_monitor
- hpe_morpheus_load_balancer_virtual_server
- hpe_morpheus_network_dhcp_server
- hpe_morpheus_network_domain (reimplemented)
- hpe_morpheus_network_firewall_rule
- hpe_morpheus_network_firewall_rule_group
- hpe_morpheus_network_router
- hpe_morpheus_network_router_bgp_neighbor
- hpe_morpheus_network_router_route

## Enhancements to existing resources

- hpe_morpheus_cloud — Added `config_azure` static config block; added 6 inventory discovery sync fields (`default_*_sync_active`)
- hpe_morpheus_instance — Added `config_azure` static config block; added `network_domain_id` attribute (change forces recreation)
- hpe_morpheus_load_balancer — Added tainting support; added `config` Dynamic attribute
- hpe_morpheus_task — Allow null `else` in conditional_workflow task
- All password/secret fields across 10 resources now enforce Sensitive + WriteOnly + PlanModifiers

## Resolved issues

- `hpe_morpheus_cloud` import fails with unknown values in static config blocks
- `hpe_morpheus_task` conditional_workflow doesn't support null else clause
- `hpe_morpheus_app_blueprint_kubernetes` YAML config not properly parsed
- `hpe_morpheus_ansible_tower_inventory` data source type assertions fail intermittently
- `hpe_morpheus_cluster` delete polling can report false failures
- Import syntax standardised to dot-notation for multi-ID resources

## Known issues

- `hpe_morpheus_cluster_hks_hvm` Destroy may return an error but the cluster will be deleted successfully, this is being investigated.
- `hpe_morpheus_instance` updates fail when removing optional fields.
  This will be addressed in a future release.
- `hpe_morpheus_instance` updates fail when removing `evars`.
  This will be addressed in a future release.
- Long running operations can fail when using username and password.
- `hpe_morpheus_instance` depending on the layout used may require one or more `volumes` to be specified,
  in these cases not specifying the correct number of `volumes` will cause instance creation to fail.
- There are intermittent issues with the provider failing to authenticate, a 500 error is returned from the Morpheus API.
  If this happens please retry the operation.  This is being investigated.
- `hpe_morpheus_datastore` when creating a datastore of type NFS the creation will silently fail if the NFS server is not reachable or the share is not accessible.
  The datastore will remain in a `provisioning` state indefinitely. Ensure the Morpheus appliance can reach the NFS server
  and that the share is accessible before creating.
- `hpe_morpheus_datastore` delete is not guaranteed to succeed.  Alletra MP HVM and Alletra MP BM datastores will delete but NFS datastores
  may fail to delete.  Always delete VMs and other resources using the datastore before deleting the datastore itself.
- `hpe_morpheus_instance` in Morpheus versions prior to 8.0.11 requires that the `root` volume is the first entry in
  the `volumes` block list

# v1.3.0 Release Notes

In this release (v1.3.0) we have added a notifier which issues a Warning if the provider version is less than the latest version available on the registry.
This can be suppressed by upgrading to the latest version or setting the environment variable `HPE_IGNORE_VERSION_CHECK`.

In this release (v1.3.0) we have added the following resource functionality:

- hpe_morpheus_cluster a generalised cluster resource has been added with support for HVM clusters, and limited update functionality
- hpe_morpheus_forms has comprehensive support for all option types
- hpe_morpheus_instance has a static `config_aws` block for AWS instances
- hpe_morpheus_instance supports Update of `network_interfaces` and `service_plan_options` for service plans that support the setting of
  options for Morpheus versions >= 8.1.2, for earlier versions changes will force a new instance to be created
- hpe_morpheus_load_balancer resource has been added
- hpe_morpheus_os_type resource has been added
- hpe_morpheus_os_type_image resource has been added

In this release (v1.3.0) we have added the following data source functionality:

- hpe_morpheus_cluster data source has been added
- hpe_morpheus_load_balancer data source has been added
- hpe_morpheus_os_type data source has been added
- hpe_morpheus_os_type_image data source has been added

## New known issues

- hpe_morpheus_cluster_hks_hvm Destroy may return an error but the cluster will be deleted successfully, this is being investigated.

## Resolved issues

- `hpe_morpheus_instance` Update of `service_plan_options` fails silently

## Known issues from previous releases

- `hpe_morpheus_instance` updates fail when removing optional fields.
  This will be addressed in a future release.
- `hpe_morpheus_instance` updates fail when removing `evars`.
  This will be addressed in a future release.
- Long running operations can fail when using username and password.
- `hpe_morpheus_instance` depending on the layout used may require one or more `volumes` to be specified,
  in these cases not specifying the correct number of `volumes` will cause instance creation to fail.
- There are intermittent issues with the provider failing to authenticate, a 500 error is returned from the Morpheus API.
  If this happens please retry the operation.  This is being investigated.
- `hpe_morpheus_datastore` when creating a datastore of type NFS the creation will silently fail if the NFS server is not reachable or the share is not accessible.
  The datastore will remain in a `provisioning` state indefinitely. Ensure the Morpheus appliance can reach the NFS server
  and that the share is accessible before creating.
- `hpe_morpheus_datastore` delete is not guaranteed to succeed.  Alletra MP HVM and Alletra MP BM datastores will delete but NFS datastores
  may fail to delete.  Always delete VMs and other resources using the datastore before deleting the datastore itself.
- `hpe_morpheus_instance` in Morpheus versions prior to 8.0.11 requires that the `root` volume is the first entry in
  the `volumes` block list

# v1.2.0 Release Notes

In this release (v1.2.0) we have added the following resource functionality:
- hpe_morpheus_instance addition and removal of volumes is now supported in Update
- hpe_morpheus_instance supports service_plan_options for use with Service Plans that accept options
- hpe_morpheus_cloud no longer requires group_id to be set
- hpe_morpheus_group supports a list of cloud-ids to associate the group with

## New known issues

- `hpe_morpheus_instance` Update of `service_plan_options` fails silently, this will be fixed in a future release

## Resolved issues

- `hpe_morpheus_datastore` data-source if a datastore with the specified name cannot be found (i.e. the corresponding
  list API request fails), the error message will indicate a 403 (Forbidden) even if the user has permission to list
  datastores.

## Known issues from previous releases

- `hpe_morpheus_instance` updates fail when removing optional fields.
  This will be addressed in a future release.
- `hpe_morpheus_instance` updates fail when removing `evars`.
  This will be addressed in a future release.
- Long running operations can fail when using username and password.
- `hpe_morpheus_instance` depending on the layout used may require one or more `volumes` to be specified,
  in these cases not specifying the correct number of `volumes` will cause instance creation to fail.
- There are intermittent issues with the provider failing to authenticate, a 500 error is returned from the Morpheus API.
  If this happens please retry the operation.  This is being investigated.
- `hpe_morpheus_datastore` when creating a datastore of type NFS the creation will silently fail if the NFS server is not reachable or the share is not accessible.
  The datastore will remain in a `provisioning` state indefinitely. Ensure the Morpheus appliance can reach the NFS server
  and that the share is accessible before creating.
- `hpe_morpheus_datastore` delete is not guaranteed to succeed.  Alletra MP HVM and Alletra MP BM datastores will delete but NFS datastores
  may fail to delete.  Always delete VMs and other resources using the datastore before deleting the datastore itself.
- `hpe_morpheus_instance` in Morpheus versions prior to 8.0.11 requires that the `root` volume is the first entry in
  the `volumes` block list
 
# v1.1.0 Release Notes

In this release (v1.1.0) we have added the following resource functionality:

- hpe_morpheus_cloud has static config schema for VMware and HVM clouds
- hpe_morpheus_instance has static config schema for VMware and HVM instances
- hpe_morpheus_task generalised task resource with static config schema for Conditional Workflow task

## New known issues

N/A

## Resolved issues

- `hpe_morpheus_cluster_hks_vsphere` scale-down and destroy issues are fixed in Morpheus version 8.1 and later

## Known issues from previous releases

- `hpe_morpheus_datastore` data-source if a datastore with the specified name cannot be found (i.e. the corresponding
  list API request fails), the error message will indicate a 403 (Forbidden) even if the user has permission to list
  datastores.  This is an API bug which is being investigated.
- `hpe_morpheus_instance` updates fail when removing optional fields.
  This will be addressed in a future release.
- `hpe_morpheus_instance` updates fail when removing `evars`.
  This will be addressed in a future release.
- Long running operations can fail when using username and password.
- `hpe_morpheus_instance` depending on the layout used may require one or more `volumes` to be specified,
  in these cases not specifying the correct number of `volumes` will cause instance creation to fail.
- There are intermittent issues with the provider failing to authenticate, a 500 error is returned from the Morpheus API.
  If this happens please retry the operation.  This is being investigated.
- `hpe_morpheus_datastore` when creating a datastore of type NFS the creation will silently fail if the NFS server is not reachable or the share is not accessible.
  The datastore will remain in a `provisioning` state indefinitely. Ensure the Morpheus appliance can reach the NFS server
  and that the share is accessible before creating.
- `hpe_morpheus_datastore` delete is not guaranteed to succeed.  Alletra MP HVM and Alletra MP BM datastores will delete but NFS datastores
  may fail to delete.  Always delete VMs and other resources using the datastore before deleting the datastore itself.
- `hpe_morpheus_instance` in Morpheus versions prior to 8.0.11 requires that the `root` volume is the first entry in
  the `volumes` block list

# v1.0.0 Release Notes

In this release (v1.0.0) we have added the following resource functionality:

- hpe_morpheus_app_blueprint_arm
- hpe_morpheus_app_blueprint_cloud_formation
- hpe_morpheus_app_blueprint_helm
- hpe_morpheus_app_blueprint_kubernetes
- hpe_morpheus_app_blueprint_terraform
- hpe_morpheus_catalog_item_app_blueprint
- hpe_morpheus_catalog_item_instance
- hpe_morpheus_catalog_item_workflow
- hpe_morpheus_cluster_layout
- hpe_morpheus_cluster_hks_hvm
- hpe_morpheus_cluster_hks_vsphere
- hpe_morpheus_cluster_package
- hpe_morpheus_contact
- hpe_morpheus_credential
- hpe_morpheus_cypher_secret
- hpe_morpheus_cypher_tfvars
- hpe_morpheus_datastore supports Alletra MP BMaaS datastores
- hpe_morpheus_environment
- hpe_morpheus_execute_schedule
- hpe_morpheus_file_template
- hpe_morpheus_form
- hpe_morpheus_identity_source_active_directory
- hpe_morpheus_identity_source_saml
- hpe_morpheus_integration_ansible
- hpe_morpheus_integration_ansible_tower
- hpe_morpheus_integration_chef
- hpe_morpheus_integration_docker_registry
- hpe_morpheus_integration_git
- hpe_morpheus_integration_puppet
- hpe_morpheus_integration_servicenow
- hpe_morpheus_integration_vro
- hpe_morpheus_instance supports VMware and BMaaS instances
- hpe_morpheus_instance_type
- hpe_morpheus_instance_type_layout
- hpe_morpheus_job_task
- hpe_morpheus_job_workflow
- hpe_morpheus_key_pair
- hpe_morpheus_license
- hpe_morpheus_network_domain
- hpe_morpheus_node_type
- hpe_morpheus_option_list_api
- hpe_morpheus_option_list_manual
- hpe_morpheus_option_list_rest
- hpe_morpheus_option_type_checkbox
- hpe_morpheus_option_type_hidden
- hpe_morpheus_option_type_number
- hpe_morpheus_option_type_password
- hpe_morpheus_option_type_radio_list
- hpe_morpheus_option_type_select_list
- hpe_morpheus_option_type_text
- hpe_morpheus_option_type_textarea
- hpe_morpheus_option_type_typeahead
- hpe_morpheus_policy has a comprehensive collection of static schema for the various supported policies
- hpe_morpheus_preseed_script
- hpe_morpheus_price
- hpe_morpheus_price_set
- hpe_morpheus_resource_pool_group
- hpe_morpheus_scale_threshold
- hpe_morpheus_script_template
- hpe_morpheus_security_package
- hpe_morpheus_setting_appliance
- hpe_morpheus_setting_backup
- hpe_morpheus_setting_guidance
- hpe_morpheus_setting_monitoring
- hpe_morpheus_setting_provisioning
- hpe_morpheus_spec_template_arm
- hpe_morpheus_spec_template_cloud_formation
- hpe_morpheus_spec_template_helm
- hpe_morpheus_spec_template_kubernetes
- hpe_morpheus_spec_template_terraform
- hpe_morpheus_task_ansible_playbook
- hpe_morpheus_task_ansible_tower
- hpe_morpheus_task_chef_bootstrap
- hpe_morpheus_task_email
- hpe_morpheus_task_groovy_script
- hpe_morpheus_task_javascript
- hpe_morpheus_task_library_script
- hpe_morpheus_task_library_template
- hpe_morpheus_task_nested_workflow
- hpe_morpheus_task_powershell_script
- hpe_morpheus_task_python_script
- hpe_morpheus_task_restart
- hpe_morpheus_task_ruby_script
- hpe_morpheus_task_shell_script
- hpe_morpheus_task_write_attributes
- hpe_morpheus_tenant
- hpe_morpheus_user_group
- hpe_morpheus_wiki_page
- hpe_morpheus_workflow_operational
- hpe_morpheus_workflow_provisioning

In this release (v1.0.0) we have added the following data-source functionality:

- hpe_morpheus_ansible_tower_inventory
- hpe_morpheus_ansible_tower_job_template
- hpe_morpheus_blueprint
- hpe_morpheus_budget
- hpe_morpheus_catalog_item_type
- hpe_morpheus_cloud_folder
- hpe_morpheus_cloud_type
- hpe_morpheus_clouds
- hpe_morpheus_cluster_type
- hpe_morpheus_contact
- hpe_morpheus_credential
- hpe_morpheus_cypher_secret
- hpe_morpheus_environments
- hpe_morpheus_execute_schedule
- hpe_morpheus_file_template
- hpe_morpheus_groups
- hpe_morpheus_images
- hpe_morpheus_integration
- hpe_morpheus_integration_git
- hpe_morpheus_instance
- hpe_morpheus_instance_type
- hpe_morpheus_job
- hpe_morpheus_key_pair
- hpe_morpheus_network_domain
- hpe_morpheus_network_group
- hpe_morpheus_network_subnet
- hpe_morpheus_networks
- hpe_morpheus_node_type
- hpe_morpheus_option_list
- hpe_morpheus_option_type
- hpe_morpheus_policies
- hpe_morpheus_power_schedule
- hpe_morpheus_price
- hpe_morpheus_price_set
- hpe_morpheus_provision_type
- hpe_morpheus_resource_pool
- hpe_morpheus_script_template
- hpe_morpheus_security_package
- hpe_morpheus_servicenow_workflow
- hpe_morpheus_spec_template
- hpe_morpheus_storage_bucket
- hpe_morpheus_storage_volume
- hpe_morpheus_storage_volume_type
- hpe_morpheus_task
- hpe_morpheus_tasks
- hpe_morpheus_tenant
- hpe_morpheus_tenants
- hpe_morpheus_user_group
- hpe_morpheus_user_groups
- hpe_morpheus_vdi_pool
- hpe_morpheus_vro_workflow
- hpe_morpheus_workflow

## New known issues

- `hpe_morpheus_cluster_hks_vsphere` has issues with scale-down, which are being investigated.
- `hpe_morpheus_cluster_hks_vsphere` destroy may not succeed, this issue is being investigated.

## Known issues from previous releases

- `hpe_morpheus_datastore` data-source if a datastore with the specified name cannot be found (i.e. the corresponding
  list API request fails), the error message will indicate a 403 (Forbidden) even if the user has permission to list
  datastores.  This is an API bug which is being investigated.
- `hpe_morpheus_instance` updates fail when removing optional fields.
  This will be addressed in a future release.
- `hpe_morpheus_instance` updates fail when removing `evars`.
  This will be addressed in a future release.
- Long running operations can fail when using username and password.
- `hpe_morpheus_instance` depending on the layout used may require one or more `volumes` to be specified,
  in these cases not specifying the correct number of `volumes` will cause instance creation to fail.
- There are intermittent issues with the provider failing to authenticate, a 500 error is returned from the Morpheus API.
  If this happens please retry the operation.  This is being investigated.
- `hpe_morpheus_datastore` when creating a datastore of type NFS the creation will silently fail if the NFS server is not reachable or the share is not accessible.
  The datastore will remain in a `provisioning` state indefinitely. Ensure the Morpheus appliance can reach the NFS server
  and that the share is accessible before creating.
- `hpe_morpheus_datastore` delete is not guaranteed to succeed.  Alletra MP HVM and Alletra MP BM datastores will delete but NFS datastores
  may fail to delete.  Always delete VMs and other resources using the datastore before deleting the datastore itself.
- `hpe_morpheus_instance` in Morpheus versions prior to 8.0.11 requires that the `root` volume is the first entry in
  the `volumes` block list

# v0.4.0 Release Notes

In this release (v0.4.0) we have added the following resource functionality:

- hpe_morpheus_image Update functionality has been added
- hpe_morpheus_instance supports multiple networks and child virtual networks
- hpe_morpheus_instance no longer requires that `ip_mode` is set to avoid forced replaces on Update
- hpe_morpheus_instance supports `timeouts`

In this release (v0.4.0) we have added the following data-source functionality:

- hpe_morpheus_image
- hpe_morpheus_policy

## New known issues

## Known issues from previous releases

- `hpe_morpheus_datastore` data-source if a datastore with the specified name cannot be found (i.e. the corresponding
  list API request fails), the error message will indicate a 403 (Forbidden) even if the user has permission to list
  datastores.  This is an API bug which is being investigated.
- `hpe_morpheus_policy` resource does not currently support the Backup Targets (`backupStorage`) policy type
  due to improper handling of the `backupStorageIds` attribute. This is an API bug which is being investigated.
- `hpe_morpheus_instance` has issues with using the same `datastore_id` with multiple volumes, please use
  a different `datastore_id` for each volume.
- `hpe_morpheus_instance` updates fail when removing optional fields.
  This will be addressed in a future release.
- `hpe_morpheus_instance` updates fail when removing `evars`.
  This will be addressed in a future release.
- Long running operations can fail when using username and password.
- `hpe_morpheus_instance` depending on the layout used may require one or more `volumes` to be specified,
  in these cases not specifying the correct number of `volumes` will cause instance creation to fail.
- There are intermittent issues with the provider failing to authenticate, a 500 error is returned from the Morpheus API.
  If this happens please retry the operation.  This is being investigated.
- `hpe_morpheus_datastore` when creating a datastore of type NFS the creation will silently fail if the NFS server is not reachable or the share is not accessible.
  The datastore will remain in a `provisioning` state indefinitely. Ensure the Morpheus appliance can reach the NFS server
  and that the share is accessible before creating.
- `hpe_morpheus_datastore` delete is not guaranteed to succeed.  AlletraMP HVM datastores will delete but NFS datastores
  may fail to delete.  Always delete VMs and other resources using the datastore before deleting the datastore itself.
- `hpe_morpheus_instance` in Morpheus versions prior to 8.0.11 requires that the `root` volume is the first entry in
  the `volumes` block list

# v0.3.0 Release Notes

In this release (v0.3.0) we have added the following resource functionality:

- hpe_morpheus_image resource has been added (Create, Delete, Read - no Update)
- hpe_morpheus_policy resource has been added (Create, Delete, Read, Update)
- hpe_morpheus_instance Update functionality has been added (The addition and removal of volumes is not yet supported)
- hpe_morpheus_service_plan `cores_per_socket` is now required
- hpe_morpheus_datastore import will now populate `resource_permissions` (`groups` only) and `tenants`

In this release (v0.3.0) we have added the following data-source functionality:

- hpe_morpheus_datastore data-source has been added

## New known issues

- hpe_morpheus_datastore data-source if a datastore with the specified name cannot be found (i.e. the corresponding
  list API request fails), the error message will indicate a 403 (Forbidden) even if the user has permission to list
  datastores.  This is an API bug which is being investigated.
- hpe_morpheus_policy resource does not currently support the Backup Targets (`backupStorage`) policy type
  due to improper handling of the `backupStorageIds` attribute. This is an API bug which is being investigated.
- hpe_morpheus_instance requires that `ip_mode` is set to avoid a forced replace on update.
  This will be addressed in a future release.
- hpe_morpheus_instance updates fail when removing optional fields.
  This will be addressed in a future release.
- hpe_morpheus_instance updates fail when removing `evars`.
  This will be addressed in a future release.
- Long running operations can fail when using username and password.

## Known issues from previous releases

- hpe_morpheus_instance has issues with using the same `datastore_id` with multiple volumes, please use
  a different `datastore_id` for each volume.
- hpe_morpheus_instance depending on the layout used may require one or more `volumes` to be specified,
  in these cases not specifying the correct number of `volumes` will cause instance creation to fail.
- There are intermittent issues with the provider failing to authenticate, a 500 error is returned from the Morpheus API.
  If this happens please retry the operation.  This is being investigated.
- hpe_morpheus_datastore when creating a datastore of type NFS the creation will silently fail if the NFS server is not reachable or the share is not accessible.
  The datastore will remain in a `provisioning` state indefinitely. Ensure the Morpheus appliance can reach the NFS server
  and that the share is accessible before creating.
- hpe_morpheus_datastore delete is not guaranteed to succeed.  AlletraMP HVM datastores will delete but NFS datastores
  may fail to delete.  Always delete VMs and other resources using the datastore before deleting the datastore itself.
- hpe_morpheus_instance only supports 1 network
- hpe_morpheus_instance in Morpheus versions prior to 8.0.11 requires that the `root` volume is the first entry in
  the `volumes` block list

# v0.2.0 Release Notes

In this release (v0.2.0) we have added the following resource functionality:

- hpe_morpheus_datastore resource has been added (Create, Delete, Read, Update)
- hpe_morpheus_service_plan Update functionality has been added
- hpe_morpheus_cloud now has a dynamic `config` block to support arbitrary cloud configuration options
- hpe_morpheus_role now supports setting a `Default Persona`

We have added the following data-source functionality:

- hpe_morpheus_datastore data-source has been added
- hpe_morpheus_role now supports reading `Default Persona` information

We have fixed the following issues:

- hpe_morpheus_user would force recreation if an attribute was updated, this has been fixed
- hpe_morpheus_network switchId is now supported
- hpe_morpheus_role data-source `Default Persona` issue has been fixed

## New known issues

- We have seen an issue with authentication for an existing user when using username/password.  The issue manifests
  as "500" errors on authentication which will not go away on retry.  It is an issue with Morpheus itself and is
  fixed from the `8.0.11` release onwards.  To work around this issue in earlier Morpheus releases
  please generate an `access_token` from the Morpheus UI (for the `morph-api` Client for example) and use
  that instead of username/password.
- hpe_morpheus_datastore when creating a datastore of type NFS the creation will silently fail if the NFS server is not reachable or the share is not accessible.
  The datastore will remain in a `provisioning` state indefinitely. Ensure the Morpheus appliance can reach the NFS server
  and that the share is accessible before creating.
- hpe_morpheus_datastore delete is not guaranteed to succeed.  AlletraMP HVM datastores will delete but NFS datastores
  may fail to delete.  Always delete VMs and other resources using the datastore before deleting the datastore itself.
- hpe_morpheus_instance only supports 1 network
- hpe_morpheus_instance in Morpheus versions prior to 8.0.11 requires that the `root` volume is the first entry in
  the `volumes` block list

## Known Issues from previous releases

- hpe_morpheus_instance has issues with using the same `datastore_id` with multiple volumes, please use
  a different `datastore_id` for each volume.
- hpe_morpheus_instance depending on the layout used may require one or more `volumes` to be specified,
  in these cases not specifying the correct number of `volumes` will cause instance creation to fail.
- There are intermittent issues with the provider failing to authenticate, a 500 error is returned from the Morpheus API.
  If this happens please retry the operation.  This is being investigated.

# v0.1.0 Release Notes

## New functionality

In this release (v0.1.0) the following resources have been added:

- hpe_morpheus_cloud for HPE HVM or HPE VME clouds
- hpe_morpheus_group
- hpe_morpheus_instance for HPE HVM or HPE VME instances (Create, Delete and Read - no Update)
- hpe_morpheus_network
- hpe_morpheus_role for Morpheus roles (user and tenant)
- hpe_morpheus_service_plan (Create, Delete and Read - no Update)
- hpe_morpheus_user (Create, Delete and Read - no Update)

In this release (v0.1.0) the following data sources have been added:

- hpe_morpheus_cloud
- hpe_morpheus_environment
- hpe_morpheus_group
- hpe_morpheus_instance_type_layout
- hpe_morpheus_network
- hpe_morpheus_role
- hpe_morpheus_service_plan

## Known Issues

- hpe_morpheus_instance has issues with using the same `datastore_id` with multiple volumes, please use
  a different `datastore_id` for each volume.
- hpe_morpheus_instance depending on the layout used may require one or more `volumes` to be specified,
  in these cases not specifying the correct number of `volumes` will cause instance creation to fail.
- hpe_morpheus_network switchId is not supported yet, prevents creating some network types, eg OVS Port Group
- hpe_morpheus_user will force recreation if an attribute is updated
- There are intermittent issues with the provider failing to authenticate, a 500 error is returned from the Morpheus API.
  If this happens please retry the operation.  This is being investigated.
