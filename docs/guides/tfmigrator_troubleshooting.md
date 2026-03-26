---
page_title: "tfmigrator: Troubleshooting"
subcategory: "Migration"
---

# tfmigrator Troubleshooting

This guide covers common issues that can appear during `tfmigrator generate`, `tfmigrator merge`, and follow-up Terraform validation and apply steps.

---

## Generate Issues

### Provider context errors during generate

Confirm `--provider-config` points to a real `.tf` file that includes your HPE provider configuration. This provider context is used when `tfmigrator` runs Terraform planning to generate configuration.

If your provider config uses variables, ensure the corresponding `variable` blocks are present in that file and pass values with `--provider-var-file` and/or `--provider-var`.
```
terraform {
  required_providers {
    morpheus = {
      source  = "gomorpheus/morpheus"
      version = "0.13.1"
    }
    hpe = {
      source  = "HPE/hpe"
      version = "= 1.2.0"
    }
  }
}


provider "morpheus" {
  url      = var.morpheus_url
  username = var.morpheus_username
  password = var.morpheus_password
}

provider "hpe" {
  morpheus {
    url      = var.morpheus_url
    username = var.morpheus_username
    password = var.morpheus_password
  }
}

variable "morpheus_url" {
  type = string
}

variable "morpheus_username" {
  type = string
}

variable "morpheus_password" {
  type      = string
  sensitive = true
}
```

### Unauthorized or invalid token errors

A failed generate can return errors similar to:

```text
# group 41596 GET failed: 401 (Unauthorized):
# {"error":"invalid_token","error_description":"Invalid access token:\n# BADTOKEN"}
```

Fix:

- Review the credentials or token used by the provider configuration.
- Confirm the URL, username, password, or token in `--provider-config` are valid for the Morpheus environment you are targeting.
- If credentials are supplied via variables, confirm the correct values are being passed through `--provider-var-file` or `--provider-var`.

### Validation errors before cleanup

During `generate`, Terraform validation errors may appear in terminal output before cleanup runs. These are often expected for raw generated configuration and can usually be ignored until the cleanup pass has completed.

If you used `--no-cleanup`, run cleanup separately before reviewing the generated file:

```bash
tfmigrator clean \
  --input ./generated/generated.tf \
  --in-place
```

### Generated file not found

If `merge` reports that the generated file cannot be found, confirm the path passed to `--generated` and ensure `generate` completed successfully first.

---

## Merge Issues

### Wrong provider source in modules

If Terraform reports an error such as:

```text
Error: Failed to query available provider packages

Could not retrieve the list of available versions for provider hashicorp/hpe:
provider registry registry.terraform.io does not have a provider named
registry.terraform.io/hashicorp/hpe
```

Fix:

Ensure `required_providers` is set correctly in the root module and in any child modules that declare providers:

```hcl
terraform {
  required_providers {
    morpheus = {
      source  = "gomorpheus/morpheus"
      version = "0.13.1"
    }
    hpe = {
      source  = "HPE/hpe"
      version = "= 1.2.0"
    }
  }
}
```

---

## hpe_morpheus_instance Issues

### nested_virtualization validation errors

Generated or merged configuration for `hpe_morpheus_instance` can fail with errors such as:

```text
Attribute config_vmware.nested_virtualization value must be one of: ["on" "off"], got: "0"
```

Fix:
- Ensure the HPE provider is at version `1.2.0` or later.
- Alternatively, set `nested_virtualization` values to `"on"` or `"off"` in the migrated configuration.


This issue is most commonly seen in generated or merged VMware instance configuration.

### Tag casing or metadata validation errors

In some environments, tag validation can fail with errors similar to:

```json
{"errors":{"metadata":"tag:
'APPLICATION' must have a valid property","APPLICATION":"must have a valid property"}}
```

Fix:

Review the generated or merged `tags` block and update tag names to match the exact casing and properties expected by the target environment.

Example adjustment:

```hcl
tags = [
  {
    name  = "APPLICATION"
    value = "ubuntu"
  },
  {
    name  = "name"
    value = "ubuntutf"
  },
]
```

If the original configuration used title-cased or additional tag names, do not assume those values will remain valid after migration. Reconcile the final tag set against what Morpheus accepts for that environment.
