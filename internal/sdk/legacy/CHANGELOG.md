## Unreleased

NOTES:

* Add support for HPE VM Essentials host maintenance and placement

## 0.5.1 (September 10, 2024)

NOTES:

* Update the policy response payload to handle power schedule policies
* Update the integrations response payload to handle chef integrations

## 0.5.0 (September 9, 2024)

NOTES:

* Update resty library from 2.12.0 to 2.14.0.
* Update cloud response payload
* Update network pool server response payload
* Update integration response payload

FEATURES:

* **Backups** Client methods for `backups` API endpoints
* **Backup Jobs** Client methods for `backup_jobs` API endpoints

## 0.4.3 (August 16, 2024)

NOTES:

* Add network pool ip address management support
* Update instance datstoreid to an interface to handle edge cases in which the value is a null value, string or integer.

## 0.4.2 (July 9, 2024)

NOTES:

* Updates to the clusters payload
* Updates to apps payload to add more details and support new actions
* Updates to instances to support additional actions (i.e. - snapshots, resize, clone, etc.).

## 0.4.1 (June 24, 2024)

FEATURES:

* **Hosts** Client methods for `hosts` API endpoints
* **Network Routers** Client methods for `network_routers` API endpoints

## 0.4.0 (May 27, 2024)

NOTES:

* Updates to instances to support add support for resize, lock, and unlock.
* Add additional payload attributes to virtual images.

FEATURES:

* **Security Groups** Client methods for `security_groups` API endpoints
* **Email Templates** Client methods for `email_templates` API endpoints
* **Storage Volumes** Client methods for `storage_volumes` API endpoints
* **Storage Volume Types** Client methods for `storage_volume_types` API endpoints
* **Task Types** Client methods for `task_types` API endpoints

## 0.3.9 (April 24, 2024)

NOTES:

* Fix issues with the `user_group` delete function and added support for newer API attributes.
* Update resty library from 2.10.0 to 2.12.0.
* Update `instance_type` and `instance_layout` resources to add support for price sets.

## 0.3.8 (January 15, 2024)

NOTES:

* Fix issues with the `user_group` delete function and added support for newer API attributes.

## 0.3.7 (November 20, 2023)

FEATURES:

* **Cluster Packages** Client methods for `cluster_packages` API endpoints
* **Forms** Client methods for `forms` API endpoints

## 0.3.6 (November 13, 2023)

NOTES:

* Update cloud, cypher, and tenant response payloads.

## 0.3.5 (November 7, 2023)

NOTES:

* Update clouds response payload along with adding support for cloud datastores and cloud resource pools.
* Update instances response payload to add support for network domains and stats.

FEATURES:

* **Cypher** Client methods for `cypher` API endpoints

## 0.3.4 (June 4, 2023)

NOTES:

* Update task response payload.

FEATURES:

* **Service Catalog** Client methods for `service_catalog` API endpoints

## 0.3.3 (May 11, 2023)

NOTES:

* Update instance and app response payloads.
* Fix an issue with the update and delete methods for resource pool management in which a cloud id was erroneously required.
* Add additional methods for managing key pairs

FEATURES:

* **Job Executions** Client methods for `job_executions` API endpoints

## 0.3.2 (April 18, 2023)

NOTES:

* Additional response payload updates.

## 0.3.1 (March 15, 2023)

NOTES:

* Update automation tasks payload.

FEATURES:

* **Resource Pool Groups** Client methods for `resource_pool_groups` API endpoints

## 0.3.0 (March 14, 2023)

NOTES:

* Update vRO payload.

FEATURES:

* **Guidance Settings** Client methods for `guidance_settings` API endpoints
* **Monitoring Settings** Client methods for `monitoring_settings` API endpoints
* **Security Scans** Client methods for `security_scans` API endpoints

## 0.2.10 (March 1, 2023)

NOTES:

* Update resty library from 2.2.0 to 2.7.0.
* Update option type response payload.
* Update roles response payload.

## 0.2.9 (February  14, 2023)


NOTES:

* Add label support for additional objects.
* Updated the return payload for networks.

FEATURES:

* **Execution Request** Client methods for `execution_request` API endpoints
* **Network Groups** Client methods for `network_groups` API endpoints
* **Oauth Clients** Client methods for `oauth_clients` API endpoints


## 0.2.8 (December 27, 2022)


NOTES:

* Updated the return payload for clouds.
* Updated the return payload for networks.

FEATURES:

* **Storage Servers** Client methods for `storage_servers` API endpoints

## 0.2.7 (December 26, 2022)

NOTES:

* Updated the return payload for clouds.
* Updated the return payload for networks.

FEATURES:

* **Network Pools** Client methods for `network_pools` API endpoints
* **Network Pool Servers** Client methods for `network_pool_servers` API endpoints
* **Network Routes** Client methods for `network_routes` API endpoints
* **Network Static Routes** Client methods for `network_static_routes` API endpoints
* **Network Subnets** Client methods for `network_subnets` API endpoints

## 0.2.6 (December 22, 2022)

NOTES:

* Updated the return payload for clouds.

FEATURES:

* **Activity** Client methods for `activity` API endpoints
* **Load Balancer Monitors** Client methods for `load_balancer_monitors` API endpoints
* **Load Balancer Pools** Client methods for `load_balancer_pools` API endpoints
* **Load Balancer Types** Client methods for `load_balancer_types` API endpoints
* **Load Balancer Virtual Servers** Client methods for `load_balancer_virtual_servers` API endpoints
* **Load Balancers** Client methods for `load_balancers` API endpoints
* **Network Proxies** Client methods for `network_proxies` API endpoints* 
* **Plugins** Client methods for `plugins` API endpoints
* **Security Packages** Client methods for `security_packages` API endpoints

