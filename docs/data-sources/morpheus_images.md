---
page_title: "hpe_morpheus_images Data Source - terraform-provider-hpe"
subcategory: "Morpheus"
description: |-
  
---
# hpe_morpheus_images (Data Source)



Finds virtual images — user-uploaded, cloud-synced and system alike — and returns the full
record for each. Its usual purpose is resolving the id to give `config_hvm.image_id` or
`config_vmware.image_id` on `hpe_morpheus_instance`.

~> **Narrow this data source server-side.** Every image the server-side filters admit is
downloaded before anything else happens. On an appliance with a large image library —
several thousand images is not unusual once cloud-synced and system images are counted —
an unnarrowed read can take a long time, and Terraform performs it on refresh, plan and
apply alike. Reach for at least one of the server-side arguments below before relying on
`filter` blocks.

## What narrows the download

The arguments in the first group are passed to the API, so they reduce what is fetched.
Everything else is applied by the provider afterwards, on whatever the first group
returned.

| Argument | Applied | Effect on the download |
|---|---|---|
| `image_id` | server-side | Smallest possible read. Prefer it when the ids are known. |
| `image_type` | server-side | Usually the most effective. Accepts several values. |
| `phrase` | server-side | Substring of the name. |
| `description` | server-side | Accepts several values. |
| `labels`, `all_labels` | server-side | |
| `system_image`, `include_system_image` | server-side | |
| `filter_type` | server-side | See the note below — the default is deliberately broad. |
| `filter` blocks | **client-side** | **Usually no effect on the download** — see below. |
| `sort`, `direction` | server-side | Affect page order only, not volume. |

!> A configuration using only `filter` blocks generally downloads the **entire** image
library before the first block is evaluated. The blocks are worth having — they take
regular expressions, which the API cannot — but they are a way to express a match, not
usually a way to fetch less.

-> **One exception.** A single-valued `filter` block on `name` whose expression *starts
with a literal* is used to narrow the fetch as well: `^ubuntu`, `ubuntu.*` and
`^rhel[0-9]` all contribute `ubuntu`, `ubuntu` and `rhel` respectively, and the server is
asked for names containing that. Expressions that begin with a pattern — `.*ubuntu`,
`[a-z]+`, or an alternation such as `ubuntu|debian` — offer nothing to narrow with, and
the whole library is fetched. Setting `phrase` yourself always takes precedence.

-> `filter_type` defaults to `All`, unlike the API, which returns only user-uploaded
images. That default is chosen so synced and system images are visible, since those are
the ones most instances are provisioned from — but it also means the default is the
broadest possible read. Setting it to `User`, `Synced` or `System` narrows the fetch
considerably where that suits.

## Choosing between this and `hpe_morpheus_image`

Use the singular [`hpe_morpheus_image`](image) to resolve one known image: it matches
`image_type` exactly, so it can distinguish images that share a name and differ only by
type.

Use this data source to discover images. Here `image_type` follows the API's aliasing —
`vmware` also matches `ovf` and `vmdk`, `virtualbox` also matches `vdi` — so a single
value can cover a hypervisor's several formats.

## Example Usage

