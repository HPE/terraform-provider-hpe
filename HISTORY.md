# Release History

Release notes for versions of this provider prior to the current release.
The current release is documented in [CHANGELOG.md](./CHANGELOG.md).

# v2.1.0 Release Notes

This release hardens the provider for day-two use. It completes the move of the tenant resource and
data sources to the plugin framework and reimplements `hpe_morpheus_images`; makes importing and
refreshing resources created outside Terraform reliable across the provider; adds support for
running as a Morpheus sub-tenant, with a `hpe_morpheus_whoami` data source to discover the caller's
tenant; and resolves a large number of state-consistency and lookup defects. OpsRamp gains
configuration-based (VMware) integrations and a management profile data source, with two breaking
schema changes listed in the [OpsRamp changes](#opsramp-changes-v210) section.

## Breaking changes

OpsRamp breaking changes are listed in the [OpsRamp changes](#opsramp-changes-v210) section below.

### `hpe_morpheus_images` reimplemented

The `hpe_morpheus_images` data source has been rewritten on the plugin framework.  It previously
returned only a list of ids and filtered a single unpaginated response, so it could neither report
anything about an image nor see one that fell beyond the first page.

It now returns the full record for every matching image, walks every page, and offers the endpoint's
complete filter set.  Existing configurations need updating:

| Before | Now |
|---|---|
| `ids` — a list of image ids | `images` — a set of full image objects |
| `source` | `filter_type` |
| `sort_ascending` | `sort` and `direction` |
| `filter` blocks | unchanged, including their regular-expression semantics |

Reading an id now means projecting it out of the set.  `images` is a set rather than a list, so it
cannot be indexed; sort the ids where one in particular is wanted:

```hcl
data "hpe_morpheus_images" "ubuntu" {
  image_type = ["qcow2", "raw"]

  filter {
    name   = "name"
    values = ["^ubuntu"]
  }
}

output "image_ids" {
  value = sort([for image in data.hpe_morpheus_images.ubuntu.images : image.id])
}
```

A set is used deliberately.  The API pages by name, names are not unique across clouds, and rows can
therefore move between pages — so no stable order can be promised, and offering a list would imply
one.  Results are de-duplicated by id for the same reason.

`filter_type` defaults to `All`, which differs from the API's own default of `User`.  Left to itself
the API returns only user-uploaded images, hiding the synced and system images most instances are
provisioned from.  Set it to `User` to restore the previous scope.

-> **This data source can be slow on a large image library.**  Every image the server-side filters
admit is downloaded before `filter` blocks are applied, and Terraform reads a data source on refresh,
plan and apply alike.  Narrow it with `image_type`, `image_id` or `phrase` wherever possible; the
data source documentation has a table of which arguments reduce the download and which do not.

### `hpe_morpheus_tenant` ported to the plugin framework

`hpe_morpheus_tenant` has been reimplemented on the Terraform plugin framework (previously
terraform-plugin-sdk/v2). The cutover is automatic: existing state is upgraded in place on the next
`terraform plan`/`apply`, with no manual `state rm` or re-import required.

The resource `id` changes type from **string to number**. State is migrated automatically by a schema
upgrade, but any configuration that consumed `hpe_morpheus_tenant.<name>.id` as a string (for example
in string interpolation) may need adjusting.

**Behavior change.**  `description` is no longer carried forward when it is omitted from the
configuration. The SDKv2 resource kept the previous value, so a description could never be unset; it
can now be cleared by removing it. This also means a configuration that omits `description` for a
tenant that has one — an imported tenant, or one whose description was set in the Morpheus UI — will
plan to clear it on the next apply. Set `description` explicitly to keep it.

The port also brings new capabilities to the resource, described under
[Enhancements to existing resources](#hpe_morpheus_tenant-supports-the-tenant-hierarchy-and-reports-more).

### `hpe_morpheus_tenant` and `hpe_morpheus_tenants` data sources reimplemented

Both tenant data sources have been rewritten on the plugin framework (previously
terraform-plugin-sdk/v2), completing the port of the `hpe_morpheus_tenant` resource above.
Attribute names are unchanged, but one type changes:

The elements of `hpe_morpheus_tenants.ids` change type from **string to number**, matching the
tenant resource's `id`. Terraform converts numbers to strings implicitly in most contexts, so plain
interpolation keeps working, but contexts that require strings must now convert explicitly — most
commonly `for_each`:

```hcl
# Before
for_each = toset(data.hpe_morpheus_tenants.all.ids)

# Now
for_each = toset([for id in data.hpe_morpheus_tenants.all.ids : tostring(id)])
```

`ids` remains a list ordered by id, ascending unless `sort_ascending = false`, and `filter` blocks
are unchanged, including their regular-expression semantics.

Both data sources also report more, and fail honestly where they previously returned nothing:

- `hpe_morpheus_tenant` now populates the full tenant read model — `description`, `enabled`,
  `subdomain`, `currency`, `external_id`, `master`, `base_role_id`, `base_role_name`, `parent_id`,
  `parent_name`, `parent_subdomain`, `instance_count`, `user_count`, `date_created` and
  `last_updated` — alongside the existing `account_name`, `account_number` and `customer_number`.
  For the master tenant, which has no parent and no base role, the `parent_*` and `base_role_*`
  attributes are `null`.
- `hpe_morpheus_tenants` gains a `tenants` attribute: a list of full tenant objects with the same
  attributes as the singular data source, in the same order as `ids`. `filter` blocks may now also
  match on `subdomain`, `currency`, `external_id`, `enabled`, `master`, `account_name`,
  `account_number` and `customer_number` in addition to `name`.
- A lookup that matches no tenant is now an error naming the problem. The SDKv2 data sources
  returned success with empty attributes, so the failure only surfaced later wherever the empty
  value was consumed. Setting neither `id` nor `name`, or an empty `name`, is likewise rejected —
  at plan time.
- `hpe_morpheus_tenants` previously fetched at most 100 tenants and silently ignored the rest; the
  cap is now 10000.

### `hpe_morpheus_task_nested_workflow` requires `operational_workflow_id`

The `operational_workflow_id` attribute on `hpe_morpheus_task_nested_workflow` is now **required**,
where it was previously optional and computed.  The Morpheus API rejects a nested
workflow task that does not reference an operational workflow, so the attribute never had a
meaningful computed value — omitting it produced a task the API would not accept.  Making it required
surfaces the mistake at `terraform plan` rather than as an apply-time API error.

Any configuration that already creates a working nested workflow task is unaffected, since a valid
task must always have supplied the workflow id.  A configuration that omitted `operational_workflow_id`
was already non-functional and must now set it explicitly.

### `hpe_morpheus_backup_job` data source no longer exposes `enabled`

The `enabled` attribute has been removed from the `hpe_morpheus_backup_job` **data source**.
The backup jobs API does not return an `enabled` field, so the attribute was always
null and could never convey a job's real state — reading it was misleading.

Any configuration that referenced `data.hpe_morpheus_backup_job.<name>.enabled` must remove that
reference.  The `hpe_morpheus_backup_job` **resource** is unaffected and still accepts `enabled`.

## New data sources

### `hpe_morpheus_whoami`

A configuration often needs to know who it is running as — which tenant, whether that tenant is the
master, what the account may do — and until now the only way was to hardcode a tenant id or user
name that differs between appliances. `hpe_morpheus_whoami` takes no arguments and reports the user
the provider is authenticated as:

```hcl
data "hpe_morpheus_whoami" "current" {}
```

It returns the user's `id`, `username`, name and email fields, status flags, `roles` and
`permissions`, `default_persona`, the `tenant` (`id` and `name`) with `tenant_id` alongside,
`is_master_account`, and `appliance_build_version`. Credential fields and the categorised `access`
object are deliberately excluded; use the `hpe_morpheus_user` data source for the latter. It is the
basis for configurations that run unchanged as the master tenant or as a sub-tenant — see
[Running as a sub-tenant](#running-as-a-sub-tenant-v210).

## OpsRamp changes (v2.1.0)

This release adds the pieces needed to declare a VMware cloud in Morpheus and its monitoring
integration in OpsRamp from one configuration, and corrects the shape of two alerting resources.
Two of these changes are breaking.

### Breaking changes

#### `hpe_opsramp_metric_alert_definition`: `entity_type` and `component` are now strings

Both were lists of strings although OpsRamp takes a single value for each. They are now plain
strings:

```hcl
# Before
entity_type = ["RESOURCE"]
component   = ["$$__name__"]

# Now
entity_type = "RESOURCE"
component   = "$$__name__"
```

`attributes`, previously required, is now optional and required only when `entity_type` is
`RESOURCE`. `no_data_condition` applies to the `STATIC_THRESHOLD` and `DYNAMIC_THRESHOLD` threshold
types and defaults to `NO_DATA_ALERT` when omitted, where it was previously required for them.

#### `hpe_opsramp_first_response_policy`: action settings are now required

Within `attribute_actions`, `run_process.process_ids` and `suppress.suppress_duration` are now
required, as is `pattern_actions.seasonality_time_frame`, which accepts `7D`, `10D`, `30D`, `60D` or
`90D`. Configurations that omit them must now supply them. Descriptions of the `learned_configuration`
and `run_immediately` settings have been corrected.

### New data source: `hpe_opsramp_management_profile`

Looks a management profile up by `name`, optionally within a `client`, and returns its `id`, `uuid`,
`type` and `description`. Its `uuid` is what a discovery profile on an `hpe_opsramp_integration`
needs.

### `hpe_opsramp_integration` supports configuration-based integrations

The resource previously covered event-based and custom integrations. It now also installs
configuration-based ones such as `VMWARE`, which connect OpsRamp to a target system: `ip_address`
names the endpoint, `credential_set` the stored credentials, and `discovery_profiles` describes what
is discovered — each with a `mgmt_profile_uuid`, a `policy` of `rules` and `actions`, an optional
`schedule`, and `scan_now` to discover immediately. Combined with `hpe_morpheus_cloud`, a VMware
cloud and its OpsRamp integration can be declared together:

```hcl
resource "hpe_opsramp_integration" "vmware" {
  display_name   = "VMware Integration"
  application    = "VMWARE"
  ip_address     = var.vcenter_url
  credential_set = hpe_opsramp_credential_set.vcenter.id

  discovery_profiles = [{
    mgmt_profile_uuid = data.hpe_opsramp_management_profile.gateway.uuid
    scan_now          = true
    policy = {
      entity_type = "ALL"
      match_type  = "ANY"
      rules       = [{ filter_type = "ANY_CLOUD_RESOURCE", resource_type = [] }]
      actions     = [{ action = "MANAGE DEVICE" }]
    }
  }]
}
```

## Running as a sub-tenant (v2.1.0)

Configurations run with sub-tenant credentials hit a class of silent failures: Morpheus accepts the
request, discards or coerces the part a sub-tenant may not set, and reports success, so the provider
read back a value it had never asked for and planned the same change on every run. Several of these
are now caught at plan time, and the appliance version check no longer depends on a permission
sub-tenants rarely hold. Where the caller's tenancy cannot be determined, the plan is not blocked.

- `visibility = "public"` is rejected at plan time when the caller is not the master tenant, on
  `hpe_morpheus_cloud`, `hpe_morpheus_image`, `hpe_morpheus_network_group`,
  `hpe_morpheus_option_list`, `hpe_morpheus_catalog_item_workflow` and
  `hpe_morpheus_workflow_operational`. Morpheus stores `private` for a sub-tenant regardless, which
  previously produced a diff that never converged. The attribute descriptions now say so.
- `hpe_morpheus_setting_whitelabel.appliance_name` is rejected at plan time for a sub-tenant. It is a
  master-tenant setting that Morpheus silently discards, so the apply failed with an inconsistent
  result.
- The `hpe_morpheus_cloud_type` data source reads `/api/zone-types`, which needs no appliance-level
  permission, so it works for sub-tenant callers. Only enabled cloud types are returned.
- The appliance version is read from `/api/whoami` rather than `/api/health`. The health endpoint
  requires the `admin-health` permission, so with a token lacking it every version-dependent check —
  such as the Morpheus 8.1.0 requirement for `hpe_morpheus_tenant.parent_id` — was silently skipped.
  `whoami` reports the same build version to any authenticated caller.

## Enhancements to existing resources

### Finding the resource pool of an HVM cluster

Morpheus creates a resource pool for every HVM cluster and binds the cluster's networks to it.
That pool is the one to provision into, but it is attached to the cluster rather than to the
cloud, so the cloud's resource-pool listing does not include it.  `hpe_morpheus_resource_pool`
resolves `name` through that listing and so could not find it, reporting only
`found 0 resourcePools`.  A user who took that at face value and created a pool by hand would
then see instance creation fail with `Invalid network`, since the networks belong to the
cluster's pool.

- `hpe_morpheus_resource_pool` now explains this when a name lookup finds nothing, and, when a
  cluster of that name exists in the cloud, gives the id of the cluster's pool in the error.  A
  lookup by `id` that does not exist reports that plainly rather than as a raw HTTP error, and
  `type` is now documented with the values Morpheus uses (`namespace` is a cluster's pool).
- `hpe_morpheus_cluster` documents `permissions.resource_pool` as the pool to provision into,
  with an example feeding it to `config_hvm.resource_pool_id`.
- `hpe_morpheus_instance` and `hpe_morpheus_instance_clone` explain an `Invalid network`
  rejection instead of passing it through bare: the error now names each requested network,
  the resource pool it belongs to and the pool the instance asked for, and points at the
  cluster's pool when they differ.  The explanation is added only on that failure and never
  replaces the API's own message.  `config_hvm.resource_pool_id` is documented as the cluster's
  pool, and the shared `resource_pool_id` description no longer says "resource group".

### `hpe_morpheus_tenant` supports the tenant hierarchy and reports more

The port of `hpe_morpheus_tenant` to the plugin framework (see
[Breaking changes](#hpe_morpheus_tenant-ported-to-the-plugin-framework)) also adds new capabilities.
`parent_id` creates a tenant under a nominated parent when authenticated as the master tenant
(changing it forces replacement); `remove_resources` de-provisions the tenant's managed instances on
destroy when set to `true` (default `false`); `destroy` now waits for the asynchronous tenant
deletion to complete. Additional read-only attributes are now populated: `master`, `parent_name`,
`parent_subdomain`, `external_id`, `base_role_name`, `instance_count`, `user_count`, `date_created`,
and `last_updated`. `currency` is validated at plan time against the appliance's live currency list,
and `subdomain` is validated for format and the not-all-numeric rule.

`parent_id` requires Morpheus 8.1.0 or later, which introduced the tenant hierarchy. On earlier
appliances — which silently ignore a nominated parent — the provider now refuses a configuration
that sets `parent_id` at plan time with a diagnostic naming the required version, rather than
applying it and then forcing a replacement on every subsequent plan. On those appliances
`parent_id`, `parent_name` and `parent_subdomain` are `null`. All other tenant functionality is
unchanged across supported Morpheus versions.

Tenant validation failures that older appliances report as HTTP 200 with `success: false` (for
example an invalid `base_role_id`) are now surfaced with the API's own message instead of a
generic "Account ID is nil" error.

The master tenant may now be updated when authenticated as the master tenant, matching what
Morpheus allows; only disabling it (`enabled = false`), assigning it a `base_role_id`, and
destroying it are refused, each with a diagnostic naming the restriction.

### `hpe_morpheus_instance` can provision from a nominated image

`config_hvm` and `config_vmware` gain `image_id`, which overrides the image configured on the
instance type layout, so one layout can serve instances that need different images.
It is create-only: changing it replaces the instance. The `hpe_morpheus_images` data source is how
the id is found.

### `hpe_morpheus_instance` reports the volume size Morpheus actually provisioned

Morpheus does not always create the disk that was asked for: when an image needs more room than
the request allows, the volume is rounded up, so asking for 10 GB with an image a little over 10 GB
yields an 11 GB disk. The provider overwrote the API's size with the requested one on every read,
so state recorded the request and the real disk was invisible to plan and refresh alike — except on
import, which read the truth and then disagreed with the state an apply had produced.

Each entry in `volumes` gains a computed `actual_size`, the size in GB that exists. `size` keeps
its meaning as the request and is now also computed, so a volume that has grown does not plan a
shrink on every run. A difference between the two is normal and produces no plan diff; a configured
`size` below the provisioned size has no effect, as Morpheus cannot shrink a disk in place; and a
disk grown outside Terraform shows up in `actual_size` on the next refresh.

### `hpe_morpheus_instance` destroy detects a failed removal and reports the reason

After deleting an instance the provider waits for it to disappear.  It treated `stopped` and
`suspended` as failures, but Morpheus writes both onto an instance that is being removed while its
servers are stopped, so they are transient during a normal teardown; on older provider versions
this occasionally failed a destroy with `reached error status: stopped`.  Meanwhile `warning`, the
status Morpheus actually sets when a removal fails, was not recognised, so a real failure was only
reported as a timeout after 45 minutes.

The wait now ignores `stopped` and `suspended`, stops on `warning`, and includes Morpheus's own
reason in the error, for example
`instance 119675: DELETE failed reached error status: warning (Unable to remove instance: ...)`.
The instance id in these messages, and in the `hpe_morpheus_image` and `hpe_morpheus_task` destroy
messages, previously printed as `{2 119675}`; it now prints as the number.

### `hpe_morpheus_security_group` group permissions are configurable

Restricting a security group to particular groups — `resource_permission_groups_all = false`
together with `resource_permission_group_ids` — is the documented usage and what the shipped
example does, yet the provider rejected the combination at plan time. Removing that
check exposed two further defects it had been hiding: the API's create endpoint ignores group
permissions, so a group created with `all = false` read back as `all = true`; and every update sent
an empty tenant-permission list when `tenant_ids` was unset, which the API treats as an instruction
to remove every permission row — including the one carrying the group permissions.

All three are fixed. The combination is accepted; group permissions are applied by a follow-up
update immediately after create; and tenant permissions are sent only when `tenant_ids` is
configured.

### `hpe_morpheus_image` no longer misses images beyond the first page

The singular `hpe_morpheus_image` data source narrowed by name server-side, but that is a SQL `like`,
so a broad name could match more images than one response holds.  An exact match falling beyond the
last page fetched was reported as not found.  Every page is now walked.

`virtio_supported` also always read null, because the field was never populated.  It now reports the
image's value.

The image sweeper used by the acceptance tests had the same defect: it requested no page size at
all, so it saw only the first page and left any test image beyond it behind.

### `hpe_morpheus_instance` no longer fails to read on an unexpected boolean value

Some instance config fields are booleans that the API almost always returns as genuine JSON booleans,
but occasionally as a string.  An unexpected string in one of these fields — for example `createUser`
as `"yes"` — previously failed the decode of the entire `GET /api/instances/{id}` response, not just
that one field, so the instance became unreadable: read, refresh and import all failed with a null
`id` and `name`.

The SDK's boolean decode now tolerates such values.  The common spellings (`yes`/`no`, `on`/`off`,
`enabled`/`disabled`, and the like) are interpreted, and any other string resolves to `false` rather
than losing the whole response.  This applies to every boolean field across the SDK, not only the
instance.  No configuration change is required.

### Morpheus data source lookups no longer fail silently or truncate results

`hpe_morpheus_policies` now walks every page before applying its filters.  It
previously fetched only the first hundred policies, so on a busy appliance a newly created policy
could fall beyond the first page and never match.  The underlying policy config is also decoded
leniently, so listing no longer fails when the API returns a field such as `maxCores` as a number in
one policy and a string in another.

`hpe_morpheus_key_pair` now looks the key pair up by its `id` argument.  It previously
read the internal resource id, which is empty during a data source read, so an `id`-only lookup fell
through to the "cannot be read without name or id" path.

`hpe_morpheus_os_type_image` retries the lookup briefly (an exponential backoff over roughly eight
seconds) to tolerate the read-after-write staleness of an image created moments earlier.

**Behavior change.**  `hpe_morpheus_instance_type` and `hpe_morpheus_storage_volume_type` now return
an error when the requested instance type or storage volume type does not exist, instead of silently
returning empty state.  A data source is expected to describe something
that exists; the previous silent-empty result left downstream references reading zero values.
Configurations that relied on the old behavior will now surface an error.

### Data sources report more

- `hpe_morpheus_cloud` — `config`, the cloud's configuration object as returned by the API. Its
  contents vary by cloud type and include values Morpheus discovers from the target system rather
  than ones supplied at creation; it is null for a cloud with no config.
- `hpe_morpheus_instance` — `storage_profile` and `create_for_multi_attach` on each volume, and
  `subnet` (`id`, `name`) on each container interface.
- `hpe_morpheus_network_router_nat` — `firewall` and `service`.
- `hpe_morpheus_network_firewall_rule_group` — `tenants` (`id`, `name`) and `visibility`.

### Plan-time validation

Several mistakes that used to surface as an apply-time API error — or as a misleading one — are now
caught at `terraform plan`:

- Data sources that look an object up by a free-text key reject an empty key. `name = ""` used to
  flow into the by-name search and come back as "not found", indistinguishable from a legitimate
  miss. The key is `name` on most data sources, `vip_name` on `hpe_morpheus_load_balancer_virtual_server`
  and `ip_address` on `hpe_morpheus_network_router_bgp_neighbor`.
- `hpe_morpheus_os_type_image` — `os_type_id` and `virtual_image_id` must be positive.
- `hpe_morpheus_network_dhcp_server` — `lease_time` must be at least 1.
- `hpe_morpheus_option_list` — `api_type` is required when `type` is `api`, and `source_url` when
  `type` is `rest` or unset.
- `hpe_morpheus_task_ansible_playbook` — `playbook` must not be empty.
- `hpe_morpheus_user_group` — `description` is limited to 255 characters; a longer value was a
  server error.
- `hpe_morpheus_image` — `min_ram` and `min_disk` cannot be negative.
- `hpe_morpheus_cluster_affinity_group` — an empty `name` is rejected instead of being sent to the
  API and reported as a 403.
- `hpe_morpheus_setting_provisioning` — `cloudinit_username` must not be empty.

## Resolved issues

### Importing and refreshing resources created outside Terraform

A resource's read runs on refresh and on `terraform import`. A Terraform-created resource
round-trips its own values, but one created in the UI or through the API and then imported exposes
every field the API omits — and several read paths handled an omitted optional field badly:
writing `null` where the schema declares a default, so the import planned an update that never
settled (or, where the attribute forces replacement, a destroy and recreate of a running VM);
dereferencing a nil pointer and crashing the provider; or failing the whole read. All three are
fixed.

- Every resource that declares a schema default fills it into state after read, exactly as the plan
  would, so an imported resource plans no change nobody made.
- `hpe_morpheus_network_router_route` (`description`, `network_mtu`), `hpe_morpheus_policy`
  (`motd.title`), `hpe_morpheus_monitoring_check` (`check_interval`), `hpe_morpheus_load_balancer`
  (tenants, and a load balancer whose optional cloud is omitted) and `hpe_morpheus_task`
  (`task_options`, `retry_delay_seconds`) no longer panic or error when the field is absent.
- `hpe_morpheus_instance` — an imported instance resolves `public_ip_type`, `is_ec2`, `layout_size`,
  `create_user` and `no_agent` to their defaults when the API omits them; previously the first
  produced a permanent diff, `is_ec2` crashed the provider, and `create_user` and `no_agent` failed
  the read outright. On appliances before 8.1.2 the resulting `network_interfaces` diff escalated
  to a replacement. An instance with no `connectionInfo` — stopped, failed or not yet
  provisioned — maps it to null instead of failing the read.
- `hpe_morpheus_cluster_namespace` — `active` is looked up on import. The single-namespace endpoint
  does not return it, so an inactive namespace imported as active.
- `hpe_morpheus_network_router_nat` — an unset `protocol` reads back as null rather than an empty
  string, so import and create agree.
- Read gaps that made an imported resource differ from its configuration:
  `hpe_morpheus_option_list_rest` did not read back `ignore_ssl_errors`,
  `inject_system_authorization_header` or `use_owner_auth`;
  `hpe_morpheus_catalog_item_workflow` and `hpe_morpheus_catalog_item_app_blueprint` did not read
  back `visibility`; `hpe_morpheus_cluster_layout` did not reconstruct its
  master and worker node pools.

### False drift and inconsistent results after apply

- `hpe_morpheus_option_list` — an unset `type` no longer plans a change on every apply; it is
  computed with the API's default of `rest`.
- `hpe_morpheus_backup_job` — `code` cannot be changed after creation, so changing it now forces
  replacement instead of silently doing nothing.
- `hpe_morpheus_backup_host` and `hpe_morpheus_backup_instance` — a configured `storage_provider_id`
  is preserved when the API omits it from the response, instead of failing with an inconsistent
  result.
- `hpe_morpheus_container_script` — a masked global script body no longer causes a perpetual diff.
- `hpe_morpheus_user` — roles Morpheus assigns automatically no longer appear as drift on
  `role_ids`; the API's list is read only on import.
- `hpe_morpheus_network_router_firewall_rule` — `protocol` and `port_range` are accepted by the API
  but not reliably returned, so the configured value is preserved and a warning is raised when the
  router did not report it. They cannot be used for drift detection; on NSX-T routers the effective
  service is selected through `application`.
- `hpe_morpheus_integration_docker_registry` — `password` is treated as write-only. The API returns
  it only masked, so comparing it produced a diff on every plan.
- `hpe_morpheus_resource_pool_group` — `tenant_ids` was sent under a key the API does not read, so
  tenants never persisted and every refresh cleared them from state.
- `hpe_morpheus_instance_type_layout` — `spec_template_ids` is sent for every layout technology. It
  was sent only for ARM, CloudFormation and Terraform layouts, so a VMware layout lost its spec
  templates and read reported them gone.
- `hpe_morpheus_cloud` — a cloud configured with `config_vmware` could not be read: the read failed
  with `failed to decode VMware configuration`, so the apply that created the cloud failed after the
  cloud existed on the appliance, and it could be neither refreshed nor imported. `cloud_type_code`
  is also now reported for clouds configured through a typed `config_*` block, not only through the
  generic `config`.

### `hpe_morpheus_setting_provisioning`

Three defects in the same create and update payload. Every
update failed with `Not found in response: ProvisioningSettings`, because update asserted on a
response envelope the endpoint never returns; create did not, which is why create worked and update
never did. `show_console_keyboard_settings` was declared and read back but never sent — and the key
the API reads on write differs from the one it returns on read, so it is now sent under the name
the API accepts. And `cloudinit_password`, `windows_password` and `pxe_root_password` were sent as
empty strings when unset, which the appliance treats as an instruction to clear the stored password;
they are now omitted unless configured.

### `hpe_morpheus_price` create and re-create

`POST /api/prices` reports a validation failure — `code: must be unique`, for example — as HTTP 200
with `success: false`. The resource treated any 200 as success, stored an id of 0, and the follow-up
read dropped the price from state, which Terraform reported as `Root object was present, but now
absent`. Such failures are now an error carrying the API's message and per-field
errors.

The usual way to hit that failure is re-creating a price with a code used before. Destroying a price
deactivates it rather than deleting it, so its `code` stays occupied, and a later create with the
same code either failed as above, inserted a duplicate row, or — where the appliance enforces the
code's uniqueness — returned a 500. Create now looks the code up first,
deactivated prices included, within the same tenant scope: a deactivated match is re-activated in
place, keeping its id and history, and becomes the resource; an active match is never adopted, and
the error says how to `terraform import` it; no match creates the price as before. `tenant_id` is
sent only when configured; it was always sent, as `0`, when omitted.

**Behavior change.** Creating a price whose `code` matches a deactivated price in the same tenant
scope re-activates that price instead of inserting a duplicate. Where several deactivated prices
share the code, the appliance chooses which is re-activated.

### `hpe_morpheus_budget` can be scoped to a specific group, cloud, or user

The budget resource exposed a `scope` but no way to point a non-account scope at a particular
entity, so a `group`, `cloud`, or `user` budget could not be expressed.  A single
`associated_resource_id` now carries that target, and the scope/id pairing is validated at plan
time: it is required when `scope` is `group`, `cloud`, or `user`, and rejected when `scope` is
`account` (which targets the whole tenant).  Omitting `scope` is treated as the `account` default,
so setting `associated_resource_id` without a scope is caught during planning rather than failing
during apply.

### `hpe_morpheus_storage_volume_type` data source exposes more attributes

The `hpe_morpheus_storage_volume_type` data source previously surfaced only `id`, `name`, `code`, and
`category`.  It now also exposes eight further scalar attributes the API returns:
`description`, `enabled`, `default_type`, `has_datastore`, `configurable_iops`, `custom_size`,
`custom_label`, and `display_order`.  The change is additive; existing configurations are unaffected.

### Other fixes

- `hpe_morpheus_cluster` data source — `name` matches case-insensitively, in line with the
  server-side query.
- `hpe_morpheus_cluster_affinity_group_member` — a parent cluster or affinity group that does not
  exist is reported as `cluster or affinity group not found` rather than as `Missing Resource State
  After Create`.
- Policies with an empty `config` — a workflow policy with no config block, for example — no longer
  fail with `unexpected end of JSON input` when created or listed, and an empty config is no longer
  mistaken for a max-memory policy.
- `hpe_morpheus_image` data source — an error while decoding the looked-up image is now reported
  instead of being discarded.
- `hpe_morpheus_job_task` — the `context_type` description explains when `appliance` is valid.

## Known issues

- `hpe_morpheus_cluster_namespace`: `name` update is not supported; changing it forces replacement.
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
