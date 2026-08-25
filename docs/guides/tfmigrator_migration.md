---
page_title: "tfmigrator"
subcategory: "Migration"
---


## What This Guide Covers

This document is a step-by-step guide for migrating Terraform configuration to HPE provider resources using `tfmigrator`. It applies both to the [Morpheus](./morpheus_to_hpe_migration.md) and [hpegl](./hpegl_to_hpe_migration.md) providers as migration sources - `tfmigrator` auto-detects the source provider in your workspace.

- Preserve your variables, modules, expressions, and overall configuration structure.
- Update resource types and schemas to HPE provider resources.
- Transform the source provider block (Morpheus or hpegl) into an `hpe` provider block.
- Run migration into a separate output directory so you can review before apply.

`tfmigrator` supports two ways of running the migration:

- **Full pipeline** - a single `migrate` command that runs every step end to end.
- **Individual steps** - run `migrate-providers`, `generate`, `merge`, and `migrate-datasources` yourself for finer control.

Both approaches are documented below.

---

## Installation

The latest release of tfmigrator can be found on [`https://github.com/HPE/terraform-provider-hpe/releases`](https://github.com/HPE/terraform-provider-hpe/releases). You can download the appropriate version for your operating system using a command line shell or a browser.

This can be useful in environments that do not allow direct access to the Internet.

### Install Scripts

See [linux](https://github.com/HPE/terraform-provider-hpe/blob/main/scripts/install-tfmigrator.sh), [windows](https://github.com/HPE/terraform-provider-hpe/blob/main/scripts/install-tfmigrator-windows.ps1), and [macOS](https://github.com/HPE/terraform-provider-hpe/blob/main/scripts/install-tfmigrator-macos.sh) install scripts that will download the latest release of tfmigrator and install it in the appropriate location for your operating system.

### Linux

The following examples use Bash on Linux (x64).

1. On a Linux operating system with Internet access, download tfmigrator from GitHub using the shell.

   ```console
   RELEASE=x.y.z
   wget -q https://github.com/HPE/terraform-provider-hpe/releases/download/v${RELEASE}/migration_tool_${RELEASE}_linux_amd64.zip
   ```

2. Extract the archive.

   ```console
   unzip migration_tool_${RELEASE}_linux_amd64.zip
   ```

3. Move the binary to a directory in your PATH.

   ```console
   sudo mv migration_tool_v${RELEASE} /usr/local/bin/tfmigrator
   chmod +x /usr/local/bin/tfmigrator
   ```

4. Verify the installation.

   ```console
   tfmigrator --version
   ```

### macOS

The following example uses Zsh (default) on macOS (Apple Silicon).

1. On a macOS operating system with Internet access, install wget with [Homebrew](https://brew.sh).

   ```console
   brew install wget
   ```

2. Download tfmigrator from GitHub using the shell.

   ```console
   RELEASE=x.y.z
   wget -q https://github.com/HPE/terraform-provider-hpe/releases/download/v${RELEASE}/migration_tool_${RELEASE}_darwin_arm64.zip
   ```

3. Extract the archive.

   ```console
   unzip migration_tool_${RELEASE}_darwin_arm64.zip
   ```

4. Move the binary to a directory in your PATH.


   ```console
   sudo mv migration_tool_v${RELEASE} /usr/local/bin/tfmigrator
   chmod +x /usr/local/bin/tfmigrator
   ```

5. Verify the installation.

   ```console
   tfmigrator --version
   ```

### Windows

The following examples use PowerShell on Windows (x64).

1. On a Windows operating system with Internet access, download tfmigrator using PowerShell.

   ```powershell
   Set-Variable -Name "RELEASE" -Value "x.y.z"
   Invoke-WebRequest https://github.com/HPE/terraform-provider-hpe/releases/download/v${RELEASE}/migration_tool_${RELEASE}_windows_amd64.zip -outfile migration_tool_${RELEASE}_windows_amd64.zip
   ```

2. Extract the archive.

   ```powershell
   Expand-Archive migration_tool_${RELEASE}_windows_amd64.zip
   cd migration_tool_${RELEASE}_windows_amd64
   ```

3. Move the binary to a directory in your PATH.

   ```powershell
   Move-Item migration_tool_v${RELEASE}.exe C:\Windows\System32\tfmigrator.exe
   ```

4. Verify the installation.

   ```powershell
   tfmigrator --version
   ```

---

## What tfmigrator Preserves

- Variable references (`var.*`)
- Expressions and interpolation
- Core block structure where possible
- Modules
- Comments in the configuration (excluding inline comments from resources with different schema)

---

## Command Summary

`tfmigrator` provides the following commands:

- `migrate` - run the full pipeline (`migrate-providers` → `generate` → `merge` → `migrate-datasources`) in one command
- `migrate-providers` - transform source (Morpheus / hpegl) provider blocks into `hpe` provider blocks
- `generate` - extract resources from state, build `generated.tf` + `import.tf`
- `merge` - merge generated resources into original config while preserving user intent
- `migrate-datasources` - rewrite `hpegl_vmaas_*` / `morpheus_*` data source blocks to `hpe_morpheus_*` equivalents
- `clean` - post-process a generated file with cleanup rules (run automatically by `generate`)

A typical migration is either:

- **Full pipeline:** `migrate` followed by `terraform init / plan / apply --refresh-only / apply` in the migrated output, or
- **Individual steps:** `migrate-providers` → `generate` → `merge` → `migrate-datasources`, then the same Terraform validation/apply flow.

Optional:

- `clean` - run if you generated raw Terraform config externally and need cleanup rules - or after running `generate --no-cleanup`

---

## Required Inputs

Before starting, ensure you have:

1. Terraform state JSON, typically from:

```bash
terraform show -json > state.json
```

2. Your original Terraform configuration directory, containing the source (`morpheus` / `hpegl`) provider configuration.

-> The `migrate` and `migrate-providers` commands generate the `hpe` provider block for you automatically from your existing configuration - you do **not** need to hand-author an `hpe` provider block for the recommended workflow. A hand-authored provider block is only required if you run `generate` on its own without `migrate-providers` (see [Step 2](#step-2-generate-migration-artifacts)).

---

## Full Pipeline (`migrate`)

The `migrate` command runs the entire migration end to end. It auto-detects the source provider
(Morpheus or hpegl) in your workspace and runs `migrate-providers` → `generate` → `merge` →
`migrate-datasources` in sequence.

### Basic command

```bash
tfmigrator migrate \
  --state state.json \
  --working-dir .
```

### Migrate flags

- `--state, -s` Path to Terraform state JSON from `terraform show -json > state.json` (required)
- `--working-dir, -w` Workspace directory containing the original Terraform configuration (default `.`)
- `--output-dir, -o` Root output directory (default `./migrated`)
- `--terraform-path` Path to the `terraform` binary (default `terraform`)
- `--var-file` Variable file(s) passed to Terraform plan during generation, repeatable (default none)
- `--var` Individual variable value(s) passed to Terraform plan in `key=value` format, repeatable (default none)
- `--original` Original Terraform file(s)/directory(ies) to include in the final merge, repeatable (defaults to `--working-dir`)
- `--dry-run` Show what would be done without writing files (default `false`)

### Migrate output

`migrate` writes into subdirectories of `--output-dir`:

- `migrated/provider-config/` - generated `hpe` provider blocks and `credentials.auto.tfvars` (from `migrate-providers`)
- `migrated/generated/` - `generated.tf` + `import.tf` (from `generate`)
- `migrated/final/` - the merged, ready-to-review configuration (from `merge`, then rewritten in place by `migrate-datasources`)

The `migrate` command derives the provider context and credentials from your existing configuration
automatically, so you do not need to hand-author a `--provider-config` file as you would when running
`generate` on its own.

### Next steps after `migrate`

```bash
cd migrated/final
terraform init
terraform fmt
terraform validate
terraform plan
# Check that everything is just imports
terraform apply --refresh-only
terraform apply
```

---

# Individual Steps

The remaining sections document each step individually. Run these when you want finer control than
the one-shot `migrate` command, for example to review the transformed provider blocks or the generated
configuration before merging.

The individual flow is: **Step 1** `migrate-providers` → **Step 2** `generate` → **Step 3** `merge`
→ **Step 4** `migrate-datasources`.

---

## Step 1: Transform Provider Blocks (`migrate-providers`)

`migrate-providers` scans your workspace for source (Morpheus / hpegl) provider configurations and
generates the equivalent `hpe` provider blocks, along with credential `variable` blocks and a
`credentials.auto.tfvars` file. This is especially important for hpegl, where the GreenLake IAM
configuration is transformed into a `morpheus { pce_identity { ... } }` block. See the
[hpegl migration guide](./hpegl_to_hpe_migration.md) for the details of that transformation.

### Basic command

```bash
tfmigrator migrate-providers \
  --working-dir . \
  --output-dir ./provider-config
```

### Migrate-providers flags

- `--working-dir, -w` Workspace directory to scan for provider configurations (default `.`)
- `--output-dir, -o` Output directory for generated provider files (default `./provider-config`)
- `--config, -c` Migration config: `auto` to detect from the workspace, an embedded config name (`hpegl`, `morpheus`), or a filesystem path (default `auto`)
- `--dry-run` Show what would be generated without writing files (default `false`)

The generated `./provider-config` directory can then be passed to `generate` via `--provider-config`,
and `credentials.auto.tfvars` supplied via `--provider-var-file`.

---

## Step 2: Generate Migration Artifacts

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

- Cleanup is run automatically at the end of `generate` unless you pass `--no-cleanup`.

### Provider context (standalone `generate` only)

`generate` runs Terraform internally to produce the configuration, so it needs a provider context via `--provider-config`.

- If you ran `migrate-providers` first (or used the full `migrate` command), this is already handled for you - point `--provider-config` at the generated `provider-config/` directory. No manual provider block is required.
- If you are running `generate` **on its own**, supply a `--provider-config` file that configures the `hpe` provider (and the source provider used for planning). This can be an existing file such as `./provider.tf` or `./main.tf`, or a temporary file created just for this step.

A minimal `--provider-config` file using variables:

```hcl
terraform {
  required_providers {
    morpheus = {
      source  = "gomorpheus/morpheus"
      version = "0.14.1"
    }
    hpe = {
      source  = "HPE/hpe"
      version = ">= 2.0.0"
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

-> If variables are used in the `--provider-config`, the respective `variable` blocks are expected to be provided in that file. The values themselves can then be passed in via `--provider-var-file` or `--provider-var`. Alternatively, for a quick throwaway context you can inline literal values in place of the `var.*` references (for example a `.provider/provider.tf` with `url = "YOURURLHERE"`, `username = "YOURUSERNAMEHERE"`, `password = "YOURPASSWORDHERE"`) and skip the `variable` blocks.

---

## Optional Step 2.5: Run `clean` Manually

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

If you write to `cleaned.tf`, use that file as the `--generated` input in Step 3.

---

## Step 3: Merge Into Original Config

### Recommended command

```bash
tfmigrator merge \
  --original .    # Will find .tf files recursively from this directory
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

## Step 4: Rewrite Data Sources (`migrate-datasources`)

`migrate-datasources` rewrites `hpegl_vmaas_*` and `morpheus_*` data source blocks to their
`hpe_morpheus_*` equivalents, updates references to them, and removes data sources that have no HPE
equivalent (for example `hpegl_vmaas_morpheus_details`).

### Basic command

```bash
tfmigrator migrate-datasources \
  --input ./migrated \
  --in-place
```

### Migrate-datasources flags

- `--input, -i` Input directory/files to process, repeatable (default `.`)
- `--output-dir, -o` Output directory for rewritten files (default `./migrated`)
- `--config, -c` Config source: `auto` (all embedded sets), a single embedded name (`hpegl`, `morpheus`), or a filesystem path (default `auto`)
- `--in-place` Modify files in place instead of writing to `--output-dir` (default `false`)
- `--dry-run` Preview changes without writing (default `false`)
- `--no-color` Disable ANSI color in the dry-run diff (default `false`)
- `--verbose, -v` Show detailed transformation info (default `false`)

Point `--input` at your merged output (for example `./migrated` from Step 3) and use `--in-place` to
rewrite the data source blocks there. The one-shot `migrate` command runs this step automatically
against `migrated/final/`.

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

## Module Handling

As part of `generate` - `tfmigrator` adds comment annotations to generated resources if they were generated from a module:
```hcl
# Original module address: module.user_policy["bob.user"].morpheus_max_cores_policy.user_policy
resource "hpe_morpheus_policy" "user_policy_module_user_policy__bob_user" {
  associated_resource_id   = 5
  associated_resource_type = "User"
  config_max_cores = {
    max_cores = "35"
  }
  description = "Max cores policy for user bob.user"
  enabled     = true
  name        = "max_cores_policy_bob.user"
  policy_type = {
    code = "maxCores"
  }
}
```

This allows it to identify which resources belong to modules on a merge and update the module accordingly.
```
=== Proposed Changes (Unified Diff) ===
--- user-policy/main.tf
+++ user-policy/main.tf
@@ -1,13 +1,19 @@
-  resource "morpheus_max_cores_policy" "user_policy" {
+  resource "hpe_morpheus_policy" "user_policy" {
     name = "max_cores_policy_${var.user_name}"
+    policy_type = {
+      code = "maxCores"
+    }
+    associated_resource_id   = var.user_id
+    associated_resource_type = "User"
+    config_max_cores = {
       max_cores = var.max_cores
+    }
     description = "Max cores policy for user ${var.user_name}"
     enabled     = var.enabled
-    scope       = "user"
-    user_id     = var.user_id
   }
   
-  resource "morpheus_python_script_task" "task" {
+  
+  resource "hpe_morpheus_task_python_script" "task" {
     name                = var.task_name
     code                = var.task_code
     labels              = var.labels
```

-> Ensure that the module configuration (located in `./user-policy` in the above example) is included as part of what is passed into `merge` as the `--original` configuration. Note that `--original` searches for Terraform configuration files recursively so passing in `--original .` while in the working directory where all modules are available as subdirectories is generally sufficient.  

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
