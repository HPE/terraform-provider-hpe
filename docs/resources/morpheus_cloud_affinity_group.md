---
page_title: "hpe_morpheus_cloud_affinity_group Resource - terraform-provider-hpe"
subcategory: "Morpheus"
description: |-
  
---
# hpe_morpheus_cloud_affinity_group (Resource)



~> **Terminology** Morpheus models these as separate things. An **instance** is the
unit you provision and manage. A **compute server** is a machine record backing it, and
may be either a guest virtual machine or a hypervisor host — both appear in the same
collection. A **host** is a compute server that hosts others: it has guests of its own
and no parent. Attributes named `compute_server_id`, `compute_servers` or `servers`
therefore refer to machine records, which in practice are usually guests; `parent_host_id`
names the hypervisor a guest runs on.

-> This resource creates a cloud-scoped affinity group, which is realised as a DRS rule on the cluster backing `pool`. To place an affinity group on an HPE Virtual Machines (HVM) cluster, use [`hpe_morpheus_cluster_affinity_group`](https://registry.terraform.io/providers/HPE/hpe/latest/docs/resources/morpheus_cluster_affinity_group) instead.

Morpheus enables affinity groups on the following cloud types:

| Cloud | Type code |
|---|---|
| VMware vCenter | `vmware` |
| VMware Cloud on AWS | `vmwareCloudAws` |
| MacStadium | `macstadium` |

-> The table records which cloud types enable affinity groups today; it is not a limit of this resource. The capability is enabled per cloud type, and a cloud contributed by a plugin can advertise it too — this resource drives the standard affinity group API rather than anything VMware-specific, so it works against such a cloud unchanged. A cloud type absent from the table may therefore still support affinity groups on a given appliance.

!> **Warning** Setting `servers` makes Terraform the sole owner of this affinity group's membership. Morpheus replaces the entire member set on every update, so any server added out-of-band — by provisioning an instance into the group, or, on VMware clouds, by vCenter DRS rule sync — is removed on the next apply. Omit `servers` entirely to manage the group's other attributes while leaving its membership untouched.

!> **Warning** On VMware clouds, Morpheus imports the vCenter DRS rules it finds on the cloud referenced by `cloud_id` as affinity groups of their own (`source = "sync"`). Those imported groups cannot be managed by this resource: every cloud sync rewrites their affinity type, active flag and membership from vCenter, so Terraform would fight the sync loop indefinitely. Use the [`hpe_morpheus_cloud_affinity_group`](https://registry.terraform.io/providers/HPE/hpe/latest/docs/data-sources/morpheus_cloud_affinity_group) data source to reference such a group instead.

-> There are two ways to manage a group's membership, and a group should use one or the other, never both. Set `servers` here to make this resource the owner of the membership: the group is created with those servers and every apply re-asserts exactly that set. Leave `servers` out to manage membership one server at a time with `hpe_morpheus_cloud_affinity_group_member`, which leaves servers added by other means alone.

-> Morpheus 9.1 refuses to create an affinity group with no members, so against 9.1 the group has to be created with `servers` set — where the members are instances managed in the same configuration, provision them first and point `servers` at their `compute_servers`. The provider checks the appliance version at plan time and refuses a plan that would create a memberless group on 9.1 or later. To use the member resource on a 9.1 group instead, bootstrap it: create the group with its first member in `servers`, then remove `servers` from the configuration — that plan is empty (the attribute is Optional+Computed, so dropping it changes nothing) and hands the membership to `hpe_morpheus_cloud_affinity_group_member` without evicting anyone. Removing the last member of a 9.1 group is refused by Morpheus; the provider leaves the server in place with a warning, and destroying the group removes it.

-> On a VMware cloud, do not rename a cloud affinity group that will later be destroyed: Morpheus 9.1 deletes the vCenter DRS rule by the group's current name while the rule keeps its original one, so the rule is left behind and later `KEEP_TOGETHER` groups containing the same VMs are refused by vCenter (Morpheus core defect MORPH-20098).

-> When `servers` is not set it shows as "known after apply" in a plan that changes anything else about the group, even when the membership has not changed. Membership can change outside Terraform, so the provider deliberately does not assume the set it last recorded is still current. Where nothing has actually changed, the apply is a no-op.

~> **Note** `pool_id` is effectively required. Morpheus rejects the creation of a cloud affinity group that has no pool, so although the schema marks `pool_id` optional, omitting it fails at apply time with `No pool specified for affinity group`. The pool must be a resource pool of type Cluster on the cloud referenced by `cloud_id`, because the group is realised as a DRS rule on that cluster.

-> Morpheus ignores `affinity_type` and `pool_id` on update, so a change to either replaces the affinity group.

-> Note that Morpheus version `8.0.10` or later is required for affinity group support.

-> An instance can alternatively be placed into a group when it is provisioned, with `config_vmware.affinity_group_id` on [`hpe_morpheus_instance`](https://registry.terraform.io/providers/HPE/hpe/latest/docs/resources/morpheus_instance). Do not mix the two approaches for the same instance: where `servers` is managed here it is authoritative, so a server this resource does not list — including one placed into the group at provision time — is removed on the next apply. `affinity_group_id` records only what was requested at creation, so `servers` here remains the source of truth for actual membership.

## Example Usage

```terraform
resource "hpe_morpheus_cloud_affinity_group" "example" {
  cloud_id      = 1
  name          = "Example Affinity Group"
  affinity_type = "KEEP_TOGETHER"

  # Morpheus requires a pool on a cloud affinity group, and it must be a
  # resource pool of type Cluster.
  pool_id = 1
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `cloud_id` (Number) The ID of the cloud.
- `name` (String) The name of the affinity group.

### Optional

- `active` (Boolean) Whether the affinity group is active.
- `affinity_type` (String) The affinity type. Valid values are KEEP_TOGETHER and KEEP_SEPARATE.
- `pool_id` (Number) The ID of the resource pool to scope the affinity group to.
- `resource_permissions` (Attributes) Resource permissions for group access. (see [below for nested schema](#nestedatt--resource_permissions))
- `servers` (Set of Number) The compute servers in this affinity group, as the API reports them.
Set it to make this resource the owner of the group's membership: the
group is created with these members and every apply re-asserts exactly
this set, so a server added by other means -- an instance provisioned
into the group, a node added to one -- plans a correction. Leave it out
to manage membership one server at a time with
hpe_morpheus_cloud_affinity_group_member (never both on one group).
Morpheus 9.1 refuses to create an affinity group with no members: a
VM-to-VM rule needs this set, and the other shape it accepts -- a
VM-to-host rule from a VM group and a host group -- is read-only on
this resource, so against 9.1 a group can only be created with this set
(the provider refuses such a plan up front on 9.1 and later).
- `tenant_ids` (Set of Number, Deprecated) Deprecated: has no effect. The Morpheus API rejects tenant assignment on affinity groups and does not apply it, so the value is retained in Terraform state only and does not reflect the appliance. Tracked as MORPH-15806.
- `visibility` (String) The visibility of the affinity group (public or private).

### Read-Only

- `host_group` (Attributes) The Host Group side of a combined VM-to-Host rule. Must be sent together with `vmGroup`, and `servers` must be empty. (see [below for nested schema](#nestedatt--host_group))
- `id` (Number) The ID of the cloud affinity group.
- `owner` (Attributes) Owning tenant account (see [below for nested schema](#nestedatt--owner))
- `source` (String) The source of the affinity group (e.g. user, sync).
- `vm_group` (Attributes) The VM Group side of a combined VM-to-Host rule. Must be sent together with `hostGroup`, and `servers` must be empty. Send `null` for both group fields to clear a combined rule. (see [below for nested schema](#nestedatt--vm_group))

<a id="nestedatt--resource_permissions"></a>
### Nested Schema for `resource_permissions`

Optional:

- `all` (Boolean) Pass true to allow access to all groups.
- `groups` (Attributes Set) Set of groups allowed access. (see [below for nested schema](#nestedatt--resource_permissions--groups))

<a id="nestedatt--resource_permissions--groups"></a>
### Nested Schema for `resource_permissions.groups`

Required:

- `id` (Number) Group ID.

Optional:

- `default` (Boolean) Whether this is the default group.



<a id="nestedatt--host_group"></a>
### Nested Schema for `host_group`

Read-Only:

- `id` (Number)
- `name` (String)


<a id="nestedatt--owner"></a>
### Nested Schema for `owner`

Read-Only:

- `id` (Number)


<a id="nestedatt--vm_group"></a>
### Nested Schema for `vm_group`

Read-Only:

- `id` (Number)
- `name` (String)

## Import

Import is supported using `cloud_id.affinity_group_id`.

```bash
terraform import hpe_morpheus_cloud_affinity_group.example 42.123
```
