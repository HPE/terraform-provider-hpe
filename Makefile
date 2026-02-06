#(C) Copyright 2025 Hewlett Packard Enterprise Development LP
#
# Note: this Makefile works with GNUMake and BSDMake
#

.PHONY: build linter lint test docs sweep build-render-tool

build:
	go build

linter:
	go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.0.2

lint:
	golangci-lint run

test:
	env TF_ACC=1 \
	go test -short -v -cover -count 1 -timeout 10m ./...

testacc:
	cd morpheus/framework && \
	env TF_ACC=1 \
	go test -v -cover -count 1 -timeout 60m ./...

testsdkv2:
	cd morpheus/sdkv2 && \
	env TF_ACC=1 \
	go test -v -cover -count 1 -timeout 60m ./...

collect-test-results:
	./scripts/collect-test-results.bash

build-render-tool:
	go build -o bin/render ./cmd/render

docs: build-render-tool
	go generate ./...
	cd tools && go generate

sweep:
	go test -v ./morpheus/utils/test/sweep/... \
	  -sweep=all -sweep-run=hpe_morpheus_datastore,hpe_morpheus_instance,hpe_morpheus_network,hpe_morpheus_policy,hpe_morpheus_user
