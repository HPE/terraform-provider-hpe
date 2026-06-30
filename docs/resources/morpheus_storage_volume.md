---
page_title: "hpe_morpheus_storage_volume Resource - terraform-provider-hpe"
subcategory: "morpheus"
description: |-
  Manages a Morpheus Storage Volume resource.
---
# hpe_morpheus_storage_volume (Resource)

Manages a Morpheus Storage Volume resource.

## Example Usage

```terraform
resource "hpe_morpheus_storage_volume" "example" {
  name              = "Example Storage Volume"
  type_code         = "hpealletraMPLUN"
  storage_server_id = 1
}
```

### Storage volume type

A storage volume requires a type, specified by **either** `type_id` **or**
`type_code` (mutually exclusive). `type_code` is recommended where available —
codes are stable across environments, whereas ids are not.

| `type_code` | Storage server | Notes |
|---|---|---|
| `3par` | HPE 3PAR / Primera | built-in (legacy) |
| `hpealletraMPLUN` | HPE Alletra Storage MP and Alletra 9000 (9060/9080) | standard LUN |
| `hpealletraMPLUN-active-pp` | HPE Alletra Storage MP and Alletra 9000 (9060/9080) | Peer Persistence (Active) |
| `hpealletraMPLUN-classic-pp` | HPE Alletra Storage MP and Alletra 9000 (9060/9080) | Peer Persistence (Classic) |

-> HPE Alletra 9000 (models 9060/9080) is managed by the same Alletra Storage MP
plugin as Alletra MP — it uses the same `hpealletraMPLUN*` types and the
`config_alletramp_bmaas` block, not the legacy `3par` type.

-> The available types depend on the storage plugins installed and the target
storage server. Discover the full set of `type_code`s and `type_id`s for your
environment with `GET /api/storage-volume-types`. Only certain storage server
types (for example 3PAR, Isilon, and Alletra Storage MP) support creating
storage volumes through this resource.

### Storage server

-> A volume must target a storage system: set `storage_server_id` (the storage
server / array) **or** `storage_group_id` (a storage group, which belongs to a
storage server). Creation routes to the storage provider through that server, so
omitting both fails with a generic `error saving volume` from the API. Look up
ids with the `hpe_morpheus_storage_server` / `hpe_morpheus_storage_servers` data
sources.

### Size

-> `max_storage` is specified in **GiB**. For example, `max_storage = 100`
creates a 100 GiB volume. (The Morpheus API stores and returns the size in bytes;
the provider converts to and from GiB.) HPE Alletra MP and Alletra 9000 volumes
must be between 1 and 65536 GiB; other storage types only require a positive
size. HPE Alletra 9000 (9060/9080) additionally enforces a 16 GiB minimum on the
array.

### Write-only configuration

