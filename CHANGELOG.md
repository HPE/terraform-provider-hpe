# v2.2.0 Release Notes

## Support

**Support policy** — Bug fixes and new features are shipped in the next release of the provider. We
do not backport changes to previous releases; if you encounter a problem, please upgrade to the
latest version before opening an issue. Security vulnerabilities may be backported to the final
release of the previous major version at the maintainers' discretion, on a case-by-case basis. The
same policy applies to the migration tool (tfmigrator).

**Support note** — HPE Morpheus functionality in this provider is fully supported by HPE subject to
the support policy above. Support for HPE OpsRamp functionality is currently provided on a
best-efforts basis.

Release notes for earlier versions are in [HISTORY.md](./HISTORY.md).

## Breaking changes

Write-only attributes require **Terraform 1.11 or later, or OpenTofu 1.11 or later**.

### Network configuration

`hpe_morpheus_network` has typed configuration blocks for HVM, OVS, NSX, Azure, GCP and OpenStack
networks. Use one block matching the network type; `config` remains available for other types.

- Configuration blocks are write-only. Replace references to `hpe_morpheus_network.x.config` with
  `hpe_morpheus_network.x.config_from_api`. Existing state is upgraded automatically.
- Specify either `type_code` or `type_id`.
- Replacement rules follow the network type's editable fields. Unsupported block attributes are
  rejected at plan time.

