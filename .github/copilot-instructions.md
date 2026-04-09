# Copilot Instructions for terraform-provider-hpe

## Build, Test, and Lint

```bash
# Build
go build

# Lint (install golangci-lint v2 first with `make linter`)
make lint

# Run all tests (short mode, skips slow acceptance tests)
make test

# Run acceptance tests (framework resources only)
make testacc

# Run acceptance tests (sdkv2 resources only)
make testsdkv2

# Run a single test
cd morpheus/framework && TF_ACC=1 go test -v -run TestAccMorpheusUserExampleOk -timeout 10m ./...

# Generate docs (requires terraform CLI on PATH)
make docs

# Clean up leftover test resources from a Morpheus appliance
make sweep
```

Acceptance tests require a live Morpheus appliance. Set provider credentials via environment variables or a provider block. Use `-short` flag to skip slow acceptance tests during local iteration.

## Architecture

### Mux Provider with SubProvider Pattern

The provider multiplexes two Terraform plugin systems through `tf6muxserver`:

1. **Plugin Framework** (primary, for new resources) — `morpheus/framework/`
2. **Plugin SDK v2** (legacy) — `morpheus/sdkv2/`

Both are served simultaneously; the mux server routes to the correct one.

New resources/data sources **must** use Plugin Framework. The sdkv2 side exists for resources not yet migrated.

### SubProvider Interface

The top-level `provider/` package defines a generic `hpeProvider` that accepts pluggable `SubProvider` implementations via `provider/subprovider/subprovider.go`. Currently, the only subprovider is `morpheus` — but the architecture supports adding others (e.g., for non-Morpheus HPE services).

Each subprovider provides its own schema block, configure logic, resources, and data sources.

### Dual SDK

- **Framework resources** use the generated OpenAPI SDK: `github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk`
- **SDKv2 resources** use the legacy hand-written SDK: `github.com/HewlettPackard/hpe-morpheus-go-sdk/legacy`

The `go.work` file references a local checkout of the SDK at `../hpe-morpheus-go-sdk/oapigen` for development.

## Key Conventions

### Framework Resource File Structure

Each framework resource lives in its own package under `morpheus/framework/resources/<name>/`:

| File | Purpose |
|---|---|
| `resource.go` | CRUD implementation, `Metadata()`, `Schema()`. When this file exceeds 250 lines, split CRUD operations into separate files: `create.go`, `read.go`, `update.go`, `delete.go`. Keep `resource.go` for `Metadata()`, `Schema()`, and shared helpers. |
| `schema_gen.go` | **Generated** schema and model struct (by `terraform-plugin-framework-generator`). Do not edit. |
| `resource_test.go` or `<name>_test.go` | Acceptance tests (larger resources may split into `create_test.go`, `update_test.go`, `import_test.go`, etc.) |
| `sweep.go` | Test resource sweeper for cleanup |
| `<name>_example.go` | `//go:generate` directives for rendering example `.tf` files |
| `example.tf.tmpl` | Template for example Terraform configs used in tests and docs |

Data sources follow the same pattern under `morpheus/framework/datasources/<name>/`.

### Framework Resource Boilerplate

Every framework resource:
- Embeds `configure.ResourceWithMorpheusConfigure` (provides `NewClient(ctx)` for API access)
- Embeds `resource.Resource`
- Uses interface compliance assertions: `var _ resource.Resource = &Resource{}`
- Names resources as `hpe_morpheus_<name>`
- Gets an API client per operation via `r.NewClient(ctx)` (not cached on the struct)

### Error Handling

- **Framework**: Accumulate into `resp.Diagnostics` with `AddError(summary, detail)`, check `HasError()` before continuing. Use `errfmt.ErrMsg(err, httpResp)` to format API errors with HTTP response bodies.
- **SDKv2**: Return `diag.FromErr(err)` or `diag.Diagnostics` from CRUD functions.

### Testing Patterns

- All acceptance tests use `t.Parallel()` and `defer testhelpers.RecordResult(t)`
- Tests skip in short mode: `if testing.Short() { t.Skip(...) }`
- Use `acctest.RandomWithPrefix(t.Name())` for unique resource names
- Use `testhelpers.ProviderBlock()` for the provider HCL config
- Use `testhelpers.RenderExample()` to render `.tf.tmpl` templates into test configs
- Each test package has `TestMain` that calls `testhelpers.WriteMergedResults()` for result aggregation
- Sweepers clean up leftover test resources; registered via `resource.AddTestSweepers()`

### Import Organization

The linter enforces `goimports` with local prefix `github.com/HPE`. Imports should be grouped as:

```go
import (
    "stdlib"

    "third-party"

    "github.com/HPE/terraform-provider-hpe/..."
)
```

### Linting Rules

Configured in `.golangci.yml` (v2 format). Notable settings:
- Max line length: 120 characters
- `gofumpt` formatting enforced
- `nlreturn` requires blank lines before returns
- `Id` and `ID` are both valid (suppressed `var-naming` initialism warning)
- `funlen` and `err113` are relaxed in test files
- `dupl` and `goconst` are relaxed in generated files

### Documentation Generation

Docs are generated with `tfplugindocs`. The pipeline:
1. Example `.tf` files live in `examples/resources/<name>/` and `examples/data-sources/<name>/`
2. Templates live in `templates/` mirroring the docs structure
3. A custom `bin/render` tool (built from `cmd/render/`) renders `.tf.tmpl` files via `//go:generate` directives
4. `make docs` runs `go generate ./...` then `go generate` in `tools/` to produce `docs/`

### Type Conversion Utilities

`utils/convert/` provides helpers for converting between SDK types and Terraform framework types. Use these consistently rather than manual conversions.
