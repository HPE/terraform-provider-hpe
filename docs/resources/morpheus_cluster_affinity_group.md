---
page_title: "hpe_morpheus_cluster_affinity_group Resource - terraform-provider-hpe"
subcategory: "Morpheus"
description: |-
  
---
# hpe_morpheus_cluster_affinity_group (Resource)



~> **Terminology** Morpheus models these as separate things. An **instance** is the
unit you provision and manage. A **compute server** is a machine record backing it, and
may be either a guest virtual machine or a hypervisor host — both appear in the same
collection. A **host** is a compute server that hosts others: it has guests of its own
and no parent. Attributes named `compute_server_id`, `compute_servers` or `servers`
therefore refer to machine records, which in practice are usually guests; `parent_host_id`
names the hypervisor a guest runs on.

-> This resource creates a cluster-scoped affinity group, which is supported on HPE Virtual Machines (HVM) clusters. To place an affinity group on a VMware, VMware Cloud on AWS or MacStadium cloud, use [`hpe_morpheus_cloud_affinity_group`](https://registry.terraform.io/providers/HPE/hpe/latest/docs/resources/morpheus_cloud_affinity_group) instead.

-> HVM is the only cluster type on which Morpheus enables affinity groups today, but that is not a limit of this resource: the capability is enabled per cluster type, and a cluster type contributed by a plugin can advertise it too, in which case this resource works against it unchanged.

-> There are two ways to manage a group's membership, and a group should use one or the other, never both. Set `servers` here to make this resource the owner of the membership: the group is created with those servers, and every apply re-asserts exactly that set. Or omit `servers` and manage the membership one server at a time with [`hpe_morpheus_cluster_affinity_group_member`](https://registry.terraform.io/providers/HPE/hpe/latest/docs/resources/morpheus_cluster_affinity_group_member), which leaves the rest of the group's members alone.

!> **Warning** Setting `servers` makes Terraform the sole owner of this affinity group's membership. Morpheus replaces the entire member set on every update, so any server added out-of-band — for example by provisioning an instance into the group — is removed on the next apply. Omit `servers` to manage the group's other attributes while leaving its membership untouched: removing it from a configuration that had it does not empty the group, the members it has at the time are kept.

-> Morpheus 9.1 refuses to create an affinity group with no members. A VM-to-VM rule needs `servers`; the other shape Morpheus accepts, a VM-to-host rule from a VM group and a host group, is read-only on this resource (`vm_group`, `host_group`), so against 9.1 the group has to be created with `servers` set — see the second example below. Where the members are instances managed in the same configuration, provision them first and point `servers` at their `compute_servers`.

~> **Note** Affinity rules on HVM are enforced by the platform's own placement loop, not applied at the moment you create or change them. That loop is the cluster's Dynamic Placement and **must be enabled on the cluster** — `config_hvm.dynamic_placement = true` on [`hpe_morpheus_cluster`](https://registry.terraform.io/providers/HPE/hpe/latest/docs/resources/morpheus_cluster). It is off by default, and with it off the rule and its membership are recorded but no virtual machine is ever placed or moved to satisfy them. With it enabled, creating a group, changing its members, or changing `affinity_type` records the intent straight away, but the affected virtual machines are only migrated when that loop next runs and finds hosts that can satisfy the rule. A `terraform apply` therefore returns once Morpheus has recorded the rule — before the virtual machines have necessarily moved to satisfy it — so do not expect placement to have converged the instant the apply completes.

-> When `servers` is not set it shows as "known after apply" in a plan that changes anything else about the group, even when the membership has not changed. Membership can change outside Terraform, so the provider deliberately does not assume the set it last recorded is still current. Where nothing has actually changed, the apply is a no-op.

-> Morpheus ignores `affinity_type` on update, so a change to it replaces the affinity group.

-> `pool_id` is assigned by Morpheus and cannot be set. It is read-only on this resource.

-> Note that Morpheus version `8.0.10` or later is required for affinity group support.

-> An instance can alternatively be placed into a group when it is provisioned, with `config_hvm.affinity_group_id` on [`hpe_morpheus_instance`](https://registry.terraform.io/providers/HPE/hpe/latest/docs/resources/morpheus_instance). Do not mix the two approaches for the same instance: where `servers` is managed here it is authoritative, so a server this resource does not list — including one placed into the group at provision time — is removed on the next apply. `affinity_group_id` records only what was requested at creation, so `servers` here remains the source of truth for actual membership.

## Example Usage

A group whose membership is managed elsewhere, with `hpe_morpheus_cluster_affinity_group_member` or by provisioning instances into it:

```terraform
resource "hpe_morpheus_cluster_affinity_group" "example" {
  cluster_id    = 1
  name          = "Example Affinity Group"
  affinity_type = "KEEP_TOGETHER"
}
```

A group that owns its membership. The servers are looked up by name; an instance managed in the same configuration can be a member through its `compute_servers` instead (`servers = hpe_morpheus_instance.web.compute_servers`):

```terraform
data "hpe_morpheus_cluster" "hvm" {
  name = "hvm-cluster"
}

# The members, looked up by name. Each is a compute server -- the machine
# record backing an instance -- of that cluster.
data "hpe_morpheus_compute_server" "member" {
  for_each = toset(["web-1", "web-2"])

  name = each.value
}

resource "hpe_morpheus_cluster_affinity_group" "members" {
  cluster_id    = data.hpe_morpheus_cluster.hvm.id
  name          = "web-keep-separate"
  affinity_type = "KEEP_SEPARATE"

  # Terraform owns the membership: the group is created with these servers and
  # every apply re-asserts exactly this set. Do not also manage members of this
  # group with hpe_morpheus_cluster_affinity_group_member.
  servers = [for s in data.hpe_morpheus_compute_server.member : s.id]
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `cluster_id` (Number) The ID of the cluster.
- `name` (String) The name of the affinity group.

### Optional

- `active` (Boolean) Whether the affinity group is active.
- `affinity_type` (String) The affinity type. Valid values are KEEP_TOGETHER and KEEP_SEPARATE.
- `description` (String, Deprecated) Deprecated: not backed by the Morpheus API. The value is retained in Terraform state only and will be removed in a future release.
- `resource_permissions` (Attributes) Resource permissions for group access. (see [below for nested schema](#nestedatt--resource_permissions))
- `servers` (Set of Number) The compute servers in this affinity group, as the API reports them.
Set it to make this resource the owner of the group's membership: the
group is created with these members and every apply re-asserts exactly
this set, so a server added by other means -- an instance provisioned
into the group, a node added to one -- plans a correction. Leave it out
to manage membership one server at a time with
hpe_morpheus_cluster_affinity_group_member (never both on one group).
Morpheus 9.1 refuses to create an affinity group with no members: a
VM-to-VM rule needs this set, and the other shape it accepts -- a
VM-to-host rule from a VM group and a host group -- is read-only on
this resource, so against 9.1 a group can only be created with this set.
- `tenant_ids` (Set of Number, Deprecated) Deprecated: has no effect. The Morpheus API rejects tenant assignment on affinity groups and does not apply it, so the value is retained in Terraform state only and does not reflect the appliance. Tracked as MORPH-15806.
- `visibility` (String) The visibility of the affinity group (public or private).

### Read-Only

- `host_group` (Attributes) The Host Group side of a combined VM-to-Host rule. Must be sent together with `vmGroup`, and `servers` must be empty. (see [below for nested schema](#nestedatt--host_group))
- `id` (Number) The ID of the cluster affinity group.
- `owner` (Attributes) Owning tenant account (see [below for nested schema](#nestedatt--owner))
- `pool_id` (Number) The ID of the resource pool associated with the affinity group.
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

Import is supported using `cluster_id.affinity_group_id`.

```bash
terraform import hpe_morpheus_cluster_affinity_group.example 42.123
```
