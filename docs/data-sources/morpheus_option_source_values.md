---
page_title: "hpe_morpheus_option_source_values Data Source - terraform-provider-hpe"
subcategory: "Morpheus"
description: |-
  Lists the values of a Morpheus option source, optionally narrowed by filter blocks. Morpheus option sources are the lookups behind the UI's dropdowns — the resource pools and networks the provisioning wizard offers, an SCVMM library share, the members of a select list — served at GET /api/options/{source}.
---
# hpe_morpheus_option_source_values (Data Source)

Lists the values of a Morpheus option source, optionally narrowed by filter blocks. Morpheus option sources are the lookups behind the UI's dropdowns — the resource pools and networks the provisioning wizard offers, an SCVMM library share, the members of a select list — served at `GET /api/options/{source}`.

Use this data source to see **everything** a source offers, to iterate over it, or to check a value
against it in a `precondition`. To require exactly one value — and fail the plan otherwise — use
[`hpe_morpheus_option_source_value`](https://registry.terraform.io/providers/HPE/hpe/latest/docs/data-sources/morpheus_option_source_value).

## What an option source is

Every dropdown in the Morpheus UI that is not a fixed list is populated by an *option source*: a
named server-side lookup that returns `{name, value}` pairs for the current context. Some read the
Morpheus database (groups, service plans); some compute a provisioning context (the resource pools
and networks a wizard step may choose from, given a cloud, layout and pool); some ask an external
system (the library shares on an SCVMM server, given credentials). They are served over the API at
`GET /api/options/{source}` and `GET /api/options/{service_type}/{source}`, take their parameters
from the query string, and return a list of options.

Until now nothing in the provider could read them. This data source reads the same lookup the UI
uses and returns every option, with the source's own attributes alongside, so a configuration can
discover, iterate and verify without a value ever being looked up by hand.

## Use case: checking a network before provisioning

When an instance is created, Morpheus checks that each requested network is selectable for the
instance's cloud, group and resource pool, and rejects the request with a bare `Invalid network`
when one is not. The usual cause is a mismatch between the network and the resource pool — for
example a pool created by hand on an HVM cloud, whose networks all belong to the cluster's own pool.
The provider now explains that error when it happens; this data source lets a configuration avoid it.

The `zoneNetworkOptions` option source is the very check Morpheus applies, exposed for reading: it
lists the networks that are selectable for a cloud, layout and pool. Reading it into a `precondition`
turns a runtime rejection into a plan-time message that says what would have worked:

```terraform
# Checking a network before provisioning: zoneNetworkOptions is the check
# Morpheus applies at create time, read here into a precondition.

data "hpe_morpheus_cloud" "hvm" {
  name = "my-hvm-cloud"
}

data "hpe_morpheus_instance_type_layout" "hvm" {
  name    = "Single KVM VM"
  version = "22.04"
}

data "hpe_morpheus_instance_type" "ubuntu" {
  name = "Ubuntu"
}

data "hpe_morpheus_service_plan" "small" {
  name                = "1 CPU, 1GB Memory"
  provision_type_code = "kvm"
}

# The cluster's pool (see hpe_morpheus_option_source_value).
data "hpe_morpheus_option_source_value" "cluster_pool" {
  source    = "zonePools"
  cloud_id  = data.hpe_morpheus_cloud.hvm.id
  layout_id = data.hpe_morpheus_instance_type_layout.hvm.id
  name      = "hvm-cluster-01"
}

# The networks selectable for that pool. zoneNetworkOptions returns an object of
# lists, so collection names the one to read. Each value is the network id in
# the form the API takes on input: "network-<id>".
data "hpe_morpheus_option_source_values" "eligible_networks" {
  source     = "zoneNetworkOptions"
  cloud_id   = data.hpe_morpheus_cloud.hvm.id
  layout_id  = data.hpe_morpheus_instance_type_layout.hvm.id
  pool_id    = data.hpe_morpheus_option_source_value.cluster_pool.value
  collection = "networks"
}

resource "hpe_morpheus_instance" "vm" {
  name             = "web-01"
  cloud_id         = data.hpe_morpheus_cloud.hvm.id
  group_id         = 1
  instance_type_id = data.hpe_morpheus_instance_type.ubuntu.id
  layout_id        = data.hpe_morpheus_instance_type_layout.hvm.id
  plan_id          = data.hpe_morpheus_service_plan.small.id

  network_interfaces = [
    { network_id = var.network_id }
  ]

  config_hvm = {
    resource_pool_id = data.hpe_morpheus_option_source_value.cluster_pool.value
  }

  lifecycle {
    precondition {
      condition = contains(
        data.hpe_morpheus_option_source_values.eligible_networks.values[*].value,
        "network-${var.network_id}",
      )
      error_message = join(" ", [
        "Network ${var.network_id} is not selectable for ${data.hpe_morpheus_option_source_value.cluster_pool.value}",
        "(${data.hpe_morpheus_option_source_value.cluster_pool.name}) in this cloud. Selectable networks:",
        join(", ", data.hpe_morpheus_option_source_values.eligible_networks.values[*].name),
      ])
    }
  }
}
```

`zoneNetworkOptions` returns an object of lists rather than a single list, so `collection` names
the one to read — `networks` here; `networkGroups`, `networkSubnets` and `networkTypes` are the
others. Its items carry the network id in the form the API takes on input, `network-<id>`, so the
comparison is against `"network-${var.network_id}"`.

## Filtering

Filter blocks narrow the list client-side after the source has been read. Blocks are ANDed; the
values within one block are ORed; values are Go regular expressions and are case-sensitive unless the
pattern says otherwise (`(?i)`). A block's `name` is the field to match: `name`, `value`, or
`attributes.<key>` for a source-specific attribute — a value that does not carry the attribute does
not match, so filtering on an attribute is also a way to keep only the values that have it.

```terraform
# Narrowing a list with filter blocks: blocks are ANDed, values within a block
# are ORed, and values are Go regular expressions.

data "hpe_morpheus_cloud" "hvm" {
  name = "my-hvm-cloud"
}

data "hpe_morpheus_instance_type_layout" "hvm" {
  name    = "Single KVM VM"
  version = "22.04"
}

# The HVM cluster pools whose name begins "prod-". attributes.<key> filters on
# a source-specific attribute; a value that lacks it does not match.
data "hpe_morpheus_option_source_values" "prod_pools" {
  source    = "zonePools"
  cloud_id  = data.hpe_morpheus_cloud.hvm.id
  layout_id = data.hpe_morpheus_instance_type_layout.hvm.id

  filter {
    name   = "attributes.providerType"
    values = ["^mvm$"]
  }

  filter {
    name   = "name"
    values = ["^prod-"]
  }
}

# The networks in this cloud whose name mentions management, either spelling.
data "hpe_morpheus_option_source_values" "management_networks" {
  source   = "networks"
  cloud_id = data.hpe_morpheus_cloud.hvm.id

  filter {
    name   = "name"
    values = ["(?i)management", "(?i)mgmt"]
  }
}

output "prod_pools" {
  value = { for v in data.hpe_morpheus_option_source_values.prod_pools.values : v.name => v.value }
}

output "management_network_ids" {
  value = data.hpe_morpheus_option_source_values.management_networks.values[*].value
}
```

For narrowing to a single value, prefer `hpe_morpheus_option_source_value` with `name` or `value`:
it fails the plan when the selection is not unique, where a filter that happens to match one value
today will silently match two tomorrow.

## Use case: discovering what a source offers

An SCVMM cloud needs `libraryShare`, a UNC path that only the VMM server can supply. Before a
configuration commits to one with the singular data source, the plural lists what there is — and the
same list drives `for_each` when one resource per value is wanted. The credentials are stored once in
Morpheus and referenced by id; they do not pass through the data source.

```terraform
# Discovering what a source offers: every library share on an SCVMM server,
# through a stored credential — and the same list driving for_each.

resource "hpe_morpheus_credential" "scvmm" {
  name     = "scvmm-prod"
  type     = "username-password"
  username = var.scvmm_username
  password = var.scvmm_password
}

data "hpe_morpheus_option_source_values" "library_shares" {
  source        = "scvmmLibraryShares"
  service_type  = "scvmm"
  credential_id = hpe_morpheus_credential.scvmm.id
  parameters    = { host = var.scvmm_host }
}

output "scvmm_library_shares" {
  value = data.hpe_morpheus_option_source_values.library_shares.values[*].value
}

# One cloud per share.
resource "hpe_morpheus_cloud" "per_share" {
  for_each = { for v in data.hpe_morpheus_option_source_values.library_shares.values : v.value => v }

  name            = "scvmm-${replace(each.key, "/[^A-Za-z0-9]+/", "-")}"
  tenant_id       = 1
  cloud_type_code = "scvmm"

  config = {
    host         = var.scvmm_host
    username     = var.scvmm_username
    password     = var.scvmm_password
    workingPath  = "c:\\Temp"
    diskPath     = "c:\\VirtualDisks"
    libraryShare = each.key
  }
}
```

## Reading any source

Any source can be read with no more than its name. Parameters the typed attributes do not cover go
in `parameters`, under the name the API uses; the full objects come back, so a source's own
attributes are visible under `attributes` (as strings — numbers and booleans in their textual form,
nested objects as JSON).

```terraform
# Reading any source: by name alone, with extra parameters under the names the
# API uses, and the source's own attributes visible under attributes.

data "hpe_morpheus_option_source_values" "groups" {
  source = "groups"
}

output "groups" {
  value = { for v in data.hpe_morpheus_option_source_values.groups.values : v.name => v.value }
}

# Another list from the object-shaped zoneNetworkOptions: the network interface
# types the wizard offers for a cloud and layout. Items here have an id but no
# value, so value is the id; code arrives under attributes.
data "hpe_morpheus_option_source_values" "network_types" {
  source     = "zoneNetworkOptions"
  cloud_id   = 1
  layout_id  = 32
  collection = "networkTypes"
}

output "network_types" {
  value = { for v in data.hpe_morpheus_option_source_values.network_types.values : v.name => v.attributes["code"] }
}

# A parameter the typed attributes do not cover goes in parameters, under the
# name the API uses. By default it is sent both in the query string and as
# config[<key>], because the server ignores a parameter in the wrong place.
data "hpe_morpheus_option_source_values" "plans_for_layout" {
  source     = "servicePlans"
  cloud_id   = 1
  parameters = { layoutId = "32", provisionTypeCode = "kvm" }
}
```

## Finding a source and its parameters

The sources above are the ones most configurations need. For any other dropdown, the Morpheus UI
itself shows the source and parameters it uses: open the browser's developer tools on the Network tab,
operate the dropdown, and look for the request to `/options/<source>?...`. The path gives `source`
(and, in the two-segment form `/options/<service_type>/<source>`, the `service_type`); the query string
gives the parameters, under the names to use in `parameters` or the typed attributes below. The UI and
the API call the same server-side lookup with the same parameters.

The provider carries a catalogue of the sources known for a Morpheus release. A `source` not in it
produces a **warning**, not an error — plugins and newer appliances add sources — and the warning
names the release the catalogue was built from.

## Parameters

Morpheus reads a source's parameters from two places, the query string (`zoneId=1`) and a
`config[<key>]` form (`config[zoneId]=1`), and different sources read different places. A parameter
sent in the wrong place is **silently ignored** — the source runs as if it had not been given, which
for a filter means an unfiltered list. By default (`parameter_style = "auto"`) the provider therefore
sends every parameter in both places; every source tested tolerates this. `parameter_style` is the
escape hatch for a source that does not: `top_level` and `config` send each parameter in one place
only.

The typed attributes are shortcuts for the common parameters, named as the rest of this provider
names things rather than as the API does: the API says *zone* where the provider says *cloud*, so the
attribute is `cloud_id` and the translation to `zoneId` happens here. Each attribute sets the key(s)
the sources read, so `cloud_id` is written once whether a source reads `zoneId`, `cloudId` or
`config[zoneId]`, and the attributes that Morpheus reads under more than one name are sent under
every name:

| Attribute | Sent as | Notes |
|---|---|---|
| `cloud_id` | `zoneId`, `cloudId` | the cloud; newer sources read the second name |
| `group_id` | `groupId`, `siteId` | sources read one or the other |
| `layout_id` | `layoutId` | the provisioning context for the wizard's sources |
| `provision_type_code` | `provisionTypeCode`, `provisionType` | alternative to `layout_id`; sources read one or the other |
| `provision_type_id` | `provisionTypeId` | the same thing by id, for the sources that read the id rather than the code |
| `pool_id` | `poolId` | `pool-<id>` or a bare id — the form `zonePools` returns |
| `plan_id` | `planId` | the service plan, for plan-scoped options |
| `server_group_id` | `serverGroupId` | the cluster, for sources scoped to one — its VMs, hosts, affinity groups, datastores |
| `provider_id` | `providerId` | the backup provider, for the backup and replication sources |
| `phrase` | `phrase` | server-side name search, where a source supports it |
| `parent_id` | `parentId`, `parentid` | for dependent sources — the load balancer for its pools, the network server for its objects; sources read one spelling or the other |
| `credential_id` | `credential.type` | a stored `hpe_morpheus_credential` |

**The tenant is always the caller's.** Morpheus sets `accountId` to the authenticated user's tenant on
every option-source call and discards any value supplied, in the query string or under `config`. A
source is therefore always evaluated for the tenant the provider authenticated as; to read it as
another tenant, authenticate as a user of that tenant. Passing `accountId` in `parameters` produces a
warning saying so rather than a silently different answer.

These are the parameters read by many sources; around two hundred distinct keys are read across the
catalogue, and anything else goes in `parameters`, keyed as the API names it. A key given both ways is
an error rather than a silent precedence rule, since the two spellings mean the same thing.

Not every parameter a source reads can be supplied from Terraform. Some sources read the record
currently open in the UI — the router whose rule is being edited, the load balancer pool being
configured — as a server-side object (`domain`, or a bare `parent`) that only the UI can pass. Such
a source cannot be given that context over the API, and reading it will return the list a fresh UI
form would show.

## What narrows the read

An option source is read in **one request** — the endpoint does not page, and does not honour
`max` — so the cost of a read is the time the server takes to compute the list plus the time to
transfer it, and both scale with how many rows the source returns. Only the parameters the server
understands reduce that; everything the provider applies afterwards works on what was already
fetched.

| Argument | Applied | Effect on the read |
|---|---|---|
| the typed attributes (`cloud_id`, `layout_id`, `provision_type_code`, `pool_id`, `server_group_id`, …) | server-side | Narrow the list — where the source reads them. The provisioning sources read the first four; which parameters any other source reads is that source's own business (see *Finding a source and its parameters*). |
| `phrase` | server-side, **some sources only** | A name search where the source implements one (`users` does); ignored by many (`networks`, `servicePlans` and `instanceTypes` do not). Harmless where ignored. |
| `parameters` | server-side | Whatever the source reads; anything else is ignored. |
| `collection` | client-side | Selects a list within an object-shaped response. The whole object is fetched. |
| `filter` blocks | **client-side** | **No effect on the read.** They take regular expressions, which the server cannot, so they are the way to express a match — not a way to fetch less. Where a source honours `phrase`, a literal from the pattern set there narrows the fetch as well. |

To give a sense of scale: on an appliance with several hundred networks and over a thousand service
plans, reading `networks` or `servicePlans` unnarrowed took a few seconds and returned around a
megabyte, while the same sources narrowed by `cloud_id` or `provision_type_code` returned a dozen rows
in under a second. Every other source tried — groups, clouds, instance types, users, the provisioning
sources — answered in under a second. Below a few hundred rows the round trip dominates and the
row count makes no measurable difference.

Two things differ from the paginated list data sources such as `hpe_morpheus_images`:

- There is no page cost to avoid. The largest read observed is a few seconds, not tens of seconds,
  and a broad read is a bandwidth cost more than a time cost. Prefer the server-side context
  parameters when a source is large and they apply, but a read that fetches a few hundred rows to
  filter is not a problem to engineer around.
- **A source that reaches an external system is priced by that system, not by row count.**
  `scvmmLibraryShares` runs a query on the VMM server for every read, and no parameter makes it
  cheaper. Data sources are read on every plan and refresh, so such a source is queried on every
  plan; keep it in a configuration that is planned when the value is needed rather than one that is
  planned often.

Each data-source block is one request. Several blocks reading the same source with different filters
make several identical reads; read once and narrow in HCL, or put the filters on one block, when that
matters.

## Sources that return an object

Most sources return a list. A few return an object of lists — `zoneNetworkOptions` returns
`networks`, `networkGroups`, `networkSubnets` and `networkTypes` (plus a few scalars) — and for those
`collection` names the list to read. Without it the plan fails naming the lists available, so an
object-shaped source can never quietly yield nothing; with it on a source that returns a plain list,
the plan fails too, so the attribute is never silently ignored.

Within an option, `value` is the source's `value`, or its `id` when it has no `value` (the
`zoneNetworkOptions` items are `{id: "network-79", name: ...}`). Either way it is already in the form
the API takes on input — `pool-2`, `network-79` — so it can be used directly. The UI's group-heading
rows (`isGroup: true`, which carry no value) are not returned.

