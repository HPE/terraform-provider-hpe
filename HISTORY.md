# Release History

Release notes for versions of this provider prior to the current release.
The current release is documented in [CHANGELOG.md](./CHANGELOG.md).

# v2.0.0 Release Notes

This is a major release.  Alongside the Morpheus support this provider already offered, it adds
support for **HPE OpsRamp**, introduces PCE (Private Cloud Enterprise) Identity authentication for
Morpheus, and closes the remaining hpegl VMaaS parity gaps.

## Breaking changes

### Write-only attributes

Six attributes are accepted by the Morpheus API on write but never returned on read.  They were
stored in Terraform state anyway, so state held values the provider could not refresh and drift was
undetectable.  They are now write-only, and each gains a `_version` companion used to signal a
change:

| Resource | Attributes now write-only |
|---|---|
| `hpe_morpheus_network_router_nat` | `action`, `firewall`, `service` |
| `hpe_morpheus_load_balancer` | `group_id`, `network_server_id` |
| `hpe_morpheus_instance_clone` | `source_instance_id` |

These attributes are no longer stored in state, so changing one on its own produces no plan diff —
increment the matching `_version` attribute to apply a change.  On `hpe_morpheus_load_balancer`,
`group_id_version` and `network_server_id_version` force replacement, because neither value can be
changed on an existing load balancer.  On `hpe_morpheus_instance_clone`,
`source_instance_id_version` has no effect, because the clone source cannot be changed after
creation; it exists so that a value may be supplied without error.

Write-only attributes require **Terraform 1.11 or later**.

`hpe_morpheus_network_router_nat.firewall` now defaults to `MATCH_INTERNAL_ADDRESS` on create.  On
update an omitted `firewall` is left out of the payload entirely, so the value already on the router
is preserved rather than overwritten.

`hpe_morpheus_instance_clone` no longer recovers `source_instance_id` from `config.cloneInstanceId`,
as a write-only attribute must be null in state.

### tfmigrator release artifacts renamed

The migration tool's release artifacts are now published as `tfmigrator_*` rather than
`migration_tool_*`.  This affects the archives, the binary inside them, and the checksum files:

| Was | Now |
|---|---|
| `migration_tool_<version>_<os>_<arch>.zip` | `tfmigrator_<version>_<os>_<arch>.zip` |
| `migration_tool_v<version>` | `tfmigrator_v<version>` |
| `migration_tool_<version>_SHA256SUMS` (+ `.sig`) | `tfmigrator_<version>_SHA256SUMS` (+ `.sig`) |

The installed binary is still called `tfmigrator`, so nothing changes once the tool is on your
PATH — only the download URL and the name of the file inside the archive.  Any automation that
fetches the archive by name needs updating.

The `install-tfmigrator` scripts have been updated to match and therefore support v2.0.0 and later.
To install an earlier version, download the `migration_tool_*` archive manually from the releases
page.

`tfmigrator --version` now reports the release it was built from.  Previously it reported a version
compiled into the source, which did not track the release.

## Deprecations

- `hpe_morpheus_instance` — `server_uuids` is deprecated in favour of the new `server_uuid`
  (String).  Morpheus assigns UUIDs strictly by position and an instance provisions exactly one
  server, so only the first element of the set was ever used and the rest were discarded silently.
  Both attributes continue to work and are mutually exclusive.
- `hpe_morpheus_task_powershell_script` and `hpe_morpheus_task_shell_script` —
  `remote_target_password` is deprecated in favour of the write-only `remote_target_password_wo`
  with `remote_target_password_wo_version`.  Morpheus returns the password as a hash, so the
  plaintext in configuration never matched the hash in state and every plan was non-empty.
- `hpe_morpheus_network_domain` — `domain_password` is deprecated in favour of the write-only
  `domain_password_wo` with `domain_password_wo_version`.  `auto_join_domain` is deprecated.
- `hpe_morpheus_cloud_affinity_group` and `hpe_morpheus_cluster_affinity_group` — `tenant_ids` is
  deprecated and has no effect, because the Morpheus API rejects tenant assignment on affinity
  groups.  `hpe_morpheus_cluster_affinity_group.description` is deprecated because it is not backed
  by the API.

