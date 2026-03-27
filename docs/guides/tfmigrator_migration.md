---
page_title: "tfmigrator: Morpheus to HPE"
subcategory: "Migration"
---


## What This Guide Covers

This document is a step-by-step guide for migrating Terraform configuration from Morpheus provider resources to HPE provider resources using `tfmigrator`.

- Preserve your variables, modules, expressions, and overall configuration structure.
- Update resource types and schemas to HPE provider resources.
- Run migration into a separate output directory so you can review before apply.

---

## Command Summary

`tfmigrator` migration flow is typically:

1. `generate` - extract resources from state, build `generated.tf` + `import.tf`
2. `merge` - merge generated resources into original config while preserving user intent
3. `terraform init / plan / apply --refresh-only / apply` - validate and apply in migrated output

Optional:

1. `clean` - run if you generated raw Terraform config externally and need cleanup rules - or after running `generate --no-cleanup`

---

## Required Inputs

Before starting, ensure you have:

1. Terraform state JSON, typically from:

```bash
terraform show -json > state.json
```

2. Your HPE provider setup in Terraform configuration (root and modules as needed).

3. Your original Terraform configuration directory.

Ensure provider setup is already in your Terraform files before running migration.

---

## Basic Provider Setup (Do This First)

Set up `required_providers` and provider blocks in your Terraform code before using `tfmigrator`. Update the `required_providers` block in root and modules as needed.

```hcl
terraform {
  required_providers {
    morpheus = {
      source  = "gomorpheus/morpheus"
      version = "0.14.1"
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

This setup is expected to exist before migration.


---

## Step 1: Generate Migration Artifacts

### Basic command

```bash
tfmigrator generate \
  --state state.json \
  --provider-config ./provider.tf
```

### Generate flags

- `--state, -s` Path to Terraform state JSON from `terraform show -json > state.json` (required)
- `--output-dir, -o` Output directory for generated files (default `./generated`)
- `--provider-config` Path to a Terraform provider context file containing `terraform` and provider blocks (and optional `variable` / `locals` blocks)
- `--provider-var-file` Optional var-file(s) for provider context, repeatable (default none)
- `--provider-var` Optional variable(s) passed to Terraform plan in `key=value` format, repeatable (default none)
- `--no-cleanup` Skip automatic cleanup pass (default `false`)
- `--dry-run` Preview output without writing files (default `false`)

### Generate output

- `generated/generated.tf`
  - Combined output for same-schema and different-schema resources
- `generated/import.tf`
  - Import blocks for migrated resources

- `--provider-config` should point to wherever your HPE provider is configured for the migration context (for example, `./provider.tf` or `./main.tf`). Alternatively - create a temporary file for this step (i.e. `.provider/provider.tf`) that follows the below format:

```hcl
terraform {
  required_providers {
    morpheus = {
      source  = "gomorpheus/morpheus"
      version = "0.14.1"
    }
    hpe = {
      source  = "HPE/hpe"
      version = "= 1.2.0"
    }
  }
}


provider "morpheus" {
  url      = "YOURURLHERE"
  username = "YOURUSERNAMEHERE"
  password = "YOURPASSWORDHERE"
}

provider "hpe" {
  morpheus {
    url      = "YOURURLHERE"
    username = "YOURUSERNAMEHERE"
    password = "YOURPASSWORDHERE"
  }
}
```
-> If variables are used in the `--provider-config` - the respective `variable` blocks are expected to be provided in that file. The values themselves can then be passed in via `--provider-var-file` or `--provider-var`.

- Cleanup is run automatically at the end of `generate` unless you pass `--no-cleanup`.

---

## Optional Step 1.5: Run `clean` Manually

Use this when:

- You ran `generate --no-cleanup`
- You want to inspect cleanup changes separately
- You are cleaning a generated file from another workflow

Recommended in-place cleanup:

```bash
tfmigrator clean \
  --input ./generated/generated.tf \
  --in-place
```

If you prefer a separate output file:

```bash
tfmigrator clean \
  --input ./generated/generated.tf \
  --output ./generated/cleaned.tf
