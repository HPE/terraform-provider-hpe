# v2.2.0 Release Notes

## Breaking changes

### `hpe_morpheus_network` — typed, write-only configuration (MORPH-17744)

A Morpheus network's behaviour is almost entirely decided by its type, and the type's configuration
fields used to be a free-form `config` map written through to state and never read back: a typo was
accepted silently, a value Morpheus would not apply after create was accepted silently, an import
left the map empty, and nothing in the schema said what any type needed. The resource is rebuilt
around what each network type declares.

- **Typed configuration blocks.** The configuration-map fields of the common types are typed,
  documented attributes: `config_ovs_port_group`, `config_hvm_data`, `config_hvm_overlay`,
  `config_hvm_private`, `config_hvm_sriov`, `config_nsxt`, `config_nsx_cloud`, `config_azure`,
  `config_gcp` and `config_openstack` (OpenStack, Huawei and Open Telekom private networks). A
  network takes at most one block, and it must be the block for its type. The free-form `config`
  remains for types without a block, such as those a network plugin adds.
- **Write-only.** The blocks and `config` are write-only: sent to Morpheus on create and update and
  never stored in state, so they can carry values Morpheus masks. What Morpheus holds is readable in
  the new computed `config_from_api`. A change is detected from a secure hash the provider keeps in
  the resource's private state — HashiCorp's recommended pattern for write-only arguments — so no
  version or trigger attribute is needed: a change to a field Morpheus edits in place plans an
  update, a change to a field it does not edit after create plans a replacement. Write-only
  attributes need **Terraform 1.11 or later, or OpenTofu 1.11 or later**.
- **`config` leaves state.** A reference to `hpe_morpheus_network.x.config` in a configuration must
  become `hpe_morpheus_network.x.config_from_api`. Existing state is upgraded automatically and the
  first plan after upgrading shows no change to the network.
- **`type_code`.** The type may be named by code (`ovsPortGroup`, `nsxtLogicalSwitch`,
  `azureNetwork`, …), which is the same on every appliance, or by `type_id`; one is required and the
  other computed. `type_id` is no longer required on its own.
- **Replacement follows the type.** Every Morpheus network type declares which of `name`, `cidr`,
  `gateway`, `vlan_id`, `dhcp_server` and the DNS attributes it edits after create; a change to any
  other is accepted by the API and ignored. The provider reads the type from the appliance at plan
  time and plans a replacement for exactly those changes, in place of the previous rule that
  `name`, `cidr` and `cidr_ipv6` always replaced the network. A built-in table of a recent Morpheus
  release's types stands in when the appliance cannot be reached during planning.
- **New attributes.** `network_server_id` (the network integration a plugin-managed type requires;
  Morpheus 9.1 requires it on the HVM port group types), `scope_id` (NSX-T transport zone),
  `vxlan_id` (OVS overlay), `parent_network_id` (custom child network) and `availability_zone`
  (Amazon subnet). Morpheus never returns the first four, so the configured values are kept in state.
- **Deprecated.** The top-level `connected_gateway` is superseded by `connected_gateway` inside
  `config_nsxt` and `config_nsx_cloud`; it still works.
- **Timing.** Create waits for the network to finish provisioning — on a cloud, or on an HVM
  cluster's hosts — before reading it back, so a failed provisioning is an error at apply rather than
  drift later. Destroy waits until the network is gone.
- **Appliance capability check.** A block attribute the appliance's own definition of the network
  type does not declare — the HVM virtual-switch layout fields on a release before 9.1, say — is
  rejected at plan time, naming the field and the appliance's version, instead of being accepted by
  Morpheus and silently ignored; an undeclared key of the free-form `config` is a warning.

The documentation describes each block, marks the attributes that replace the network, and gives a
validated example for the OVS port group, HVM data, NSX-T, OpenStack, Azure, Google, Amazon and
child network types.

## New data sources

### `hpe_morpheus_option_source_value` and `hpe_morpheus_option_source_values` (MORPH-17643)

Morpheus option sources are the enumerations behind the UI's dropdowns: values Morpheus discovers
from a cloud or computes for a provisioning context, such as the resource pools and networks the
provisioning wizard offers, an SCVMM library share, or the members of a select list. They were not
reachable from Terraform, so a configuration that needed one — an SCVMM cloud's `libraryShare`, or
the resource pool Morpheus creates for an HVM cluster — had to be written from a value obtained by
hand.

