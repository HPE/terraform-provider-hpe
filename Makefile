#(C) Copyright 2025-2026 Hewlett Packard Enterprise Development LP
#
# Note: this Makefile works with GNUMake and BSDMake
#

.PHONY: build linter lint test test-json docs sweep build-render-tool

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
#
# The cluster package is the binding constraint: it runs four tests
# sequentially, each provisioning a three-worker cluster with an explicit 55m
# create timeout, so its worst case is 4 x 55m = 220m. At the previous 120m the
# package was killed before its own tests could time out, which replaced the
# per-test diagnostics with a bare `panic: test timed out`.
TEST_TIMEOUT ?= 240m

build:
	go build

linter:
	go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.11.3

# The vendored Morpheus SDK (internal/sdk/{oapigen,legacy}) is generated /
# hand-written third-party code that is not subject to the provider's lint
# rules. It is still type-checked as a dependency, but excluded from the lint
# target set (linting ~9k generated files is both wrong and prohibitively slow).
lint:
	golangci-lint run $$(go list -f '{{.Dir}}' ./... | grep -v '/internal/sdk')

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
