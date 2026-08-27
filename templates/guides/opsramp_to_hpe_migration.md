---
page_title: "OpsRamp to HPE"
subcategory: "Migration"
---

# Migrating from OpsRamp Provider to HPE Provider

This guide describes how to migrate existing Terraform configurations from the standalone OpsRamp provider (`HPE/opsramp`) to the unified HPE provider (`HPE/hpe`). The migration involves two changes: updating the provider declaration and updating all resource/data source addresses.

## Prerequisites

- Terraform >= 1.5.0 (required for `import` blocks)
- HPE provider >= 2.0.0
- Access to your existing OpsRamp resource IDs (from state or the OpsRamp API)

## Step 1: Update the Provider Declaration

The OpsRamp provider is now embedded as a nested configuration block within the HPE provider.

**Before (standalone OpsRamp provider):**

```HCL
terraform {
  required_providers {
    opsramp = {
      source  = "HPE/opsramp"
      version = ">= 1.5.0"
    }
  }
}

provider "opsramp" {
  client_id     = "abcdefghijklmnopqrstuvwxyz123456"
  client_secret = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890ab"
  endpoint      = "tenant.api.pov.opsramp.com"
  tenant        = "abcdefgh-1234-5678-90ab-cdefghijklmn"
}
```

**After (unified HPE provider):**

```HCL
terraform {
  required_providers {
    hpe = {
      source  = "HPE/hpe"
      version = ">= 2.0.0"
    }
  }
}

provider "hpe" {
  opsramp {
    client_id     = "abcdefghijklmnopqrstuvwxyz123456"
    client_secret = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890ab"
    endpoint      = "tenant.api.pov.opsramp.com"
    tenant        = "abcdefgh-1234-5678-90ab-cdefghijklmn"
  }
}
```

OpsRamp credentials are now specified inside an `opsramp` block within the `provider "hpe"` declaration. All other provider attributes remain unchanged.

## Step 2: Update Resource and Data Source Addresses

All resource and data source type names have been prefixed with `hpe_` to reflect the new provider namespace. The naming convention changes from `opsramp_<resource>` to `hpe_opsramp_<resource>`.

### Resources

| OpsRamp Provider (old) | HPE Provider (new) |
|---|---|
| `opsramp_alert_correlation_policy` | `hpe_opsramp_alert_correlation_policy` |
| `opsramp_alert_escalation_policy` | `hpe_opsramp_alert_escalation_policy` |
| `opsramp_alert_prediction_policy` | `hpe_opsramp_alert_prediction_policy` |
| `opsramp_client` | `hpe_opsramp_client` |
| `opsramp_credential_set` | `hpe_opsramp_credential_set` |
| `opsramp_custom_integration` | `hpe_opsramp_custom_integration` |
| `opsramp_device_group` | `hpe_opsramp_device_group` |
| `opsramp_first_response_policy` | `hpe_opsramp_first_response_policy` |
| `opsramp_integration` | `hpe_opsramp_integration` |
| `opsramp_integration_app` | `hpe_opsramp_integration_app` |
| `opsramp_integration_config` | `hpe_opsramp_integration_config` |
| `opsramp_integration_event` | `hpe_opsramp_integration_event` |
| `opsramp_kb_article` | `hpe_opsramp_kb_article` |
| `opsramp_kb_category` | `hpe_opsramp_kb_category` |
| `opsramp_log_alert_definition` | `hpe_opsramp_log_alert_definition` |
| `opsramp_management_profile` | `hpe_opsramp_management_profile` |
| `opsramp_metric_alert_definition` | `hpe_opsramp_metric_alert_definition` |
| `opsramp_permission_set` | `hpe_opsramp_permission_set` |
| `opsramp_resource` | `hpe_opsramp_resource` |
| `opsramp_role` | `hpe_opsramp_role` |
| `opsramp_schedule_maintenance` | `hpe_opsramp_scheduled_maintenance` |
| `opsramp_script` | `hpe_opsramp_script` |
| `opsramp_script_category` | `hpe_opsramp_script_category` |
| `opsramp_servicedesk_business_impact` | `hpe_opsramp_servicedesk_business_impact` |
| `opsramp_servicedesk_category` | `hpe_opsramp_servicedesk_category` |
| `opsramp_servicedesk_urgency` | `hpe_opsramp_servicedesk_urgency` |
| `opsramp_servicemap` | `hpe_opsramp_servicemap` |
| `opsramp_servicemap_link` | `hpe_opsramp_servicemap_link` |
| `opsramp_site` | `hpe_opsramp_site` |
| `opsramp_user` | `hpe_opsramp_user` |
| `opsramp_user_group` | `hpe_opsramp_user_group` |

### Data Sources

| OpsRamp Provider (old) | HPE Provider (new) |
|---|---|
| `data.opsramp_custom_event_alert_source` | `data.hpe_opsramp_custom_event_alert_source` |
| `data.opsramp_resource_lookup` | `data.hpe_opsramp_resource_lookup` |
| `data.opsramp_role` | `data.hpe_opsramp_role` |
| `data.opsramp_servicedesk_business_impact` | `data.hpe_opsramp_servicedesk_business_impact` |
| `data.opsramp_servicedesk_category` | `data.hpe_opsramp_servicedesk_category` |
| `data.opsramp_servicedesk_urgency` | `data.hpe_opsramp_servicedesk_urgency` |
| `data.opsramp_tenant` | `data.hpe_opsramp_tenant` |

## Step 3: Import Existing Resources

To migrate resources that were previously managed by the OpsRamp provider into the HPE provider, use Terraform [import blocks](https://developer.hashicorp.com/terraform/language/import). This allows you to adopt existing infrastructure without recreating it.

~> **Note:** The [Bulk Import](https://developer.hashicorp.com/terraform/language/v1.14.x/import/bulk?page=import&page=bulk) feature requires provider support for the `List` RPC, which the HPE provider does not currently implement. Bulk import support is planned for a future release.

### Using Import Blocks

Add `import` blocks to your configuration specifying the resource ID and the new HPE provider address:

```HCL
import {
  id = "125623"
  to = hpe_opsramp_resource.my_server
}

import {
  id = "DGP-abc12345-def6-7890-ghij-klmnopqrstuv"
  to = hpe_opsramp_device_group.production
}
```

### Migration Workflow

1. Remove the old `opsramp` provider from `required_providers` and add `hpe`.
2. Replace `provider "opsramp" { ... }` with the nested `provider "hpe" { opsramp { ... } }` block.
3. Update all resource and data source type names from `opsramp_*` to `hpe_opsramp_*`.
4. Add `import` blocks for each resource, using the IDs from your existing state (`terraform state list` / `terraform state show`).
5. Run `terraform init` to install the HPE provider.
6. Run `terraform plan` to verify the import plan shows no unexpected changes.
7. Run `terraform apply` to import the resources into the new state.
8. Confirm with `terraform plan` that no further changes are detected.

-> **Tip:** You can extract resource IDs from your existing OpsRamp state file before switching providers using `terraform state show <resource_address>` to view the `id` attribute of each resource.