Two data sources now read any option source. `hpe_morpheus_option_source_values` returns the full
list, with a repeatable `filter` block on `name`, `value` or a source-specific attribute;
`hpe_morpheus_option_source_value` returns exactly one and fails the plan when the source returns
none or several, naming the candidates. Both take the source name, an optional `service_type` for
sources defined on a cloud-specific service, a `parameters` map, and typed shortcuts for the common
parameters (`cloud_id`, `layout_id`, `provision_type_code`, `pool_id`, `group_id`, `credential_id`,
...), named as this provider names things — the translation to the API's `zoneId` is the data
source's job. Values come back in the form the API takes on input (`pool-2`, `network-79`), with
every other key the source returned under `attributes`.

The two sources behind the provisioning wizard are the motivating case. `zonePools` lists the pools
selectable for a cloud and layout — including an HVM cluster's own pool, which the cloud's
resource-pool listing omits — with `providerType` alongside; `zoneNetworkOptions` lists the networks
selectable for a pool, which is the check Morpheus applies when an instance is created, so a
`lifecycle.precondition` can turn a runtime `Invalid network` into a plan-time message.

Three server behaviours are handled so they do not become silent mistakes: a parameter sent in the
wrong position is ignored rather than rejected, so each parameter is sent in both positions by
default; sources that need a provisioning context return an empty list rather than an error without
one, and the provider says so; and the server returns one generic error for an unknown source and
for bad parameters alike, so the provider carries a catalogue of the sources known for a Morpheus
release, warns on a name outside it, and says which cause is more likely. A fourth is documented and
warned about: the server evaluates every source for the caller's tenant and discards any `accountId`
sent, so there is no tenant attribute, and passing one in `parameters` warns rather than silently
answering for the wrong tenant. Raw credentials must not be passed in `parameters`; sources that
authenticate to a target system take a stored credential via `credential_id`.

## Improvements

### Write-only values: changes are detected automatically; the `*_version` markers are deprecated (MORPH-18005)

A write-only attribute such as `password_wo` is never stored in state, so until now the provider
could not see it change: a companion `*_version` marker had to be bumped by hand for a new value to
be sent, and a forgotten bump left the appliance silently out of date. The provider now keeps a
digest of each configured write-only value in the resource's private state and compares it on
every plan, so **changing the value is enough** — the plan shows the change as
`~ write_only_revision = <n> -> (known after apply)` and the apply sends the new value.

This release migrates the resources whose write-only attributes are secrets sent on update:

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

Nested attributes are named by their dotted path in the import ID (`"123#config_vmware.password=<sha256>"`);
a map or object value is digested as sha256 of its JSON form (`sha256(jsonencode(v))`), a number as its text.

For the create-only attributes, the deprecated marker used to be the *only* way to ask for a replacement
("increment to recreate"); now the change itself is detected and plans the replacement, bumping the marker
still does, and **removing** the marker from the configuration is an in-place update, never a replacement.
For a migrated resource:

- **`*_wo_version` is deprecated.** It still works for the whole 2.2 line — changing it re-sends the
  value as before — but every plan warns that it is no longer needed. Remove the attribute from
  your configuration: the removal is one in-place update, after which the plan is empty. The
  attributes will be removed in the 2.6 release.
- **A new read-only `write_only_revision`** records the time of the last apply that sent a
  write-only value. It is the visible trace of "a write-only value was re-sent on this apply".
- **Upgrading from 2.1:** the first plan after upgrading shows one in-place update per resource
  (`+ write_only_revision`), which re-sends the values the appliance already has and records their
  digests. This is a one-time cost; every later plan is clean. It cannot be avoided: the digest can
  only be taken when the configured value is next seen at apply.
- **Import:** the plan after `terraform import` (or an `import` block) is empty. Without more
  information the provider cannot know whether the configured value matches the appliance, so it is
  taken as applied and a change to it is only detected once the resource is next updated for another
  reason — Terraform warns about this on every plan until then. To have changes detected from the
  start, put the sha256 digest of each configured write-only value in the import ID:
  `"123#password_wo=<sha256 hex>&linux_password_wo=<sha256 hex>"` (Terraform's `sha256()` function
  produces it in an import block). Migration tooling that knows the values uses this form.
