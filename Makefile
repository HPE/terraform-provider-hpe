#(C) Copyright 2025 Hewlett Packard Enterprise Development LP
#
# Note: this Makefile works with GNUMake and BSDMake
#

.PHONY: build linter lint test

build:
	go build ./cmd/...

linter:
	go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.0.2

lint:
	golangci-lint run

test:
	env TF_ACC=1 \
	go test -short -v -cover -count 1 -timeout 10m ./...

testacc:
	env TF_ACC=1 \
	go test -v -cover -count 1 -timeout 10m ./...

generate-docs:
	go generate ./...
	cd tools; go generate