See the [network resource documentation](https://registry.terraform.io/providers/HPE/hpe/latest/docs/resources/morpheus_network)
for block names, fields and examples.

### Licence management

`hpe_morpheus_license.key` is write-only. The attribute name is unchanged; refresh removes the key
from state automatically.

- `id` and the import ID are now the key's first eight characters, not `1`.
- Licence details are read back, with unspecified values reported as `null` rather than empty
  strings, zero or false. Update expressions that rely on those values.
- A changed key is installed in place. A missing key is reinstalled on the next apply.
- Destroy uninstalls only the matching licence. Use `terraform state rm` to stop managing it
  without uninstalling it.

### Instances and clones

- `hpe_morpheus_instance`: `evars` can be added, changed and removed in place. Removing the block
  removes the managed environment variables. Changing `ports` replaces the instance.
- `hpe_morpheus_instance_clone`: refresh excludes CD/DVD drives from `volumes`. Remove any CD/DVD
  entry previously added to the configuration to suppress drift; otherwise Terraform plans a new
  disk. Volume renames and reordering are rejected at plan time.

## Deprecations

- Write-only `*_version` and `*_wo_version` markers are deprecated and scheduled for removal in
  **2.6**. See [Write-only values](#write-only-values) for upgrade instructions.
- Network `connected_gateway`: use the field inside `config_nsxt` or `config_nsx_cloud`.
- Cloud `config_aws.backup_provider` has no effect on Morpheus 9.x; use `backup_provider_ids`.
- Network domain `auto_join_domain` (already a no-op) and `domain_password`, and shell/PowerShell task
  `remote_target_password`, remain accepted until their removal in 2.6.

## New resources and data sources

- **`hpe_morpheus_plugin`** resource and data source: install, update, configure and remove plugins.
  Supply `jar_path` or `jar_url`, optionally pinned with `jar_sha256`. Changed jars trigger updates;
  replacing a plugin with a different plugin code is rejected. Look up installed plugins by `id`
  or `code`.
- **`hpe_morpheus_option_source_value`** and **`hpe_morpheus_option_source_values`**: query the values
  behind Morpheus dropdowns, including selectable networks and HVM cluster resource pools. The
  singular lookup requires exactly one match; the plural lookup supports filters. Use
  `credential_id` for authentication, not credentials in `parameters`.

## OpsRamp changes

- Added `hpe_opsramp_template` resource and data source for G2 monitoring templates using
  `metric_groups` (SNMP, JMX, OpenTelemetry and remote-code execution). Generation-1 templates
  are not supported.
- Added `hpe_opsramp_device_management_policy` resource and data source. Policies select resources
  and assign monitoring settings; create and update also apply the policy.
- Added `hpe_opsramp_metric` and `hpe_opsramp_adapter_type` data sources for template configuration.
- Added import support for alert escalation policies and device groups.
- Fixed scheduled-maintenance drift caused by equivalent timestamps with different timezone offsets.

## Enhancements to existing resources

### Write-only values

Changing a write-only value now triggers an update without a manual version bump.
`write_only_revision` records the change in the plan.

- Remove deprecated markers from your configuration. This causes one in-place update, not a
  replacement. Changing a create-only value still replaces the resource.
- The first plan after upgrading from 2.1 shows an in-place `write_only_revision` update per
  migrated resource. Applying records digests and may re-send configured values.
- After import, change detection starts with the next resource update. To enable it immediately,
  include each value's SHA-256 in the import ID, for example
  `"123#password_wo=<sha256 hex>&linux_password_wo=<sha256 hex>"`.

| Resource | Write-only attribute(s) | Deprecated marker(s) |
|---|---|---|
| `hpe_morpheus_user` | `password_wo`, `linux_password_wo`, `windows_password_wo` | `*_wo_version` |
| `hpe_morpheus_backup_host` | `ssh_password_wo` | `ssh_password_wo_version` |
| `hpe_morpheus_image` | `ssh_password_wo` | `ssh_password_wo_version` |
| `hpe_morpheus_load_balancer_monitor` | `monitor_password_wo` | `monitor_password_wo_version` |
| `hpe_morpheus_network_pool_server` | `service_password_wo` | `service_password_wo_version` |
| `hpe_morpheus_network_router_bgp_neighbor` | `password_wo` | `password_wo_version` |
| `hpe_morpheus_provisioning_license` | `license_key_wo` | `license_key_wo_version` |
| `hpe_morpheus_storage_bucket` | `access_key`, `secret_key` | `access_key_version`, `secret_key_version` |
| `hpe_morpheus_storage_server` | `service_password_wo` | `service_password_wo_version` |
| `hpe_morpheus_cloud` | `config_vmware.password`, `config_azure.client_secret` | `config_vmware.password_version`, `config_azure.client_secret_version` |
| `hpe_morpheus_cluster` | `server.ssh_password_wo` (create-only: a change **replaces** the cluster) | `server.ssh_password_wo_version` |
| `hpe_morpheus_load_balancer` | `group_id`, `network_server_id` (create-only: a change **replaces**) | `group_id_version`, `network_server_id_version` |
| `hpe_morpheus_instance_clone` | `source_instance_id` (create-only: a change **replaces**) | `source_instance_id_version` |
| `hpe_morpheus_subnet` | `config` (create-only: a change **replaces**) | `config_version` |
| `hpe_morpheus_storage_volume` | `config`, `config_alletramp_bmaas` (create-only: a change **replaces**) | `config_wo_version`, `config_alletramp_bmaas_wo_version` |
| `hpe_morpheus_network_router_nat` | `action`, `firewall`, `service` (updated in place) | `action_version`, `firewall_version`, `service_version` |
| `hpe_morpheus_load_balancer_pool` | `partition` (updated in place) | — (never had one) |
| `hpe_morpheus_instance_node` | `ssh_password` (create-only: a change **replaces**) | — (never had one) |

Network domain and shell/PowerShell task resources also use automatic change detection for
`domain_password_wo` and `remote_target_password_wo`. Their move to the plugin framework includes
automatic state upgrades. A task still using `remote_target_password` may plan one update.

Appliance-side secret changes are detected where Morpheus supplies a fingerprint: images, load
balancer monitors, BGP neighbours, storage servers, network pool servers, VMware/Azure cloud
credentials, instance nodes and clusters. Detection starts after an apply sets the secret. Terraform
restores the configured value, or warns if doing so would require replacement. User, storage bucket,
backup host and provisioning licence secrets remain configuration-only.

### New fields and lookup options

- **Instances and clones:** `config_hvm.user_data` supplies cloud-init data at first boot. Map keys
  merge with Morpheus-generated data; lists replace it. Changing the value replaces the instance.
  Instance `evars` also support `export` and `masked`.
- **HVM clusters:** added `virt_switch_name`, `dns1`–`dns3`, `search_domain`, `ntp1`, `ntp2`, `proxy`
  and `no_proxy` under `config_hvm`. These require cluster layout version 2.0 or later; changes
  replace the cluster. Documentation covers layout upgrades and virtual-switch requirements for
  scale-out. Set `overlay_interface_name` for HVM 1.x layouts on Morpheus 9.0.x to avoid host-build
  failures; HVM 2.0 does not use it.
- **Users:** `require_password_change` forces a password change at the next UI sign-in whenever the
  provider sets the password. API sign-in is unaffected.
- **Clouds:** `backup_provider_ids` assigns backup integrations. Before Morpheus 9.0.2, the
  configured value is retained because the API does not return it.
- **Network pool servers:** added `app_id`, required for phpIPAM and rejected for other types.
- **Appliance settings:** added `no_proxy` for the comma-separated proxy bypass list.
- **Networks:** added `network_server_id`, `scope_id`, `vxlan_id`, `parent_network_id` and
  `availability_zone`. Morpheus 9.1 HVM port groups require `network_server_id`.
- **Lookups:** network data sources accept `cloud_id` and `cluster_id`; datastore data sources
  accept `cluster_id`; instance-type layouts return `instance_type` (`id`, `name`, `code`).
- **Read-only fields:** affinity groups return `host_group`, `vm_group` and `owner`; compute-server
  data sources return `guest_agent_status`; snapshots return `max_storage` and their data source's
  `snapshot_files` include `file_size`; user data sources return `tenant_id`.
- **Storage-volume lookups:** added `include_tenants` and `tenant_id` filters.
- **Whitelabel settings:** added `security_banner`, `external_legal_text`, `supports_themes` and
  read-only `account`. Tenant `currency` is validated against ISO 4217 at plan time.

### Clearing optional values

Removing a previously configured optional attribute now clears it where the API supports this.
Tracking starts with the first apply under this version; apply once before removing attributes.
Booleans and values you never configured are retained. Required, immutable and non-clearable fields
have resource-specific rules; see the provider documentation under *Removing an attribute from the
configuration*.

## Resolved issues

### Provisioning and updates

- Instances and clones no longer resize on unrelated updates. Instance `timeouts`-only changes
  make no API call; `labels`-only and `plan_id`-only changes are applied correctly.
- ISO CD/DVD drives no longer cause instance volume drift.
- Networks apply supported rename and IPv6 changes. Create and destroy wait for completion.
- Instance and clone `Invalid network` errors list selectable networks.
- Service plans no longer send zero for unset `max_disks`, `max_cpu`, `max_cores` or `sort_order`.
  Repair affected existing plans by setting `max_disks` explicitly once.
- Appliance-settings updates no longer fail with `Not found in response: ApplianceSettings` and
  now send changes to `api_allowed_origins`, `registration_enabled` and both default role IDs.
- Security-group rules validate required source, destination and port fields at plan time. IDs of
  1,000,000 or more work correctly.
- Users can clear `linux_username` and `windows_username`.
- Storage-volume `max_storage` changes resize the volume; a name change replaces it.
- Storage servers, load balancer pools and instance nodes send their previously omitted
  `service_password_wo`, `partition` and `ssh_password` values. Storage buckets and network pool
  servers send credentials only when changed.
- Image updates with `tenant_ids` no longer crash the provider.
- Failed creates report the original error without additional unknown-state errors.
- Corrected OVS, Azure and Amazon network examples.

### State and deletion

- Fixed unnecessary replacements caused by unset create-only attributes on images, load balancer
  monitors, network routes, firewall rule groups, OS type images, policies and storage volumes.
- Image, datastore, policy and task IDs remain known during updates, avoiding replacements of
  dependent resources.
- Network `tenant_ids` updates no longer force replacement. Invalid public visibility is rejected
  at plan time. Monitor password-version changes no longer replace load balancer monitors.
- Refused deletes now report the API error and retain the resource in state.
- Cluster destroy waits for removal, preventing name conflicts on recreation
  (`timeouts.delete` defaults to 45 minutes).
- Instance destroy uses the instance API to remove the VM on Morpheus 8.1.x, not just its records.
  A warning is issued if deletion falls back to the server API.
- Fixed refresh-on-404 handling across 25 framework resources, allowing Terraform to recreate
  objects deleted outside Terraform.

### Authentication, tenancy and lookups

- Fixed username/password authentication failures after token expiry during long-running operations.
- Sub-tenant role `multitenant`/`multitenant_locked` settings and policy `tenants` settings outside
  `Group` scope are rejected at plan time. Role `cloud_permissions` are preserved when omitted
  from the API response.
- Storage bucket, node type, cloud folder and Ansible Tower job template data sources report missing
  matches instead of returning empty state or crashing.
- Updated the bundled SDK for the 9.1.0 API, including omitted sub-objects and empty numeric values.
- `tfmigrator` reports newer releases. Set `HPE_IGNORE_VERSION_CHECK` to disable the check.

## Known issues

- Cluster namespace renames require replacement. HKS HVM cluster destroy may report an error after
  successful deletion.
- Instance layouts may require a specific number of volumes. Before Morpheus 8.0.11, put the root
  volume first.
- Retry authentication if Morpheus returns an intermittent HTTP 500.
- NFS datastore creation can remain in `provisioning` if the server or share is inaccessible.
  Deletion can also fail; remove dependent resources first.
- Sub-tenant role imports cannot recover cloud permissions hidden by the API.
