---
page_title: "hpe_morpheus_network_pool_server Resource - terraform-provider-hpe"
subcategory: "morpheus"
description: |-
  Manages a Morpheus Network Pool Server (IPAM integration) resource.
---
# hpe_morpheus_network_pool_server (Resource)

Manages a Morpheus Network Pool Server (IPAM integration) resource.

This resource supports multiple IPAM pool server types, including Infoblox, Bluecat,
phpIPAM, SolarWinds, and EfficientIP (SOLIDserver). The pool server type is selected with
either `type_code` (a stable, environment-independent string) or `type_id` (a numeric ID).
These two attributes are mutually exclusive, and exactly one must be set. Using `type_code`
is recommended. All types share common connection attributes, while some attributes only
apply to specific types.

## Pool Server Type Codes

The `type_code` value is resolved by the Morpheus API against the registered pool server
types. The codes below are the most common; some are provided by built-in integrations and
others by installed plugins (for example, EfficientIP requires the
`morpheus-efficientip-plugin`).

| IPAM System | `type_code` | Availability |
|-------------|-------------|--------------|
| Infoblox | `infoblox` | Built-in / plugin |
| BlueCat | `bluecat` | Built-in / plugin |
| phpIPAM | `phpipam` | Built-in / plugin |
| SolarWinds | `solarwindsipam` | Built-in / plugin |
| EfficientIP (SOLIDserver) | `solidserver` | Plugin |
| NetBox | `netbox` | Plugin |
| Micetro | `micetro` | Plugin |
| Morpheus (internal) | `morpheus` | Built-in |

The authoritative, environment-specific list of available type codes can be retrieved from
the Morpheus API at `GET /api/networks/pool-server-types`, since the available codes depend
on which integrations and plugins are installed.

## Attribute Applicability by Type

| Attribute | Infoblox | Bluecat | phpIPAM | SolarWinds | EfficientIP |
|-----------|:--------:|:-------:|:-------:|:----------:|:-----------:|
| `name` | Yes | Yes | Yes | Yes | Yes |
| `type_code` | Yes | Yes | Yes | Yes | Yes |
| `type_id` | Yes | Yes | Yes | Yes | Yes |
| `enabled` | Yes | Yes | Yes | Yes | Yes |
| `service_url` | Yes | Yes | Yes | Yes | Yes |
| `service_username` | Yes | Yes | Yes | Yes | Yes |
| `service_password_wo` | Yes | Yes | Yes | Yes | Yes |
| `credential_id` | Yes | Yes | Yes | Yes | Yes |
| `ignore_ssl` | Yes | Yes | Yes | Yes | Yes |
| `inventory_existing` | Yes | Yes | Yes | Yes | Yes |
| `service_throttle_rate` | Yes | Yes | Yes | Yes | Yes |
| `network_filter` | Yes | Yes | Yes | - | - |
| `zone_filter` | Yes | - | - | - | - |
| `tenant_match` | Yes | - | - | - | - |
| `service_mode` | Yes | - | - | - | - |

## Example Usage

### Infoblox

```terraform
# Infoblox Network Pool Server
#
# Applicable attributes for Infoblox:
#   name, type_code, enabled, service_url, service_username, service_password_wo,
#   service_password_wo_version, credential_id, ignore_ssl, inventory_existing,
#   network_filter, zone_filter, tenant_match, service_mode, service_throttle_rate
resource "hpe_morpheus_network_pool_server" "infoblox" {
  name                        = "Infoblox IPAM"
  type_code                   = "infoblox"
  enabled                     = true
  service_url                 = "https://infoblox.example.com/wapi/v2.12"
  service_username            = "admin"
  service_password_wo         = "changeme"
  service_password_wo_version = 1
  ignore_ssl                  = true
  network_filter              = "10.0.0.0/8"
  zone_filter                 = "example.com"
  tenant_match                = ".*"
  service_mode                = "static"
  service_throttle_rate       = 0
}
```

### Bluecat

```terraform
# Bluecat Network Pool Server
#
# Applicable attributes for Bluecat:
#   name, type_code, enabled, service_url, service_username, service_password_wo,
#   service_password_wo_version, credential_id, ignore_ssl, inventory_existing,
#   network_filter, service_throttle_rate
resource "hpe_morpheus_network_pool_server" "bluecat" {
  name                        = "Bluecat IPAM"
  type_code                   = "bluecat"
  enabled                     = true
  service_url                 = "https://bluecat.example.com/api"
  service_username            = "admin"
  service_password_wo         = "changeme"
  service_password_wo_version = 1
  ignore_ssl                  = false
  network_filter              = "192.168.0.0/16"
  service_throttle_rate       = 50
}
```

### phpIPAM