## OpsRamp support

This provider now serves HPE OpsRamp resources and data sources alongside Morpheus, configured with
an `opsramp` block in the provider configuration.  See the
[OpsRamp to HPE migration guide](./docs/guides/opsramp_to_hpe_migration.md) for moving an existing
OpsRamp provider configuration across.

## PCE Identity authentication

The `morpheus` provider block accepts two new mutually exclusive blocks, so Morpheus connection
details can be obtained from GreenLake rather than configured by hand:

- `pce_identity` — Connected PCE, using GLCS IAM and scoped by GreenLake Space
- `pce_disconnected_identity` — Disconnected PCE, using GLP IAM and scoped by GreenLake Workspace

Each accepts either GreenLake API client credentials — `client_id`, `client_secret` and an issuer
URL, which is `issuer_url` in the Connected block and `token_issuer_url` in the Disconnected one —
or a pre-generated `iam_token`.  Neither can be combined with `url`, `username`, `password`,
`access_token` or `tenant_subdomain`.

## New resources

In this release (v2.0.0) we have added the following resources:

### Morpheus

- hpe_morpheus_cloud_affinity_group
- hpe_morpheus_cloud_affinity_group_member
- hpe_morpheus_cluster_affinity_group_member
- hpe_morpheus_instance_node
- hpe_morpheus_network_router_firewall_rule_group

### OpsRamp

- hpe_opsramp_alert_correlation_policy
- hpe_opsramp_alert_escalation_policy
- hpe_opsramp_alert_prediction_policy
- hpe_opsramp_client
- hpe_opsramp_credential_set
- hpe_opsramp_custom_integration
- hpe_opsramp_device_group
- hpe_opsramp_first_response_policy
- hpe_opsramp_integration
- hpe_opsramp_integration_app
- hpe_opsramp_integration_config
- hpe_opsramp_integration_event
- hpe_opsramp_kb_article
- hpe_opsramp_kb_category
- hpe_opsramp_log_alert_definition
- hpe_opsramp_management_profile
- hpe_opsramp_metric_alert_definition
- hpe_opsramp_permission_set
- hpe_opsramp_resource
- hpe_opsramp_role
- hpe_opsramp_scheduled_maintenance
- hpe_opsramp_script
- hpe_opsramp_script_category
- hpe_opsramp_servicedesk_business_impact
- hpe_opsramp_servicedesk_category
- hpe_opsramp_servicedesk_urgency
- hpe_opsramp_servicemap
- hpe_opsramp_servicemap_link
- hpe_opsramp_site
- hpe_opsramp_user
- hpe_opsramp_user_group

## New data sources

In this release (v2.0.0) we have added the following data sources:

### Morpheus

- hpe_morpheus_cloud_affinity_group
- hpe_morpheus_cloud_affinity_groups
- hpe_morpheus_cluster_affinity_groups
- hpe_morpheus_clusters
- hpe_morpheus_compute_server
- hpe_morpheus_compute_servers
- hpe_morpheus_instance_disk_type
- hpe_morpheus_instance_storage_controller
- hpe_morpheus_network_interface_type
- hpe_morpheus_network_proxy
- hpe_morpheus_network_router_firewall_rule_group
- hpe_morpheus_network_server_group

### OpsRamp

- hpe_opsramp_custom_event_alert_source
- hpe_opsramp_resource_lookup
- hpe_opsramp_role
- hpe_opsramp_servicedesk_business_impact
- hpe_opsramp_servicedesk_category
- hpe_opsramp_servicedesk_urgency
- hpe_opsramp_tenant

## Enhancements to existing resources

- `hpe_morpheus_instance` — added `wait_for_ip_address`, an opt-in wait that holds the apply until
  at least one container reports a usable address, rather than recording the `0.0.0.0` placeholder
  Morpheus returns for a container that has not reported yet.  On expiry it warns and continues.
  Added `config_vmware.affinity_group_id` and `config_hvm.affinity_group_id` to
  place an instance into an affinity group at provision time; create-only, and rejected alongside
  `config_hvm.kvm_host_id`.  Added the computed `compute_servers` and `container_id`,
  and the new `server_uuid`.
