#(C) Copyright 2025 Hewlett Packard Enterprise Development LP
#
# Note: this Makefile works with GNUMake and BSDMake
#

.PHONY: build linter lint test docs docs-experimental experimental sweep

build:
	go build

experimental:
	go build -tags=experimental

linter:
	go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.0.2

lint:
	golangci-lint run

test:
	env TF_ACC=1 \
	go test -short -v -cover -count 1 -timeout 10m ./...

testacc:
	env TF_ACC=1 \
	go test -v -cover -count 1 -timeout 60m ./...

docs:
	go generate ./...
	cd tools && go generate

docs-experimental:
	rm -rf templates-combined-temp
	mkdir -p templates-combined-temp
	cp -r ./templates/* templates-combined-temp
	cp -r ./templates-experimental/* templates-combined-temp
	go generate -tags=experimental ./...
	cd tools && env GOFLAGS="-tags=experimental" go generate
	rm -rf templates-combined-temp

sweep:
	go test -v ./internal/subproviders/morpheus/test/sweep/... \
	  -sweep=all -sweep-run=hpe_morpheus_cloud,hpe_morpheus_datastore,hpe_morpheus_instance,hpe_morpheus_network,hpe_morpheus_user
