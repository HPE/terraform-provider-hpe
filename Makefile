#(C) Copyright 2025 Hewlett Packard Enterprise Development LP
#
# Note: this Makefile works with GNUMake and BSDMake
#

.PHONY: build linter lint test docs sweep build-render-tool

# Usage: make sweep SWEEP=resource_name SWEEP_SYSTEMS=systemname
SWEEP ?= all
SWEEP_SYSTEMS ?= zodiac,feature
SWEEP_RUN_ARGS = $(if $(filter all,$(SWEEP)),,-sweep-run=$(SWEEP))

build:
	go build

linter:
	go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.0.2

lint:
	golangci-lint run

test:
	env TF_ACC=1 \
	go test -v -cover -count 1 -timeout 60m ./...

unit-tests:
	# exclude the framework and sdkv2 packages and the
	# terraform-provider-hpe/morpheus package (but NOT its subpackages)
	pkgs=$$(go list ./... | grep -Ev 'sdkv2/(resources|datasources)|framework/(resources|datasources)|terraform-provider-hpe/morpheus$$'); \
	go test -v -count=1 -short -skip "TestAcc*" $$pkgs

collect-test-results:
	./scripts/collect-test-results.bash

build-render-tool:
	go build -o bin/render ./cmd/render

docs: build-render-tool
	go generate ./...
	cd tools && go generate

sweep:
	# Find only resource packages with sweep_test.go and run sweepers in those
	pkgs=$$(find ./morpheus/framework/resources -name sweep_test.go -exec dirname {} \; | sort -u); \
	go test -v -run '^$$' $$pkgs -sweep=$(SWEEP_SYSTEMS) $(SWEEP_RUN_ARGS)
