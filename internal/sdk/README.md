# Vendored Morpheus SDK (`internal/sdk`)

This directory contains the Morpheus SDK, vendored **in-tree** as part of the
`github.com/HPE/terraform-provider-hpe` module. It was previously the external
`github.com/HewlettPackard/hpe-morpheus-go-sdk` repository (now archived).

| Path | Package | Origin | Built? | Edit? |
|------|---------|--------|--------|-------|
| `oapigen/` | `sdk` | **Generated** by OpenAPI Generator from the Morpheus OpenAPI spec, then **pruned** to what the provider reaches | Yes — the import path | No — regenerated |
| `sdkfull/` | `sdkfull` | The same generator output, **complete** (every operation and model), kept as a reference | **No** — every file carries `//go:build sdkfull` | No — regenerated |
| `legacy/`  | `morpheus` | Hand-written SDK (used by the SDKv2 resources/data sources) | Yes | Yes — maintained here |

## Import paths

```go
import sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen" // framework resources
import morpheus "github.com/HPE/terraform-provider-hpe/internal/sdk/legacy" // sdkv2 resources
```

`sdkfull` is never imported. It exists so the whole API is visible in-tree; the
code you call always comes from `oapigen`.

## Regenerating `oapigen` and `sdkfull`

Both are generated from the Morpheus OpenAPI specification by an internal
code-generation pipeline (the same pipeline that produces the `schema_gen.go`
resource/data-source schemas) and delivered into this repo together via a
manual task.

Do not hand-edit files under `oapigen/` or `sdkfull/` — your changes will be
overwritten on the next regeneration.

## `oapigen` is pruned; `sdkfull` is the complete copy

The generator produces a client for the *entire* Morpheus API (≈1 M lines, 9.5 k
files). The provider uses roughly a third of the operations, and compiling all
of it as one package cost ~8 GB of RAM and most of the build time. So:

- **`oapigen/`** — the import path — is **pruned to what the provider reaches**:
  every model and API operation that provider code, its unit tests, its
  acceptance tests or its sweepers reference, plus everything those
  transitively need. Unused operations and models are removed; the per-model
  `Nullable<X>` wrappers the generator emits are removed unless a
  `nullable: true` property actually uses them. This is purely a compile-time
  saving (~3× faster, ~4× less memory); runtime behaviour is identical, because
  Go cannot reach a type that nothing references.
- **`sdkfull/`** is the **same generator output, complete**, with two mechanical
  changes: the package clause is `package sdkfull`, and every file opens with
  `//go:build sdkfull`. That tag is set by nothing, so `go build ./...`,
  `go vet ./...`, `go test ./...`, `go list ./...`, gopls and `go mod tidy` all
  skip the directory — it costs the build nothing. It is the in-repo reference
  for the whole API, it is what the internal tooling's SDK index is built from,
  and it is the source the pruned tree is restored from. To type-check it on
  purpose: `go vet -tags sdkfull ./internal/sdk/sdkfull`.

Consequences for day-to-day work:

- **If the operation or model you need is missing from `oapigen`**, look in
  `sdkfull/` — it is there; it was pruned, not left out of the spec. The normal
  fix is the pipeline: its `sdk-to-provider` task delivers the full SDK into
  `oapigen`, you write your code (tests included), then `prune-sdk` cuts it back
  down and you commit the result. Anything your new code or tests reference is
  kept from then on.
- **Without the pipeline** (a standalone clone): copy the `api_*.go` and
  `model_*.go` files you need from `sdkfull/` into `oapigen/`, change their
  package clause back to `package sdk` and delete the `//go:build sdkfull`
  line; the compiler will tell you which models are still missing. The next
  pipeline run normalises whatever you copied.
- **Do not hand-write SDK code**, and do not add files to `oapigen/` that have
  no counterpart in `sdkfull/` — the next regeneration would drop them.
- The pipeline's `check-sdk-prune` task verifies that the committed `oapigen`
  compiles everywhere (all packages and test binaries, sweepers included, on
  Linux and Windows), is fully pruned (a prune would change nothing), and that
  `sdkfull/` is a faithful, still-non-buildable superset of it.

## Reproducibility

The pipeline is deterministic: given the same OpenAPI spec commit and its
pinned tools (the generator, the formatters, and the Go toolchain whose
`go/format` shapes every file), it regenerates both trees byte-for-byte, which
has been verified end-to-end. The trees themselves carry no record of which
spec commit they were generated from; when that matters, it is stated in the
pull request that regenerated them.

Two properties of the prune are worth knowing when you add code:

- The reachability analysis loads the provider as `GOOS=linux GOARCH=amd64`
  with every `*_test.go` and `-tags sweep`, regardless of the machine it runs
  on. **SDK calls must live in portable, untagged code (or `sweep`-tagged
  code).** A file behind any other `//go:build` tag, or a platform-specific
  file, would be invisible to the analysis, so the pruner refuses to run at all
  if it finds one — the fix is to move the SDK use, or to have the pipeline
  enable the tag.
- Formatting is pinned to the Go version in this module's `go` directive.
  Bumping that directive may reformat the SDK trees on the next regeneration;
  that diff is expected and mechanical.

## Tooling notes

- This whole tree is **excluded from `make lint`** (generated/third-party code).
  `oapigen/` and `legacy/` are still built and type-checked as part of the
  module; `sdkfull/` is not (see above).
- There is no `go.work` or `replace` directive — these are normal in-module packages.
