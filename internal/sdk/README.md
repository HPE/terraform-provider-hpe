# Vendored Morpheus SDK (`internal/sdk`)

This directory contains the Morpheus SDK, vendored **in-tree** as part of the
`github.com/HPE/terraform-provider-hpe` module. It was previously the external
`github.com/HewlettPackard/hpe-morpheus-go-sdk` repository (now archived).

| Path | Package | Origin | Edit? |
|------|---------|--------|-------|
| `oapigen/` | `sdk` | **Generated** by OpenAPI Generator from the Morpheus OpenAPI spec | No — regenerated |
| `legacy/`  | `morpheus` | Hand-written SDK (used by the SDKv2 resources/data sources) | Yes — maintained here |

## Import paths

```go
import sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen" // framework resources
import morpheus "github.com/HPE/terraform-provider-hpe/internal/sdk/legacy" // sdkv2 resources
```

## Regenerating `oapigen`

`oapigen` is generated from the Morpheus OpenAPI specification by an internal
code-generation pipeline (the same pipeline that produces the `schema_gen.go`
resource/data-source schemas) and delivered into this repo via a manual task.

Do not hand-edit files under `oapigen/` — your changes will be overwritten on the
next regeneration.

## Tooling notes

- This tree is **excluded from `make lint`** (generated/third-party code) but is still
  built and type-checked as part of the module.
- There is no `go.work` or `replace` directive — it is a normal in-module package.
