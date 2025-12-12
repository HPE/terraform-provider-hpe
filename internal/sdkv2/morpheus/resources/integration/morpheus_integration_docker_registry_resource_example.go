// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package integration

//go:generate go run ../../../../../cmd/render -out examples/resources/morpheus_integration_docker_registry/resource.tf morpheus_integration_docker_registry_resource.tf.tmpl Name tfexampledockerregistry Enabled true Url https://index.docker.io/v1/ Username admin Password password123