```terraform
# phpIPAM Network Pool Server
#
# Applicable attributes for phpIPAM:
#   name, type_code, enabled, service_url, service_username, service_password_wo,
#   service_password_wo_version, credential_id, ignore_ssl, inventory_existing,
#   network_filter, service_throttle_rate
resource "hpe_morpheus_network_pool_server" "phpipam" {
  name                        = "phpIPAM"
  type_code                   = "phpipam"
  enabled                     = true
  service_url                 = "https://phpipam.example.com/api/app"
  service_username            = "admin"
  service_password_wo         = "changeme"
  service_password_wo_version = 1
  ignore_ssl                  = false
  network_filter              = "172.16.0.0/12"
  service_throttle_rate       = 0
}
```

### SolarWinds

```terraform
# SolarWinds Network Pool Server
#
# Applicable attributes for SolarWinds:
#   name, type_code, enabled, service_url, service_username, service_password_wo,
#   service_password_wo_version, credential_id, ignore_ssl, inventory_existing,
#   service_throttle_rate
resource "hpe_morpheus_network_pool_server" "solarwinds" {
  name                        = "SolarWinds IPAM"
  type_code                   = "solarwindsipam"
  enabled                     = true
  service_url                 = "https://solarwinds.example.com:17778/SolarWinds/InformationService/v3/Json"
  service_username            = "admin"
  service_password_wo         = "changeme"
  service_password_wo_version = 1
  ignore_ssl                  = true
  service_throttle_rate       = 100
}
```

### EfficientIP (SOLIDserver)

```terraform
# EfficientIP (SOLIDserver) Network Pool Server
#
# EfficientIP is provided by the external morpheus-efficientip-plugin and uses the
# "solidserver" type code. The plugin must be installed for this type to be available.
#
# Applicable attributes for EfficientIP:
#   name, type_code, enabled, service_url, service_username, service_password_wo,
#   service_password_wo_version, credential_id, ignore_ssl, inventory_existing,
#   service_throttle_rate
resource "hpe_morpheus_network_pool_server" "efficientip" {
  name                        = "EfficientIP IPAM"
  type_code                   = "solidserver"
  enabled                     = true
  service_url                 = "https://solidserver.example.com"
  service_username            = "admin"
  service_password_wo         = "changeme"
  service_password_wo_version = 1
  ignore_ssl                  = true
  inventory_existing          = false
  service_throttle_rate       = 0
}
```

### Using a Stored Credential

```terraform
# Network Pool Server using a stored credential
#
# Instead of providing service_username and service_password_wo inline,
# reference a stored credential by ID. credential_id conflicts with
# service_username and service_password_wo.
resource "hpe_morpheus_network_pool_server" "with_credential" {
  name          = "Infoblox with Credential"
  type_code     = "infoblox"
  enabled       = true
  service_url   = "https://infoblox.example.com/wapi/v2.12"
  credential_id = 42
  ignore_ssl    = true
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) The name of the network pool server.

### Optional

> **NOTE**: [Write-only arguments](https://developer.hashicorp.com/terraform/language/resources/ephemeral#write-only-arguments) are supported in Terraform 1.11 and later.

- `credential_id` (Number) The ID of a stored credential to use for authentication. Conflicts with service_username and service_password_wo.
- `enabled` (Boolean) Whether the network pool server is enabled.
- `ignore_ssl` (Boolean) Whether to ignore SSL certificate errors.
- `inventory_existing` (Boolean) Whether to inventory existing networks from the pool server. Applies to all pool server types.
- `network_filter` (String) Filter expression for which networks to sync from the pool server.
- `service_mode` (String) The service mode (e.g. "static" or "dhcp"). Applies to Infoblox.
- `service_password_wo` (String, Sensitive, [Write-only](https://developer.hashicorp.com/terraform/language/resources/ephemeral#write-only-arguments)) The service password for authentication. Conflicts with credential_id.
- `service_password_wo_version` (Number) Service password version. Used to determine if service_password_wo has been updated.
- `service_throttle_rate` (Number) Rate limit (in milliseconds) for API requests to the pool server.
- `service_url` (String) The service URL for the IPAM integration.
- `service_username` (String) The service username for authentication. Conflicts with credential_id.
- `tenant_match` (String) Tenant matching expression for multi-tenancy (Infoblox only).
- `type_code` (String) The network pool server type code (e.g. "infoblox", "bluecat", "phpipam", "solarwindsipam", "solidserver"). Mutually exclusive with type_id.
- `type_id` (Number) The ID of the network pool server type. Mutually exclusive with type_code.
- `zone_filter` (String) Filter expression for which DNS zones to sync (Infoblox only).

### Read-Only

- `id` (Number) The ID of the network pool server.
- `status` (String) The current status of the network pool server.

## Import

Network pool servers can be imported using the pool server ID, e.g.

```shell
terraform import hpe_morpheus_network_pool_server.example 42
```
