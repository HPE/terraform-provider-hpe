---
page_title: "hpe_morpheus_option_source_value Data Source - terraform-provider-hpe"
subcategory: "Morpheus"
description: |-
  Reads exactly one value of a Morpheus option source and fails the plan when the source returns none or several. Morpheus option sources are the lookups behind the UI's dropdowns — the resource pools and networks the provisioning wizard offers, an SCVMM library share, the members of a select list — served at GET /api/options/{source}.
---
# hpe_morpheus_option_source_value (Data Source)

Reads exactly one value of a Morpheus option source and fails the plan when the source returns none or several. Morpheus option sources are the lookups behind the UI's dropdowns — the resource pools and networks the provisioning wizard offers, an SCVMM library share, the members of a select list — served at `GET /api/options/{source}`.

Use this data source when a configuration needs **one** value that only Morpheus can supply and the
plan should fail if it cannot be found — the resource pool of a named HVM cluster, or the library
share of an SCVMM server. To list what a source offers, or to iterate over it, use
[`hpe_morpheus_option_source_values`](https://registry.terraform.io/providers/HPE/hpe/latest/docs/data-sources/morpheus_option_source_values).

## What an option source is

Every dropdown in the Morpheus UI that is not a fixed list is populated by an *option source*: a
named server-side lookup that returns `{name, value}` pairs for the current context. Some read the
Morpheus database (groups, service plans); some compute a provisioning context (the resource pools
and networks a wizard step may choose from, given a cloud, layout and pool); some ask an external
system (the library shares on an SCVMM server, given credentials). They are served over the API at
`GET /api/options/{source}` and `GET /api/options/{service_type}/{source}`, take their parameters
from the query string, and return a list of options.

Until now nothing in the provider could read them. A value such as a cluster's resource pool or an
SCVMM library share had to be looked up by hand — in the UI, the CLI or with `curl` — and typed into
the configuration as a literal. This data source reads the same lookup the UI uses, so the
configuration can resolve the value itself and verify it at plan time.

## Use case: provisioning into an HVM cluster

Morpheus creates a resource pool for every HVM cluster and binds the cluster's networks to it. That
pool is the one to provision into — but it is attached to the cluster rather than to the cloud, so
the cloud's resource-pool listing does not return it, and neither does `hpe_morpheus_resource_pool`
when looked up by name. A practitioner who cannot find the pool may conclude it does not exist and
create one by hand; an instance pointed at that pool then fails with `Invalid network`, because the
networks belong to the cluster's pool.

The provisioning wizard does not have this problem because it uses the `zonePools` option source,
which lists the pools that are selectable for a cloud and layout, cluster pools included. Reading the
same source gives the configuration the right pool, by the cluster's name, in exactly the form the
instance takes:

```terraform
# Provisioning into an HVM cluster: read the cluster's own resource pool — the
# one its networks belong to — from the zonePools source, by the cluster's name.

data "hpe_morpheus_cloud" "hvm" {
  name = "my-hvm-cloud"
}

# The layout supplies the provisioning context. Without it (or a provision type)
# zonePools returns an empty list rather than an error.
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

# The cluster's pool, selected by the cluster's name. value is already in the
# form the instance takes ("pool-<id>"), and the source reports the pool's
# provider type alongside.
data "hpe_morpheus_option_source_value" "cluster_pool" {
  source    = "zonePools"
  cloud_id  = data.hpe_morpheus_cloud.hvm.id
  layout_id = data.hpe_morpheus_instance_type_layout.hvm.id
  name      = "hvm-cluster-01"
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

  # With the generic config block instead of config_hvm, the pool's provider
  # type is needed too, and the source has it:
  #
  #   config = {
  #     resourcePoolId   = data.hpe_morpheus_option_source_value.cluster_pool.value
  #     poolProviderType = data.hpe_morpheus_option_source_value.cluster_pool.attributes["providerType"]
  #   }
}
```

`zonePools` needs a provisioning context — `layout_id` or `provision_type_code` — and returns an
empty list rather than an error without one; the data source says so when that happens. The
[`hpe_morpheus_cluster`](https://registry.terraform.io/providers/HPE/hpe/latest/docs/data-sources/morpheus_cluster)
data source also exposes the pool, as `permissions.resource_pool.id`, when the cluster is already
looked up for another reason.

## Use case: declaring an SCVMM cloud

An SCVMM cloud requires `libraryShare`, the UNC path of a library share on the VMM server. Only the
server can say which shares exist, so the value cannot be written from documentation. The UI
discovers it by asking the server once it has credentials — the `scvmmLibraryShares` option source —
and that lookup works before the cloud exists, so the first cloud can be declared without obtaining
the share by hand:

```terraform
# Declaring an SCVMM cloud: resolve the required libraryShare from the VMM
# server itself, through a stored credential, before the cloud exists.

# The SCVMM credentials, stored once in Morpheus. The secret lives here and is
# referenced by id below; it does not pass through the data source, so nothing
# sensitive is recorded in its state.
resource "hpe_morpheus_credential" "scvmm" {
  name     = "scvmm-prod"
  type     = "username-password"
  username = var.scvmm_username
  password = var.scvmm_password
}

# The library share. The source is defined on the SCVMM service, so service_type
# selects it, and it reads the VMM server address from the host parameter. When
# the server has several shares, narrow with name — they are UNC paths — or the
# plan fails listing them.
data "hpe_morpheus_option_source_value" "library_share" {
  source        = "scvmmLibraryShares"
  service_type  = "scvmm"
  credential_id = hpe_morpheus_credential.scvmm.id
  parameters    = { host = var.scvmm_host }
}

resource "hpe_morpheus_cloud" "scvmm" {
  name            = "scvmm-prod"
  tenant_id       = 1
  cloud_type_code = "scvmm"

  config = {
    host         = var.scvmm_host
    username     = var.scvmm_username
    password     = var.scvmm_password
    workingPath  = "c:\\Temp"
    diskPath     = "c:\\VirtualDisks"
    libraryShare = data.hpe_morpheus_option_source_value.library_share.value
  }
}
```

Three things in that example carry over to any source that reaches an external system:

- **The source is defined on a cloud-specific service**, so `service_type` selects the two-segment
  path `/api/options/scvmm/scvmmLibraryShares`. Without it Morpheus looks on its base service, does
  not find the source, and returns a generic server error; the data source warns when the catalogue
  says a source needs a `service_type`.
- **Credentials are referenced, not passed.** Data-source arguments are recorded in state, so raw
  credentials in `parameters` would be too. The sources that authenticate to a target system accept
  a stored credential instead: store it with `hpe_morpheus_credential` and reference it with
  `credential_id`.
- **The source's own parameters** (`host` here) go in `parameters` under the name the API uses.

## Selecting one of several values

A source often returns several values. `name` and `value` select one — both are also outputs, so
whichever is not used to select comes back — and selecting by `value` is a way to confirm at plan
time that a value the configuration already holds is one the source offers:

```terraform
# Selecting by value: confirm at plan time that a value the configuration
# already holds is one the source offers, and learn its name.

data "hpe_morpheus_cloud" "hvm" {
  name = "my-hvm-cloud"
}

data "hpe_morpheus_instance_type_layout" "hvm" {
  name    = "Single KVM VM"
  version = "22.04"
}

# If pool-2 is not selectable for this cloud and layout, the plan fails and
# names the pools that are.
data "hpe_morpheus_option_source_value" "pool_2" {
  source    = "zonePools"
  cloud_id  = data.hpe_morpheus_cloud.hvm.id
  layout_id = data.hpe_morpheus_instance_type_layout.hvm.id
  value     = "pool-2"
}

output "pool_2" {
  value = {
    name          = data.hpe_morpheus_option_source_value.pool_2.name
    provider_type = data.hpe_morpheus_option_source_value.pool_2.attributes["providerType"]
    cluster_type  = data.hpe_morpheus_option_source_value.pool_2.attributes["type"]
  }
}
```

When several values remain after selection, the plan fails and names them:

```
Error: multiple matches for option source "zonePools"

Expected exactly one value, found 2:
    hvm-cluster-01 (pool-2)
    hvm-cluster-02 (pool-7)

Narrow the result with name or value, or use hpe_morpheus_option_source_values to select
from the list.
```

When none remains, the plan fails and lists what the source did return, so the fix is visible:

```
Error: no match for option source "zonePools"

None of the 2 values the source returned has name "hvm-cluster-03". The values are:
    hvm-cluster-01 (pool-2)
    hvm-cluster-02 (pool-7)
```

## Finding a source and its parameters

The two sources above are the ones most configurations need. For any other dropdown, the Morpheus UI
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
| `name`, `value` | client-side | Select one row from the fetched list. **No effect on the read.** Where a source honours `phrase`, setting it to the name narrows the fetch as well. |

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
  select one is not a problem to engineer around.
- **A source that reaches an external system is priced by that system, not by row count.**
  `scvmmLibraryShares` runs a query on the VMM server for every read, and no parameter makes it
  cheaper. Data sources are read on every plan and refresh, so such a source is queried on every
  plan; keep it in a configuration that is planned when the value is needed rather than one that is
  planned often.

Each data-source block is one request. Several blocks reading the same source with different
selections make several identical reads; read once and select in HCL, or with `filter` blocks on a
single `hpe_morpheus_option_source_values` block, when that matters.

## Sources that return an object

Most sources return a list. A few return an object of lists — `zoneNetworkOptions` returns
`networks`, `networkGroups`, `networkSubnets` and `networkTypes` (plus a few scalars) — and for those
`collection` names the list to read. Without it the plan fails naming the lists available, so an
object-shaped source can never quietly yield nothing; with it on a source that returns a plain list,
the plan fails too, so the attribute is never silently ignored.

Within an option, `value` is the source's `value`, or its `id` when it has no `value` (the
`zoneNetworkOptions` items are `{id: "network-79", name: ...}`). Either way it is already in the form
the API takes on input — `pool-2`, `network-79` — so it can be used directly.

## Empty results

Some sources answer an incomplete request with an empty list rather than an error. The provisioning
wizard's sources are the known cases: `zonePools` and `zoneNetworkOptions` return nothing without a
`layout_id` or `provision_type_code`. When such a source returns nothing and the context parameter
is absent, the no-match error says so:

```
Error: no match for option source "zonePools"

The source returned no values. zonePools requires a provisioning context: set layout_id or
provision_type_code in addition to cloud_id. Without one the server returns an empty list rather
than an error.
```

## Errors and warnings

| Diagnostic | Meaning | What to do |
|---|---|---|
| Warning *option source not in the catalogue* | The `source` is not among the sources known for the catalogue's Morpheus release, or it is defined only under a `service_type` that is not set | Check the spelling; set `service_type` if the warning names one; otherwise proceed — a plugin or newer appliance may serve it |
| Error *returned HTTP 500 (Looks like the server threw a gasket)* | Morpheus's generic failure. It is returned alike for an unknown source and for a known source given missing or malformed parameters; the message says which the catalogue thinks more likely, and shows the request | Fix the name or `service_type` if indicated; otherwise compare the parameters sent against what the UI sends for the same dropdown |
| Error *returns an object, not a list; set collection to one of …* | The source is object-shaped | Set `collection` to one of the listed names |
| Error *returns a list, so collection … does not apply* | `collection` was set for a list-shaped source | Remove `collection` |
| Error *no match for option source* | Nothing matched `name`/`value`, or the source returned nothing | The message lists the values returned, or names the missing context parameter |
| Error *multiple matches for option source* | Several values remain after selection | Narrow with `name` or `value`, or use the plural data source |
| Error *conflicting parameters* | A key was set both by a typed attribute and in `parameters` | Use one |

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
- `group_id` (Number) The group (site) to evaluate the source for. Sent as both `groupId` and `siteId`, since sources read one or the other.
- `layout_id` (Number) The instance type layout that gives a provisioning source its context — the provisioning wizard's sources (`zonePools`, `zoneNetworkOptions`) return an empty list without one. Sent as `layoutId`.
- `name` (String) The display label of the value to select. When set, only values with exactly this name are considered. Set on output.
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
- `value` (String) The value to select, e.g. `pool-2`. When set, only values with exactly this value are considered. On output, the value to use in configuration: the source's `value`, or its `id` when it has no `value`.

### Read-Only

- `attributes` (Map of String) Every other key the source returned for the option, as strings (numbers and booleans in their textual form, nested objects as JSON). Source-specific: `zonePools` returns `providerType` and `type`, `zoneNetworkOptions` returns `dhcpServer` and `pool`.