- **An unset write-only attribute is not managed.** Leaving `linux_password_wo` out of the
  configuration is not a change and never triggers a re-send.
- **`hpe_morpheus_cluster`:** Morpheus cannot change a cluster's SSH password after creation, so a
  change to `server.ssh_password_wo` now plans a replacement of the cluster instead of being silently
  ignored (before, neither the value nor its marker had any effect after create). Bumping the
  deprecated marker plans the same replacement; removing it does not.
- **`hpe_morpheus_network_router_nat`:** a change to `action`, `firewall` or `service` now plans an
  update (before, the plan was empty and the change only went out when something else changed).
- **Not migrated:** the three SDKv2 resources with a write-only attribute — `hpe_morpheus_network_domain`
  (`domain_password_wo`), `hpe_morpheus_task_shell_script` and `hpe_morpheus_task_powershell_script`
  (`remote_target_password_wo`) — cannot keep digests (the SDKv2 has no private state for provider
  code), so their `*_wo_version` markers remain the way to re-send the value and are **not** deprecated.
  Their descriptions say so. They will move to the digest mechanism when the resources move to the
  plugin framework.
- **Behaviour corrections found on the way:** `hpe_morpheus_load_balancer_monitor` no longer
  *replaces* the monitor when `monitor_password_wo_version` changes (the marker had `RequiresReplace`
  although the API updates the password in place); `hpe_morpheus_storage_server` now actually sends
  `service_password_wo` on update (it read the write-only value from the plan, where it is always
  null); `hpe_morpheus_storage_bucket` and `hpe_morpheus_network_pool_server` send their keys /
  password only when they changed rather than on every update; and an in-place update of
  `hpe_morpheus_image` with `tenant_ids` set no longer crashes the provider (a wrong set conversion
  behind a `log.Fatal`); `hpe_morpheus_load_balancer_pool` never sent `partition` and
  `hpe_morpheus_instance_node` never sent `ssh_password` (both read the write-only value from the plan,
  where it is always null) — they now do.

### `hpe_morpheus_network` data source: `cloud_id` selects among per-cloud records

Morpheus keeps one network record per cloud. When two clouds cover the same account — a cloud scoped
to one Azure resource group or AWS region alongside one covering the whole subscription or all
regions, or a vSphere cloud alongside the NSX-T cloud for the same vCenter — every cloud discovers
the network on its next refresh and keeps its own record of it, under the same name, so a lookup by
`name` alone found several and failed. The data source gains `cloud_id`: it filters a lookup by
name, is checked against a lookup by id, and is returned on every read. The ambiguity error now
names each record's cloud. The documentation explains the per-cloud model and warns against
deleting another cloud's record to tidy up, since that deletes the network itself.

### `hpe_morpheus_network` applies a rename and IPv6 changes in place where the type allows

The update request never carried `name`, `cidr_ipv6`, `gateway_ipv6`, `netmask_ipv6`, the IPv6
DNS and DHCP attributes, `pool_ipv6_id` or `ipv4enabled`/`ipv6enabled`, so a change to any of them
that was planned in place did nothing on the appliance. They are sent now; `name` only when it
changes, because the API also sets a blank display name from it.

### `hpe_morpheus_network` reports a delete Morpheus refused

Morpheus answers a refused network delete with HTTP 200 and `success: false` — for example when a
network plugin's host-side clean-up fails — and the provider read the status code alone, so the
resource was removed from state while the network remained. The answer is read now: a refusal is
retried briefly and then reported with Morpheus's own message.

### `Invalid network` on instance create and clone now says which networks would have been accepted (MORPH-17643)

Version 2.1 made `hpe_morpheus_instance` and `hpe_morpheus_instance_clone` explain an `Invalid
network` rejection — which network, which resource pool it belongs to, which pool the instance asked
for. The explanation now also answers the question that follows: which networks *would* Morpheus
have accepted? It asks the appliance's own selectability check (the `zoneNetworkOptions` option
source, run for the same cloud, group, layout and resource pool the validation used) and lists the
result, naming the scope so it can be reproduced with `hpe_morpheus_option_source_values`. When the
answer is empty — a resource pool with no selectable network at all, which is what the incident that
motivated 2.1's change turned out to be — the error says so as a fact rather than inviting the user
to choose another network. A clone reads the cloud, group and layout it inherits from its source
instance. The list is added only on that failure, is capped at fifteen networks, and a failed
lookup leaves the rest of the explanation as it was.