```terraform
# Narrow server-side wherever possible. Every image the filters admit is
# fetched before the provider does anything with it, and a large estate makes
# that expensive: on an appliance holding around 7000 images an unfiltered read
# takes well over a minute, on every plan and every apply.
#
# image_type accepts several values and some are aliases: "vmware" also matches
# ovf and vmdk, "virtualbox" also matches vdi. Anything else, such as "qcow2",
# "raw" or "iso", is matched exactly.
data "hpe_morpheus_images" "hvm_bootable" {
  image_type = ["qcow2", "raw"]
}

# filter_type decides what is considered before any other filter. It defaults to
# "All"; left to itself the API returns only images a user uploaded, hiding the
# synced and system images most instances are actually provisioned from.
data "hpe_morpheus_images" "user_uploaded" {
  image_type  = ["qcow2", "raw"]
  filter_type = "User"
}

# Filter blocks are applied by the provider once the matching images have been
# fetched, so they can express patterns the API cannot. Values are Go regular
# expressions and are unanchored, so anchor them when that is what you mean.
#
# Combine them with the server-side arguments rather than relying on them alone:
# the arguments decide how much is fetched, the blocks only decide what survives.
data "hpe_morpheus_images" "ubuntu_active" {
  image_type = ["qcow2", "raw"]

  filter {
    name   = "name"
    values = ["^ubuntu"]
  }

  filter {
    name   = "status"
    values = ["active"]
  }
}

# The id of an image to provision from, for config_hvm.image_id or
# config_vmware.image_id on hpe_morpheus_instance. images is a set, so it cannot
# be indexed; sort the ids to pick one deterministically.
output "ubuntu_image_ids" {
  value = sort([for image in data.hpe_morpheus_images.ubuntu_active.images : image.id])
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `all_labels` (String) Filter by label, requiring every given label to be present rather than any.
- `description` (Set of String) Filter by description. An image matches if its description matches any of the values.

Matching is a case-insensitive SQL `like` performed by the server, so `%` is the wildcard — `ubuntu%`, not `ubuntu.*`. For pattern matching use a `filter` block, which takes regular expressions.
- `direction` (String) Sort direction used while paging, `asc` (the default) or `desc`.
- `filter` (Block Set) Filter block. Repeat to apply multiple filters (all ANDed together). Filter values are case-sensitive and support Go regular expressions (https://regex101.com/).

These are applied by the provider after every matching image has been fetched, so they can express patterns the API cannot. Prefer the top-level arguments where they suffice: those are applied by the server and fetch less. (see [below for nested schema](#nestedblock--filter))
- `filter_type` (String) Which images to consider, before any other filter is applied.

Defaults to `All`, which differs from the platform's own default of `User`. The API returns only user-uploaded images unless told otherwise, hiding synced and system images — including the cloud-provided images most instances are actually provisioned from.

`System` returns only system images, `Synced` only images synced from a cloud, `User` only those uploaded by a user.

Ignored when `system_image` or `include_system_image` is set: those take precedence server-side.
- `image_id` (Set of Number) Filter by image ID. Returns the images with these IDs, and is the cheapest way to fetch a known set.
- `image_type` (Set of String) Filter by image type code. An image matches if it is any of the values given.

Some values are aliases rather than exact codes: `vmware` also matches `ovf` and `vmdk`, and `virtualbox` also matches `vdi`. Any other value is matched exactly, so `qcow2`, `raw` and `iso` mean precisely themselves.

This is the usual way to narrow by hypervisor: `["qcow2", "raw"]` for HVM, `["vmware"]` for VMware, `["iso"]` for bare metal.
- `include_system_image` (Boolean) Include system images alongside non-system ones. System images are excluded unless this is `true`.

Takes precedence over `filter_type`, and is itself overridden by `system_image`.
- `labels` (String) Filter by label. Images carrying any of the given labels match.
- `phrase` (String) Filter by a substring of the image name, matched case-insensitively by the server.

There is deliberately no top-level `name` argument: it would be a server-side wildcard match sitting beside a `filter` block that takes regular expressions, which is two pattern syntaxes on one data source. Use `phrase` for a substring, or a `filter` block on `name` for a regular expression.
- `sort` (String) Property the server sorts by while paging. Defaults to `name`.

This does not order the result: `images` is a set, so it has no order. It still affects which records fall on which page, so it is not inert.
- `system_image` (Boolean) Filter on whether an image is a system image. `true` returns only system images, `false` only non-system images.

Takes precedence over both `filter_type` and `include_system_image`.

### Read-Only

- `images` (Attributes Set) The images matching the supplied filters.

This is a set, so it has no order and cannot be indexed. Ordering is not offered because the server pages by a non-unique key, so a stable order cannot be promised. (see [below for nested schema](#nestedatt--images))

<a id="nestedblock--filter"></a>
### Nested Schema for `filter`

Required:

- `name` (String) The field to filter on. Valid names are: name, description, image_type, status, visibility.
- `values` (Set of String) The filter values. An image matches the block if the chosen field matches ANY value (Go regular expression).

Expressions are unanchored, so `Test` matches anywhere in the value. Anchor explicitly with `^Test` or `^Test$` when that is what is meant.


<a id="nestedatt--images"></a>
### Nested Schema for `images`

Read-Only:

- `auto_join_domain` (Boolean) Auto Join Domain?
- `cloud_init` (Boolean) Cloud Init Enabled?
- `config_azure` (Attributes) Azure Reference Virtual Image Parameters (see [below for nested schema](#nestedatt--images--config_azure))
- `console_keymap` (String)
- `description` (String) A description for the virtual image
- `external_id` (String)
- `fips_enabled` (Boolean) FIPS enabled?
- `force_customization` (Boolean) Force Guest Customization?
- `id` (Number)
- `image_type` (String) Code of image type. eg. vmware, ami, etc.
- `install_agent` (Boolean) Install Agent?
- `labels` (Set of String) Array of label strings, can be used for filtering.
- `min_disk` (Number) Minimum disk size the image requires, in GB.

The API reports this in bytes; the provider converts it, so a value of `10` means 10GB. A volume smaller than this is rejected at provision time rather than grown.
- `min_ram` (Number) Minimum memory the image requires, in GB.

The API reports this in bytes; the provider converts it, so a value of `4` means 4GB.
- `name` (String) A name for the virtual image
- `os_type_id` (Number)
- `owner_id` (Number) Owner of the image
- `raw_size` (Number) Size of the image on disk, in bytes.

Unlike `min_disk` and `min_ram` this is left as the API reports it. A volume that clears `min_disk` but falls short of this is rounded up at provision time.
- `ssh_key` (String) SSH Key
- `ssh_username` (String) SSH Username
- `status` (String)
- `storage_provider_id` (Number)
- `sysprep` (Boolean) Sysprep Enabled?
- `system_image` (Boolean) Is created by system?
- `tags` (Attributes Set) Metadata tags, Array of objects having a name and value (see [below for nested schema](#nestedatt--images--tags))
- `tenants` (Attributes Set) (see [below for nested schema](#nestedatt--images--tenants))
- `trial_version` (Boolean) Is Trial Version?
- `uefi` (Boolean) UEFI enabled?
- `user_data` (String) Cloud-Init User Data, a bash script
- `user_defined` (Boolean) Is defined by an user?
- `user_uploaded` (Boolean) Is uploaded by an user?
- `virtio_supported` (Boolean) VirtIO Drivers Loaded?
- `visibility` (String) private or public
- `vm_tools_installed` (Boolean) VM Tools Installed?

<a id="nestedatt--images--config_azure"></a>
### Nested Schema for `images.config_azure`

Read-Only:

- `offer` (String) The name of the offer in the Azure Marketplace
- `publisher` (String) The name of the publisher in the Azure Marketplace
- `sku` (String) The name of the sku in the Azure Marketplace
- `version` (String) The name of the version in the Azure Marketplace


<a id="nestedatt--images--tags"></a>
### Nested Schema for `images.tags`

Read-Only:

- `name` (String)
- `value` (String)


<a id="nestedatt--images--tenants"></a>
### Nested Schema for `images.tenants`

Read-Only:

- `id` (Number)
- `name` (String)