-> The `config` and `config_alletramp_bmaas` arguments are
[write-only arguments](https://developer.hashicorp.com/terraform/language/resources/ephemeral#write-only-arguments)
(Terraform 1.11 and later): their values are sent to the Morpheus API on create
but are never stored in Terraform state. Because Terraform cannot detect changes
to a write-only value, each is paired with a companion version argument
(`config_wo_version` and `config_alletramp_bmaas_wo_version`). Increment the
version argument whenever you change the corresponding configuration to recreate
the volume with the new values.

#### Alletra MP BMaaS volume

The typed `config_alletramp_bmaas` block configures an HPE Alletra Storage MP
Bare Metal (BMaaS) volume. It applies to both HPE Alletra MP and HPE Alletra
9000 (9060/9080), which share the same plugin and API. `compute_server_id` and
`instance_ids` are mutually exclusive.

```terraform
resource "hpe_morpheus_storage_volume" "alletramp_bmaas" {
  name              = "Example Alletra MP BMaaS Volume"
  type_code         = "hpealletraMPLUN"
  storage_server_id = 1
  max_storage       = 30 # GiB

  # config_alletramp_bmaas is a write-only block: its values are sent to the API
  # on create but are never stored in Terraform state. Because Terraform cannot
  # detect changes to a write-only value, increment config_alletramp_bmaas_wo_version
  # whenever you change this block to recreate the volume with the new configuration.
  config_alletramp_bmaas = {
    datastore_id      = 5
    shared            = false
    compute_server_id = 10
  }
  config_alletramp_bmaas_wo_version = 1
}
```

#### Generic write-only config

For storage volume types without a typed config block, use the generic `config`
map with the storage plugin's native keys.

```terraform
resource "hpe_morpheus_storage_volume" "generic" {
  name              = "Example Storage Volume"
  type_code         = "hpealletraMPLUN"
  storage_server_id = 1
  max_storage       = 30 # GiB

  # config is a generic, write-only configuration map for storage volume types
  # that do not have a typed config block. Like config_alletramp_bmaas, its
  # values are sent to the API on create but are never stored in state.
  # Increment config_wo_version to recreate the volume with the new
  # configuration. The keys are the storage plugin's native config keys.
  config = {
    hpe_storage_datastore = 5
  }
  config_wo_version = 1
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) The name of the storage volume.

### Optional

> **NOTE**: [Write-only arguments](https://developer.hashicorp.com/terraform/language/resources/ephemeral#write-only-arguments) are supported in Terraform 1.11 and later.

- `config` (Dynamic, [Write-only](https://developer.hashicorp.com/terraform/language/resources/ephemeral#write-only-arguments)) Generic write-only configuration options for the storage volume, varies based on the storage volume type. Only sent to the API on create; not stored in state. Increment config_wo_version to apply a change.
- `config_alletramp_bmaas` (Attributes, [Write-only](https://developer.hashicorp.com/terraform/language/resources/ephemeral#write-only-arguments)) Alletra MP BMaaS storage volume configuration. This is a write-only attribute; its values are not stored in state. Increment config_alletramp_bmaas_wo_version to apply a change. (see [below for nested schema](#nestedatt--config_alletramp_bmaas))
- `config_alletramp_bmaas_wo_version` (Number) Version trigger for the write-only config_alletramp_bmaas attribute. Increment whenever config_alletramp_bmaas changes to recreate the volume with the new configuration.
- `config_wo_version` (Number) Version trigger for the write-only config attribute. Increment whenever config changes to recreate the volume with the new configuration.
- `max_storage` (Number) The storage volume size in GiB. HPE Alletra MP and Alletra 9000 volumes must be between 1 and 65536 GiB.
- `provision_type` (String) Provision type for storage volume types that support it.
- `storage_group_id` (Number) The ID of the storage group.
- `storage_server_id` (Number) The ID of the storage server.
- `type_code` (String) The storage volume type code, which is more stable across environments than type_id (e.g. "3par", "hpealletraMPLUN", "hpealletraMPLUN-active-pp", "hpealletraMPLUN-classic-pp"). Mutually exclusive with type_id.
- `type_id` (Number) The ID of the storage volume type. Mutually exclusive with type_code.

### Read-Only

- `id` (Number) The ID of the storage volume.
- `status` (String) The status of the storage volume.
- `wwn` (String)

<a id="nestedatt--config_alletramp_bmaas"></a>
### Nested Schema for `config_alletramp_bmaas`

Required:

- `datastore_id` (Number, [Write-only](https://developer.hashicorp.com/terraform/language/resources/ephemeral#write-only-arguments)) ID of the Alletra MP BMaaS data store (pool) in which to create the volume.

Optional:

- `compute_server_id` (Number, [Write-only](https://developer.hashicorp.com/terraform/language/resources/ephemeral#write-only-arguments)) Compute server ID to export a non-shared volume to.
- `instance_ids` (List of Number, [Write-only](https://developer.hashicorp.com/terraform/language/resources/ephemeral#write-only-arguments)) List of instance IDs to export a shared volume to.
- `remote_copy_target_id` (String, [Write-only](https://developer.hashicorp.com/terraform/language/resources/ephemeral#write-only-arguments)) Remote copy (replication) target ID. Required for replicated LUN volume types.
- `shared` (Boolean, [Write-only](https://developer.hashicorp.com/terraform/language/resources/ephemeral#write-only-arguments)) Whether the volume is shared (multi-attach).
- `use_existing_volume_set` (Boolean, [Write-only](https://developer.hashicorp.com/terraform/language/resources/ephemeral#write-only-arguments)) Whether to add the volume to an existing, exported volume set rather than creating a new one.
- `volume_set_id` (String, [Write-only](https://developer.hashicorp.com/terraform/language/resources/ephemeral#write-only-arguments)) ID of an existing volume set to add the volume to.
- `volume_set_name` (String, [Write-only](https://developer.hashicorp.com/terraform/language/resources/ephemeral#write-only-arguments)) Base name for a new volume set (a unique suffix is always appended).

## Import

Import is supported using the resource ID.

```bash
terraform import hpe_morpheus_storage_volume.example 123
```