```

If you write to `cleaned.tf`, use that file as the `--generated` input in Step 2.

---

## Step 2: Merge Into Original Config

### Recommended command

```bash
tfmigrator merge \
  --original . \  # Will find .tf files recursively from this directory
```

### Merge flags

- `--generated, -g` Path to generated Terraform file (default `./generated/generated.tf`)
- `--original, -o` Original Terraform configuration file(s) and/or directories, repeatable
- `--output-dir, -d` Output directory for merged configuration (default `./migrated`)
- `--imports-file` Explicit imports file path (auto-discovered from generated file location when omitted)
- `--var-file` Var-file(s) to copy into output directory, repeatable (default none)
- `--dry-run` Preview changes only (default `false`)
- `--no-color` Disable ANSI color in diff display (default `false`)
- `--verbose, -v` More detailed logs (default `false`)

### What merge does

- Same-schema resources are renamed while keeping existing expressions/variables where possible.
- Different-schema resources are merged from generated output.
- Result is written to your `--output-dir` so you can review before applying.

---

## Dry-Run Workflow (Recommended)

Run both phases in dry-run mode first:

```bash
tfmigrator generate \
  --state state.json \
  --provider-config ./provider.tf \
  --dry-run

# After an actual generation
tfmigrator merge \
  --original . \
  --dry-run
```

What to review in merge dry-run output:

1. `Merge Results` counts
2. `Merge Transformations` lines
3. Unified diff section (`Proposed Changes`)


---

## Dry-Run Output Example

Example excerpt from a successful `merge --dry-run` run:

```text
Merging configurations...

=== Merge Results ===

Resources added: 0
Resources merged: 2
Modules updated: 0

Dry run completed - no files were modified

=== Merge Transformations ===
~ Updated same-schema module resource morpheus_python_script_task.task -> hpe_morpheus_task_python_script (preserved original attributes)
~ Updated module resource morpheus_max_cores_policy.user_policy -> hpe_morpheus_policy (affects all instances)

=== Proposed Changes (Unified Diff) ===
--- user-policy/main.tf
+++ user-policy/main.tf
@@
- resource "morpheus_max_cores_policy" "user_policy" {
+ resource "hpe_morpheus_policy" "user_policy" {
    description = "Max cores policy for user ${var.user_name}"
    enabled     = var.enabled
-   user_id     = var.user_id
+   associated_resource_id   = var.user_id
+   associated_resource_type = "User"
  }
```

- In terminal output, removed lines are red and added lines are green.

---

## Apply Workflow

After dry-run looks correct:

1. Run merge without `--dry-run`
2. Move to migrated directory
3. Run standard Terraform validation/apply flow

```bash
tfmigrator merge \
  --original . \
  --generated ./generated/generated.tf \
  --output-dir ./migrated

cd ./migrated
terraform init
terraform plan
# Check that everything is just imports
terraform apply --refresh-only
terraform apply
```

---

## Common Patterns

### Using provider var files during generate

```bash
tfmigrator generate \
  --state state.json \
  --provider-config ./provider.tf \
  --provider-var-file terraform.tfvars \
  --provider-var-file secrets.auto.tfvars \
  --output-dir ./generated
```

Note this currently requires the used `variable` blocks to be in the provider.tf file.

### Copying var files into migrated output during merge

```bash
tfmigrator merge \
  --original . \
  --generated ./generated/generated.tf \
  --var-file terraform.tfvars \
  --var-file env/dev.tfvars \
  --output-dir ./migrated
```

### Explicit imports file (if auto-discovery is not desired)

```bash
tfmigrator merge \
  --original . \
  --generated ./generated/generated.tf \
  --imports-file ./generated/import.tf \
  --output-dir ./migrated
```

---

## Troubleshooting

For common `generate`, `merge`, and resource-specific migration issues, see [tfmigrator Troubleshooting](./tfmigrator_troubleshooting.md).

---

## What tfmigrator Preserves

- Variable references (`var.*`)
- Expressions and interpolation
- Core block structure where possible
- Modules
- Comments in the configuration (excluding inline comments from resources with different schema)
