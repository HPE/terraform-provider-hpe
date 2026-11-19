# v2.3.0 Release Notes

## Breaking changes

### `hpe_opsramp_alert_correlation_policy.algorithm_correlation.alerts_time_window`

`alerts_time_window` is now a number (`Int64`) instead of a string. Update configurations that
quoted this value; for example, change `alerts_time_window = "30"` to
`alerts_time_window = 30`. Existing state is upgraded automatically — the resource implements
`ResourceWithUpgradeState` to convert a stored string value to a number — so no `terraform state`
surgery is needed.

### `organization_matching_type` drops `EXCLUDE`

`hpe_opsramp_alert_correlation_policy`, `hpe_opsramp_alert_escalation_policy`,
`hpe_opsramp_alert_prediction_policy` and `hpe_opsramp_first_response_policy` accept only `ALL` and
`INCLUDE` for `organization_matching_type`; `EXCLUDE` is no longer a valid value. A configuration
that set `organization_matching_type = "EXCLUDE"` must be changed to scope the policy with `ALL` or
with `INCLUDE` and an explicit `included_clients` list.

### `hpe_morpheus_setting_appliance`: `stats_retainment_period` applied; credentials need Terraform 1.11

The resource is rebuilt on the plugin framework (see *Improvements*). Existing state is upgraded in
place, and configurations carry over, with two exceptions:

- `stats_retainment_period` now takes effect: the previous resource accepted it but never sent it.
  A configuration that sets it changes the appliance's stats retention on the first apply after
  upgrading. It accepts 30, 60 or 90.
- `smtp_password`, `proxy_password` and `currency_provider_api_key` keep their names but are now
  write-only. A configuration that sets any of them, or the new `twilio_auth_token`, needs
  **Terraform 1.11 or later, or OpenTofu 1.11 or later**.

## Deprecations

- `hpe_morpheus_setting_backup`: `retention_days` is deprecated in favour of `retention_count` and
  scheduled for removal in **2.6**. Despite its name the value has always been the number of
  successful backups to retain, not a number of days; `retention_count` names it correctly. Both are
  sent to the same setting, so changing `retention_days = N` to `retention_count = N` plans no
  change. A configuration that sets `retention_days` now gets a deprecation warning.
- `hpe_morpheus_network_router_nat` data source: `action` and `firewall` are deprecated and scheduled
  for removal in **2.6** (MORPH-14596). The Morpheus NAT read API does not return either value, so
  they have always been null; a configuration that references them now gets a deprecation warning.
  The `hpe_morpheus_network_router_nat` resource is unaffected — it already treats both as write-only.

## New functions

### `morpheus_sha512_crypt` (MORPH-20062)

`provider::hpe::morpheus_sha512_crypt(password, salt, rounds)` derives deterministic
Linux SHA-512 crypt password hashes for cloud-init without an external command or
an appliance connection. Terraform 1.8 or later is required. Salt and rounds are
explicit; 5000 rounds preserves the output of `openssl passwd -6`. Passwords are
limited to 4096 bytes and rounds to 1000 through 1000000 to bound computation.
Sensitive arguments retain sensitivity, but sensitive hashes can still be stored
in Terraform state. The compatibility round count is not a security recommendation.

## New resources

### `hpe_morpheus_appliance_setup`: first-time setup of a new appliance (MORPH-20057)

A freshly installed appliance could only be initialised through the setup wizard in a browser or
a script against its unauthenticated setup endpoint, so a Terraform-built appliance needed a
provisioner before the provider could log in. The new `hpe_morpheus_appliance_setup` resource waits
for the appliance's API to come up (a first boot takes a good twenty minutes; the create timeout
defaults to 40 minutes), then creates the master tenant and its first administrator in one call,
installing `license_key_wo` when it is set. The password and licence key are write-only. An
appliance already set up is adopted with `already_initialised = true` and a warning; a refresh that
finds the appliance needing setup again (reinstalled behind the same URL) removes the resource so
the next apply sets it up again. Destroying the resource leaves the appliance as it is. Import takes
the appliance URL.

### `hpe_morpheus_image`: local image uploads (MORPH-20059)

Images now accept `file_path` and optional `file_sha256` instead of `url`. Local files are streamed
without buffering the whole image, verified through the file listing, and updated in place without
replacing the image ID. Existing same-name, same-size files can be adopted after ordinary image
import. Local content changes are detected by checksum, while remote drift detection uses filename
and byte count because Morpheus does not expose remote checksums. Single-file management removes
other files after verification; replacing a same-name file removes it before uploading, so a failed
upload may temporarily leave the image without usable content. Removing `file_path` stops file
management without deleting its content. Existing URL and metadata-only configurations still work.

### `hpe_morpheus_network_pool_ip` resource and data source (MORPH-15704)

A single IP address can now be reserved from a network pool. Until now the provider could define
pools (`hpe_morpheus_network_pool`, `hpe_morpheus_ip_pool_ipv4`) but not take an address out of one;
that needed the API or the UI.

- **`hpe_morpheus_network_pool_ip`**: give `network_pool_id` and a `hostname` and the pool's IPAM
  allocates its next free address, returned in `ip_address`; set `ip_address` to reserve a specific
  one, and an address outside the pool's ranges is refused with Morpheus's own message. `hostname`,
  `ref_type` and `ref_id` (for example `ComputeServer` and a server id) change in place; the address,
  `description`, `domain_id` and `create_dns` (DNS A and PTR records in a network domain) replace the
  reservation. Destroying it releases the address. Import with `<network_pool_id>:<id>`.
- **`hpe_morpheus_network_pool_ip` data source**: read a reservation by `id`, `ip_address` or
  `hostname` within a pool.

Verified against the built-in Morpheus pool type and an Infoblox-backed pool.

### `hpe_morpheus_virtual_switch` resource and virtual switch data sources (MORPH-17745)

Morpheus 9.1.0 virtual switches on HVM layout-2.0 clusters can now be managed from Terraform. A
virtual switch is a per-cluster switch with traffic **segments** and per-host **uplinks**, realised
on the hosts as Linux bridges or OVS.

