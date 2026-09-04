#(C) Copyright 2025-2026 Hewlett Packard Enterprise Development LP
#
# Note: this Makefile works with GNUMake and BSDMake
#

.PHONY: build linter lint lint-ci test test-json docs sweep build-render-tool

# Usage: make sweep SWEEP=resource_name SWEEP_SYSTEMS=systemname SWEEP_PREFIX=prefix
# SWEEP_PREFIX optionally overrides the resource-name prefix the sweeper matches
# (default: TestAccMorpheus). Leave unset for normal test cleanup.
SWEEP ?= all
SWEEP_SYSTEMS ?= all
SWEEP_PREFIX ?=
SWEEP_RUN_ARGS = $(if $(filter all,$(SWEEP)),,-sweep-run=$(SWEEP))

# Per-package timeout for the acceptance test targets (`test`, `test-json`).
# Acceptance tests provision real infrastructure; on slower backends a single
# instance can take ~25m to provision, and packages such as the instance
# resource run several such tests in parallel (with multi-step update tests
# provisioning twice), so the default needs generous headroom. Override per run,
# e.g. `make test TEST_TIMEOUT=180m`.
TEST_TIMEOUT ?= 120m

build:
	go build

linter:
	go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.11.3

# The vendored Morpheus SDK (internal/sdk/{oapigen,legacy}) is generated /
# hand-written third-party code that is not subject to the provider's lint
# rules. It is still type-checked as a dependency, but excluded from the lint
# target set (linting ~9k generated files is both wrong and prohibitively slow).
LINT_DIRS = $$(packages="$$(go list -f '{{.ImportPath}} {{.Dir}}' ./...)" && printf '%s\n' "$$packages" | awk 'NF && $$1 !~ /\/internal\/sdk(\/|$$)/ { print $$2 }' | tr '\n' ' ')

lint:
	set -e; \
	dirs="$(LINT_DIRS)"; \
	test -n "$$dirs" || { echo "no lint targets found (go list produced nothing)" >&2; exit 1; }; \
	golangci-lint run $$dirs

# The set CI runs (see .github/workflows/lint.yaml). The SSA/fact-based linters
# (unused, staticcheck, govet) dominate golangci-lint's peak memory and pushed
# the ~7 GB CI runner into an OOM kill (exit 143) on the full package set. gosec
# has been removed due to timeout issues. CI keeps the cheap AST linters and
# drops those three and gosec, which cuts peak memory and runs ~3x faster.
# `make lint` above keeps FULL coverage locally; this target reproduces exactly
# what CI runs. Keep the --disable list in sync with the lint workflow.
lint-ci:
	set -e; \
	dirs="$(LINT_DIRS)"; \
	test -n "$$dirs" || { echo "no lint targets found (go list produced nothing)" >&2; exit 1; }; \
	golangci-lint run --disable=gosec,unused,staticcheck,govet $$dirs

test:
	pkgs=$$(go list ./... | grep -v '/internal/sdk'); \
	env TF_ACC=1 \
	go test -v -cover -count 1 -timeout $(TEST_TIMEOUT) $$pkgs

# Same as `test` but emits machine-readable `go test -json` on stdout (and
# nothing else, so the stream stays valid JSON). Used by the nightly runner to
# capture a full, parseable log including API traces. The leading `@` keeps the
# recipe command out of stdout.
test-json:
	@pkgs=$$(go list ./... | grep -v '/internal/sdk'); \
	env TF_ACC=1 \
	go test -json -cover -count 1 -timeout $(TEST_TIMEOUT) $$pkgs

unit-tests:
	# exclude the framework and sdkv2 packages, the
	# terraform-provider-hpe/morpheus package (but NOT its subpackages), and the
	# vendored SDK under internal/sdk (generated/third-party, no provider tests)
	pkgs=$$(go list ./... | grep -Ev 'sdkv2/(resources|datasources)|framework/(resources|datasources)|terraform-provider-hpe/morpheus$$|/internal/sdk'); \
	go test -v -count=1 -short -skip "TestAcc*" $$pkgs

collect-test-results:
	./scripts/collect-test-results.bash

build-render-tool:
	go build -o bin/render ./cmd/render

docs: build-render-tool
	go generate ./...
	cd tools && go generate

# -sweep-allow-failures keeps the run going when an individual sweeper fails.
# Without it the first failure aborts the whole sweep in map-iteration order,
# so unrelated resources are left on the appliance and the leak compounds.
sweep:
	env TF_ACC_SWEEP_PREFIX=$(SWEEP_PREFIX) \
	go test -v -tags sweep ./morpheus/testhelpers/sweep/... -sweep=$(SWEEP_SYSTEMS) -sweep-allow-failures $(SWEEP_RUN_ARGS)