## Empty results

An empty `values` list is a legitimate answer for most sources and is not an error. Some sources,
though, answer an *incomplete* request with an empty list rather than an error: the provisioning
wizard's `zonePools` and `zoneNetworkOptions` return nothing without a `layout_id` or
`provision_type_code`. When such a source returns nothing and the context parameter is absent, the
data source raises a warning naming it:

```
Warning: option source "zonePools" returned no values

zonePools requires a provisioning context: set layout_id or provision_type_code in addition to
cloud_id. Without one the server returns an empty list rather than an error.
```

## Errors and warnings

| Diagnostic | Meaning | What to do |
|---|---|---|
| Warning *option source not in the catalogue* | The `source` is not among the sources known for the catalogue's Morpheus release, or it is defined only under a `service_type` that is not set | Check the spelling; set `service_type` if the warning names one; otherwise proceed — a plugin or newer appliance may serve it |
| Warning *option source … returned no values* | A wizard source was read without its provisioning context | Set `layout_id` or `provision_type_code` |
| Error *returned HTTP 500 (Looks like the server threw a gasket)* | Morpheus's generic failure. It is returned alike for an unknown source and for a known source given missing or malformed parameters; the message says which the catalogue thinks more likely, and shows the request | Fix the name or `service_type` if indicated; otherwise compare the parameters sent against what the UI sends for the same dropdown |
| Error *returns an object, not a list; set collection to one of …* | The source is object-shaped | Set `collection` to one of the listed names |
| Error *returns a list, so collection … does not apply* | `collection` was set for a list-shaped source | Remove `collection` |
| Error *conflicting parameters* | A key was set both by a typed attribute and in `parameters` | Use one |
| Error *invalid regular expression … for filter …* | A filter value does not compile as a Go regular expression | Fix the pattern; literal dots and brackets need escaping |