- **`hpe_morpheus_virtual_switch`**: `cluster_id`, `name` (at most 16 characters), `base_type`
  (`general`, `iscsi` or `sdn`; immutable), `bond_mode`, `mtu`, `description`; a `segments` map
  (`traffic_type`, `vlan_id`, addresses, per-host `host_ips`, `mtu`) and `uplinks` (`server_id`,
  `nic1`, optional `nic2`). Rename, description and bond changes, and adding, editing or removing
  segments apply in place through the appliance's rename, update, composite-edit and segment
  endpoints; a change to a segment's `traffic_type` or `mtu`, or to `base_type`, replaces the
  switch. Every mutation is asynchronous on the appliance and the resource waits for its process
  and for the switch to leave `syncing`. The plan-time rules Morpheus enforces (management segment
  needs an address, `iscsi` is single-NIC and exclusive of other data segments, one untagged segment
  per data type, no duplicate VLANs, no uplink removal) are reported before apply.
- **`hpe_morpheus_virtual_switch`** and **`hpe_morpheus_virtual_switches`** data sources: read a
  switch by id or name, or list a cluster's switches.
- **`hpe_morpheus_cluster_available_nics`** data source: the spare NICs of each host of a cluster,
  so uplinks can be chosen by data rather than by guess.

Requires Morpheus 9.1.0; gated on the `virtual_switch` test capability.

## Improvements

### `hpe_morpheus_network`: HVM networks on a virtual switch (MORPH-18378)