- `hpe_morpheus_cluster_affinity_group` — reworked.  The shipped resource could not set an affinity
  type and could not manage membership at all; it now supports `affinity_type`, `pool_id`, `servers`
  and `source`.
- `hpe_morpheus_load_balancer` — added `enabled`, to activate or disable a load balancer on create
  and update.
- `hpe_morpheus_network` — added `connected_gateway`, the provider ID of a connected NSX-T Tier-1
  gateway.
- `hpe_morpheus_network_router` — `provider_id` is available again, for configurations that
  reference it when building dependent resources.
- `hpe_morpheus_network_router_nat` — added `translated_ports`.
- `hpe_morpheus_network_router_route` — added `priority`, which forces replacement to match the
  API's behaviour.
- `hpe_morpheus_compute_server` and `hpe_morpheus_compute_servers` — report `parent_host_id` and
  `parent_host_name`, the hypervisor host a guest runs on; the plural data source can filter on it.
- `hpe_morpheus_cluster_affinity_group` (data source) — now exposes `servers`, `source`,
  `tenant_ids` and `resource_permissions`.
- `hpe_morpheus_network_dhcp_server` (data source) — added `provider_id`.
- `hpe_morpheus_network_pool` (data source) — added `display_name`.

## Resolved issues

- `hpe_morpheus_instance` — unrelated computed attributes no longer churn the plan; a small edit
  used to show `connection_info`, `labels` and every computed field of `network_interfaces` and
  `volumes` as `(known after apply)`.  An imported instance no longer plans changes
  nobody made, which on appliances before 8.1.2 escalated to a replacement of a running VM.
  Instances created by the hpegl provider no longer fail to read: string-encoded `noAgent` is
  handled, and an absent `nestedVirtualization` is treated as optional rather than an error.
- `hpe_morpheus_subnet` — `resource_permission_groups_all` was sent under a request key the Morpheus
  API does not read, so the setting was silently dropped; an explicitly configured `pool_id` was
  overwritten with `null` when the API response omitted the pool.
- `hpe_morpheus_network_domain` — `public_zone`, `visibility` and `active` are now sent on update,
  so changing them takes effect; `auto_join_domain` is preserved on import and `tenant_id` is read
  back.
- `hpe_morpheus_network_router_firewall_rule` — creating a rule without a `description` no longer
  fails; the required format of `parent_id` is documented.
- `hpe_morpheus_network_router` — BGP neighbor configuration is read correctly on import, and API
  flags returned as JSON booleans are handled alongside the `on`/`off` strings.
- `hpe_morpheus_os_type_image` — inconsistent `os_type_id` after apply.
- `hpe_morpheus_tenant` — `currency` is validated against the supported ISO codes.
- `hpe_morpheus_option_type` — rows and description are validated at plan time.
- `hpe_morpheus_service_plan` and `hpe_morpheus_datastore` — name lookups now work on Private Cloud
  appliances.  Both looked the object up by name and then re-fetched it by id, and it was that
  second request that failed.
- `hpe_morpheus_load_balancer` — resources no longer return state inconsistent with the plan, and
  the sweepers keep up with leaked NSX-T load balancer services.
- `hpe_morpheus_instance_clone` — clone failures are now detected and reported.  Create polled only
  for the clone's name and never inspected the `cloning` process, so a server-side failure was
  indistinguishable from a slow clone and the diagnostic reported `<nil>`.  The documented examples
  used block syntax for `volumes` and `network_interfaces`, which are list nested attributes, so
  copying them produced `Blocks of type "volumes" are not expected here`; they now use
  list-of-object syntax.
- hpegl to hpe migration — several Read-vs-plan consistency errors that blocked `terraform apply`
  with `import` blocks are fixed, covering instance, instance clone, load balancer and its profiles,
  monitors, pool and virtual server, network, network router, BGP neighbor, NAT rule, static route
  and firewall rule group.
- API tracing (`MORPHEUS_API_HTTPTRACE`) now redacts the `Authorization` header, the appliance
  password sent to `/oauth/token`, and the `access_token` and `refresh_token` in the response.
  These traces are captured in test output and kept as build artifacts, so credentials were being
  written to logs in clear text.