## Secrets

Never pass a raw credential in `parameters`. Data sources cannot have write-only arguments, so
anything in `parameters` is recorded in state. The sources that need to authenticate to a target
system (SCVMM, Amazon, OpenStack, …) accept a stored credential: create it with
`hpe_morpheus_credential` and pass its id as `credential_id`.

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `source` (String) The option source name, as the API names it, e.g. `zonePools`, `zoneNetworkOptions`, `groups`, `scvmmLibraryShares`. A name not in the provider's catalogue of sources known for Morpheus 9.0.2 produces a warning, not an error: plugins and newer appliances add sources.

### Optional

- `cloud_id` (Number) The cloud to evaluate the source in. Sent as `zoneId`, the API's name for a cloud (and as `cloudId`, which some newer sources read instead).
- `collection` (String) For a source that returns an object of lists rather than a single list — `zoneNetworkOptions` returns `networks`, `networkGroups`, `networkSubnets` and `networkTypes` — the list to read. Required for such a source; the error names the lists available. Not allowed for a source that returns a list.
- `credential_id` (Number) The id of an `hpe_morpheus_credential` for sources that authenticate to a target system (SCVMM, Amazon, OpenStack). Sent as `credential.type`. This is the route for secrets; never put them in `parameters`.
- `filter` (Block Set) Filter block, applied client-side after the source has been read. Repeat to apply several (all must match). Filter values are Go regular expressions and are case-sensitive. (see [below for nested schema](#nestedblock--filter))
- `group_id` (Number) The group (site) to evaluate the source for. Sent as both `groupId` and `siteId`, since sources read one or the other.
- `layout_id` (Number) The instance type layout that gives a provisioning source its context — the provisioning wizard's sources (`zonePools`, `zoneNetworkOptions`) return an empty list without one. Sent as `layoutId`.
- `parameter_style` (String) Where each parameter is placed: `auto` (the default) sends it both in the query string and as `config[<key>]`; `top_level` and `config` send it in one place only. Use them only for a source that misbehaves when a parameter is duplicated.
- `parameters` (Map of String) Parameters for the source, keyed as the API names them (`zoneId`, `layoutId`, ...). Morpheus reads parameters from two places — the query string and `config[<key>]` — and a parameter in the wrong place is silently ignored rather than rejected, so by default each key is sent in both (see `parameter_style`). The typed attributes below are shortcuts for the common keys; a key given both ways is an error.
- `parent_id` (Number) The id of the object a dependent source is scoped to — the load balancer for its pools, the network server for its objects. Sent as `parentId` (and `parentid`, which some sources read).
- `phrase` (String) A search phrase, for sources that narrow by name server-side. Sent as `phrase`.
- `plan_id` (Number) The service plan, for sources whose options depend on it (plan-scoped sizing and storage options). Sent as `planId`.
- `pool_id` (String) The resource pool to scope a provisioning source to, as `pool-<id>` or a bare id — the form `zonePools` returns as `value`. Sent as `poolId`.
- `provider_id` (Number) The backup provider, for the backup and replication sources scoped to one. Sent as `providerId`.
- `provision_type_code` (String) The provision type code (`kvm`, `vmware`, ...) as an alternative to `layout_id`. Sent as both `provisionTypeCode` and `provisionType`, since sources read one or the other.
- `provision_type_id` (Number) The provision type by id, for the sources that read the id rather than the code (`zonePools` reads both; `storageControllerTypes` reads only the id). Sent as `provisionTypeId`.
- `server_group_id` (Number) The cluster (server group), for sources scoped to one — its VMs, hosts, affinity groups, datastores. Sent as `serverGroupId`.
- `service_type` (String) Selects the two-segment path `/api/options/{service_type}/{source}`, which is how a source defined on a cloud-specific service is reached — for example `scvmm` for `scvmmLibraryShares`. A source that needs one and is called without it fails with a generic server error; the provider warns when the catalogue says so.

### Read-Only

- `values` (Attributes List) The values the source returned, in the source's order, after the filter blocks. Header rows the UI uses as group labels are not included. (see [below for nested schema](#nestedatt--values))

<a id="nestedblock--filter"></a>
### Nested Schema for `filter`

Required:

- `name` (String) The field to filter on: `name`, `value`, or `attributes.<key>` for a source-specific attribute, e.g. `attributes.providerType`.
- `values` (Set of String) The patterns. A value matches the block if the field matches ANY of them. A value that lacks the field does not match.


<a id="nestedatt--values"></a>
### Nested Schema for `values`

Read-Only:

- `attributes` (Map of String) Every other key the source returned for the option, as strings (numbers and booleans in their textual form, nested objects as JSON). Source-specific: `zonePools` returns `providerType` and `type`, `zoneNetworkOptions` returns `dhcpServer` and `pool`.
- `name` (String) The display label.
- `value` (String) The value to use in configuration — the source's `value`, or its `id` when it has no `value` (as `zoneNetworkOptions` items do). Already in the form the API takes on input, e.g. `pool-2`, `network-79`.