On a Morpheus 9.1 HVM cluster with virtual switches, the network resource binds each HVM type to
its switch: a standard network's (`ovsPortGroup`) `switch_id` is the switch's libvirt network name
(`hvmnetN`, from the `hvm-vswitch-or-router` option source by the switch's name), a private network
takes `config_hvm_private.virtual_switch_id`, and a data network a segment of a switch through
`config_hvm_data.virtual_switch_segment_id`. The `switch_id` description and the docs of both
resources now say so, together with what Morpheus enforces across them — the switch must be in
the network's cluster, a network VLAN or trunk may not collide with a segment VLAN on the switch,
and a switch with networks cannot be deleted or have its MTU changed — with an example that creates
a switch and both kinds of network on it. Two changes here make that example apply in one go
(the third it relies on, finding a cluster's hidden network server by name, is MORPH-20649 below):

- A data network on a virtual switch segment waits, before the create, until Morpheus has
  recorded the segment's host interface, which it does on the cluster's inventory sync (every few
  minutes); without it Morpheus fails the create with a bare 500. The wait is bounded at twelve
  minutes.
- `hpe_morpheus_virtual_switch` keeps the ids of its existing segments known when the switch is
  updated; previously every edit of the switch planned a network built on one of its segments for
  replacement — and destroyed it.

### `hpe_morpheus_network_pool`: every IP range of a pool, with `ip_address_ranges` (MORPH-15713)

A Morpheus pool can have any number of IP ranges, but the resource could only express one: `ip_ranges`
was a single `{ starting_address, ending_address }` object. Read kept the first range and dropped the
rest, and — because Morpheus deletes any range missing from an update — applying the resource to a pool
that had several ranges silently deleted the others. The new `ip_address_ranges` attribute is a list of
`{ id, starting_address, ending_address }` holding every range. Each range is edited in place by its
`id`, so reservations in it (including `hpe_morpheus_network_pool_ip` resources) survive a change to its
addresses; a range removed from the list is deleted, and when Morpheus refuses because addresses in it
are still reserved, the provider reports which range rather than returning an inconsistent result. The
`hpe_morpheus_network_pool` data source exposes the same list.

`ip_ranges` is deprecated and is removed in the 2.6 release. Moving a configuration to the same values
under `ip_address_ranges` plans no change. While `ip_ranges` is still in use on a pool that has several
ranges, the provider refuses the update instead of deleting the other ranges.

### `hpe_morpheus_setting_appliance` moves to the plugin framework

The appliance settings resource is rebuilt on the plugin framework. The type name and the attribute
names are unchanged, `id` included (still always `"1"`), and it gains the settings the previous
resource did not have (see below). Existing state is upgraded in place: the first plan after
upgrading shows one `~ write_only_revision` update, which re-sends the configured credentials, plus
any change to `stats_retainment_period` the configuration asks for (see *Breaking changes*).

- **The credentials are write-only.** `smtp_password`, `proxy_password` and
  `currency_provider_api_key` are sent to Morpheus and no longer stored in state. A change is
  detected from a digest kept in the resource's private state and plans an update through
  `write_only_revision`, the computed attribute every resource with write-only values has: the time
  they were last sent to Morpheus. A credential changed on the appliance outside Terraform is detected
  from the hash Morpheus reports for it, and re-sent. Before, a configured password planned a change
  on every run (the resource compared a plain SHA-256 of it with the appliance's HMAC), and the
  currency key was kept in state in clear.
- **Only the settings that change are sent.** A first apply no longer switches tenant registration
  and Docker privileged mode off when the configuration does not set them, and a plan lists only the
  settings that change.
- **`stats_retainment_period` takes effect.** The previous resource declared it but never sent or
  read it. It is now sent when it changes and read back, and a value other than 30, 60 or 90 fails
  when planning.
- **Every setting the endpoint takes.** The resource also manages the password policy
  (`password_min_length`, `password_min_upper_case`, `password_min_numbers`,
  `password_min_symbols`), browser sessions (`user_browser_session_timeout`,
  `user_browser_session_warning`), user deactivation (`expire_password_after_days`,
  `disable_user_after_attempts`, `disable_user_after_days_inactive`,
  `warn_user_days_before_disable`), Twilio (`twilio_account_sid`, `twilio_sms_from`,
  `twilio_auth_token`), the sync intervals (`cloud_sync_interval_seconds`,
  `cluster_sync_interval_seconds`), the other retention periods (`usage_retainment_period`,
  `invoice_retainment_period`, `incident_retainment_period`, `reports_retainment_period`), the HTTP
  task hosts (`http_blacklist_hosts`, `http_approvelist_hosts`), the Agent and SSH flags
  (`no_agent`, `agent_ssl_verify`, `disable_ssh_password_auth`), `default_locale`,
  `default_vdi_gateway_id`, `max_option_list_size`, `dashboards_to_display` and `enable_murph`.
  Each is read from the appliance and sent only when configured. The four user-deactivation
  settings are named after the appliance settings page, with the unit where it is days.
  `twilio_auth_token` is write-only like the other credentials; Morpheus reports no hash
  for it, so a change made to it outside Terraform is not detected.
- **`registration_enabled` changes in place** rather than replacing the resource; the outcome on the
  appliance is the same.
- **Import and types.** Import takes any ID, optionally with the digests of the configured
  credentials (`1#smtp_password=<sha256>`); without them the configured credentials are taken as
  already applied, so a later change to one is only detected once the settings are next updated for
  another reason. `smtp_port`, `proxy_port`, `default_role_id` and `default_user_role_id` are numbers
  in state; quoted values in a configuration keep working.

### `hpe_morpheus_setting_backup` moves to the plugin framework; `retention_count`

The backup settings resource is rebuilt on the plugin framework. The type name and the attribute
names are unchanged, `id` included (still always `"1"`), and existing state is upgraded in place: an
unset default storage bucket or schedule, which the previous resource stored as `0`, becomes null, and
the plan after upgrading is clean.

- **`retention_count`** is the new name for the retention setting: the number of successful backups
  to retain. `retention_days` is kept as a deprecated alias of the same setting (see *Deprecations*);
  the two must not both be set, and whichever is set is mirrored into the other in state, so a
  configuration can move from one to the other without a change being planned. Both must be at
  least 1: Morpheus stores a `0` but treats it as unset where it reads the retention, so `0` and
  negative values, which the previous resource passed through, are refused when planning.
- **Only the settings that change are sent.** A first apply sends the configured settings only, so
  an attribute left out of the configuration keeps the appliance's value instead of being reset, and
  a plan lists only the settings that change.
- **A refused update fails the apply.** Morpheus answers a default schedule or storage bucket that
  does not exist with a success flag of `false`; the resource reports Morpheus' message instead of
  recording the settings as applied. A setting Morpheus accepts but does not apply fails the apply
  too, naming the setting and the value the appliance reports, instead of Terraform's generic
  "inconsistent result" error.
- **`backup_appliance` is guarded for the master tenant.** Morpheus keeps one appliance backup for
  the whole appliance: a request from any other tenant to enable it is accepted but applied to the
  master tenant's appliance backup, which that tenant then reads as `false`, and a request to disable
  it is ignored. The resource refuses `backup_appliance = true` at plan time when the provider is not
  configured as the master tenant (with a warning instead, when it cannot determine the tenant). The
  attribute descriptions now say which settings are appliance-wide, which is per tenant, and which is
  the master tenant's.
- **Import takes any ID** (for example `1`); the settings are a singleton.

### `hpe_morpheus_license`: the `key` description says what the attribute is

The `key` attribute's description says what the attribute is — the licence key, sent to the
appliance and never read back — and leaves how the provider tracks a change (write-only, the
digest in private state, the in-place update) to the resource description, so that documentation
generated from the provider specification for other tools does not inherit Terraform's state and
plan vocabulary for an argument. No behaviour changes.
## Bug fixes

### `hpe_morpheus_instance`: tags Morpheus adds itself no longer fail the apply or show as drift (MORPH-20973)

Morpheus does not always hold only the tags an instance was created with: tag policies,
provisioning workflows (`addTags`), layout option types and the cloud's own tag sync (vCenter tags
on VMware, for example) add tags to a new instance, or overwrite a value. The provider stored the
API's tag list verbatim after create and update, so the apply failed with *Provider produced
inconsistent result after apply … `.tags`: length changed from 4 to 6* (or *was cty.SetVal([…]),
but now null* when Morpheus held none yet), Terraform tainted the healthy, running instance, and the
next apply destroyed and recreated it only to fail the same way. A configured `tags = []` failed
even when Morpheus added nothing, since an empty list was stored as null.

`tags` is now the set of tags the resource **manages**, matched by name:

- a read keeps, of the tags Morpheus holds, only the configured ones, with Morpheus's value — so a
  value Morpheus overwrites and a tag it removes are still drift, while a tag it adds is not. A value
  that differs only in case is taken as the configured one: a VMware cloud reports a tag in the case
  of the vCenter tag it matched (`prod` as `Prod`), which is not drift. Tag names must be unique
  ignoring case;
- after a create or update a configured tag Morpheus does not (yet) report, or reports with another
  value, is kept as configured, with a warning naming the difference, so the apply completes; the
  next refresh shows what Morpheus holds;
- an update sends only the difference between the previously managed tags and the plan
  (`addTags` / `removeTags`) instead of replacing the whole list, so tags Morpheus added survive;
- an omitted `tags` stays null and `tags = []` stays an empty set; either removes the tags the
  resource previously managed and nothing else;
- the new computed **`tags_all`** lists every tag Morpheus holds for the instance.

On import `tags` is not populated — the resource cannot know which of Morpheus's tags a
configuration means to manage — and `tags_all` is; configuring tags after an import adds them
through `addTags` and manages them from then on.

Thanks to the customer who diagnosed the defect and proposed a fix.

### Compute server lookups work on an appliance with an external Kubernetes cluster

`hpe_morpheus_compute_servers`, and `hpe_morpheus_compute_server` looked up by name, failed with
*decoding failed … cannot parse "+0000" as "Z07:00"* on any appliance that had registered an
external Kubernetes cluster. Morpheus reports the time a server's statistics were sampled with a
`Z` offset for most servers and with `+0000` for the cluster's node; the second form is not a
valid RFC 3339 timestamp, and one such server made the whole server list unreadable, whatever was
being looked up. The provider now reads that timestamp as text. Neither data source exposes it, so
their attributes are unchanged.

### `hpe_morpheus_cluster`: an External Kubernetes cluster's token and kubeconfig no longer reach state (MORPH-19991)

For an External Kubernetes cluster the appliance returns `config.apiToken` — a cluster-admin
service-account token — and `config.serviceAccess` — the kubeconfig carrying it — in clear, and
the `hpe_morpheus_cluster` data source copied the whole map into its `config` attribute, as did
the resource on import. Reading such a cluster, even just for its id, put live cluster
credentials into `terraform.tfstate` and plan output, unmarked. Both now withhold credential
keys from `config` before it is stored: `apiToken`, `serviceAccess`, `serviceToken`, `kubeToken`,
`kubeConfig`, `joinCommand`, and any key whose name ends in `token`, `password`, `passwd`,
`secret`, `passphrase`, `privateKey`, `kubeconfig` or `serviceAccess`, at any depth. The keys are
removed rather than masked, so nothing about them is in state; the rest of `config` is unchanged.
A cluster imported with an earlier release keeps what it imported until it is removed from state
and imported again; the resource warns on every refresh while a credential key remains. A
configuration that writes such a key into the resource's own `config` is stored as written — that
is what a typed block with write-only attributes is for.

### `hpe_morpheus_network_server` data source: a cluster's hidden network server is found by name (MORPH-20649)

An HVM cluster's own network server — `<cluster> HVM Network Server`, the one every network on
the cluster must name in `network_server_id` — is hidden from the ordinary listing, and the data
source listed without hidden servers and then re-read its match by id, which Morpheus refuses for
a hidden server. A lookup by name now includes hidden servers and uses the listing's own record,
so `data "hpe_morpheus_network_server" "hvm" { name = "<cluster> HVM Network Server" }` resolves.
An exact name is unambiguous, so there is no flag to set. A lookup by `id` that the appliance
refuses — it answers 404 for a hidden server regardless of `showHidden` — is likewise satisfied from
the listing, so hidden servers resolve either way. The data source's `config` description now
says what it holds: the configuration as Morpheus returns it, with only the top-level
`service_password` and `service_token` masked.

### `hpe_morpheus_setting_backup`: the default storage bucket and schedule are applied (MORPH-20485)

`default_backup_storage_bucket_id` and `default_backup_schedule_id` were sent beside the
`backupSettings` object in the request instead of inside it, where the appliance reads them, so
both were accepted and silently ignored: the setting applied without error, the appliance kept no
default bucket or schedule, and every plan proposed the same change again. Both now travel inside
`backupSettings` and take effect on the first apply. A configuration that set either attribute
sees the appliance change accordingly after upgrading.
### Affinity groups on Morpheus 9.1 (MORPH-20061, MORPH-20083)

Morpheus 9.1 refuses an affinity group with no servers (its new rule-shape validation: a group is a
VM-to-VM rule with at least one server, or a VM-to-Host rule), on create and when the last server is
removed, so the provider could no longer create an affinity group against 9.1 at all.

`servers` on `hpe_morpheus_cluster_affinity_group` and `hpe_morpheus_cloud_affinity_group` is settable
again: set, this resource owns the membership — the group is created with those servers and every
apply re-asserts exactly that set; unset, nothing is sent on create and the live membership is
re-asserted on update, so the `*_affinity_group_member` resources and provisioning-time placement are
not evicted. A group should use one of the two, never both. On a 9.1 or later appliance both group
resources refuse **at plan time** a create or replacement with no `servers`, with a message naming the
fix, instead of failing mid-apply with the API's error (or, on some 9.1 releases, a bare 403); on
older appliances an empty group is still created, and the check fails open when the appliance version
cannot be read. A member-managed group is bootstrapped on 9.1 by creating it with its first member in
`servers` and then dropping `servers` from the configuration, which plans empty. Destroying a member
that is the group's last server on 9.1 no longer fails: the server is left in place with a warning,
so `terraform destroy` of a group and its members works (the group's own destroy removes the
membership).

Known issue (Morpheus core, MORPH-20098): on a VMware cloud, renaming a cloud affinity group and then
destroying it leaves the DRS rule behind in vCenter under its original name, after which later
`KEEP_TOGETHER` groups containing the same VMs are refused by vCenter. Until that is fixed, avoid
renaming VMware cloud affinity groups that will later be destroyed.

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

### `hpe_morpheus_license` — the licence key is write-only and leaves state

The appliance licence key — a few thousand characters that unlock the product — used to be stored in
Terraform state, where anyone with access to the state file could read it. The resource is rebuilt
on the plugin framework so that it no longer is.

- **`key` is write-only.** It is sent to the appliance on create and whenever it changes, and is
  never stored in state or plan, so an ephemeral value may be used. **Nothing changes in your
  configuration** — the attribute keeps its name — but write-only attributes need **Terraform 1.11
  or later, or OpenTofu 1.11 or later**.
- **Existing state is upgraded automatically.** The first refresh after upgrading drops the stored
  key from state and the plan is empty.
- **A change to the key is detected without any version or trigger attribute.** Morpheus reports a
  SHA-256 hash of each installed licence's key; the provider keeps the same digest in the resource's
  private state and compares it with the configured key on every plan. A changed key plans an
  in-place update that installs the new key; surrounding whitespace in the key is ignored, as it is
  by Morpheus. The key is only re-sent when it actually changed — a non-stacking install deactivates
  every other installed licence, so an unchanged key is never posted again.
- **Import by key id leaves a clean plan.** `terraform import hpe_morpheus_license.x <keyId>` (the
  first eight characters of the key) records the digest from the appliance, so the configured key is
  recognised at once and a later change to it is detected from the first plan. No digest needs to be
  put in the import ID.
- **A licence uninstalled or replaced outside Terraform now plans a *create*** (the refresh removes
  the resource from state, as every framework resource does when its object is gone), where the
  previous resource planned an in-place update. The apply re-installs the configured key either way.
- **A detail the licence does not specify is null, not a zero value.** The previous resource wrote
  `""`, `0` or `false` for anything the appliance did not report; the rebuilt resource, like every
  other framework resource, writes null. Concretely: `end_date` is null for a perpetual licence (was
  `""`); `max_instances`, `max_memory` and `max_storage` are null when unlimited (was `0` — Morpheus
  itself reports an unlimited legacy limit as null, and on the newer standard limits `0` means
  *nothing allowed*); `free_trial` is null (was always `false`: current Morpheus releases do not
  report it); `hard_limit`, `support_level` and `account_name` are null when the licence does not
  specify them. An expression that compared one of these with `""`, `0` or `false` should compare
  with `null`, or wrap the attribute in `coalesce()`.
- Destroy still uninstalls the licence matching `key` and no other.

## New resources

### `hpe_morpheus_plugin` resource and data source (MORPH-18591)

Morpheus plugins — the jars behind *Administration → Integrations → Plugins*, whether one of the
appliance's prepackaged set, a release from the Morpheus Exchange or a locally built one — could
only be installed through the UI or a hand-rolled multipart upload. The `hpe_morpheus_plugin`
resource installs, updates, enables, configures and removes them.

The jar comes from a local path (`jar_path`) or a URL (`jar_url`); exactly one is given. A URL is
downloaded by Terraform and uploaded to the appliance, so the appliance itself needs no route to the
exchange: `jar_url` can be an exchange `/api/packages/pull?code=…&version=…` link, which is what the
exchange records as each release's `downloadURL`. `jar_sha256` pins the artifact: when set — from
`filesha256()` for a local jar, or the `shaHash` the exchange publishes beside each release — the
jar is verified before upload and a change to it re-uploads; when a local jar is rebuilt in place,
the provider hashes it at plan time, so the rebuild shows as an update without any trigger
attribute. An install waits for the plugin to report `loaded` and fails with the plugin's own
message when it reports `error`. Uploading a newer jar updates the installed plugin in place
(Morpheus matches on the plugin's `code`); a jar whose code differs from the managed plugin's is
refused, because Morpheus would install it alongside rather than replace it. `enabled` and the
plugin's option-type values in `config` go through the plugin update endpoint; `code`, `name`,
`version`, `status`, `providers`, `is_prepackaged` and the plugin's URLs are read back. Destroying
the resource uninstalls the plugin. A plugin installed by other means is imported by id.

The `hpe_morpheus_plugin` data source reads one plugin by `id` or by its `code`, for configurations
that refer to a prepackaged plugin or one installed outside Terraform.

A plugin's providers come back from the API in a different order on every read, and several may
share a type and a name (a plugin's UI tabs often do), so `providers` is a list the provider sorts
by type, then name.

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

### `hpe_morpheus_cluster_layouts`: find a cluster layout by its version (MORPH-17748)

Morpheus 9.1.0 introduces cluster layout 2.0 for HVM clusters, and from then on the layout a cluster
is built from decides what the cluster can do: only a layout whose `cluster_version` is 2.0 or
greater builds on virtual switches and takes the 2.0 network settings described under
`hpe_morpheus_cluster` below. A configuration therefore needs to select a layout by that version,
rather than by an id or a display name that differs from one appliance to the next.

The new `hpe_morpheus_cluster_layouts` data source lists the cluster layouts, with a repeatable
`filter` block over `name`, `code`, `cluster_version`, `compute_version` and `provision_type`. Each
value is a Go regular expression, matched unanchored — `2.0` also matches `12.0`, so an exact version
is written `^2\.0$`; a layout matches a block when the field matches any of the block's values, and
every block must match. Filtering is client-side, as in `hpe_morpheus_clusters`. Each layout is
returned with its scalar attributes, so

```hcl
data "hpe_morpheus_cluster_layouts" "hvm_2" {
  filter {
    name   = "cluster_version"
    values = ["^2\\.0$"]
  }
  filter {
    name   = "code"
    values = ["^mvm-"]
  }
}
```

selects the HVM 2.0 layout, and `one(data.hpe_morpheus_cluster_layouts.hvm_2.cluster_layouts).id`
supplies its `id` to `hpe_morpheus_cluster`'s `layout_id`. `cluster_layouts` is a set, and a regular
expression can select several layouts; `one()` fails the plan unless exactly one matched, so a
selection that has become ambiguous is caught at plan time rather than building a cluster from the
wrong layout.

`cluster_version` is the version of the cluster software a layout deploys: the Kubernetes version
for a Kubernetes layout (`1.33.x`), the HVM cluster version for an HVM layout (`1.3`, `2.0`). The
singular `hpe_morpheus_cluster_layout` data source already exposed it; the `hpe_morpheus_cluster_layout`
resource now reads it back as a read-only attribute too. Morpheus does not accept the field when a
layout is created or updated, so a custom layout cannot be made a 2.0 layout through the API.

## Improvements

### Cloud-init user data on HVM instances and clones (MORPH-19115)

`hpe_morpheus_instance` and `hpe_morpheus_instance_clone` gain `config_hvm.user_data`: cloud-init
user data (cloud-config YAML or a plain script) applied once, at the VM's first boot, on top of the
user data Morpheus generates. Your `bootcmd`/`runcmd` entries run before Morpheus's own, so a guest
network fix-up (for example a netplan drop-in keeping a second NIC off the default route) can be in
place before the agent is installed. Map keys merge; list values replace Morpheus's. The value is
provision-time only, so changing it replaces the instance. Previously the only way to pass user data
was the untyped `config` block, which cannot be combined with `config_hvm`. Only the HVM provisioning
path honours per-instance user data as a typed attribute; the generic `config` block remains the
route for plugin-based clouds.
### Changes made to a secret on the appliance are detected (MORPH-19005)

A write-only secret — a password, key or client secret the provider sends to Morpheus but never
stores in state — used to be tracked only against the configuration: a change in HCL planned an
update, but a change made on the appliance (in the Morpheus UI, through the API or by another tool)
went unnoticed, and Terraform never restored the configured value.

Morpheus returns a fingerprint of each stored secret beside its masked value. The provider now
records that fingerprint after it sends a value and compares it on every refresh. When the secret
has been changed on the appliance, the next plan reports it and re-sends the configured value; for a
secret that cannot be changed in place (a cluster's or instance node's host SSH password) the plan
carries a warning instead, since a re-send would replace the resource. Nothing changes in the
configuration or the state file; detection starts with the first apply that sets the secret, so an
imported resource has no baseline until then.

Covered: `hpe_morpheus_image` (`ssh_password_wo`), `hpe_morpheus_load_balancer_monitor`
(`monitor_password_wo`), `hpe_morpheus_network_router_bgp_neighbor` (`password_wo`),
`hpe_morpheus_storage_server` and `hpe_morpheus_network_pool_server` (`service_password_wo`),
`hpe_morpheus_cloud` (`config_vmware.password`, `config_azure.client_secret`),
`hpe_morpheus_instance_node` (`ssh_password`) and `hpe_morpheus_cluster` (`server.ssh_password_wo`).
Morpheus returns no fingerprint for `hpe_morpheus_user`, `hpe_morpheus_storage_bucket`,
`hpe_morpheus_backup_host` and `hpe_morpheus_provisioning_license`, which keep the configuration-only
behaviour.

On Morpheus 9.1.x a load balancer monitor password changed through the API is not actually stored
(the stored-secret fingerprint is unchanged), so no drift arises from that path; the detection
mechanism is unchanged and is exercised live on resources whose secret the appliance does store.

Acceptance tests that create a network pool server or a storage server now need a reachable
service: Morpheus 9.1.2 verifies the connection and the credentials on create (`TF_ACC_IPAM_*`,
`TF_ACC_STORAGE_SERVER_*`); they skip when it is not configured. Because Morpheus 9.1.x also
re-validates those credentials on update, an out-of-band password rotation cannot be staged against
an ordinary integration, so the network pool server and storage server write-only and drift
acceptance tests are currently skipped pending an IPAM or storage service with rotatable
credentials; set `TF_ACC_IPAM_ROTATABLE=1` or `TF_ACC_STORAGE_SERVER_ROTATABLE=1` to re-enable them
against such a service.

### `hpe_morpheus_cluster`: the network settings of an HVM cluster layout 2.0 (MORPH-17748)

An HVM cluster built from a layout 2.0 (Morpheus 9.1.0) is networked over the hosts' virtual switches
rather than the bond-and-bridge scheme of the earlier layouts, and carries its own DNS, NTP and proxy
settings. `config_hvm` gains the corresponding optional attributes: `virt_switch_name` — the virtual
switch every host must already carry, as Morpheus refuses to add a host without it — and `dns1`,
`dns2`, `dns3`, `search_domain`, `ntp1`, `ntp2`, `proxy` and `no_proxy`.

Morpheus silently ignores these settings on a layout below 2.0. Before creating or updating a
cluster the provider therefore reads the selected layout's `cluster_version` and, when any of them
is set for an older layout, refuses with a message naming the layout's version (when the version
cannot be determined, the operation is allowed through). A change to `virt_switch_name` replaces the
cluster, since switches are attached to each host as it is added. A change to any of the other
settings is an in-place update that Morpheus reconciles onto the hosts through its daily
drift-correction sync, so the hosts converge asynchronously after the apply, and a value that is
cleared is not rolled back on the hosts by itself. The resource documentation carries a complete
2.0 example and explains why a change of `layout_id` replaces the cluster rather than upgrading it.

### New read-only and query attributes across several resources and data sources (MORPH-18721)

Schemas were refreshed against the current Morpheus API, surfacing attributes the provider did not
previously expose:

- `hpe_morpheus_cloud_affinity_group` and `hpe_morpheus_cluster_affinity_group` (and their singular
  and plural data sources) gain read-only `host_group`, `vm_group` and `owner` objects describing a
  group's rule shape and owning tenant.
- The `hpe_morpheus_compute_server` and `hpe_morpheus_compute_servers` data sources gain the computed
  `guest_agent_status`.
- The `hpe_morpheus_instance_snapshot` resource and data source gain the computed `max_storage`
  (total snapshot size in bytes), and the data source's `snapshot_files` gain `file_size` (bytes).
- The `hpe_morpheus_storage_volumes` data source gains the optional `include_tenants` and `tenant_id`
  query inputs (master-tenant list filters).
- The `hpe_morpheus_user` data source gains the computed `tenant_id` (the owning tenant).
- The `hpe_morpheus_setting_whitelabel` resource gains the settable `security_banner`,
  `external_legal_text` and `supports_themes`, and the read-only `account`.

### `hpe_morpheus_cloud`: backup providers are assigned at the cloud level (MORPH-18721)

Morpheus 9.x assigns backup providers to a cloud rather than through a per-cloud-type backup mode.
The new optional-and-computed `backup_provider_ids` (a set of integration ids) is sent on create and
update and read back; left unset, Morpheus assigns the internal backup provider by default and that
default is read into state. On Morpheus releases before 9.0.2 the API accepts and stores the value but
does not return it, so there the configured value is kept in state rather than read back. The existing
`config_aws.backup_provider` (stored as `config.backupMode`)
is **deprecated**: it has no effect on Morpheus 9.x — the provider no longer sends it, and its value
is retained in Terraform state only. Use `backup_provider_ids` instead.

### `hpe_morpheus_cluster` and `hpe_morpheus_tenant`: refreshed descriptions and currency validation (MORPH-18721)

`hpe_morpheus_cluster`'s `config_hvm.overlay_interface_name` now carries a full description explaining
when the overlay (VXLAN) Open vSwitch bridge interface is needed, which HVM layouts use it, and the
Morpheus 9.0.x build caveat. `hpe_morpheus_tenant`'s `currency` is validated at plan time against the
set of ISO 4217 currency codes Morpheus accepts, so an unsupported code is rejected before apply
rather than by the API.

### `hpe_morpheus_user`: `require_password_change` issues a temporary password (MORPH-18374)

A user can be created with a password they must change at their first sign-in to the Morpheus UI,
the way an administrator hands out a temporary password. With `require_password_change = true` the
provider marks the password expired whenever it sets it — on create, when the attribute is turned on
for an existing user (their current password expires), and again whenever `password_wo` changes —
and Morpheus redirects the user's next UI sign-in to the change-password page (signing in to the API
with the password is not affected). The user's own change clears the flag, which the read-only
`password_expired` reflects; `require_password_change` itself is not read back, so it never
re-expires a password the user has changed. Left out or `false`, nothing is sent and the behaviour
is as before.

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
- **The three SDKv2 resources moved too (MORPH-18860).** `hpe_morpheus_network_domain`
  (`domain_password_wo`), `hpe_morpheus_task_shell_script` and `hpe_morpheus_task_powershell_script`
  (`remote_target_password_wo`) are rebuilt on the plugin framework and use the same digest
  mechanism; their `*_wo_version` markers are deprecated like the others. Type names and attributes
  are unchanged, and existing state is upgraded in place: the first plan after upgrading shows one
  `~ write_only_revision` update (and, for a task that still sets the deprecated
  `remote_target_password`, that attribute — see below). Configurations written for the previous resources
  keep working: the deprecated `auto_join_domain` (a no-op), `domain_password` and
  `remote_target_password` are still accepted and honoured, and are removed in 2.6 with the markers.
  One difference for `remote_target_password`: the previous resource read the API's password *hash*
  back into it, which the plugin framework does not allow; it now holds the configured value, so a
  configuration that still sets it plans one update after upgrading (sending the configured
  password), and a configuration that has moved to `remote_target_password_wo` plans nothing. An
  unrelated update of a network domain no longer re-sends `domain_password_wo` — the provider sends
  the API's mask value instead, because Morpheus clears the stored password when the field is omitted
  on update (the deprecated `domain_password` is still sent on every update, as before). A task's
  `script_content` is still compared ignoring surrounding whitespace, so a heredoc with a trailing
  newline does not plan a change.
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

### `hpe_morpheus_resource_pool_group`: `tenant_ids` no longer plans a change on every run (MORPH-8207)

Morpheus always lists the tenant that owns a resource pool group in the group's `tenants`, and the
resource wrote that list straight into state. A configuration that did not list the owning tenant —
including one that left `tenant_ids` out altogether — planned `- <owner>` after every apply, and
applying it changed nothing. The owning tenant is now left out of state unless it is listed
explicitly, in which case it is kept; every other tenant Morpheus reports is kept and any it no
longer reports is dropped, so a tenant added or removed in the UI still shows as drift. Importing a
group yields the tenants it is shared with. The owning tenant is the one the provider is running
as, so the same holds for a sub-tenant. If the provider cannot determine it, the read keeps the
full list and warns that a plan may show the owner being removed.

### A resource that fails after Morpheus created it no longer buries the real error under "invalid result object" diagnostics (MORPH-18935)

When a create succeeded at the API but a later step failed — a cluster whose provisioning poll
hit the authentication failure above, for example — the provider marked the resource tainted so
the next apply would replace it, but left every computed attribute unknown in the state it
returned. Terraform rejects unknown values after an apply even when the provider reports an error,
so the practitioner saw a run of *"Provider returned invalid result object after apply …
attribute `cluster_type_code` is unknown"* (and `cpu_placement_mode`, `service_url`, `uuid`, …) on
top of the failure that mattered. The tainting helper used by every resource now nulls those
values first, so the output is the real error and the "partially created, marked tainted" warning,
and nothing else. The next plan still proposes the replacement as before.

### Username and password authentication no longer fails after 24 minutes (MORPH-18753)

With `username`/`password` provider configuration, any single plan or apply that ran for longer than
about 24 minutes — a cluster build, host onboarding, a virtual switch fan-out, anything that polls a
Morpheus process — failed part-way with *"could not authenticate with the Morpheus API using the
username … verify that the credentials are correct: 401"*, although the credentials were fine. The
access token the provider obtains at the start lives 24 minutes, and the provider did log in again
when a request was refused — but the new login carried the expired token in its own `Authorization`
header (the bundled SDK adds the cached header to every request, the token request included), and
the Morpheus token endpoint refuses any request that carries a bearer token, valid or expired (on
9.0 and 9.1 alike). The old token is now dropped before logging in again, concurrent requests
refused with the same token share a single new login instead of each logging in, and a request with
a body that is retried after the new login sends its body again. Access-token authentication was
never affected.

### A resource deleted outside Terraform is recreated instead of failing the refresh (MORPH-18584)

When an object this provider manages was deleted outside Terraform — through the Morpheus UI, by
another automation, or by cleanup — the next `terraform plan` or `apply` failed with an error such
as *"read cluster resource / cluster 4 GET failed: 404"* instead of noticing the object was gone.
The only way forward was `terraform state rm`, which practitioners should never need to reach for.
Twenty-five resources had this flaw: `hpe_morpheus_backup_host`, `hpe_morpheus_backup_instance`,
`hpe_morpheus_cloud`, `hpe_morpheus_cluster`, `hpe_morpheus_datastore`, `hpe_morpheus_group`,
`hpe_morpheus_image`, the load balancer resources (`hpe_morpheus_load_balancer`, its monitor, pool
and profile), the network resources (`hpe_morpheus_network`, its DHCP server, and the network
router with its BGP neighbour, firewall rule, NAT and route sub-resources), `hpe_morpheus_os_type`,
`hpe_morpheus_os_type_image`, `hpe_morpheus_policy`, `hpe_morpheus_role`, `hpe_morpheus_service_plan`,
`hpe_morpheus_task` and `hpe_morpheus_user`. Read now answers the API's 404 by removing the resource
from state and logging that it did, so the next plan proposes recreating it — the behaviour the
other resources already had. A resource vanishing immediately after create or update is still a hard
error, since that is a genuine failure rather than drift. A unit rule test now holds every framework
resource to this convention so the regression cannot return.

### `hpe_morpheus_cluster`: `config_hvm.overlay_interface_name` is documented (MORPH-18589)

The attribute had no description. It names the host interface used for the overlay (VXLAN) Open
vSwitch bridge; Morpheus treats it as optional and HVM 2.0 layouts do not use it, but on Morpheus
9.0.x with an HVM 1.x layout omitting it makes every host fail to build at the "HVM Host Build OVS
Interface" step. The documentation now says so and recommends setting it on 1.x layouts. The
attribute stays optional, as in the API.

### `hpe_morpheus_instance`: destroy deletes through the instance API, so the virtual machine is removed on Morpheus 8.1.x (MORPH-18513)

Since 2.0.0 the resource destroyed an instance by stopping and force-deleting its server(s) through
the server API, then deleting the instance. On Morpheus 8.1.x that path removes the Morpheus records
but can leave the virtual machine behind in the hypervisor: the appliance waits on the instance
removal before queuing the VM deletion, that wait times out after two hours, and the deletion is
never queued — while Terraform, seeing the records gone, reports the destroy complete. The destroy
now deletes the instance through the instance API, the path the Morpheus UI uses, which tears the
virtual machine down through the cloud's provisioning service, and then confirms the server(s) are
gone. Because Morpheus refuses an instance delete while the instance is still provisioning, resizing,
restarting or cloning — and honours a forced delete only for users with the *Provisioning: Force
Delete* permission — the provider waits for the instance to settle first. The server path is kept
as a fallback for an instance Morpheus still refuses to delete, or whose removal fails, and a warning
says when it was used. The acceptance-test sweeper follows the same order.

### `hpe_morpheus_instance`: a volume's `controller_mount_point` is now sent to Morpheus (MORPH-18424)

Setting `controller_mount_point` on a volume — for example `"-1:1:6:0"` from
`hpe_morpheus_instance_storage_controller`, to put a data disk on a new SCSI VMware Paravirtual
controller at bus 1 while the OS disk stays on the LSI Logic controller — had no effect: the value
was never included in the create or update request, so Morpheus placed the disk on the existing
controller and the apply failed with *"Provider produced inconsistent result after apply …
controller_mount_point: was "-1:1:6:0", but now "62383:0:13:1""*. The value is now sent, and the
controller is created. Because Morpheus identifies a controller by bus and type rather than by the
`-1` placeholder id, state keeps the configured `-1:…` form once the request has been satisfied, so
there is no diff on the next plan; a disk that has really moved to another controller is still
reported as drift.

### `hpe_morpheus_license`: the installed licence is read back, drift is detected, and destroy is documented (MORPH-18310)

The resource's read called `GET /api/license` and discarded the answer, so it never noticed when the
configured key was no longer the installed one: after an appliance rebuild (first-time setup installs
a trial licence) or a key installed through the UI, `terraform plan` said *No changes* and the only
remedy was a `-replace`. Morpheus reports each installed licence's key id (the key's first eight
characters) and the SHA-256 of its key, so the resource now matches the configured key against the
installed licences on every read and plans an in-place re-install when it is missing. The licence's
details are read back as new read-only attributes — `product_tier`, `account_name`, `start_date`,
`end_date`, `max_instances`, `max_memory`, `max_storage`, `hard_limit`, `free_trial`, `multi_tenant`,
`support_level` — and `id` is now the key id rather than the constant `1` (the key id stays the same
for as long as the same key is used, unlike Morpheus's internal record id, which changes each time a
key is re-installed); existing state is updated on the next refresh. Destroy uninstalls only the
licence matching `key` (it used to deactivate every installed licence), the documentation says so,
and gives `terraform state rm` as the way to stop managing a licence without uninstalling it. Import
takes the key id.

### `hpe_morpheus_cluster`: destroying a cluster waits for Morpheus to finish removing it (MORPH-18379)

Morpheus removes a cluster in a background job, and the resource meant to wait for that job before
reporting the destroy complete — but its wait ended on the first poll whether or not the cluster
was still there. A `-/+` replacement, or a `destroy` followed at once by an `apply`, therefore
raced the job and failed with *"Cluster name must be unique"*. The destroy now returns only once
the cluster can no longer be read, bounded by `timeouts.delete` (45 minutes by default). The
deletes of `hpe_morpheus_image`, `hpe_morpheus_os_type` and `hpe_morpheus_task` had the same wait;
Morpheus deletes those synchronously, so nothing changes for them in practice.

### `hpe_morpheus_setting_appliance`: the proxy bypass list can be managed (MORPH-18311)

The resource configured the appliance's proxy (`proxy_host`, `proxy_port`, …) but had no attribute
for its **No Proxy** list, so a proxied appliance sent its calls to lab-internal hosts — HVM hosts,
distributed workers, other appliances — through the proxy, and the only remedy was the UI, which
Terraform then neither read nor preserved. The new `no_proxy` attribute takes the comma-separated
hosts, domains or CIDR ranges the appliance reaches directly (for example
`10.0.0.0/8,localhost,127.0.0.1`) and is written and read like the other proxy settings.
### `hpe_morpheus_image`, `hpe_morpheus_datastore`, `hpe_morpheus_policy`, `hpe_morpheus_task`: an in-place update no longer replaces the resources that reference them (MORPH-18309)

These four resources planned their `id` as *known after apply* on every in-place update — the one
computed attribute that cannot change was missing the plan modifier the other resources carry. A
dependant that references the id through a create-only attribute (an instance's
`config_hvm.image_id`, for instance) saw it as unknown and was planned for **destroy and
recreate**: adding `user_data` to an image would have replaced the running VM booted from it. The
known id is now kept in the plan, so an image, datastore, policy or task update is just that.
### `hpe_morpheus_service_plan`: a plan created without `max_disks` no longer rejects every instance (MORPH-18341)

Creating a plan with `max_disks` (or `max_cpu`, `max_cores`, `sort_order`) left out sent the
attribute as `0`, because the create request tested only whether the value was null and an unset
optional-and-computed attribute is *unknown* at that point. Morpheus stored `maxDisks: 0`, and with
its default `minDisks` of 1 every instance provisioned on the plan failed with *"Number of volumes
must be between 1 and 0 for the selected plan"*. An attribute left out of the configuration is now
left out of the request, as the update already did, so Morpheus keeps it unset — the same as a plan
created in the UI. Plans already created with the bad value are repaired by setting `max_disks`
explicitly once.

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