One line of 2.1's text is corrected: a network that belongs to no resource pool passes Morpheus's
pool filter for every pool, so the explanation no longer says such a network was excluded by the
requested pool; it points at the network's type and its group and tenant visibility instead.

## Bug fixes

### `hpe_morpheus_storage_volume`: resizing works, and a rename replaces the volume (MORPH-18280)

The resource's update sent the new name and size to `PUT /api/storage-volumes/{id}`, which never
applied either — the API has no rename for a volume, and from Morpheus 9.1.0 that endpoint transfers
the volume to another tenant. A change to `max_storage` now calls the resize endpoint, and `name`
cannot be changed after the volume is created: a change replaces the volume, which the plan shows.

### The bundled Morpheus SDK follows the 9.1.0 API (MORPH-18280)

The SDK is regenerated from the 9.1.0 API description. No resource attribute changes; the provider
reads the sub-objects the API can now omit (an instance's plan, a policy's or task's type, a
router's cloud, a cloud's type, a storage server's type, a tenant's account name) as unset rather
than failing, and sends the inputs the API accepts in several forms (a role's default persona, a
tenant's role and parent, a user's roles and key pair, a licence's type and images) in the form it
did before.


### `hpe_morpheus_network`: removing an optional attribute clears it (MORPH-17990)

Removing an optional attribute from a network's configuration used to do one of two wrong things.
An attribute Morpheus also computes (`gateway`, the DNS servers, `pool_id`, `description`, the IPv6
addressing, `labels`, …) was silently kept: Terraform planned no change, and the value stayed on the
network. An attribute Morpheus does not compute (`cidr`, `cidr_ipv6`) was planned to `null` but
never sent, so every apply failed with *Provider produced inconsistent result after apply*. Both are
fixed: a removed attribute is cleared on the network and reads back empty, `cidr` included. The
provider records at each apply which attributes the configuration set, so a value the user never
set — a display name Morpheus filled in, a gateway a cloud sync reported — is left alone; booleans
keep their value. Removing `tenant_ids` revokes every tenant but the owning one, which Morpheus
always keeps. A clear of an attribute the network type does not edit after create is planned as a
replacement, like a change to it, and `zone_pool_id` cannot be removed from the HVM port group types,
whose pool is the network's cluster. The record starts with the first apply made by this version:
an attribute removed before that apply is kept once more — apply, then remove.

### Removing an optional attribute clears it on fifteen more resources (MORPH-17990)

The same fix, resource by resource, for the attributes practitioners edit: `description` on
`hpe_morpheus_subnet`, `hpe_morpheus_cluster`, `hpe_morpheus_security_group`,
`hpe_morpheus_monitoring_check`, `hpe_morpheus_monitoring_group`, `hpe_morpheus_vdi_pool`,
`hpe_morpheus_vdi_gateway`, `hpe_morpheus_vdi_app` and `hpe_morpheus_instance` (and
`hpe_morpheus_tenant`, `hpe_morpheus_role`, whose removal already worked and is now covered by
tests); `labels` on `hpe_morpheus_group`, `hpe_morpheus_subnet`, `hpe_morpheus_cluster` and
`hpe_morpheus_instance`; `code` and `location` on `hpe_morpheus_group`, and its `cloud_ids`;
`pool_id` and `resource_permission_group_ids` on `hpe_morpheus_subnet`;
`resource_permission_group_ids` on `hpe_morpheus_security_group`; `name` on
`hpe_morpheus_security_group_rule`; `gateway_url` on `hpe_morpheus_vdi_gateway`; `retention_count`
and `schedule_id` on `hpe_morpheus_backup_job`; `path`, `ssh_username` and `storage_provider_id` on
`hpe_morpheus_backup_host`; `subdomain`, `account_number`, `account_name` and `customer_number` on
`hpe_morpheus_tenant`. Each of these was either a failed apply (*Provider produced inconsistent
result after apply*) or a silently ignored removal. The bundled Morpheus SDK sends null for them;
the cluster API rejects an empty label list, so an empty list goes under the key it does accept.
The provider index documents the behaviour and which resources have it.

### Removing an optional attribute clears it on fourteen more resources (MORPH-17990)

The same fix for `hpe_morpheus_user` (`first_name`, `last_name`, `linux_username`,
`windows_username`, `linux_key_pair_id`), `hpe_morpheus_service_plan` (`description`,
`price_set_ids`), `hpe_morpheus_image` (`labels`, `min_disk`, `os_type_id`, `ssh_username`,
`user_data`, `tenant_ids`), `hpe_morpheus_cloud` (`code`, `labels`, `location`, `appliance_url`,
`data_center_name`, `external_id`, `import_existing_vms`, `keyboard_layout`),
`hpe_morpheus_load_balancer`, `hpe_morpheus_load_balancer_profile` and
`hpe_morpheus_load_balancer_pool` (`description`), `hpe_morpheus_load_balancer_monitor`
(`description`, `send_data`, `send_type`, `receive_data`, `receive_code`, `monitor_destination`,
`extra_config`), `hpe_morpheus_load_balancer_virtual_server` (`description`, `vip_hostname`,
`ssl_cert`, `ssl_server_cert`, `pool_id`), `hpe_morpheus_network_router_nat` (`description`,
`source_network`, `destination_network`, `translated_network`, `translated_ports`),
`hpe_morpheus_network_router_firewall_rule` (`description`, `application`, `port_range`),
`hpe_morpheus_network_router_firewall_rule_group` (`description`, `tenant_ids`),
`hpe_morpheus_network_router_bgp_neighbor` (`description`, `forwarding_address`,
`protocol_address`, `route_filtering_in`, `route_filtering_out`) and
`hpe_morpheus_network_router_route`, whose `name`, `description` and `priority` cannot be edited
after create: removing one replaces the route, as a change does. Cleared cloud configuration
values, which Morpheus stores as empty strings, now read back as unset rather than as `""`. Three
attributes stay as they are when removed, because the Morpheus API cannot clear them:
`hpe_morpheus_image` `min_ram` (an empty minimum memory is rejected with a server error) and
`storage_provider_id`, and `hpe_morpheus_load_balancer_monitor` `monitor_username` (both set only
when the resource is created). The provider index lists the resources that clear a removed
attribute.

### Removing an optional attribute clears it on nineteen more resources (MORPH-17990)

The same fix completes the provider: `description` on `hpe_morpheus_budget`,
`hpe_morpheus_certificate`, `hpe_morpheus_deployment`, `hpe_morpheus_power_schedule`,
`hpe_morpheus_policy` (and its `tenants`), `hpe_morpheus_network_group` and
`hpe_morpheus_cluster_namespace` (and their `tenant_ids`), `hpe_morpheus_provisioning_license` (and
its `virtual_images` and `tenants`), `hpe_morpheus_network_firewall_rule_group` (and its
`tenant_ids`) and `hpe_morpheus_network_firewall_rule` (and its `application` and `profile` lists);
`description`, `source_url` and `api_type` on `hpe_morpheus_option_list`; `description`,
`service_host`, `service_url`, `service_username` and `tenants` on `hpe_morpheus_storage_server`;
`service_url`, `service_username`, `network_filter`, `zone_filter` and `tenant_match` on
`hpe_morpheus_network_pool_server`; every optional text attribute of
`hpe_morpheus_os_type` (`description`, `os_family`, `os_version`, `cloud_init_version`, `category`,
`os_codename`, `os_name`, `vendor`); `code`, `result_type`, `retry_count`, `retry_delay_seconds`
and `labels` on `hpe_morpheus_task`; `category`, `run_as_user` and `labels` on
`hpe_morpheus_container_script`; `dns_domain`, `gateway`, `netmask` and `subnet_address` on
`hpe_morpheus_network_pool`; `retention_days` on `hpe_morpheus_storage_bucket`; and
`appliance_name`, `primary_color` and `secondary_color` on `hpe_morpheus_setting_whitelabel`, whose
logos and favicon were already reset when removed. The bundled Morpheus SDK sends null for them,
or an empty list. The fourteen remaining resources have nothing that can be cleared, and their
documentation is unchanged; the provider index lists every resource that clears a removed
attribute. A container script's empty label list now reads back as unset rather than as `[]`.

### `hpe_morpheus_security_group_rule`: rules with ids of a million or more could not be read

The rule id was sent to the API as a floating-point number, which for ids of 1,000,000 and above
was formatted in exponent notation (`/rules/1.962589e+06`) and not found: every read, update and
delete of such a rule failed. The id is now sent as an integer. The rule sweeper also passed the
security group and rule ids in the wrong order and could not remove leftover test rules.

### A delete Morpheus refuses is now an error, not a success (MORPH-17987)

Morpheus renders a delete it accepted as a request but refused as an operation — a network whose
host-side removal failed, a resource still in use — as **HTTP 200 with a body of `{"success": false,
"msg": ...}`**. Every resource judged the DELETE by status code alone, so a refused delete removed the
resource from Terraform state while it remained on the appliance, and the next plan tried to create
it again and failed on a name or VLAN collision.

- Every framework resource's delete now reads the body: a refusal fails the destroy with Morpheus's
  own message (`Morpheus refused the delete: <msg>; <field>: <error>`, or "no reason given"), and the
  resource stays in state to be retried or removed by hand. `hpe_morpheus_network` keeps its three
  retries for transient host-side failures. A 404 is still treated as already deleted.
- The same applies to the SDKv2 resources through the legacy API client, and to any other API call
  answered that way: a 2xx whose body says `success: false` is a failure carrying the message.
- Also fixed on the way: removing an environment variable from `hpe_morpheus_instance` reported
  success on a refusal.

### `hpe_morpheus_security_group_rule` reports a missing source, destination or port range at plan (MORPH-18279)

A rule that left out an attribute Morpheus requires planned cleanly and failed at apply with the API's
400 (`"source": "Source CIDR required"`, `"portRange": "Port Range required"`). The requirement is
conditional — `source`/`destination` when the type is `cidr` (the default), `port_range` unless the
protocol is `icmp` — so the attributes stay optional (an ICMP rule has no ports; a group or tier source
has no CIDR) and the plan now fails on the attribute with the condition spelled out.

### `hpe_morpheus_network` (MORPH-17908, MORPH-17909, MORPH-17910)

- Any change to a network with `tenant_ids` unset was planned as a replacement, because the
  attribute's unconditional replacement rule compared an unknown planned value with state. Tenant
  permissions are applied in place by Morpheus, so `tenant_ids` no longer replaces the network at all
  (MORPH-17908).
- `visibility = "public"` on a network whose `tenant_ids` did not include the master tenant failed
  after apply with an inconsistent result, because Morpheus creates such a network as `private`. The
  provider now reports this at plan time (MORPH-17909).
- The OVS port group and Azure examples in the documentation were invalid HCL (unquoted string
  values); the Amazon example placed `assignPublicIp` and `availabilityZone` in `config`, where
  Morpheus never read them — they are the top-level `assign_public_ip` and `availability_zone`
  (MORPH-17910).

### `hpe_morpheus_user` (MORPH-18061)

- Removing a configured `linux_username` or `windows_username` from the configuration failed with
  `resource was not updated`. The update sent the removed username as an explicit null, which the
  Morpheus API rejects with an internal error (MORPH-18060). It is now sent as an empty string, which
  Morpheus accepts and stores as unset, so the username is cleared and the following plan is empty.

### Updating a resource no longer replaces it because of an optional attribute left unset

Thirteen optional attributes on seven resources are set when the resource is created and replace it
when they change: `description` and `storage_provider_id` on `hpe_morpheus_image`,
`monitor_username` on `hpe_morpheus_load_balancer_monitor`, `group_layer` on
`hpe_morpheus_network_firewall_rule_group`, `name`, `description` and `priority` on
`hpe_morpheus_network_router_route`, `cloud_id` and `provision_type_id` on
`hpe_morpheus_os_type_image`, `associated_resource_id` on `hpe_morpheus_policy`, and `type_id`,
`type_code` and `provision_type` on `hpe_morpheus_storage_volume`. When one of them was left out of
the configuration, every update of the resource — a rename, say — was planned as a replacement: the
plan showed the attribute as known after apply, and the replacement rule counted that as a change.
Renaming an image created without a `description` destroyed and re-created the image. An
unconfigured attribute now keeps its value from state, so only a change the practitioner makes
replaces the resource. A provider-wide unit test keeps the attribute shape that caused this from
returning.

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
