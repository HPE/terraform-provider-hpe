---
page_title: "Affinity Groups"
subcategory: "Guides"
---

## What This Guide Covers

Affinity groups tell Morpheus to keep a set of guests on the same hypervisor, or
to keep them apart. This guide covers how to express that in Terraform with the
HPE provider: creating a group, getting guests into it, and proving where they
actually landed.

It also covers what changed in v2.0.0, because membership works differently from
v1.6.0 and existing configuration will need editing.

- [The two group types](#the-two-group-types)
- [Membership is a resource](#membership-is-a-resource)
- [Two ways to join a group](#two-ways-to-join-a-group)
- [Ordering matters](#ordering-matters)
- [Proving placement](#proving-placement)
- [Platform differences](#platform-differences)
- [What forces replacement](#what-forces-replacement)
- [Migrating from v1.6.0](#migrating-from-v160)
- [tenant_ids is deprecated](#tenant_ids-is-deprecated)

---

## The Two Group Types

Affinity groups belong either to a cloud or to a cluster, and the provider has a
separate resource for each. They are not interchangeable — pick the one that
matches where your guests live.

| Resource | Belongs to | Typical use |
|---|---|---|
| `hpe_morpheus_cloud_affinity_group` | a cloud | VMware |
| `hpe_morpheus_cluster_affinity_group` | a cluster | HVM |

```hcl
resource "hpe_morpheus_cluster_affinity_group" "web" {
  cluster_id    = var.cluster_id
  name          = "web-tier"
  affinity_type = "KEEP_TOGETHER"
  active        = true
}
```

```hcl
resource "hpe_morpheus_cloud_affinity_group" "web" {
  cloud_id      = var.cloud_id
  pool_id       = var.pool_id
  name          = "web-tier"
  affinity_type = "KEEP_SEPARATE"
  active        = true
}
```

`affinity_type` is `KEEP_TOGETHER` or `KEEP_SEPARATE`.

---

## Membership Is a Resource

The group's `servers` attribute is **read-only**. It reports who is in the group;
it does not set it. To put a guest in a group, declare a membership resource —
one per member.

```hcl
resource "hpe_morpheus_cluster_affinity_group_member" "app" {
  cluster_id        = var.cluster_id
  affinity_group_id = hpe_morpheus_cluster_affinity_group.web.id
  server_id         = one(hpe_morpheus_instance.app.compute_servers)
}
```

The cloud equivalent is `hpe_morpheus_cloud_affinity_group_member`, taking
`cloud_id` instead of `cluster_id`.

Note `server_id` is a **compute server**, not an instance. An instance may have
more than one; `one()` is correct for a single-node instance, otherwise index the
list explicitly.

### Why membership is separate

The API treats a supplied server list as authoritative: anything absent from it
is removed from the group. A group that named its own members would therefore
evict any guest that joined by another route — an instance provisioned directly
into the group, or a node added to a cluster — on the next apply, silently.

One resource per member avoids that. Terraform manages the memberships you
declared and leaves everything else alone.

---

## Two Ways to Join a Group

**After provisioning**, with a membership resource, as above. This works on any
existing guest.

**At provision time**, with `affinity_group_id` on the instance's typed config
block. Morpheus places the guest against the group's existing membership as it
builds it.

```hcl
resource "hpe_morpheus_instance" "app" {
  name             = "app-01"
  cloud_id         = var.cloud_id
  group_id         = var.group_id
  instance_type_id = var.instance_type_id
  layout_id        = var.layout_id
  plan_id          = var.plan_id

  network_interfaces = [{ network_id = var.network_id }]

  config_hvm = {
    resource_pool_id  = var.resource_pool_id
    affinity_group_id = hpe_morpheus_cluster_affinity_group.web.id
  }
}
```

For VMware the typed block is `config_vmware`, which also takes
`vmware_folder_id`:

```hcl
resource "hpe_morpheus_instance" "app" {
  name             = "app-01"
  cloud_id         = var.cloud_id
  group_id         = var.group_id
  instance_type_id = var.instance_type_id
  layout_id        = var.layout_id
  plan_id          = var.plan_id

  network_interfaces = [{ network_id = var.network_id }]

  config_vmware = {
    resource_pool_id  = var.resource_pool_id
    vmware_folder_id  = var.vmware_folder_id
    affinity_group_id = hpe_morpheus_cloud_affinity_group.web.id
  }
}
```

Do not set `kvm_host_id` alongside `affinity_group_id`. The provider rejects the
combination, because Morpheus applies the pinned host and then lets the affinity
group override it, discarding the pin without warning.

---

## Ordering Matters

If you use `affinity_group_id`, and the group's other members are created in the
same configuration, **add an explicit `depends_on`**.

```hcl
resource "hpe_morpheus_cluster_affinity_group_member" "seed" {
  cluster_id        = var.cluster_id
  affinity_group_id = hpe_morpheus_cluster_affinity_group.web.id
  server_id         = one(hpe_morpheus_instance.seed.compute_servers)
}

resource "hpe_morpheus_instance" "app" {
  # ...

  config_hvm = {
    resource_pool_id  = var.resource_pool_id
    affinity_group_id = hpe_morpheus_cluster_affinity_group.web.id
  }

  # Without this, Terraform may build this instance before the seed's
  # membership exists. An empty group gives Morpheus nothing to place the
  # guest with or away from, and the result is down to luck.
  depends_on = [hpe_morpheus_cluster_affinity_group_member.seed]
}
```

The instance references the *group*, so Terraform knows the group must exist
first — but nothing tells it the group must be **populated** first. Membership is
a sibling resource, not a dependency of the instance, so Terraform is free to
create them in either order.

This applies to **both** platforms. On VMware an empty group is worse than
useless: provisioning into one can fail at power on with `Error running vm`, and
the failed instance can be awkward to remove.

A group with one member is equally pointless for placement — there is nothing to
be placed with or away from — so seed the group before the guest whose placement
you care about.

---

## Proving Placement

Membership is recorded whether or not anything was placed, so `servers` does not
tell you the rule was honoured. The hypervisor does.

```hcl
data "hpe_morpheus_compute_server" "app" {
  id         = one(hpe_morpheus_instance.app.compute_servers)
  depends_on = [hpe_morpheus_instance.app]
}

data "hpe_morpheus_compute_server" "seed" {
  id         = one(hpe_morpheus_instance.seed.compute_servers)
  depends_on = [hpe_morpheus_instance.app]
}

output "placement" {
  value = (
    data.hpe_morpheus_compute_server.app.parent_host_name ==
    data.hpe_morpheus_compute_server.seed.parent_host_name
    ? "together on ${data.hpe_morpheus_compute_server.app.parent_host_name}"
    : "separate: ${data.hpe_morpheus_compute_server.app.parent_host_name} vs ${data.hpe_morpheus_compute_server.seed.parent_host_name}"
  )
}
```

`parent_host_name` is the hypervisor the guest is running on. The `depends_on`
keeps the reads until after the instance exists.

---

## Platform Differences

The two platforms apply rules at different times, which changes what you see
immediately after `terraform apply`.

| | VMware | HVM |
|---|---|---|
| `KEEP_TOGETHER` | at placement | at placement |
| `KEEP_SEPARATE` | at placement | **after placement**, by the cluster |
| Mechanism | DRS rule | placement loop, appears as a `Move Server` process |

On HVM, `KEEP_SEPARATE` is not applied as the guest is built. The cluster
notices the violated rule afterwards and moves a guest, which has been observed
to take anywhere from about ninety seconds to six minutes. An apply that returns
with both guests still co-located has not failed — check again shortly.

HVM also requires the cluster to have dynamic placement enabled. With it off,
membership is recorded but nothing is placed or moved:

```
GET /api/clusters/<id>  ->  .cluster.config.dynamicPlacementMode
```

A cluster needs at least two hypervisors for `KEEP_SEPARATE` to mean anything.

---

## What Forces Replacement

Two attributes cannot be changed in place.

| Attribute | Changing it |
|---|---|
| `affinity_type` on the group | replaces the group |
| `affinity_group_id` on the instance | **replaces the instance** |

The second is the one to watch: editing or removing `affinity_group_id` on a
running instance destroys and rebuilds it.

To move guests from one rule to another without rebuilding them, leave the
original group in place, deactivate it, and add a second group with membership
resources:

```hcl
resource "hpe_morpheus_cluster_affinity_group" "together" {
  cluster_id    = var.cluster_id
  name          = "web-together"
  affinity_type = "KEEP_TOGETHER"

  # Deactivated rather than deleted: an instance still names it in
  # affinity_group_id, and removing it would rebuild that instance.
  active = false
}

resource "hpe_morpheus_cluster_affinity_group" "apart" {
  cluster_id    = var.cluster_id
  name          = "web-apart"
  affinity_type = "KEEP_SEPARATE"
  active        = true
}

resource "hpe_morpheus_cluster_affinity_group_member" "app" {
  cluster_id        = var.cluster_id
  affinity_group_id = hpe_morpheus_cluster_affinity_group.apart.id
  server_id         = one(hpe_morpheus_instance.app.compute_servers)
}
```

`active` updates in place.

---

## Migrating From v1.6.0

In v1.6.0 the group named its own members, usually with `ignore_changes` to stop
Terraform removing guests that joined by other routes:

```hcl
# v1.6.0 -- no longer valid
resource "hpe_morpheus_cluster_affinity_group" "web" {
  cluster_id    = var.cluster_id
  name          = "web-tier"
  affinity_type = "KEEP_TOGETHER"

  servers = hpe_morpheus_instance.seed.compute_servers

  lifecycle {
    ignore_changes = [servers]
  }
}
```

`servers` is read-only in v2.0.0, so this configuration will not apply. Replace
it with a membership resource per guest, and drop the `lifecycle` block — the
behaviour it worked around is gone:

```hcl
resource "hpe_morpheus_cluster_affinity_group" "web" {
  cluster_id    = var.cluster_id
  name          = "web-tier"
  affinity_type = "KEEP_TOGETHER"
  active        = true
}

resource "hpe_morpheus_cluster_affinity_group_member" "seed" {
  cluster_id        = var.cluster_id
  affinity_group_id = hpe_morpheus_cluster_affinity_group.web.id
  server_id         = one(hpe_morpheus_instance.seed.compute_servers)
}
```

If any instance in the same configuration uses `affinity_group_id`, add the
`depends_on` described in [Ordering matters](#ordering-matters). The old
configuration got that ordering implicitly, because the instance referenced a
group whose `servers` referenced the seed. Splitting membership out breaks that
chain, and without the explicit dependency placement becomes a race.

---

## `tenant_ids` Is Deprecated

`tenant_ids` on both affinity group resources is deprecated and **has no
effect**. The Morpheus API rejects tenant assignment on affinity groups and does
not apply it, so the provider no longer sends it.

Because the attribute is optional and computed, a value set in configuration is
retained in Terraform state — so state will show tenants that do not exist on the
appliance. Do not rely on it. Remove it from configuration.

Tenant access to an affinity group has to be arranged outside Terraform until the
appliance defect is resolved.