- Resources with a `Dynamic` `config` attribute — the provider no longer panics when a value inside
  `config` is not known until apply, for example `config = { templateId = var.image_id }`.  A literal
  worked, and so did arithmetic between literals, because Terraform folds those at parse time, so the
  crash appeared only once a value was deferred — taking an id from a `.tfvars` file was enough.
  Unknown is now treated as absent, as null already was.
- `hpe_morpheus_image` — reading an image with two or more tenants failed with `Duplicate Set
  Element`, whether or not anything was duplicated.  The tenant objects were built without a known
  state, so their `name` and `id` were discarded and every tenant became identical.

## Known issues

- `hpe_morpheus_cluster_namespace`: `active` is not supported on import. `name` update is not supported.
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

# v1.6.0 Release Notes

In this release (v1.6.0) we have added the following resources:

### Networking

- hpe_morpheus_load_balancer_profile

In this release (v1.6.0) we have added the following data sources:

- hpe_morpheus_backup_host
- hpe_morpheus_backup_instance
- hpe_morpheus_backup_job
- hpe_morpheus_certificate
- hpe_morpheus_cluster_affinity_group
- hpe_morpheus_cluster_layout
- hpe_morpheus_cluster_namespace
- hpe_morpheus_container_script
- hpe_morpheus_datastore_types
- hpe_morpheus_datastores
- hpe_morpheus_deployment
- hpe_morpheus_load_balancer_profile
- hpe_morpheus_monitoring_alert
- hpe_morpheus_monitoring_check
- hpe_morpheus_monitoring_check_type
- hpe_morpheus_monitoring_group
- hpe_morpheus_network_pool_server
- hpe_morpheus_network_pool_server_type
- hpe_morpheus_network_router_firewall_rule
- hpe_morpheus_network_router_firewall_rule_groups
- hpe_morpheus_network_router_nat
- hpe_morpheus_network_router_type
- hpe_morpheus_provisioning_license
- hpe_morpheus_security_group
- hpe_morpheus_security_group_rule
- hpe_morpheus_security_groups
- hpe_morpheus_storage_server
- hpe_morpheus_storage_servers
- hpe_morpheus_storage_volume (reimplemented)
- hpe_morpheus_storage_volumes
- hpe_morpheus_subnet_type
- hpe_morpheus_vdi_app
- hpe_morpheus_vdi_gateway

## Enhancements to existing resources

- hpe_morpheus_instance — Added `config_bmaas` static config block for HPE bare metal (BMaaS) instances; added `host_name`, `labels` and `server_uuids`; added `user_group` and `storage_profile`
- hpe_morpheus_datastore — Added `config_alletramp_bmaas` static config block for HPE Alletra MP Bare Metal (BMaaS) datastores; the datastore type `code` is now preferred (`type_id` is optional and resolved from the code)
- hpe_morpheus_storage_volume — Added a typed write-only `config` block (the resource now uses `WriteOnly` and a `Dynamic` `config` attribute)
- hpe_morpheus_network_pool_server — Added EfficientIP (SOLIDserver) support
- hpe_morpheus_role — Added `tenant_copies`
- hpe_morpheus_cluster_affinity_group and hpe_morpheus_cluster_namespace — `resource_permissions.groups` now supports a per-group `default` flag
- Added `visibility`, `resource_permissions` and `tenant_ids` support to several resources introduced in v1.4.0 (hpe_morpheus_cluster_affinity_group, hpe_morpheus_cluster_namespace, hpe_morpheus_network_group, hpe_morpheus_network_router, hpe_morpheus_power_schedule, hpe_morpheus_provisioning_license)

## Resolved issues