## 0.2.5 (December 12, 2022)

NOTES:

* Update workflow payloads to return the task sets data.

## 0.2.4 (December 2, 2022)

NOTES:

* Update execute schedule data payload.
* Fix invalid policy data payload.
* Update task data payload.

## 0.2.3 (November 14, 2022)

NOTES:

* Update cluster resource naming policy data type for auto resolve naming conflict setting.

## 0.2.2 (November 10, 2022)

NOTES:

* Fix user role creation response payload struct
* Fix delayed delete policy response payload struct

## 0.2.1 (October 27, 2022)

NOTES:

* Update policies code
* Update data structs

## 0.2.0 (October 7, 2022)

NOTES:

* Added additional tests and examples for existing resources
* Updated the response payloads for existing resources

FEATURES:

* **Alerts** Client methods for `alerts` API endpoints
* **Check Apps** Client methods for `check_apps` API endpoints
* **Deployments** Client methods for `deployments` API endpoints

## 0.1.10 (September 15, 2022)

* **Archives** Client methods for `archive_buckets` API endpoints

## 0.1.9 (September 13, 2022)

FEATURES:

* **Boot Scripts** Client methods for `boot_scripts` API endpoints
* **Preseed Scripts** Client methods for `preseed_scripts` API endpoints

## 0.1.8 (September 12, 2022)

FEATURES:

* **Users** Client methods for `users` API endpoints
* **User Groups** Client methods for `user_groups` API endpoints
* **Whitelabel Settings** Client methods for `whitelabel_settings` API endpoints

## 0.1.7 (September 8, 2022)

FEATURES:

* **Appliance Settings** Client methods for `appliance_settings` API endpoints
* **Backup Settings** Client methods for `backup_settings` API endpoints
* **Cluster Types** Client methods for `cluster_types` API endpoints
* **License** Client methods for `license` API endpoints
* **Provisioning Settings** Client methods for `provisioning_settings` API endpoints
* **VDI Allocations** Client methods for `vdi_allocations` API endpoints
* **VDI Apps** Client methods for `vdi_aps` API endpoints

## 0.1.6 (September 8, 2022)

Retracted Release

## 0.1.5 (August 26, 2022)

FEATURES:

* **Budgets** Client methods for `budgets` API endpoints
* **Credentials** Client methods for `credentials` API endpoints
* **Identity Sources** Client methods for `user_sources` API endpoints
* **Scale Thresholds** Client methods for `scale_thresholds` API endpoints
* **Software Licenses** Client methods for `software_licenses` API endpoints
* **Storage Buckets** Client methods for `storage_buckets` API endpoints
* **VDI Gateways** Client methods for `vdi_gateways` API endpoints
* **VDI Pools** Client methods for `vdi_pools` API endpoints

## 0.1.4 (August 3, 2022)

FEATURES:

* **Check Groups** Client methods for `check_groups` API endpoints
* **Checks** Client methods for `checks` API endpoints
* **Contacts** Client methods for `contacts` API endpoints
* **Incidents** Client methods for `incidents` API endpoints
* **Jobs** Client methods for `jobs` API endpoints
* **Monitoring Apps** Client methods for `monitoring_apps` API endpoints
* **Node Types** Client methods for `node_types` API endpoints
* **Policies** Client methods for `policies` API endpoints
* **Prices** Client methods for `prices` API endpoints
* **Price Sets** Client methods for `price_sets` API endpoints
* **Provision Types** Client methods for `provision_types` API endpoints
* **Roles** Client methods for `roles` API endpoints
* **Service Plans** Client methods for `service_plans` API endpoints
* **Wikis** Client methods for `wikis` API endpoints

## 0.3 (May 17, 2021)

FEATURES:

* **Instance Layouts** Client methods for `instance_layouts` API endpoints
* **Instance Types** Client methods for `instance_types` API endpoints
* **Plans** Client methods for `plans` API endpoints
* **Resource Pools** Client methods for `resource_pools` API endpoints

## 0.2 (May 7, 2021)

FEATURES:

* **Accounts** Client methods for `accounts` (Tenants) API endpoints
* **Apps** Client methods for `apps` API endpoints
* **Blueprints** Client methods for `blueprints` API endpoints
* **Catalog Items** Client methods for `catalog_items` API endpoints
* **Clusters** Client methods for `clusters` API endpoints
* **Environments** Client methods for `environments` API endpoints
* **Option Lists** Client methods for `option_lists` API endpoints
* **Option Types** Client methods for `option_types` API endpoints
* **Task Sets** Client methods for `task_sets` API endpoints
* **Tasks** Client methods for `tasks` API endpoints

## 0.1 (November 27, 2019)

NOTES:

* This is a **BETA** version of the Morpheus Go SDK.

FEATURES:

* **Client** Base client for API communication with remote morpheus appliance.
* **Groups:** Client methods for `groups` API endpoints.
* **Clouds:** Client methods for `clouds` API endpoints.
* **Instances:** Client methods for `instances` API endpoints.
* **Network:** Client methods for `networks` API endpoints.
* **NetworkDomains:** Client methods for `network-domains` API endpoints.
* **Setup:** Client methods for `setup` API endpoints.
* **Whoami:** Client methods for `whoami` API endpoints.
