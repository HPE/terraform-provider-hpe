---
page_title: "hpegl to HPE"
subcategory: "Migration"
---

# Migrate hpegl VMaaS Provider Resources to HPE Provider Resources

To migrate resources from the hpegl (`hpe/hpegl`) VMaaS provider to the HPE Terraform provider
(`HPE/hpe`) the basic principle is to import the existing resources created with the hpegl provider
into the HPE provider. HashiCorp documentation on import can be found
[here](https://developer.hashicorp.com/terraform/language/import).

-> [tfmigrator](./tfmigrator_migration.md) has been developed in order to assist with migration. By
following the documented process, this tooling ports `hpegl_vmaas_*` resources to their corresponding
`hpe_morpheus_*` resources - including handling for the provider block transformation, modules, and
variables. See the [tfmigrator guide](./tfmigrator_migration.md) for the full command workflow (both
the one-shot `migrate` command and the individual steps).<br><br>
The section on [Bulk Import](https://developer.hashicorp.com/terraform/language/v1.14.x/import/bulk?page=import&page=bulk)
requires provider support for the `List` RPC. The HPE provider does not currently support `List`,
this is something that we are considering for a future release.<br><br>
Unlike the Morpheus provider, **every** hpegl VMaaS resource has a different schema in the HPE
provider (hpegl used the Terraform Plugin SDK v2; the HPE provider uses the Plugin Framework). Every
resource therefore requires Terraform-generated configuration and manual post-migration review. The
process documented for [tfmigrator](./tfmigrator_migration.md) handles this generation for you.

## Provider Block Transformation

The largest difference between an hpegl configuration and an HPE configuration is the provider block.
hpegl authenticates through the GreenLake IAM broker, while the HPE provider connects to Morpheus
either through GreenLake IAM or directly. The [`migrate-providers`](./tfmigrator_migration.md) step
performs this transformation automatically (including emitting credential `variable` blocks and a
`credentials.auto.tfvars` file), but the mapping is summarised here for reference.

The generated block depends on the `iam_version` attribute of the hpegl provider (falling back to the
`HPEGL_IAM_VERSION` environment variable):

### GreenLake IAM - Connected (`iam_version = "glcs"`)

```HCL
# Before (hpegl)
provider "hpegl" {
  iam_version      = "glcs"
  iam_service_url  = "https://client.greenlake.hpe.com/..."
  vmaas {
    location   = "HPE-Data-Center"
    space_name = "my-space"
    broker_url = "https://vmaas-broker.example.com"
  }
}

# After (hpe)
provider "hpe" {
  morpheus {
    pce_identity {
      location      = "HPE-Data-Center"
      space         = "my-space"
      broker_url    = "https://vmaas-broker.example.com"
      issuer_url    = "https://client.greenlake.hpe.com/..."
      client_id     = var.hpe_pce_client_id
      client_secret = var.hpe_pce_client_secret
    }
  }
}
```

### GreenLake IAM - Disconnected (`iam_version = "glp"`)

```HCL
# After (hpe)
provider "hpe" {
  morpheus {
    pce_disconnected_identity {
      location         = "HPE-Data-Center"
      workspace_id     = "my-space"
      broker_url       = "https://vmaas-broker.example.com"
      token_issuer_url = "https://client.greenlake.hpe.com/..."
      client_id        = var.hpe_pce_client_id
      client_secret    = var.hpe_pce_client_secret
    }
  }
}
```

### Direct-connect (`vmaas.morpheus_url` present)

When the hpegl `vmaas` block sets `morpheus_url`, the provider connects directly to Morpheus with no
GreenLake IAM broker. In that case the credentials land directly inside the `morpheus` block and no
`pce_identity` sub-block is emitted:

```HCL
# After (hpe)
provider "hpe" {
  morpheus {
    url          = "https://morpheus.example.com"
    access_token = var.hpe_morpheus_access_token
  }
}
```

## Resource Definitions

All hpegl VMaaS resources have a different schema in the HPE provider, so Terraform is used to
generate the HPE resource definition (this is how [tfmigrator](./tfmigrator_migration.md) handles
these resources). In addition to the generated resource definitions, an `import` block must be
supplied for each resource. After the import ensure that a `terraform plan` shows no changes; if
there are changes then the resource definitions need to be adjusted to match the existing resources.

| hpegl Provider Resource Name             | HPE Provider Resource Name                     |
|------------------------------------------|------------------------------------------------|
| hpegl_vmaas_instance                     | hpe_morpheus_instance                          |
| hpegl_vmaas_instance_clone               | hpe_morpheus_instance_clone                    |
| hpegl_vmaas_network                      | hpe_morpheus_network                           |
| hpegl_vmaas_dhcp_server                  | hpe_morpheus_network_dhcp_server               |
| hpegl_vmaas_router                       | hpe_morpheus_network_router                    |
| hpegl_vmaas_router_bgp_neighbor          | hpe_morpheus_network_router_bgp_neighbor       |
| hpegl_vmaas_router_nat_rule              | hpe_morpheus_network_router_nat                |
| hpegl_vmaas_router_route                 | hpe_morpheus_network_router_route              |
| hpegl_vmaas_router_firewall_rule_group   | hpe_morpheus_network_router_firewall_rule_group |
| hpegl_vmaas_load_balancer                | hpe_morpheus_load_balancer                     |
| hpegl_vmaas_load_balancer_monitor        | hpe_morpheus_load_balancer_monitor             |
| hpegl_vmaas_load_balancer_pool           | hpe_morpheus_load_balancer_pool                |
| hpegl_vmaas_load_balancer_profile        | hpe_morpheus_load_balancer_profile             |
| hpegl_vmaas_load_balancer_virtual_server | hpe_morpheus_load_balancer_virtual_server      |

### Composite Import IDs

Most resources import with a simple integer ID. The following resources are sub-resources whose HPE
equivalent requires a **composite** import ID rather than a bare integer. tfmigrator constructs these
automatically from the resource state; they are listed here so the generated `import` blocks make
sense during review.

| hpegl Provider Resource Name             | Import ID format                |
|------------------------------------------|---------------------------------|
| hpegl_vmaas_router_bgp_neighbor          | `{router_id}.{id}`              |
| hpegl_vmaas_router_nat_rule              | `{router_id}.{id}`              |
| hpegl_vmaas_router_route                 | `{router_id}.{id}`              |
| hpegl_vmaas_router_firewall_rule_group   | `{router_id}.{id}`              |
| hpegl_vmaas_dhcp_server                  | `{network_server_id}.{id}`      |
| hpegl_vmaas_load_balancer_monitor        | `{lb_id}.{id}`                  |
| hpegl_vmaas_load_balancer_virtual_server | `{lb_id}.{id}`                  |
| hpegl_vmaas_load_balancer_pool           | `{lb_id}/{id}`                  |
| hpegl_vmaas_load_balancer_profile        | `{lb_id}/{id}`                  |

## Data Source Definitions

Data sources are migrated in the same way, by updating the data source type. The following table
lists the hpegl provider data sources and their HPE equivalents.

| hpegl Provider Data Source Name                   | HPE Provider Data Source Name              |
|---------------------------------------------------|--------------------------------------------|
| hpegl_vmaas_cloud                                 | hpe_morpheus_cloud                         |
| hpegl_vmaas_group                                 | hpe_morpheus_group                         |
| hpegl_vmaas_plan                                  | hpe_morpheus_service_plan                  |
| hpegl_vmaas_layout                                | hpe_morpheus_instance_type_layout          |
| hpegl_vmaas_network                               | hpe_morpheus_network                       |
| hpegl_vmaas_datastore                             | hpe_morpheus_datastore                     |
| hpegl_vmaas_template                              | hpe_morpheus_image                         |
| hpegl_vmaas_environment                           | hpe_morpheus_environment                   |
| hpegl_vmaas_network_type                          | hpe_morpheus_network_type                  |
| hpegl_vmaas_network_pool                          | hpe_morpheus_network_pool                  |
| hpegl_vmaas_network_domain                        | hpe_morpheus_network_domain                |
| hpegl_vmaas_router                                | hpe_morpheus_network_router                |
| hpegl_vmaas_edge_cluster                          | hpe_morpheus_network_edge_cluster          |
| hpegl_vmaas_transport_zone                        | hpe_morpheus_network_transport_zone        |
| hpegl_vmaas_load_balancer                         | hpe_morpheus_load_balancer                 |
| hpegl_vmaas_load_balancer_monitor                 | hpe_morpheus_load_balancer_monitor         |
| hpegl_vmaas_load_balancer_pool                    | hpe_morpheus_load_balancer_pool            |
| hpegl_vmaas_load_balancer_profile                 | hpe_morpheus_load_balancer_profile         |
| hpegl_vmaas_load_balancer_virtual_server          | hpe_morpheus_load_balancer_virtual_server  |
| hpegl_vmaas_dhcp_server                           | hpe_morpheus_network_dhcp_server           |
| hpegl_vmaas_resource_pool                         | hpe_morpheus_resource_pool                 |
| hpegl_vmaas_power_schedule                        | hpe_morpheus_power_schedule                |
| hpegl_vmaas_cloud_folder                          | hpe_morpheus_cloud_folder                  |
| hpegl_vmaas_network_proxy                         | hpe_morpheus_network_proxy                 |
| hpegl_vmaas_network_interface                     | hpe_morpheus_network_interface_type        |
| hpegl_vmaas_instance_storage_controller           | hpe_morpheus_instance_storage_controller   |
| hpegl_vmaas_lb_pool_member_group                  | hpe_morpheus_network_server_group          |
| hpegl_vmaas_instance_disk_type                    | hpe_morpheus_instance_disk_type            |
| hpegl_vmaas_load_balancer_virtual_server_ssl_cert | hpe_morpheus_certificate                   |

### Removed Data Sources

The following hpegl data source has no HPE equivalent and is removed during migration:

- `hpegl_vmaas_morpheus_details` - the provider-chain broker data source; its role is replaced by the
  provider block transformation performed by [`migrate-providers`](./tfmigrator_migration.md).

## Migration Workflow

The end-to-end command workflow - installation, the one-shot `tfmigrator migrate` command, and the
individual `migrate-providers` -> `generate` -> `merge` -> `migrate-datasources` -> `cleanup-providers` steps - is documented
in the [tfmigrator guide](./tfmigrator_migration.md). tfmigrator auto-detects the hpegl provider in
your workspace, so the same commands apply here. For common issues, see
[tfmigrator Troubleshooting](./tfmigrator_troubleshooting.md).