- `hpe_morpheus_network_router` import and `enable_bgp` handling; conflicting `group_id` and tenant permissions are now rejected at plan time; NAT protocol round-trip and NAT/firewall/BGP handling
- `hpe_morpheus_vdi_pool` now requires `max_idle` >= `min_idle`; fixed `idle_timeout` and `max_session_timeout` handling
- `hpe_morpheus_vdi_gateway` `gateway_url` is now optional; `description` is limited to 255 characters
- `hpe_morpheus_vdi_app` `launch_prefix` is now required
- `hpe_morpheus_storage_volume` validates `max_storage` for the selected `type_id` at plan time, and requires a storage server or storage group while rejecting a conflicting export target; changing the write-only `config` now correctly recreates the volume
- `hpe_morpheus_container_script` name/phase validation and CRLF handling
- `hpe_morpheus_instance` Read panic on a 404 from the environment-variables endpoint; volume matching on resize; externally-attached SAN volumes are excluded from volume tracking; `stopped` is accepted as a valid create status
- `hpe_morpheus_setting_whitelabel` now validates hex color format
- `hpe_morpheus_subnet` unknown `cidr`, `netmask` and `subnet_address` on create
- `hpe_morpheus_cloud` write-only secrets dropped on create and update
- `hpe_morpheus_node_type` template ID list read-back
- `hpe_morpheus_spec_template_*` `source_type` input validation
- `hpe_morpheus_app_blueprint_kubernetes` Read panic
- The underlying transport error is now surfaced by the legacy client

## Known issues

- `hpe_morpheus_cluster_namespace`: `active` is not supported on import. `name` update is not supported.
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

# v1.5.0 Release Notes

In this release (v1.5.0) we have added the following resources:

### Networking

- hpe_morpheus_load_balancer_pool

### Compute & Instances

- hpe_morpheus_instance_clone
- hpe_morpheus_instance_power_state
- hpe_morpheus_instance_snapshot

### Backup

- hpe_morpheus_backup_host
- hpe_morpheus_backup_instance

In this release (v1.5.0) we have added the following data sources:

- hpe_morpheus_backup
- hpe_morpheus_backup_type
- hpe_morpheus_instance_snapshot
- hpe_morpheus_load_balancer_pool
- hpe_morpheus_network_edge_cluster
- hpe_morpheus_network_pool
- hpe_morpheus_network_server
- hpe_morpheus_network_transport_zone
- hpe_morpheus_network_type

## Enhancements to existing resources

- hpe_morpheus_setting_whitelabel — `header_logo`, `footer_logo`, `login_logo` and `favicon` now take a local image file path and are uploaded to Morpheus (png/jpg/svg for logos; ico/png for the favicon). Remote URLs are not supported.
- hpe_morpheus_instance — Added `subnet_id` to `network_interfaces` (required by some clouds, e.g. Azure); added a computed `status` attribute and out-of-band deletion detection (the instance is removed from state when the underlying VM no longer exists)
- hpe_morpheus_service_plan (data source) — Added a `cloud_id` filter (region/zone)
- hpe_morpheus_identity_source_saml — Exposed the SP metadata (`entity_id`, `acs_url`) as computed outputs
- hpe_morpheus_option_list_rest — Added `inject_system_authorization_header` and `use_owner_auth`
- Option type resources and hpe_morpheus_form — Reject a self-referential `dependent_field` (circular `dependsOnCode`) at plan time

## Resolved issues

- `hpe_morpheus_setting_whitelabel` accepted `header_logo`, `footer_logo`, `login_logo` and `favicon` on apply but returned them as null on refresh; the logos are now uploaded via the whitelabel images endpoint and the configured path is preserved in state
- `hpe_morpheus_form` dropped attributes on read for the secGroup field
- `hpe_morpheus_form` cross-field leakage between option_type reads
- `hpe_morpheus_form` no attribute generated for plain cascade keys
- `hpe_morpheus_price_set` type enum and create/update failure handling
- Appliance URLs with trailing slashes now authenticate correctly against Morpheus 9.0
- New resources re-read via GET after POST and taint their state on failure, improving reliability

## Known issues

- `hpe_morpheus_cluster_namespace`: `active` is not supported on import. `name` update is not supported.
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
- hpe_morpheus_provisioning_license

### VDI

- hpe_morpheus_vdi_app
- hpe_morpheus_vdi_gateway
- hpe_morpheus_vdi_pool

### Compute & Cluster

- hpe_morpheus_cluster_affinity_group
- hpe_morpheus_cluster_namespace

### Identity & Governance

- hpe_morpheus_setting_whitelabel

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
- hpe_morpheus_security_group
- hpe_morpheus_security_groups

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

- `hpe_morpheus_cluster_namespace`: `active` is not supported on import. `name` update is not supported.
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
