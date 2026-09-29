# v2.1.0 Release Notes

This release hardens the provider for day-two use. It completes the move of the tenant resource and
data sources to the plugin framework and reimplements `hpe_morpheus_images`; makes importing and
refreshing resources created outside Terraform reliable across the provider; adds support for
running as a Morpheus sub-tenant, with a `hpe_morpheus_whoami` data source to discover the caller's
tenant; and resolves a large number of state-consistency and lookup defects. OpsRamp gains
configuration-based (VMware) integrations and a management profile data source, with two breaking
schema changes listed in the [OpsRamp changes](#opsramp-changes) section.

Release notes for earlier versions are in [HISTORY.md](./HISTORY.md).

## Breaking changes

OpsRamp breaking changes are listed in the [OpsRamp changes](#opsramp-changes) section below.

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
[Running as a sub-tenant](#running-as-a-sub-tenant).

## OpsRamp changes

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

## Running as a sub-tenant

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
