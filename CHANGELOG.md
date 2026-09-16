# v2.1.0 Release Notes

## Breaking changes

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

### `hpe_morpheus_tenant` and `hpe_morpheus_tenants` data sources reimplemented

Both tenant data sources have been rewritten on the plugin framework (previously
terraform-plugin-sdk/v2), completing the port begun with the `hpe_morpheus_tenant` resource in
v2.0.0. Attribute names are unchanged, but one type changes:

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
where it was previously optional and computed.  The Morpheus API rejects a nested workflow task that
does not reference an operational workflow, so the attribute never had a meaningful computed value —
omitting it produced a task the API would not accept.  Making it required surfaces the mistake at
`terraform plan` rather than as an apply-time API error.

Any configuration that already creates a working nested workflow task is unaffected, since a valid
task must always have supplied the workflow id.  A configuration that omitted `operational_workflow_id`
was already non-functional and must now set it explicitly.

## Enhancements to existing resources

### `hpe_morpheus_image` no longer misses images beyond the first page

The singular `hpe_morpheus_image` data source narrowed by name server-side, but that is a SQL `like`,
so a broad name could match more images than one response holds.  An exact match falling beyond the
last page fetched was reported as not found.  Every page is now walked.

`virtio_supported` also always read null, because the field was never populated.  It now reports the
image's value.

### Test image sweeper no longer misses images

The image sweeper requested no page size at all, so it saw only the first page and left any test
image beyond it behind.  Which images those were depended on what else existed at the time, since the
API sorts by name.

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

`hpe_morpheus_policies` now walks every page before applying its filters.  It previously fetched only
the first hundred policies, so on a busy appliance a newly created policy could fall beyond the first
page and never match.  The underlying policy config is also decoded leniently, so listing no longer
fails when the API returns a field such as `maxCores` as a number in one policy and a string in
another.

`hpe_morpheus_key_pair` now looks the key pair up by its `id` argument.  It previously read the
internal resource id, which is empty during a data source read, so an `id`-only lookup fell through to
the "cannot be read without name or id" path.

`hpe_morpheus_os_type_image` retries the lookup briefly (an exponential backoff over roughly eight
seconds) to tolerate the read-after-write staleness of an image created moments earlier.

**Behavior change.**  `hpe_morpheus_instance_type` and `hpe_morpheus_storage_volume_type` now return
an error when the requested instance type or storage volume type does not exist, instead of silently
returning empty state.  A data source is expected to describe something that exists; the previous
silent-empty result left downstream references reading zero values.  Configurations that relied on the
old behavior will now surface an error.

## Resolved issues

### `hpe_morpheus_budget` can be scoped to a specific group, cloud, or user

The budget resource exposed a `scope` but no way to point a non-account scope at a particular
entity, so a `group`, `cloud`, or `user` budget could not be expressed.  A single
`associated_resource_id` now carries that target, and the scope/id pairing is validated at plan
time: it is required when `scope` is `group`, `cloud`, or `user`, and rejected when `scope` is
`account` (which targets the whole tenant).  Omitting `scope` is treated as the `account` default,
so setting `associated_resource_id` without a scope is caught during planning rather than failing
during apply.

# v2.0.0 Release Notes

This is a major release.  Alongside the Morpheus support this provider already offered, it adds
support for **HPE OpsRamp**, introduces PCE (Private Cloud Enterprise) Identity authentication for
Morpheus, and closes the remaining hpegl VMaaS parity gaps.

Release notes for earlier versions are in [HISTORY.md](./HISTORY.md).

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

### hpe_morpheus_tenant ported to the plugin framework

`hpe_morpheus_tenant` has been reimplemented on the Terraform plugin framework (previously
terraform-plugin-sdk/v2). The cutover is automatic: existing state is upgraded in place on the next
`terraform plan`/`apply`, with no manual `state rm` or re-import required.

The resource `id` changes type from **string to number**. State is migrated automatically by a schema
upgrade, but any configuration that consumed `hpe_morpheus_tenant.<name>.id` as a string (for example
in string interpolation) may need adjusting.

New capabilities: `parent_id` creates a tenant under a nominated parent when authenticated as
the master tenant (changing it forces replacement); `remove_resources` de-provisions the tenant's
managed instances on destroy when set to `true` (default `false`); `destroy` now waits for the
asynchronous tenant deletion to complete. Additional read-only attributes are now populated:
`master`, `parent_name`, `parent_subdomain`, `external_id`, `base_role_name`, `instance_count`,
`user_count`, `date_created`, and `last_updated`. `currency` is validated at plan time against the
appliance's live currency list, and `subdomain` is validated for format and the not-all-numeric rule.

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

`description` can now be cleared by removing it from the configuration. The SDKv2 resource carried
the previous description forward when it was omitted, so it could never be unset.

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
  Both attributes continue to work and are mutually exclusive (MORPH-15552).
- `hpe_morpheus_task_powershell_script` and `hpe_morpheus_task_shell_script` —
  `remote_target_password` is deprecated in favour of the write-only `remote_target_password_wo`
  with `remote_target_password_wo_version`.  Morpheus returns the password as a hash, so the
  plaintext in configuration never matched the hash in state and every plan was non-empty
  (MORPH-14024).
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
details can be obtained from GreenLake rather than configured by hand (MORPH-15611):

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
  Morpheus returns for a container that has not reported yet.  On expiry it warns and continues
  (MORPH-12804).  Added `config_vmware.affinity_group_id` and `config_hvm.affinity_group_id` to
  place an instance into an affinity group at provision time; create-only, and rejected alongside
  `config_hvm.kvm_host_id` (MORPH-15596).  Added the computed `compute_servers` and `container_id`,
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
  `volumes` as `(known after apply)` (MORPH-14919).  An imported instance no longer plans changes
  nobody made, which on appliances before 8.1.2 escalated to a replacement of a running VM.
  Instances created by the hpegl provider no longer fail to read: string-encoded `noAgent` is
  handled, and an absent `nestedVirtualization` is treated as optional rather than an error.
- `hpe_morpheus_subnet` — `resource_permission_groups_all` was sent under a request key the Morpheus
  API does not read, so the setting was silently dropped; an explicitly configured `pool_id` was
  overwritten with `null` when the API response omitted the pool (MORPH-14001).
- `hpe_morpheus_network_domain` — `public_zone`, `visibility` and `active` are now sent on update,
  so changing them takes effect; `auto_join_domain` is preserved on import and `tenant_id` is read
  back (MORPH-8836/MORPH-10305).
- `hpe_morpheus_network_router_firewall_rule` — creating a rule without a `description` no longer
  fails; the required format of `parent_id` is documented.
- `hpe_morpheus_network_router` — BGP neighbor configuration is read correctly on import, and API
  flags returned as JSON booleans are handled alongside the `on`/`off` strings.
- `hpe_morpheus_os_type_image` — inconsistent `os_type_id` after apply (MORPH-13276).
- `hpe_morpheus_tenant` — `currency` is validated against the supported ISO codes (MORPH-10304).
- `hpe_morpheus_option_type` — rows and description are validated at plan time
  (MORPH-7445/MORPH-8853).
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
  Unknown is now treated as absent, as null already was (MORPH-16244).
- `hpe_morpheus_image` — reading an image with two or more tenants failed with `Duplicate Set
  Element`, whether or not anything was duplicated.  The tenant objects were built without a known
  state, so their `name` and `id` were discarded and every tenant became identical (MORPH-16245).

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
